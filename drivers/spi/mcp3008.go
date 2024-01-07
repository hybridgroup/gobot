package spi

import (
	"fmt"
	"strconv"
)

// MCP3008DriverMaxChannel is the number of channels of this A/D converter.
const MCP3008DriverMaxChannel = 8

// MCP3008Driver is a driver for the MCP3008 A/D converter.
type MCP3008Driver struct {
	*Driver
}

// NewMCP3008Driver creates a new Gobot Driver for MCP3008Driver A/D converter
//
// Params:
//
//	a *Adaptor - the Adaptor to use with this Driver
//
// Optional params:
//
//	 spi.WithBusNumber(int):  bus to use with this driver
//		spi.WithChipNumber(int): chip to use with this driver
//	 spi.WithMode(int):    	 mode to use with this driver
//	 spi.WithBitCount(int):   number of bits to use with this driver
//	 spi.WithSpeed(int64):    speed in Hz to use with this driver
func NewMCP3008Driver(a Connector, options ...func(Config)) *MCP3008Driver {
	d := &MCP3008Driver{
		Driver: NewDriver(a, "MCP3008"),
	}
	for _, option := range options {
		option(d)
	}
	return d
}

// Read reads the current analog data for the desired channel.
func (d *MCP3008Driver) Read(channel int) (int, error) {
	if channel < 0 || channel > MCP3008DriverMaxChannel-1 {
		return 0, fmt.Errorf("Invalid channel '%d' for read", channel)
	}

	tx := make([]byte, 3)
	tx[0] = 0x01
	tx[1] = byte(8+channel) << 4
	tx[2] = 0x00

	rx := make([]byte, 3)

	if err := d.connection.ReadCommandData(tx, rx); err != nil || len(rx) != 3 {
		return 0, err
	}

	result := int((rx[1]&0x3))<<8 + int(rx[2])

	return result, nil
}

// AnalogRead returns value from analog reading of specified pin
func (d *MCP3008Driver) AnalogRead(pin string) (int, error) {
	channel, _ := strconv.Atoi(pin)
	return d.Read(channel)
}
