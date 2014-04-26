package gobotNeurosky

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

type NeuroskyInterface interface {
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

func NewNeurosky(adaptor *NeuroskyAdaptor) *NeuroskyDriver {
	d := new(NeuroskyDriver)
	d.Events = make(map[string]chan interface{})
	d.Events["Extended"] = make(chan interface{})
	d.Events["Signal"] = make(chan interface{})
	d.Events["Attention"] = make(chan interface{})
	d.Events["Meditation"] = make(chan interface{})
	d.Events["Blink"] = make(chan interface{})
	d.Events["Wave"] = make(chan interface{})
	d.Events["EEG"] = make(chan interface{})
	d.Adaptor = adaptor
	d.Commands = []string{}
	return d
}

func (me *NeuroskyDriver) Init() bool { return true }
func (me *NeuroskyDriver) Start() bool {
	go func() {
		for {
			var buff = make([]byte, int(2048))
			_, err := me.Adaptor.sp.Read(buff[:])
			if err != nil {
				panic(err)
			} else {
				me.parse(bytes.NewBuffer(buff))
			}
		}
	}()
	return true
}
func (me *NeuroskyDriver) Halt() bool { return true }

func (me *NeuroskyDriver) parse(buf *bytes.Buffer) {
	for buf.Len() > 2 {
		b1, _ := buf.ReadByte()
		b2, _ := buf.ReadByte()
		if b1 == BT_SYNC && b2 == BT_SYNC {
			length, _ := buf.ReadByte()
			var payload = make([]byte, int(length))
			buf.Read(payload)
			//checksum, _ := buf.ReadByte()
			buf.Next(1)
			me.parsePacket(payload)
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
	eeg := EEG{}
	eeg.Delta = me.parse3ByteInteger(data[0:3])
	eeg.Theta = me.parse3ByteInteger(data[3:6])
	eeg.LoAlpha = me.parse3ByteInteger(data[6:9])
	eeg.HiAlpha = me.parse3ByteInteger(data[9:12])
	eeg.LoBeta = me.parse3ByteInteger(data[12:15])
	eeg.HiBeta = me.parse3ByteInteger(data[15:18])
	eeg.LoGamma = me.parse3ByteInteger(data[18:21])
	eeg.MidGamma = me.parse3ByteInteger(data[21:25])
	return eeg
}

func (me *NeuroskyDriver) parse3ByteInteger(data []byte) int {
	b1 := int(data[0])
	b2 := int(data[1])
	b3 := int(data[2])
	bigEndianInteger := ((b1 << 16) | (((1 << 16) - 1) & (b2 << 8)) | ((1<<8)-1)&b3)
	return bigEndianInteger
}
