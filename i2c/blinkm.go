package gobotI2C

import (
	"fmt"
	"github.com/hybridgroup/gobot"
)

type BlinkM struct {
	gobot.Driver
	Adaptor I2cInterface
}

func NewBlinkM(a I2cInterface) *BlinkM {
	w := new(BlinkM)
	w.Adaptor = a
	w.Events = make(map[string]chan interface{})
	w.Commands = []string{
		"RgbC",
		"FadeC",
		"ColorC",
		"FirmwareVersionC",
	}
	return w
}

func (self *BlinkM) Start() bool {
	self.Adaptor.I2cStart(0x09)
	self.Adaptor.I2cWrite([]uint16{uint16([]byte("o")[0])})
	self.Rgb(0, 0, 0)
	return true
}
func (self *BlinkM) Init() bool { return true }
func (self *BlinkM) Halt() bool { return true }

func (self *BlinkM) Rgb(r byte, g byte, b byte) {
	self.Adaptor.I2cWrite([]uint16{uint16([]byte("n")[0])})
	self.Adaptor.I2cWrite([]uint16{uint16(r), uint16(g), uint16(b)})
}

func (self *BlinkM) Fade(r byte, g byte, b byte) {
	self.Adaptor.I2cWrite([]uint16{uint16([]byte("c")[0])})
	self.Adaptor.I2cWrite([]uint16{uint16(r), uint16(g), uint16(b)})
}

func (self *BlinkM) FirmwareVersion() string {
	self.Adaptor.I2cWrite([]uint16{uint16([]byte("Z")[0])})
	data := self.Adaptor.I2cRead(uint16(2))
	if len(data) != 2 {
		return ""
	}
	return fmt.Sprintf("%v.%v", data[0], data[1])
}

func (self *BlinkM) Color() []byte {
	self.Adaptor.I2cWrite([]uint16{uint16([]byte("g")[0])})
	data := self.Adaptor.I2cRead(uint16(3))
	if len(data) != 3 {
		return make([]byte, 0)
	}
	return []byte{byte(data[0]), byte(data[1]), byte(data[2])}
}
