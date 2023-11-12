package gpio

import (
	"log"

	"gobot.io/x/gobot/v2"
)

// RelayDriver represents a digital relay
type RelayDriver struct {
	pin        string
	name       string
	connection DigitalWriter
	high       bool
	Inverted   bool
	gobot.Commander
}

// NewRelayDriver return a new RelayDriver given a DigitalWriter and pin.
//
// Adds the following API Commands:
//
//	"Toggle" - See RelayDriver.Toggle
//	"On" - See RelayDriver.On
//	"Off" - See RelayDriver.Off
func NewRelayDriver(a DigitalWriter, pin string) *RelayDriver {
	d := &RelayDriver{
		name:       gobot.DefaultName("Relay"),
		pin:        pin,
		connection: a,
		high:       false,
		Inverted:   false,
		Commander:  gobot.NewCommander(),
	}

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
func (d *RelayDriver) Start() error { return nil }

// Halt implements the Driver interface
func (d *RelayDriver) Halt() error { return nil }

// Name returns the RelayDrivers name
func (d *RelayDriver) Name() string { return d.name }

// SetName sets the RelayDrivers name
func (d *RelayDriver) SetName(n string) { d.name = n }

// Pin returns the RelayDrivers name
func (d *RelayDriver) Pin() string { return d.pin }

// Connection returns the RelayDrivers Connection
func (d *RelayDriver) Connection() gobot.Connection {
	if conn, ok := d.connection.(gobot.Connection); ok {
		return conn
	}

	log.Printf("%s has no gobot connection\n", d.name)
	return nil
}

// State return true if the relay is On and false if the relay is Off
func (d *RelayDriver) State() bool {
	if d.Inverted {
		return !d.high
	}
	return d.high
}

// On sets the relay to a high state.
func (d *RelayDriver) On() error {
	newValue := byte(1)
	if d.Inverted {
		newValue = 0
	}
	if err := d.connection.DigitalWrite(d.Pin(), newValue); err != nil {
		return err
	}

	d.high = !d.Inverted

	return nil
}

// Off sets the relay to a low state.
func (d *RelayDriver) Off() error {
	newValue := byte(0)
	if d.Inverted {
		newValue = 1
	}
	if err := d.connection.DigitalWrite(d.Pin(), newValue); err != nil {
		return err
	}

	d.high = d.Inverted

	return nil
}

// Toggle sets the relay to the opposite of it's current state
func (d *RelayDriver) Toggle() error {
	if d.State() {
		return d.Off()
	}

	return d.On()
}
