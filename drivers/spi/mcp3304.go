package spi

import (
	"errors"
	"strconv"

	"gobot.io/x/gobot"
)

// MCP3304DriverMaxChannel is the number of channels of this A/D converter.
const MCP3304DriverMaxChannel = 8

// MCP3304Driver is a driver for the MCP3304 A/D converter.
type MCP3304Driver struct {
	name       string
	connector  Connector
	connection Connection
	Config
	gobot.Commander
}

// NewMCP3304Driver creates a new Gobot Driver for MCP3304Driver A/D converter
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
func NewMCP3304Driver(a Connector, options ...func(Config)) *MCP3304Driver {
	d := &MCP3304Driver{
		name:      gobot.DefaultName("MCP3304"),
		connector: a,
		Config:    NewConfig(),
	}
	for _, option := range options {
		option(d)
	}
	return d
}

// Name returns the name of the device.
func (d *MCP3304Driver) Name() string { return d.name }

// SetName sets the name of the device.
func (d *MCP3304Driver) SetName(n string) { d.name = n }

// Connection returns the Connection of the device.
func (d *MCP3304Driver) Connection() gobot.Connection { return d.connection.(gobot.Connection) }

// Start initializes the driver.
func (d *MCP3304Driver) Start() (err error) {
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
func (d *MCP3304Driver) Halt() (err error) {
	d.connection.Close()
	return
}

// Read reads the current analog data for the desired channel.
func (d *MCP3304Driver) Read(channel int) (result int, err error) {
	if channel < 0 || channel > MCP3304DriverMaxChannel-1 {
		return 0, errors.New("Invalid channel for read")
	}

	tx := make([]byte, 3)
	tx[0] = 0x0c + (byte(channel) >> 1)
	tx[1] = (byte(channel) & 0x01) << 7
	tx[2] = 0x00

	rx := make([]byte, 3)

	err = d.connection.Tx(tx, rx)
	if err == nil && len(rx) == 3 {
		result = int((rx[1]&0xf))<<8 + int(rx[2])
	}

	return result, err
}

// AnalogRead returns value from analog reading of specified pin, scaled to 0-1023 value.
func (d *MCP3304Driver) AnalogRead(pin string) (value int, err error) {
	channel, _ := strconv.Atoi(pin)
	value, err = d.Read(channel)
	if err != nil {
		value = int(gobot.ToScale(gobot.FromScale(float64(value), 0, 4095), 0, 1023))
	}

	return
}
