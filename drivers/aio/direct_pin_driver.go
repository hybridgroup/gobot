package aio

import (
	"gobot.io/x/gobot"
)

// DirectPinDriver represents a AIO pin
type DirectPinDriver struct {
	name       string
	pin        string
	connection gobot.Connection
	gobot.Commander
}

// NewDirectPinDriver return a new DirectPinDriver given a Connection and pin.
//
// Adds the following API Command:
// 	"AnalogRead" - See DirectPinDriver.AnalogRead
func NewDirectPinDriver(a gobot.Connection, pin string) *DirectPinDriver {
	d := &DirectPinDriver{
		name:       "DirectPin",
		connection: a,
		pin:        pin,
		Commander:  gobot.NewCommander(),
	}

	d.AddCommand("AnalogRead", func(params map[string]interface{}) interface{} {
		val, err := d.AnalogRead()
		return map[string]interface{}{"val": val, "err": err}
	})

	return d
}

// Name returns the DirectPinDrivers name
func (d *DirectPinDriver) Name() string { return d.name }

// SetName sets the DirectPinDrivers name
func (d *DirectPinDriver) SetName(n string) { d.name = n }

// Pin returns the DirectPinDrivers pin
func (d *DirectPinDriver) Pin() string { return d.pin }

// Connection returns the DirectPinDrivers Connection
func (d *DirectPinDriver) Connection() gobot.Connection { return d.connection }

// Start implements the Driver interface
func (d *DirectPinDriver) Start() (err error) { return }

// Halt implements the Driver interface
func (d *DirectPinDriver) Halt() (err error) { return }

// AnalogRead reads the current analog reading of the pin
func (d *DirectPinDriver) AnalogRead() (val int, err error) {
	if reader, ok := d.Connection().(AnalogReader); ok {
		return reader.AnalogRead(d.Pin())
	}
	err = ErrAnalogReadUnsupported
	return
}
