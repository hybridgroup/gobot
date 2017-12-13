package spi

import (
	"errors"
	"strconv"

	"gobot.io/x/gobot"
)

// MCP3002DriverMaxChannel is the number of channels of this A/D converter.
const MCP3002DriverMaxChannel = 2

// MCP3002Driver is a driver for the MCP3002 A/D converter.
type MCP3002Driver struct {
	name       string
	connector  Connector
	connection Connection
}

// NewMCP3002Driver creates a new Gobot Driver for MCP3002 A/D converter
//
// Params:
//      a *Adaptor - the Adaptor to use with this Driver
//
func NewMCP3002Driver(a Connector) *MCP3002Driver {
	d := &MCP3002Driver{
		name:      gobot.DefaultName("MCP3002"),
		connector: a,
	}
	return d
}

// Name returns the name of the device.
func (d *MCP3002Driver) Name() string { return d.name }

// SetName sets the name of the device.
func (d *MCP3002Driver) SetName(n string) { d.name = n }

// Connection returns the Connection of the device.
func (d *MCP3002Driver) Connection() gobot.Connection { return d.connection.(gobot.Connection) }

// Start initializes the driver.
func (d *MCP3002Driver) Start() (err error) {
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
func (d *MCP3002Driver) Halt() (err error) {
	d.connection.Close()
	return
}

// Read reads the current analog data for the desired channel.
func (d *MCP3002Driver) Read(channel int) (result int, err error) {
	if channel < 0 || channel > MCP3002DriverMaxChannel-1 {
		return 0, errors.New("Invalid channel for read")
	}

	tx := make([]byte, 2)
	tx[0] = 0x68 + (byte(channel) << 4)
	tx[1] = 0x00

	rx := make([]byte, 2)

	err = d.connection.Tx(tx, rx)
	if err == nil && len(rx) == 2 {
		result = int(((rx[0] & 0x3) << 8) + rx[1])
	}

	return result, err
}

// AnalogRead returns value from analog reading of specified pin
func (d *MCP3002Driver) AnalogRead(pin string) (value int, err error) {
	channel, _ := strconv.Atoi(pin)
	value, err = d.Read(channel)

	return
}
