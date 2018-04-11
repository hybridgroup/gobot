package spi

import (
	"errors"
	"strconv"

	"gobot.io/x/gobot"
)

// MCP3008DriverMaxChannel is the number of channels of this A/D converter.
const MCP3008DriverMaxChannel = 8

// MCP3008Driver is a driver for the MCP3008 A/D converter.
type MCP3008Driver struct {
	name       string
	connector  Connector
	connection Connection
	Config
	gobot.Commander
}

// NewMCP3008Driver creates a new Gobot Driver for MCP3008Driver A/D converter
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
func NewMCP3008Driver(a Connector, options ...func(Config)) *MCP3008Driver {
	d := &MCP3008Driver{
		name:      gobot.DefaultName("MCP3008"),
		connector: a,
		Config:    NewConfig(),
	}
	for _, option := range options {
		option(d)
	}
	return d
}

// Name returns the name of the device.
func (d *MCP3008Driver) Name() string { return d.name }

// SetName sets the name of the device.
func (d *MCP3008Driver) SetName(n string) { d.name = n }

// Connection returns the Connection of the device.
func (d *MCP3008Driver) Connection() gobot.Connection { return d.connection.(gobot.Connection) }

// Start initializes the driver.
func (d *MCP3008Driver) Start() (err error) {
	bus := d.GetBusOrDefault(d.connector.GetSpiDefaultBus())
	chip := d.GetChipOrDefault(d.connector.GetSpiDefaultChip())
	mode := d.GetModeOrDefault(d.connector.GetSpiDefaultMode())
	bits := d.GetBitsOrDefault(d.connector.GetSpiDefaultBits())
	maxSpeed := d.GetSpeedOrDefault(d.connector.GetSpiDefaultMaxSpeed())

	d.connection, err = d.connector.GetSpiConnection(bus, chip, mode, bits, maxSpeed)
	if err != nil {
		return err
	}
	return nil
}

// Halt stops the driver.
func (d *MCP3008Driver) Halt() (err error) {
	d.connection.Close()
	return
}

// Read reads the current analog data for the desired channel.
func (d *MCP3008Driver) Read(channel int) (result int, err error) {
	if channel < 0 || channel > MCP3008DriverMaxChannel-1 {
		return 0, errors.New("Invalid channel for read")
	}

	tx := make([]byte, 3)
	tx[0] = 0x01
	tx[1] = byte(8+channel) << 4
	tx[2] = 0x00

	rx := make([]byte, 3)

	err = d.connection.Tx(tx, rx)
	if err == nil && len(rx) == 3 {
		result = int((rx[1]&0x3))<<8 + int(rx[2])
	}

	return result, err
}

// AnalogRead returns value from analog reading of specified pin
func (d *MCP3008Driver) AnalogRead(pin string) (value int, err error) {
	channel, _ := strconv.Atoi(pin)
	value, err = d.Read(channel)

	return
}
