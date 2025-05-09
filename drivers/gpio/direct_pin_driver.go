package gpio

import (
	"strconv"

	"gobot.io/x/gobot/v2"
)

// DirectPinDriver represents a GPIO pin
type DirectPinDriver struct {
	*driver
}

// NewDirectPinDriver return a new DirectPinDriver given a Connection and pin.
//
// Supported options:
//
//	"WithName"
//
// Adds the following API Commands:
//
//	"DigitalRead" - See DirectPinDriver.DigitalRead
//	"DigitalWrite" - See DirectPinDriver.DigitalWrite
//	"PwmWrite" - See DirectPinDriver.PwmWrite
//	"ServoWrite" - See DirectPinDriver.ServoWrite
func NewDirectPinDriver(a gobot.Connection, pin string, opts ...interface{}) *DirectPinDriver {
	d := &DirectPinDriver{
		driver: newDriver(a, "DirectPin", append(opts, withPin(pin))...),
	}

	d.AddCommand("DigitalRead", func(_ map[string]interface{}) interface{} {
		val, err := d.DigitalRead()
		return map[string]interface{}{"val": val, "err": err}
	})
	//nolint:forcetypeassert // ok here
	d.AddCommand("DigitalWrite", func(params map[string]interface{}) interface{} {
		level, _ := strconv.Atoi(params["level"].(string))
		return d.DigitalWrite(byte(level))
	})
	//nolint:forcetypeassert // ok here
	d.AddCommand("PwmWrite", func(params map[string]interface{}) interface{} {
		level, _ := strconv.Atoi(params["level"].(string))
		return d.PwmWrite(byte(level))
	})
	//nolint:forcetypeassert // ok here
	d.AddCommand("ServoWrite", func(params map[string]interface{}) interface{} {
		level, _ := strconv.Atoi(params["level"].(string))
		return d.ServoWrite(byte(level))
	})

	return d
}

// Off turn off pin
func (d *DirectPinDriver) Off() error {
	return d.digitalWrite(d.driverCfg.pin, byte(0))
}

// On turn on pin
func (d *DirectPinDriver) On() error {
	return d.digitalWrite(d.driverCfg.pin, byte(1))
}

// DigitalRead returns the current digital state of the pin
func (d *DirectPinDriver) DigitalRead() (int, error) {
	return d.digitalRead(d.driverCfg.pin)
}

// DigitalWrite writes to the pin. Acceptable values are 1 or 0
func (d *DirectPinDriver) DigitalWrite(level byte) error {
	return d.digitalWrite(d.driverCfg.pin, level)
}

// PwmWrite writes the 0-254 value to the specified pin
func (d *DirectPinDriver) PwmWrite(level byte) error {
	return d.pwmWrite(d.driverCfg.pin, level)
}

// ServoWrite writes value to the specified pin
func (d *DirectPinDriver) ServoWrite(level byte) error {
	return d.servoWrite(d.driverCfg.pin, level)
}
