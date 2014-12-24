package neurosky

import (
	"bytes"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*NeuroskyDriver)(nil)

const BTSync byte = 0xAA

// Extended code
const CodeEx byte = 0x55

// POOR_SIGNAL quality 0-255
const CodeSignalQuality byte = 0x02

// ATTENTION eSense 0-100
const CodeAttention byte = 0x04

// MEDITATION eSense 0-100
const CodeMeditation byte = 0x05

// BLINK strength 0-255
const CodeBlink byte = 0x16

// RAW wave value: 2-byte big-endian 2s-complement
const CodeWave byte = 0x80

// ASIC EEG POWER 8 3-byte big-endian integers
const CodeAsicEEG byte = 0x83

type NeuroskyDriver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}

type EEG struct {
	Delta    int
	Theta    int
	LoAlpha  int
	HiAlpha  int
	LoBeta   int
	HiBeta   int
	LoGamma  int
	MidGamma int
}

// NewNeuroskyDriver creates a NeuroskyDriver by name
// and adds the following events:
//
//   extended - user's current extended level
//   signal - shows signal strength
//   attention - user's current attention level
//   meditation - user's current meditation level
//   blink - user's current blink level
//   wave - shows wave data
//   eeg - showing eeg data
func NewNeuroskyDriver(a *NeuroskyAdaptor, name string) *NeuroskyDriver {
	n := &NeuroskyDriver{
		name:       name,
		connection: a,
		Eventer:    gobot.NewEventer(),
	}

	n.AddEvent("extended")
	n.AddEvent("signal")
	n.AddEvent("attention")
	n.AddEvent("meditation")
	n.AddEvent("blink")
	n.AddEvent("wave")
	n.AddEvent("eeg")
	n.AddEvent("error")

	return n
}
func (n *NeuroskyDriver) Connection() gobot.Connection { return n.connection }
func (n *NeuroskyDriver) Name() string                 { return n.name }

// adaptor returns neurosky adaptor
func (n *NeuroskyDriver) adaptor() *NeuroskyAdaptor {
	return n.Connection().(*NeuroskyAdaptor)
}

// Start creates a go routine to listen from serial port
// and parse buffer readings
func (n *NeuroskyDriver) Start() (errs []error) {
	go func() {
		for {
			buff := make([]byte, 1024)
			_, err := n.adaptor().sp.Read(buff[:])
			if err != nil {
				gobot.Publish(n.Event("error"), err)
			} else {
				n.parse(bytes.NewBuffer(buff))
			}
		}
	}()
	return
}

// Halt stops neurosky driver (void)
func (n *NeuroskyDriver) Halt() (errs []error) { return }

// parse converts bytes buffer into packets until no more data is present
func (n *NeuroskyDriver) parse(buf *bytes.Buffer) {
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
func (n *NeuroskyDriver) parsePacket(buf *bytes.Buffer) {
	for buf.Len() > 0 {
		b, _ := buf.ReadByte()
		switch b {
		case CodeEx:
			gobot.Publish(n.Event("extended"), nil)
		case CodeSignalQuality:
			ret, _ := buf.ReadByte()
			gobot.Publish(n.Event("signal"), ret)
		case CodeAttention:
			ret, _ := buf.ReadByte()
			gobot.Publish(n.Event("attention"), ret)
		case CodeMeditation:
			ret, _ := buf.ReadByte()
			gobot.Publish(n.Event("meditation"), ret)
		case CodeBlink:
			ret, _ := buf.ReadByte()
			gobot.Publish(n.Event("blink"), ret)
		case CodeWave:
			buf.Next(1)
			var ret = make([]byte, 2)
			buf.Read(ret)
			gobot.Publish(n.Event("wave"), int16(ret[0])<<8|int16(ret[1]))
		case CodeAsicEEG:
			ret := make([]byte, 25)
			i, _ := buf.Read(ret)
			if i == 25 {
				gobot.Publish(n.Event("eeg"), n.parseEEG(ret))
			}
		}
	}
}

// parseEEG returns data converted into EEG map
func (n *NeuroskyDriver) parseEEG(data []byte) EEG {
	return EEG{
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

func (n *NeuroskyDriver) parse3ByteInteger(data []byte) int {
	return ((int(data[0]) << 16) |
		(((1 << 16) - 1) & (int(data[1]) << 8)) |
		(((1 << 8) - 1) & int(data[2])))
}
