package i2c

import (
	"fmt"
	"github.com/hybridgroup/gobot"
)

type BlinkMDriver struct {
	gobot.Driver
	Adaptor I2cInterface
}

func NewBlinkMDriver(a I2cInterface, name string) *BlinkMDriver {
	return &BlinkMDriver{
		Driver: gobot.Driver{
			Name: name,
			Commands: []string{
				"RgbC",
				"FadeC",
				"ColorC",
				"FirmwareVersionC",
			},
		},
		Adaptor: a,
	}
}

func (b *BlinkMDriver) Start() bool {
	b.Adaptor.I2cStart(0x09)
	b.Adaptor.I2cWrite([]uint16{uint16([]byte("o")[0])})
	b.Rgb(0, 0, 0)
	return true
}
func (b *BlinkMDriver) Init() bool { return true }
func (b *BlinkMDriver) Halt() bool { return true }

func (b *BlinkMDriver) Rgb(red byte, green byte, blue byte) {
	b.Adaptor.I2cWrite([]uint16{uint16([]byte("n")[0])})
	b.Adaptor.I2cWrite([]uint16{uint16(red), uint16(green), uint16(blue)})
}

func (b *BlinkMDriver) Fade(red byte, green byte, blue byte) {
	b.Adaptor.I2cWrite([]uint16{uint16([]byte("c")[0])})
	b.Adaptor.I2cWrite([]uint16{uint16(red), uint16(green), uint16(blue)})
}

func (b *BlinkMDriver) FirmwareVersion() string {
	b.Adaptor.I2cWrite([]uint16{uint16([]byte("Z")[0])})
	data := b.Adaptor.I2cRead(uint16(2))
	if len(data) != 2 {
		return ""
	}
	return fmt.Sprintf("%v.%v", data[0], data[1])
}

func (b *BlinkMDriver) Color() []byte {
	b.Adaptor.I2cWrite([]uint16{uint16([]byte("g")[0])})
	data := b.Adaptor.I2cRead(uint16(3))
	if len(data) != 3 {
		return make([]byte, 0)
	}
	return []byte{byte(data[0]), byte(data[1]), byte(data[2])}
}
