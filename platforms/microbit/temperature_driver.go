package microbit

import (
	"bytes"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
)

// TemperatureDriver is the Gobot driver for the Microbit's built-in thermometer
type TemperatureDriver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}

const (
	// BLE services
	temperatureService = "e95d6100251d470aa062fa1922dfa9a8"

	// BLE characteristics
	temperatureCharacteristic = "e95d9250251d470aa062fa1922dfa9a8"

	// Temperature event
	Temperature = "temperature"
)

// NewTemperatureDriver creates a Microbit TemperatureDriver
func NewTemperatureDriver(a ble.BLEConnector) *TemperatureDriver {
	n := &TemperatureDriver{
		name:       gobot.DefaultName("Microbit Temperature"),
		connection: a,
		Eventer:    gobot.NewEventer(),
	}

	n.AddEvent(Temperature)

	return n
}

// Connection returns the BLE connection
func (b *TemperatureDriver) Connection() gobot.Connection { return b.connection }

// Name returns the Driver Name
func (b *TemperatureDriver) Name() string { return b.name }

// SetName sets the Driver Name
func (b *TemperatureDriver) SetName(n string) { b.name = n }

// adaptor returns BLE adaptor
func (b *TemperatureDriver) adaptor() ble.BLEConnector {
	return b.Connection().(ble.BLEConnector)
}

// Start tells driver to get ready to do work
func (b *TemperatureDriver) Start() (err error) {
	// subscribe to temperature notifications
	b.adaptor().Subscribe(temperatureCharacteristic, func(data []byte, e error) {
		var l int8
		buf := bytes.NewBuffer(data)
		val, _ := buf.ReadByte()
		l = int8(val)

		b.Publish(b.Event(Temperature), l)
	})

	return
}

// Halt stops Temperature driver (void)
func (b *TemperatureDriver) Halt() (err error) {
	return
}
