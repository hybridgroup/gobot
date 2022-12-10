package spi

import (
	"errors"
	"strconv"

	"gobot.io/x/gobot"
)

// MCP3204DriverMaxChannel is the number of channels of this A/D converter.
const MCP3204DriverMaxChannel = 4

// MCP3204Driver is a driver for the MCP3204 A/D converter.
type MCP3204Driver struct {
	*Driver
}

// NewMCP3204Driver creates a new Gobot Driver for MCP3204Driver A/D converter
//
// Params:
//      a *Adaptor - the Adaptor to use with this Driver
//
// Optional params:
//      spi.WithBus(int):    	bus to use with this driver
//     	spi.WithChip(int):    	chip to use with this driver
//      spi.WithMode(int):    	mode to use with this driver
//      spi.WithBits(int):    	number of bits to use with this driver
//      spi.WithSpeed(int64):   speed in Hz to use with this driver
//
func NewMCP3204Driver(a Connector, options ...func(Config)) *MCP3204Driver {
	d := &MCP3204Driver{
		Driver: NewDriver(a, "MCP3204"),
	}
	for _, option := range options {
		option(d)
	}
	return d
}

// Read reads the current analog data for the desired channel.
func (d *MCP3204Driver) Read(channel int) (result int, err error) {
	if channel < 0 || channel > MCP3204DriverMaxChannel-1 {
		return 0, errors.New("Invalid channel for read")
	}

	tx := make([]byte, 3)
	tx[0] = 0x06 + (byte(channel) >> 2)
	tx[1] = (byte(channel) & 0x03) << 6
	tx[2] = 0x00

	rx := make([]byte, 3)

	err = d.connection.Tx(tx, rx)
	if err == nil && len(rx) == 3 {
		result = int((rx[1]&0xf))<<8 + int(rx[2])
	}

	return result, err
}

// AnalogRead returns value from analog reading of specified pin, scaled to 0-1023 value.
func (d *MCP3204Driver) AnalogRead(pin string) (value int, err error) {
	channel, _ := strconv.Atoi(pin)
	value, err = d.Read(channel)
	if err != nil {
		value = int(gobot.ToScale(gobot.FromScale(float64(value), 0, 4095), 0, 1023))
	}

	return
}
