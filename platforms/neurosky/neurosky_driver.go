package neurosky

import (
	"bytes"

	"gobot.io/x/gobot"
)

const (
	// BTSync is the sync code
	BTSync byte = 0xAA

	// CodeEx Extended code
	CodeEx byte = 0x55

	// CodeSignalQuality POOR_SIGNAL quality 0-255
	CodeSignalQuality byte = 0x02

	// CodeAttention ATTENTION eSense 0-100
	CodeAttention byte = 0x04

	// CodeMeditation MEDITATION eSense 0-100
	CodeMeditation byte = 0x05

	// CodeBlink BLINK strength 0-255
	CodeBlink byte = 0x16

	// CodeWave RAW wave value: 2-byte big-endian 2s-complement
	CodeWave byte = 0x80

	// CodeAsicEEG ASIC EEG POWER 8 3-byte big-endian integers
	CodeAsicEEG byte = 0x83

	// Extended event
	Extended = "extended"

	// Signal event
	Signal = "signal"

	// Attention event
	Attention = "attention"

	// Meditation event
	Meditation = "meditation"

	// Blink event
	Blink = "blink"

	// Wave event
	Wave = "wave"

	// EEG event
	EEG = "eeg"

	// Error event
	Error = "error"
)

// Driver is the Gobot Driver for the Mindwave
type Driver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}

// EEGData is the EEG raw data returned from the Mindwave
type EEGData struct {
	Delta    int
	Theta    int
	LoAlpha  int
	HiAlpha  int
	LoBeta   int
	HiBeta   int
	LoGamma  int
	MidGamma int
}

// NewDriver creates a Neurosky Driver
// and adds the following events:
//
//   extended - user's current extended level
//   signal - shows signal strength
//   attention - user's current attention level
//   meditation - user's current meditation level
//   blink - user's current blink level
//   wave - shows wave data
//   eeg - showing eeg data
func NewDriver(a *Adaptor) *Driver {
	n := &Driver{
		name:       "Neurosky",
		connection: a,
		Eventer:    gobot.NewEventer(),
	}

	n.AddEvent(Extended)
	n.AddEvent(Signal)
	n.AddEvent(Attention)
	n.AddEvent(Meditation)
	n.AddEvent(Blink)
	n.AddEvent(Wave)
	n.AddEvent(EEG)
	n.AddEvent(Error)

	return n
}

// Connection returns the Driver's connection
func (n *Driver) Connection() gobot.Connection { return n.connection }

// Name returns the Driver name
func (n *Driver) Name() string { return n.name }

// SetName sets the Driver name
func (n *Driver) SetName(name string) { n.name = name }

// adaptor returns neurosky adaptor
func (n *Driver) adaptor() *Adaptor {
	return n.Connection().(*Adaptor)
}

// Start creates a go routine to listen from serial port
// and parse buffer readings
func (n *Driver) Start() (err error) {
	go func() {
		for {
			buff := make([]byte, 1024)
			_, err := n.adaptor().sp.Read(buff[:])
			if err != nil {
				n.Publish(n.Event("error"), err)
			} else {
				n.parse(bytes.NewBuffer(buff))
			}
		}
	}()
	return
}

// Halt stops neurosky driver (void)
func (n *Driver) Halt() (err error) { return }

// parse converts bytes buffer into packets until no more data is present
func (n *Driver) parse(buf *bytes.Buffer) {
	for buf.Len() > 2 {
		b1, _ := buf.ReadByte()
		b2, _ := buf.ReadByte()
		if b1 == BTSync && b2 == BTSync {
			length, _ := buf.ReadByte()
			payload := make([]byte, length)
			buf.Read(payload)
			//checksum, _ := buf.ReadByte()
			buf.Next(1)
			n.parsePacket(bytes.NewBuffer(payload))
		}
	}
}

// parsePacket publishes event according to data parsed
func (n *Driver) parsePacket(buf *bytes.Buffer) {
	for buf.Len() > 0 {
		b, _ := buf.ReadByte()
		switch b {
		case CodeEx:
			n.Publish(n.Event("extended"), nil)
		case CodeSignalQuality:
			ret, _ := buf.ReadByte()
			n.Publish(n.Event("signal"), ret)
		case CodeAttention:
			ret, _ := buf.ReadByte()
			n.Publish(n.Event("attention"), ret)
		case CodeMeditation:
			ret, _ := buf.ReadByte()
			n.Publish(n.Event("meditation"), ret)
		case CodeBlink:
			ret, _ := buf.ReadByte()
			n.Publish(n.Event("blink"), ret)
		case CodeWave:
			buf.Next(1)
			var ret = make([]byte, 2)
			buf.Read(ret)
			n.Publish(n.Event("wave"), int16(ret[0])<<8|int16(ret[1]))
		case CodeAsicEEG:
			ret := make([]byte, 25)
			i, _ := buf.Read(ret)
			if i == 25 {
				n.Publish(n.Event("eeg"), n.parseEEG(ret))
			}
		}
	}
}

// parseEEG returns data converted into EEG map
func (n *Driver) parseEEG(data []byte) EEGData {
	return EEGData{
		Delta:    n.parse3ByteInteger(data[0:3]),
		Theta:    n.parse3ByteInteger(data[3:6]),
		LoAlpha:  n.parse3ByteInteger(data[6:9]),
		HiAlpha:  n.parse3ByteInteger(data[9:12]),
		LoBeta:   n.parse3ByteInteger(data[12:15]),
		HiBeta:   n.parse3ByteInteger(data[15:18]),
		LoGamma:  n.parse3ByteInteger(data[18:21]),
		MidGamma: n.parse3ByteInteger(data[21:25]),
	}
}

func (n *Driver) parse3ByteInteger(data []byte) int {
	return ((int(data[0]) << 16) |
		(((1 << 16) - 1) & (int(data[1]) << 8)) |
		(((1 << 8) - 1) & int(data[2])))
}
