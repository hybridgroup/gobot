package neurosky

import (
	"bytes"
	"github.com/hybridgroup/gobot"
)

const BT_SYNC byte = 0xAA
const CODE_EX byte = 0x55             // Extended code
const CODE_SIGNAL_QUALITY byte = 0x02 // POOR_SIGNAL quality 0-255
const CODE_ATTENTION byte = 0x04      // ATTENTION eSense 0-100
const CODE_MEDITATION byte = 0x05     // MEDITATION eSense 0-100
const CODE_BLINK byte = 0x16          // BLINK strength 0-255
const CODE_WAVE byte = 0x80           // RAW wave value: 2-byte big-endian 2s-complement
const CODE_ASIC_EEG byte = 0x83       // ASIC EEG POWER 8 3-byte big-endian integers

type NeuroskyDriver struct {
	gobot.Driver
	Adaptor *NeuroskyAdaptor
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

func NewNeuroskyDriver(a *NeuroskyAdaptor) *NeuroskyDriver {
	return &NeuroskyDriver{
		Driver: gobot.Driver{
			Events: map[string]chan interface{}{
				"Extended":   make(chan interface{}),
				"Signal":     make(chan interface{}),
				"Attention":  make(chan interface{}),
				"Meditation": make(chan interface{}),
				"Blink":      make(chan interface{}),
				"Wave":       make(chan interface{}),
				"EEG":        make(chan interface{}),
			},
		},
		Adaptor: a,
	}
}

func (n *NeuroskyDriver) Init() bool { return true }
func (n *NeuroskyDriver) Start() bool {
	go func() {
		for {
			var buff = make([]byte, int(2048))
			_, err := n.Adaptor.sp.Read(buff[:])
			if err != nil {
				panic(err)
			} else {
				n.parse(bytes.NewBuffer(buff))
			}
		}
	}()
	return true
}
func (n *NeuroskyDriver) Halt() bool { return true }

func (n *NeuroskyDriver) parse(buf *bytes.Buffer) {
	for buf.Len() > 2 {
		b1, _ := buf.ReadByte()
		b2, _ := buf.ReadByte()
		if b1 == BT_SYNC && b2 == BT_SYNC {
			length, _ := buf.ReadByte()
			var payload = make([]byte, int(length))
			buf.Read(payload)
			//checksum, _ := buf.ReadByte()
			buf.Next(1)
			n.parsePacket(payload)
		}
	}
}

func (me *NeuroskyDriver) parsePacket(data []byte) {
	buf := bytes.NewBuffer(data)
	for buf.Len() > 0 {
		b, _ := buf.ReadByte()
		switch b {
		case CODE_EX:
			gobot.Publish(me.Events["Extended"], nil)
		case CODE_SIGNAL_QUALITY:
			ret, _ := buf.ReadByte()
			gobot.Publish(me.Events["Signal"], ret)
		case CODE_ATTENTION:
			ret, _ := buf.ReadByte()
			gobot.Publish(me.Events["Attention"], ret)
		case CODE_MEDITATION:
			ret, _ := buf.ReadByte()
			gobot.Publish(me.Events["Meditation"], ret)
		case CODE_BLINK:
			ret, _ := buf.ReadByte()
			gobot.Publish(me.Events["Blink"], ret)
		case CODE_WAVE:
			buf.Next(1)
			var ret = make([]byte, 2)
			buf.Read(ret)
			gobot.Publish(me.Events["Wave"], ret)
		case CODE_ASIC_EEG:
			var ret = make([]byte, 25)
			n, _ := buf.Read(ret)
			if n == 25 {
				gobot.Publish(me.Events["EEG"], me.parseEEG(ret))
			}
		}
	}
}

func (me *NeuroskyDriver) parseEEG(data []byte) EEG {
	return EEG{
		Delta:    me.parse3ByteInteger(data[0:3]),
		Theta:    me.parse3ByteInteger(data[3:6]),
		LoAlpha:  me.parse3ByteInteger(data[6:9]),
		HiAlpha:  me.parse3ByteInteger(data[9:12]),
		LoBeta:   me.parse3ByteInteger(data[12:15]),
		HiBeta:   me.parse3ByteInteger(data[15:18]),
		LoGamma:  me.parse3ByteInteger(data[18:21]),
		MidGamma: me.parse3ByteInteger(data[21:25]),
	}
}

func (me *NeuroskyDriver) parse3ByteInteger(data []byte) int {
	return ((int(data[0]) << 16) | (((1 << 16) - 1) & (int(data[1]) << 8)) | (((1 << 8) - 1) & int(data[2])))
}
