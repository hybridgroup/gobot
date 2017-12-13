package spi

import (
	"errors"
	"strconv"

	"gobot.io/x/gobot"
)

// MCP3202DriverMaxChannel is the number of channels of this A/D converter.
const MCP3202DriverMaxChannel = 2

// MCP3202Driver is a driver for the MCP3202 A/D converter.
type MCP3202Driver struct {
	name       string
	connector  Connector
	connection Connection
}

// NewMCP3202Driver creates a new Gobot Driver for MCP3202Driver A/D converter
//
// Params:
//      a *Adaptor - the Adaptor to use with this Driver
//
func NewMCP3202Driver(a Connector) *MCP3202Driver {
	d := &MCP3202Driver{
		name:      gobot.DefaultName("MCP3202"),
		connector: a,
	}
	return d
}

// Name returns the name of the device.
func (d *MCP3202Driver) Name() string { return d.name }

// SetName sets the name of the device.
func (d *MCP3202Driver) SetName(n string) { d.name = n }

// Connection returns the Connection of the device.
func (d *MCP3202Driver) Connection() gobot.Connection { return d.connection.(gobot.Connection) }

// Start initializes the driver.
func (d *MCP3202Driver) Start() (err error) {
	bus := d.connector.GetSpiDefaultBus()
	mode := d.connector.GetSpiDefaultMode()
	maxSpeed := d.connector.GetSpiDefaultMaxSpeed()
	d.connection, err = d.connector.GetSpiConnection(bus, mode, maxSpeed)
	if err != nil {
		return err
	}
	return nil
}

// Halt stops the driver.
func (d *MCP3202Driver) Halt() (err error) {
	d.connection.Close()
	return
}

// Read reads the current analog data for the desired channel.
func (d *MCP3202Driver) Read(channel int) (result int, err error) {
	if channel < 0 || channel > MCP3202DriverMaxChannel-1 {
		return 0, errors.New("Invalid channel for read")
	}

	tx := make([]byte, 3)
	tx[0] = 0x01
	tx[1] = 0xa0 + byte(channel)<<6
	tx[2] = 0x00

	rx := make([]byte, 3)

	err = d.connection.Tx(tx, rx)
	if err == nil && len(rx) == 3 {
		result = int(((rx[1] & 0xf) << 8) + rx[2])
	}

	return result, err
}

// AnalogRead returns value from analog reading of specified pin, scaled to 0-1023 value.
func (d *MCP3202Driver) AnalogRead(pin string) (value int, err error) {
	channel, _ := strconv.Atoi(pin)
	value, err = d.Read(channel)
	if err != nil {
		value = int(gobot.ToScale(gobot.FromScale(float64(value), 0, 4095), 0, 1023))
	}

	return
}
