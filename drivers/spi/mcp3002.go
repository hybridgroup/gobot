package spi

import (
	"fmt"
	"strconv"
)

// MCP3002DriverMaxChannel is the number of channels of this A/D converter.
const MCP3002DriverMaxChannel = 2

// MCP3002Driver is a driver for the MCP3002 A/D converter.
type MCP3002Driver struct {
	*Driver
}

// NewMCP3002Driver creates a new Gobot Driver for MCP3002 A/D converter
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
func NewMCP3002Driver(a Connector, options ...func(Config)) *MCP3002Driver {
	d := &MCP3002Driver{
		Driver: NewDriver(a, "MCP3002"),
	}
	for _, option := range options {
		option(d)
	}
	return d
}

// Read reads the current analog data for the desired channel.
func (d *MCP3002Driver) Read(channel int) (int, error) {
	if channel < 0 || channel > MCP3002DriverMaxChannel-1 {
		return 0, fmt.Errorf("Invalid channel '%d' for read", channel)
	}

	tx := make([]byte, 2)
	tx[0] = 0x68 + (byte(channel) << 4)
	tx[1] = 0x00

	rx := make([]byte, 2)

	if err := d.connection.ReadCommandData(tx, rx); err != nil {
		return 0, err
	}

	return int((rx[0]&0x3))<<8 + int(rx[1]), nil
}

// AnalogRead returns value from analog reading of specified pin
func (d *MCP3002Driver) AnalogRead(pin string) (int, error) {
	channel, _ := strconv.Atoi(pin)
	return d.Read(channel)
}
