package i2c

import (
	"fmt"
)

const blinkmDefaultAddress = 0x09

// BlinkMDriver is a Gobot Driver for a BlinkM LED
type BlinkMDriver struct {
	*Driver
}

// NewBlinkMDriver creates a new BlinkMDriver.
//
// Params:
//
//	c Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
func NewBlinkMDriver(c Connector, options ...func(Config)) *BlinkMDriver {
	b := &BlinkMDriver{
		Driver: NewDriver(c, "BlinkM", blinkmDefaultAddress),
	}
	b.afterStart = b.initialize

	for _, option := range options {
		option(b)
	}

	//nolint:forcetypeassert // ok here
	b.AddCommand("Rgb", func(params map[string]interface{}) interface{} {
		red := byte(params["red"].(float64))
		green := byte(params["green"].(float64))
		blue := byte(params["blue"].(float64))
		return b.Rgb(red, green, blue)
	})

	//nolint:forcetypeassert // ok here
	b.AddCommand("Fade", func(params map[string]interface{}) interface{} {
		red := byte(params["red"].(float64))
		green := byte(params["green"].(float64))
		blue := byte(params["blue"].(float64))
		return b.Fade(red, green, blue)
	})

	b.AddCommand("FirmwareVersion", func(_ map[string]interface{}) interface{} {
		version, err := b.FirmwareVersion()
		return map[string]interface{}{"version": version, "err": err}
	})

	b.AddCommand("Color", func(_ map[string]interface{}) interface{} {
		color, err := b.Color()
		return map[string]interface{}{"color": color, "err": err}
	})

	return b
}

// Rgb sets color using r,g,b params
func (b *BlinkMDriver) Rgb(red byte, green byte, blue byte) error {
	if _, err := b.connection.Write([]byte("n")); err != nil {
		return err
	}
	_, err := b.connection.Write([]byte{red, green, blue})
	return err
}

// Fade removes color using r,g,b params
func (b *BlinkMDriver) Fade(red byte, green byte, blue byte) error {
	if _, err := b.connection.Write([]byte("c")); err != nil {
		return err
	}
	_, err := b.connection.Write([]byte{red, green, blue})
	return err
}

// FirmwareVersion returns version with MAYOR.minor format
func (b *BlinkMDriver) FirmwareVersion() (string, error) {
	if _, err := b.connection.Write([]byte("Z")); err != nil {
		return "", err
	}
	data := []byte{0, 0}
	read, err := b.connection.Read(data)
	if read != 2 || err != nil {
		return "", err
	}
	return fmt.Sprintf("%v.%v", data[0], data[1]), nil
}

// Color returns an array with current rgb color
func (b *BlinkMDriver) Color() ([]byte, error) {
	if _, err := b.connection.Write([]byte("g")); err != nil {
		return []byte{}, err
	}
	data := []byte{0, 0, 0}
	read, err := b.connection.Read(data)
	if read != 3 || err != nil {
		return []byte{}, err
	}
	return []byte{data[0], data[1], data[2]}, nil
}

func (b *BlinkMDriver) initialize() error {
	if _, err := b.connection.Write([]byte("o")); err != nil {
		return err
	}
	return nil
}
