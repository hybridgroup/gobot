package spi

import (
	"fmt"
	"strconv"

	"gobot.io/x/gobot/v2"
)

// MCP3304DriverMaxChannel is the number of channels of this A/D converter.
const MCP3304DriverMaxChannel = 8

// MCP3304Driver is a driver for the MCP3304 A/D converter.
type MCP3304Driver struct {
	*Driver
}

// NewMCP3304Driver creates a new Gobot Driver for MCP3304Driver A/D converter
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
func NewMCP3304Driver(a Connector, options ...func(Config)) *MCP3304Driver {
	d := &MCP3304Driver{
		Driver: NewDriver(a, "MCP3304"),
	}
	for _, option := range options {
		option(d)
	}
	return d
}

// Read reads the current analog data for the desired channel.
func (d *MCP3304Driver) Read(channel int) (int, error) {
	if channel < 0 || channel > MCP3304DriverMaxChannel-1 {
		return 0, fmt.Errorf("Invalid channel '%d' for read", channel)
	}

	tx := make([]byte, 3)
	tx[0] = 0x0c + (byte(channel) >> 1)
	tx[1] = (byte(channel) & 0x01) << 7
	tx[2] = 0x00

	rx := make([]byte, 3)

	if err := d.connection.ReadCommandData(tx, rx); err != nil || len(rx) != 3 {
		return 0, err
	}

	result := int((rx[1]&0xf))<<8 + int(rx[2])

	return result, nil
}

// AnalogRead returns value from analog reading of specified pin, scaled to 0-1023 value.
func (d *MCP3304Driver) AnalogRead(pin string) (int, error) {
	channel, _ := strconv.Atoi(pin)
	value, err := d.Read(channel)
	if err != nil {
		return 0, err
	}

	return int(gobot.ToScale(gobot.FromScale(float64(value), 0, 4095), 0, 1023)), err
}
