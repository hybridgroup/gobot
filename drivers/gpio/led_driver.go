package gpio

import (
	"log"

	"gobot.io/x/gobot/v2"
)

// LedDriver represents a digital Led
type LedDriver struct {
	pin        string
	name       string
	connection DigitalWriter
	high       bool
	gobot.Commander
}

// NewLedDriver return a new LedDriver given a DigitalWriter and pin.
//
// Adds the following API Commands:
//
//	"Brightness" - See LedDriver.Brightness
//	"Toggle" - See LedDriver.Toggle
//	"On" - See LedDriver.On
//	"Off" - See LedDriver.Off
func NewLedDriver(a DigitalWriter, pin string) *LedDriver {
	d := &LedDriver{
		name:       gobot.DefaultName("LED"),
		pin:        pin,
		connection: a,
		high:       false,
		Commander:  gobot.NewCommander(),
	}

	d.AddCommand("Brightness", func(params map[string]interface{}) interface{} {
		level := byte(params["level"].(float64)) //nolint:forcetypeassert // ok here
		return d.Brightness(level)
	})

	d.AddCommand("Toggle", func(params map[string]interface{}) interface{} {
		return d.Toggle()
	})

	d.AddCommand("On", func(params map[string]interface{}) interface{} {
		return d.On()
	})

	d.AddCommand("Off", func(params map[string]interface{}) interface{} {
		return d.Off()
	})

	return d
}

// Start implements the Driver interface
func (d *LedDriver) Start() error { return nil }

// Halt implements the Driver interface
func (d *LedDriver) Halt() error { return nil }

// Name returns the LedDrivers name
func (d *LedDriver) Name() string { return d.name }

// SetName sets the LedDrivers name
func (d *LedDriver) SetName(n string) { d.name = n }

// Pin returns the LedDrivers name
func (d *LedDriver) Pin() string { return d.pin }

// Connection returns the LedDrivers Connection
func (d *LedDriver) Connection() gobot.Connection {
	if conn, ok := d.connection.(gobot.Connection); ok {
		return conn
	}

	log.Printf("%s has no gobot connection\n", d.name)
	return nil
}

// State return true if the led is On and false if the led is Off
func (d *LedDriver) State() bool {
	return d.high
}

// On sets the led to a high state.
func (d *LedDriver) On() error {
	if err := d.connection.DigitalWrite(d.Pin(), 1); err != nil {
		return err
	}
	d.high = true
	return nil
}

// Off sets the led to a low state.
func (d *LedDriver) Off() error {
	if err := d.connection.DigitalWrite(d.Pin(), 0); err != nil {
		return err
	}
	d.high = false
	return nil
}

// Toggle sets the led to the opposite of it's current state
func (d *LedDriver) Toggle() error {
	if d.State() {
		return d.Off()
	}
	return d.On()
}

// Brightness sets the led to the specified level of brightness
func (d *LedDriver) Brightness(level byte) error {
	if writer, ok := d.connection.(PwmWriter); ok {
		return writer.PwmWrite(d.Pin(), level)
	}
	return ErrPwmWriteUnsupported
}
