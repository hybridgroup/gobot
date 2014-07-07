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
		Driver: *gobot.NewDriver(
			name,
			"BlinkMDriver",
			a.(gobot.AdaptorInterface),
		),
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

func (b *BlinkMDriver) adaptor() I2cInterface {
	return b.Driver.Adaptor().(I2cInterface)
}

func (b *BlinkMDriver) Start() bool {
	b.adaptor().I2cStart(0x09)
	b.adaptor().I2cWrite([]byte("o"))
	b.Rgb(0, 0, 0)
	return true
}
func (b *BlinkMDriver) Init() bool { return true }
func (b *BlinkMDriver) Halt() bool { return true }

func (b *BlinkMDriver) Rgb(red byte, green byte, blue byte) {
	b.adaptor().I2cWrite([]byte("n"))
	b.adaptor().I2cWrite([]byte{red, green, blue})
}

func (b *BlinkMDriver) Fade(red byte, green byte, blue byte) {
	b.adaptor().I2cWrite([]byte("c"))
	b.adaptor().I2cWrite([]byte{red, green, blue})
}

func (b *BlinkMDriver) FirmwareVersion() string {
	b.adaptor().I2cWrite([]byte("Z"))
	data := b.adaptor().I2cRead(2)
	if len(data) != 2 {
		return ""
	}
	return fmt.Sprintf("%v.%v", data[0], data[1])
}

func (b *BlinkMDriver) Color() []byte {
	b.adaptor().I2cWrite([]byte("g"))
	data := b.adaptor().I2cRead(3)
	if len(data) != 3 {
		return make([]byte, 0)
	}
	return []byte{data[0], data[1], data[2]}
}
