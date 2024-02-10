package microbit

import (
	"bytes"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble"
)

const (
	// temperatureService = "e95d6100251d470aa062fa1922dfa9a8"
	temperatureChara = "e95d9250251d470aa062fa1922dfa9a8"

	TemperatureEvent = "temperature"
)

// TemperatureDriver is the Gobot driver for the Microbit's built-in thermometer
type TemperatureDriver struct {
	*ble.Driver
	gobot.Eventer
}

// NewTemperatureDriver creates a Microbit TemperatureDriver
func NewTemperatureDriver(a gobot.BLEConnector, opts ...ble.OptionApplier) *TemperatureDriver {
	d := &TemperatureDriver{
		Eventer: gobot.NewEventer(),
	}
	d.Driver = ble.NewDriver(a, "Microbit Temperature", d.initialize, nil, opts...)

	d.AddEvent(TemperatureEvent)

	return d
}

// initialize tells driver to get ready to do work
func (d *TemperatureDriver) initialize() error {
	// subscribe to temperature notifications
	return d.Adaptor().Subscribe(temperatureChara, func(data []byte) {
		var l int8
		buf := bytes.NewBuffer(data)
		val, _ := buf.ReadByte()
		l = int8(val)

		d.Publish(d.Event(TemperatureEvent), l)
	})
}
