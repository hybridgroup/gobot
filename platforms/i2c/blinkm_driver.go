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
	b := &BlinkMDriver{
		Driver: gobot.Driver{
			Name:     name,
			Commands: make(map[string]func(map[string]interface{}) interface{}),
		},
		Adaptor: a,
	}

	b.Driver.AddCommand("FirmwareVersion", func(params map[string]interface{}) interface{} {
		return b.FirmwareVersion()
	})
	b.Driver.AddCommand("Color", func(params map[string]interface{}) interface{} {
		return b.Color()
	})
	b.Driver.AddCommand("Rgb", func(params map[string]interface{}) interface{} {
		red := byte(params["red"].(float64))
		green := byte(params["green"].(float64))
		blue := byte(params["blue"].(float64))
		b.Rgb(red, green, blue)
		return nil
	})
	b.Driver.AddCommand("Fade", func(params map[string]interface{}) interface{} {
		red := byte(params["red"].(float64))
		green := byte(params["green"].(float64))
		blue := byte(params["blue"].(float64))
		b.Fade(red, green, blue)
		return nil
	})

	return b
}

func (b *BlinkMDriver) Start() bool {
	b.Adaptor.I2cStart(0x09)
	b.Adaptor.I2cWrite([]byte("o"))
	b.Rgb(0, 0, 0)
	return true
}
func (b *BlinkMDriver) Init() bool { return true }
func (b *BlinkMDriver) Halt() bool { return true }

func (b *BlinkMDriver) Rgb(red byte, green byte, blue byte) {
	b.Adaptor.I2cWrite([]byte("n"))
	b.Adaptor.I2cWrite([]byte{red, green, blue})
}

func (b *BlinkMDriver) Fade(red byte, green byte, blue byte) {
	b.Adaptor.I2cWrite([]byte("c"))
	b.Adaptor.I2cWrite([]byte{red, green, blue})
}

func (b *BlinkMDriver) FirmwareVersion() string {
	b.Adaptor.I2cWrite([]byte("Z"))
	data := b.Adaptor.I2cRead(2)
	if len(data) != 2 {
		return ""
	}
	return fmt.Sprintf("%v.%v", data[0], data[1])
}

func (b *BlinkMDriver) Color() []byte {
	b.Adaptor.I2cWrite([]byte("g"))
	data := b.Adaptor.I2cRead(3)
	if len(data) != 3 {
		return make([]byte, 0)
	}
	return []byte{data[0], data[1], data[2]}
}
