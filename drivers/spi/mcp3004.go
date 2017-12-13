package spi

import (
	"strconv"

	"gobot.io/x/gobot"
)

// MCP3004Driver is a driver for the MCP3008 A/D converter.
type MCP3004Driver struct {
	name       string
	connector  Connector
	connection Connection
}

// NewMCP3004Driver creates a new Gobot Driver for MCP3004Driver A/D converter
//
// Params:
//      a *Adaptor - the Adaptor to use with this Driver
//
func NewMCP3004Driver(a Connector) *MCP3004Driver {
	d := &MCP3004Driver{
		name:      gobot.DefaultName("MCP3008"),
		connector: a,
	}
	return d
}

// Name returns the name of the device.
func (d *MCP3004Driver) Name() string { return d.name }

// SetName sets the name of the device.
func (d *MCP3004Driver) SetName(n string) { d.name = n }

// Connection returns the Connection of the device.
func (d *MCP3004Driver) Connection() gobot.Connection { return d.connection.(gobot.Connection) }

// Start initializes the driver.
func (d *MCP3004Driver) Start() (err error) {
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
func (d *MCP3004Driver) Halt() (err error) {
	d.connection.Close()
	return
}

// Read reads the current analog data for the desired channel.
func (d *MCP3004Driver) Read(channel int) (int, error) {
	tx := make([]byte, 3)
	tx[0] = 0x01
	tx[1] = 0x80 + (byte(channel) << 4)
	tx[2] = 0x00

	d.connection.Tx(tx, nil)

	return 0, nil
}

// AnalogRead returns value from analog reading of specified pin
func (d *MCP3004Driver) AnalogRead(pin string) (value int, err error) {
	channel, _ := strconv.Atoi(pin)
	value, err = d.Read(channel)

	return
}
