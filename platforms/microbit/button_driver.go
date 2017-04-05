package microbit

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
)

// ButtonDriver is the Gobot driver for the Microbit's built-in buttons
type ButtonDriver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}

const (
	// BLE services
	buttonService = "e95d9882251d470aa062fa1922dfa9a8"

	// BLE characteristics
	buttonACharacteristic = "e95dda90251d470aa062fa1922dfa9a8"
	buttonBCharacteristic = "e95dda91251d470aa062fa1922dfa9a8"

	// ButtonA event
	ButtonA = "buttonA"

	// ButtonB event
	ButtonB = "buttonB"
)

// NewButtonDriver creates a Microbit ButtonDriver
func NewButtonDriver(a ble.BLEConnector) *ButtonDriver {
	n := &ButtonDriver{
		name:       gobot.DefaultName("Microbit Button"),
		connection: a,
		Eventer:    gobot.NewEventer(),
	}

	n.AddEvent(ButtonA)
	n.AddEvent(ButtonB)

	return n
}

// Connection returns the BLE connection
func (b *ButtonDriver) Connection() gobot.Connection { return b.connection }

// Name returns the Driver Name
func (b *ButtonDriver) Name() string { return b.name }

// SetName sets the Driver Name
func (b *ButtonDriver) SetName(n string) { b.name = n }

// adaptor returns BLE adaptor
func (b *ButtonDriver) adaptor() ble.BLEConnector {
	return b.Connection().(ble.BLEConnector)
}

// Start tells driver to get ready to do work
func (b *ButtonDriver) Start() (err error) {
	// subscribe to button A notifications
	b.adaptor().Subscribe(buttonACharacteristic, func(data []byte, e error) {
		b.Publish(b.Event(ButtonA), data)
	})

	// subscribe to button B notifications
	b.adaptor().Subscribe(buttonBCharacteristic, func(data []byte, e error) {
		b.Publish(b.Event(ButtonB), data)
	})

	return
}

// Halt stops LED driver (void)
func (b *ButtonDriver) Halt() (err error) {
	return
}
