package i2c

import (
	"fmt"

	"github.com/hybridgroup/gobot"
)

var _ gobot.DriverInterface = (*BlinkMDriver)(nil)

type BlinkMDriver struct {
	gobot.Driver
}

// NewBlinkMDriver creates a new BlinkMDriver with specified name.
//
// Adds the following API commands:
//	Rgb - sets RGB color
//	Fade - fades the RGB color
//	FirmwareVersion - returns the version of the current Frimware
//	Color - returns the color of the LED.
func NewBlinkMDriver(a I2cInterface, name string) *BlinkMDriver {
	b := &BlinkMDriver{
		Driver: *gobot.NewDriver(
			name,
			"BlinkMDriver",
			a.(gobot.AdaptorInterface),
		),
	}

	b.AddCommand("Rgb", func(params map[string]interface{}) interface{} {
		red := byte(params["red"].(float64))
		green := byte(params["green"].(float64))
		blue := byte(params["blue"].(float64))
		b.Rgb(red, green, blue)
		return nil
	})
	b.AddCommand("Fade", func(params map[string]interface{}) interface{} {
		red := byte(params["red"].(float64))
		green := byte(params["green"].(float64))
		blue := byte(params["blue"].(float64))
		b.Fade(red, green, blue)
		return nil
	})
	b.AddCommand("FirmwareVersion", func(params map[string]interface{}) interface{} {
		return b.FirmwareVersion()
	})
	b.AddCommand("Color", func(params map[string]interface{}) interface{} {
		return b.Color()
	})

	return b
}

// adaptor returns I2C adaptor
func (b *BlinkMDriver) adaptor() I2cInterface {
	return b.Adaptor().(I2cInterface)
}

// Start writes start bytes and resets color
func (b *BlinkMDriver) Start() error {
	b.adaptor().I2cStart(0x09)
	b.adaptor().I2cWrite([]byte("o"))
	b.Rgb(0, 0, 0)
	return nil
}

// Halt returns true if device is halted successfully
func (b *BlinkMDriver) Halt() error { return nil }

// Rgb sets color using r,g,b params
func (b *BlinkMDriver) Rgb(red byte, green byte, blue byte) {
	b.adaptor().I2cWrite([]byte("n"))
	b.adaptor().I2cWrite([]byte{red, green, blue})
}

// Fade removes color using r,g,b params
func (b *BlinkMDriver) Fade(red byte, green byte, blue byte) {
	b.adaptor().I2cWrite([]byte("c"))
	b.adaptor().I2cWrite([]byte{red, green, blue})
}

// FirmwareVersion returns version with MAYOR.minor format
func (b *BlinkMDriver) FirmwareVersion() string {
	b.adaptor().I2cWrite([]byte("Z"))
	data := b.adaptor().I2cRead(2)
	if len(data) != 2 {
		return ""
	}
	return fmt.Sprintf("%v.%v", data[0], data[1])
}

// Color returns an array with current rgb color
func (b *BlinkMDriver) Color() []byte {
	b.adaptor().I2cWrite([]byte("g"))
	data := b.adaptor().I2cRead(3)
	if len(data) != 3 {
		return make([]byte, 0)
	}
	return []byte{data[0], data[1], data[2]}
}
