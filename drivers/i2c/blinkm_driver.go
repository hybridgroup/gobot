package i2c

import (
	"fmt"

	"gobot.io/x/gobot"
)

const blinkmAddress = 0x09

// BlinkMDriver is a Gobot Driver for a BlinkM LED
type BlinkMDriver struct {
	name       string
	connection I2c
	gobot.Commander
}

// NewBlinkMDriver creates a new BlinkMDriver.
//
// Adds the following API commands:
//	Rgb - sets RGB color
//	Fade - fades the RGB color
//	FirmwareVersion - returns the version of the current Frimware
//	Color - returns the color of the LED.
func NewBlinkMDriver(a I2c) *BlinkMDriver {
	b := &BlinkMDriver{
		name:       "BlinkM",
		connection: a,
		Commander:  gobot.NewCommander(),
	}

	b.AddCommand("Rgb", func(params map[string]interface{}) interface{} {
		red := byte(params["red"].(float64))
		green := byte(params["green"].(float64))
		blue := byte(params["blue"].(float64))
		return b.Rgb(red, green, blue)
	})
	b.AddCommand("Fade", func(params map[string]interface{}) interface{} {
		red := byte(params["red"].(float64))
		green := byte(params["green"].(float64))
		blue := byte(params["blue"].(float64))
		return b.Fade(red, green, blue)
	})
	b.AddCommand("FirmwareVersion", func(params map[string]interface{}) interface{} {
		version, err := b.FirmwareVersion()
		return map[string]interface{}{"version": version, "err": err}
	})
	b.AddCommand("Color", func(params map[string]interface{}) interface{} {
		color, err := b.Color()
		return map[string]interface{}{"color": color, "err": err}
	})

	return b
}

// Name returns the Name for the Driver
func (b *BlinkMDriver) Name() string                 { return b.name }

// SetName sets the Name for the Driver
func (b *BlinkMDriver) SetName(n string)             { b.name = n }

// Connection returns the connection for the Driver
func (b *BlinkMDriver) Connection() gobot.Connection { return b.connection.(gobot.Connection) }

// Start starts the Driver up, and writes start command
func (b *BlinkMDriver) Start() (err error) {
	if err := b.connection.I2cStart(blinkmAddress); err != nil {
		return err
	}
	if err := b.connection.I2cWrite(blinkmAddress, []byte("o")); err != nil {
		return err
	}
	return
}

// Halt returns true if device is halted successfully
func (b *BlinkMDriver) Halt() (err error) { return }

// Rgb sets color using r,g,b params
func (b *BlinkMDriver) Rgb(red byte, green byte, blue byte) (err error) {
	if err = b.connection.I2cWrite(blinkmAddress, []byte("n")); err != nil {
		return
	}
	err = b.connection.I2cWrite(blinkmAddress, []byte{red, green, blue})
	return
}

// Fade removes color using r,g,b params
func (b *BlinkMDriver) Fade(red byte, green byte, blue byte) (err error) {
	if err = b.connection.I2cWrite(blinkmAddress, []byte("c")); err != nil {
		return
	}
	err = b.connection.I2cWrite(blinkmAddress, []byte{red, green, blue})
	return
}

// FirmwareVersion returns version with MAYOR.minor format
func (b *BlinkMDriver) FirmwareVersion() (version string, err error) {
	if err = b.connection.I2cWrite(blinkmAddress, []byte("Z")); err != nil {
		return
	}
	data, err := b.connection.I2cRead(blinkmAddress, 2)
	if len(data) != 2 || err != nil {
		return
	}
	return fmt.Sprintf("%v.%v", data[0], data[1]), nil
}

// Color returns an array with current rgb color
func (b *BlinkMDriver) Color() (color []byte, err error) {
	if err = b.connection.I2cWrite(blinkmAddress, []byte("g")); err != nil {
		return
	}
	data, err := b.connection.I2cRead(blinkmAddress, 3)
	if len(data) != 3 || err != nil {
		return []byte{}, err
	}
	return []byte{data[0], data[1], data[2]}, nil
}
