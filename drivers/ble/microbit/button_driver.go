package microbit

import (
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble"
)

const (
	// buttonService = "e95d9882251d470aa062fa1922dfa9a8"
	buttonAChara = "e95dda90251d470aa062fa1922dfa9a8"
	buttonBChara = "e95dda91251d470aa062fa1922dfa9a8"

	ButtonAEvent = "buttonA"
	ButtonBEvent = "buttonB"
)

// ButtonDriver is the Gobot driver for the Microbit's built-in buttons
type ButtonDriver struct {
	*ble.Driver
	gobot.Eventer
}

// NewButtonDriver creates a new driver
func NewButtonDriver(a gobot.BLEConnector, opts ...ble.OptionApplier) *ButtonDriver {
	d := &ButtonDriver{
		Eventer: gobot.NewEventer(),
	}

	d.Driver = ble.NewDriver(a, "Microbit Button", d.initialize, nil, opts...)

	d.AddEvent(ButtonAEvent)
	d.AddEvent(ButtonBEvent)

	return d
}

// initialize tells driver to get ready to do work
func (d *ButtonDriver) initialize() error {
	// subscribe to button A notifications
	if err := d.Adaptor().Subscribe(buttonAChara, func(data []byte) {
		d.Publish(d.Event(ButtonAEvent), data)
	}); err != nil {
		return err
	}

	// subscribe to button B notifications
	return d.Adaptor().Subscribe(buttonBChara, func(data []byte) {
		d.Publish(d.Event(ButtonBEvent), data)
	})
}
