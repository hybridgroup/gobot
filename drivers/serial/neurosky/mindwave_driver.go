package neurosky

import (
	"bytes"
	"log"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/serial"
)

type mindWaveSerialAdaptor interface {
	gobot.Adaptor
	serial.SerialReader
}

const (
	// mindWaveBTSync is the sync code
	mindWaveBTSync byte = 0xAA
	// mindWaveCodeEx Extended code
	mindWaveCodeEx byte = 0x55
	// mindWaveCodeSignalQuality POOR_SIGNAL quality 0-255
	mindWaveCodeSignalQuality byte = 0x02
	// mindWaveCodeAttention ATTENTION eSense 0-100
	mindWaveCodeAttention byte = 0x04
	// mindWaveCodeMeditation MEDITATION eSense 0-100
	mindWaveCodeMeditation byte = 0x05
	// mindWaveCodeBlink BLINK strength 0-255
	mindWaveCodeBlink byte = 0x16
	// mindWaveCodeWave RAW wave value: 2-byte big-endian 2s-complement
	mindWaveCodeWave byte = 0x80
	// mindWaveCodeAsicEEG ASIC EEG POWER 8 3-byte big-endian integers
	mindWaveCodeAsicEEG byte = 0x83

	ExtendedEvent   = "extended"
	SignalEvent     = "signal"
	AttentionEvent  = "attention"
	MeditationEvent = "meditation"
	BlinkEvent      = "blink"
	WaveEvent       = "wave"
	EEGEvent        = "eeg"
	ErrorEvent      = "error"
)

// MindWaveDriver is the Gobot driver for the Neurosky MindWave Sensor
type MindWaveDriver struct {
	*serial.Driver
	gobot.Eventer
}

// MindWaveEEGData is the EEG raw data returned from the sensor
type MindWaveEEGData struct {
	Delta    int
	Theta    int
	LoAlpha  int
	HiAlpha  int
	LoBeta   int
	HiBeta   int
	LoGamma  int
	MidGamma int
}

// NewMindWaveDriver creates a driver for Neurosky MindWave
// and adds the following events:
//
//	extended - user's current extended level
//	signal - shows signal strength
//	attention - user's current attention level
//	meditation - user's current meditation level
//	blink - user's current blink level
//	wave - shows wave data
//	eeg - showing eeg data
func NewMindWaveDriver(a mindWaveSerialAdaptor, opts ...serial.OptionApplier) *MindWaveDriver {
	d := &MindWaveDriver{
		Eventer: gobot.NewEventer(),
	}
	d.Driver = serial.NewDriver(a, "MindWave", d.initialize, nil, opts...)

	d.AddEvent(ExtendedEvent)
	d.AddEvent(SignalEvent)
	d.AddEvent(AttentionEvent)
	d.AddEvent(MeditationEvent)
	d.AddEvent(BlinkEvent)
	d.AddEvent(WaveEvent)
	d.AddEvent(EEGEvent)
	d.AddEvent(ErrorEvent)

	return d
}

// initialize creates a go routine to listen from serial port and parse buffer readings
// TODO: stop the go routine gracefully on Halt()
func (d *MindWaveDriver) initialize() error {
	go func() {
		for {
			buff := make([]byte, 1024)
			_, err := d.adaptor().SerialRead(buff)
			if err != nil {
				d.Publish(d.Event("error"), err)
			} else {
				if err := d.parse(bytes.NewBuffer(buff)); err != nil {
					panic(err)
				}
			}
		}
	}()
	return nil
}

// parse converts bytes buffer into packets until no more data is present
func (d *MindWaveDriver) parse(buf *bytes.Buffer) error {
	for buf.Len() > 2 {
		b1, _ := buf.ReadByte()
		b2, _ := buf.ReadByte()
		if b1 == mindWaveBTSync && b2 == mindWaveBTSync {
			length, _ := buf.ReadByte()
			payload := make([]byte, length)
			if _, err := buf.Read(payload); err != nil {
				return err
			}
			// checksum, _ := buf.ReadByte()
			buf.Next(1)
			if err := d.parsePacket(bytes.NewBuffer(payload)); err != nil {
				panic(err)
			}
		}
	}

	return nil
}

// parsePacket publishes event according to data parsed
func (d *MindWaveDriver) parsePacket(buf *bytes.Buffer) error {
	for buf.Len() > 0 {
		b, _ := buf.ReadByte()
		switch b {
		case mindWaveCodeEx:
			d.Publish(d.Event(ExtendedEvent), nil)
		case mindWaveCodeSignalQuality:
			ret, _ := buf.ReadByte()
			d.Publish(d.Event(SignalEvent), ret)
		case mindWaveCodeAttention:
			ret, _ := buf.ReadByte()
			d.Publish(d.Event(AttentionEvent), ret)
		case mindWaveCodeMeditation:
			ret, _ := buf.ReadByte()
			d.Publish(d.Event(MeditationEvent), ret)
		case mindWaveCodeBlink:
			ret, _ := buf.ReadByte()
			d.Publish(d.Event(BlinkEvent), ret)
		case mindWaveCodeWave:
			buf.Next(1)
			ret := make([]byte, 2)
			if _, err := buf.Read(ret); err != nil {
				return err
			}
			d.Publish(d.Event(WaveEvent), int16(ret[0])<<8|int16(ret[1]))
		case mindWaveCodeAsicEEG:
			ret := make([]byte, 25)
			i, _ := buf.Read(ret)
			if i == 25 {
				d.Publish(d.Event(EEGEvent), d.parseEEG(ret))
			}
		}
	}

	return nil
}

// parseEEG returns data converted into EEG map
func (d *MindWaveDriver) parseEEG(data []byte) MindWaveEEGData {
	return MindWaveEEGData{
		Delta:    d.parse3ByteInteger(data[0:3]),
		Theta:    d.parse3ByteInteger(data[3:6]),
		LoAlpha:  d.parse3ByteInteger(data[6:9]),
		HiAlpha:  d.parse3ByteInteger(data[9:12]),
		LoBeta:   d.parse3ByteInteger(data[12:15]),
		HiBeta:   d.parse3ByteInteger(data[15:18]),
		LoGamma:  d.parse3ByteInteger(data[18:21]),
		MidGamma: d.parse3ByteInteger(data[21:25]),
	}
}

func (d *MindWaveDriver) parse3ByteInteger(data []byte) int {
	return ((int(data[0]) << 16) |
		(((1 << 16) - 1) & (int(data[1]) << 8)) |
		(((1 << 8) - 1) & int(data[2])))
}

func (d *MindWaveDriver) adaptor() mindWaveSerialAdaptor {
	if a, ok := d.Connection().(mindWaveSerialAdaptor); ok {
		return a
	}

	log.Printf("%s has no Neurosky serial connector\n", d.Name())
	return nil
}
