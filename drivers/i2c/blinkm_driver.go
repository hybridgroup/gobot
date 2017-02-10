package i2c

import (
	"fmt"

	"gobot.io/x/gobot"
)

const blinkmAddress = 0x09

// BlinkMDriver is a Gobot Driver for a BlinkM LED
type BlinkMDriver struct {
	name       string
	connector  Connector
	connection Connection
	Config
	gobot.Commander
}

// NewBlinkMDriver creates a new BlinkMDriver.
//
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewBlinkMDriver(a Connector, options ...func(Config)) *BlinkMDriver {
	b := &BlinkMDriver{
		name:      gobot.DefaultName("BlinkM"),
		Commander: gobot.NewCommander(),
		connector: a,
		Config:    NewConfig(),
	}

	for _, option := range options {
		option(b)
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
func (b *BlinkMDriver) Name() string { return b.name }

// SetName sets the Name for the Driver
func (b *BlinkMDriver) SetName(n string) { b.name = n }

// Connection returns the connection for the Driver
func (b *BlinkMDriver) Connection() gobot.Connection { return b.connection.(gobot.Connection) }

// Start starts the Driver up, and writes start command
func (b *BlinkMDriver) Start() (err error) {
	bus := b.GetBusOrDefault(b.connector.GetDefaultBus())
	address := b.GetAddressOrDefault(blinkmAddress)

	b.connection, err = b.connector.GetConnection(address, bus)
	if err != nil {
		return
	}

	if _, err := b.connection.Write([]byte("o")); err != nil {
		return err
	}
	return
}

// Halt returns true if device is halted successfully
func (b *BlinkMDriver) Halt() (err error) { return }

// Rgb sets color using r,g,b params
func (b *BlinkMDriver) Rgb(red byte, green byte, blue byte) (err error) {
	if _, err = b.connection.Write([]byte("n")); err != nil {
		return
	}
	_, err = b.connection.Write([]byte{red, green, blue})
	return
}

// Fade removes color using r,g,b params
func (b *BlinkMDriver) Fade(red byte, green byte, blue byte) (err error) {
	if _, err = b.connection.Write([]byte("c")); err != nil {
		return
	}
	_, err = b.connection.Write([]byte{red, green, blue})
	return
}

// FirmwareVersion returns version with MAYOR.minor format
func (b *BlinkMDriver) FirmwareVersion() (version string, err error) {
	if _, err = b.connection.Write([]byte("Z")); err != nil {
		return
	}
	data := []byte{0, 0}
	read, err := b.connection.Read(data)
	if read != 2 || err != nil {
		return
	}
	return fmt.Sprintf("%v.%v", data[0], data[1]), nil
}

// Color returns an array with current rgb color
func (b *BlinkMDriver) Color() (color []byte, err error) {
	if _, err = b.connection.Write([]byte("g")); err != nil {
		return
	}
	data := []byte{0, 0, 0}
	read, err := b.connection.Read(data)
	if read != 3 || err != nil {
		return []byte{}, err
	}
	return []byte{data[0], data[1], data[2]}, nil
}
