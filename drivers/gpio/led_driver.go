package gpio

import (
	"gobot.io/x/gobot/v2"
)

// LedDriver represents a digital Led
type LedDriver struct {
	*driver
	high bool
}

// NewLedDriver return a new LedDriver given a DigitalWriter and pin.
//
// Supported options:
//
//	"WithName"
//
// Adds the following API Commands:
//
//	"Brightness" - See LedDriver.Brightness
//	"Toggle" - See LedDriver.Toggle
//	"On" - See LedDriver.On
//	"Off" - See LedDriver.Off
func NewLedDriver(a DigitalWriter, pin string, opts ...interface{}) *LedDriver {
	//nolint:forcetypeassert // no error return value, so there is no better way
	d := &LedDriver{
		driver: newDriver(a.(gobot.Connection), "LED", append(opts, withPin(pin))...),
	}
	//nolint:forcetypeassert // ok here
	d.AddCommand("Brightness", func(params map[string]interface{}) interface{} {
		level := byte(params["level"].(float64))
		return d.Brightness(level)
	})

	d.AddCommand("Toggle", func(_ map[string]interface{}) interface{} {
		return d.Toggle()
	})

	d.AddCommand("On", func(_ map[string]interface{}) interface{} {
		return d.On()
	})

	d.AddCommand("Off", func(_ map[string]interface{}) interface{} {
		return d.Off()
	})

	return d
}

// State return true if the led is On and false if the led is Off
func (d *LedDriver) State() bool {
	return d.high
}

// On sets the led to a high state.
func (d *LedDriver) On() error {
	if err := d.digitalWrite(d.driverCfg.pin, 1); err != nil {
		return err
	}
	d.high = true
	return nil
}

// Off sets the led to a low state.
func (d *LedDriver) Off() error {
	if err := d.digitalWrite(d.driverCfg.pin, 0); err != nil {
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
	return d.pwmWrite(d.driverCfg.pin, level)
}
