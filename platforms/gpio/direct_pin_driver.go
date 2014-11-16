package gpio

import (
	"strconv"

	"github.com/hybridgroup/gobot"
)

var _ gobot.DriverInterface = (*DirectPinDriver)(nil)

// Represents a raw GPIO pin
type DirectPinDriver struct {
	gobot.Driver
}

// NewDirectPinDriver return a new DirectPinDriver given a DirectPin, name and pin.
//
// Adds the following API Commands:
// 	"DigitalRead" - See DirectPinDriver.DigitalRead
// 	"DigitalWrite" - See DirectPinDriver.DigitalWrite
// 	"AnalogRead" - See DirectPinDriver.AnalogRead
// 	"AnalogWrite" - See DirectPinDriver.AnalogWrite
// 	"PwmWrite" - See DirectPinDriver.PwmWrite
// 	"ServoWrite" - See DirectPinDriver.ServoWrite
func NewDirectPinDriver(a DirectPin, name string, pin string) *DirectPinDriver {
	d := &DirectPinDriver{
		Driver: *gobot.NewDriver(
			name,
			"DirectPinDriver",
			a.(gobot.AdaptorInterface),
			pin,
		),
	}

	d.AddCommand("DigitalRead", func(params map[string]interface{}) interface{} {
		return d.DigitalRead()
	})
	d.AddCommand("DigitalWrite", func(params map[string]interface{}) interface{} {
		level, _ := strconv.Atoi(params["level"].(string))
		d.DigitalWrite(byte(level))
		return nil
	})
	d.AddCommand("AnalogRead", func(params map[string]interface{}) interface{} {
		return d.AnalogRead()
	})
	d.AddCommand("AnalogWrite", func(params map[string]interface{}) interface{} {
		level, _ := strconv.Atoi(params["level"].(string))
		d.AnalogWrite(byte(level))
		return nil
	})
	d.AddCommand("PwmWrite", func(params map[string]interface{}) interface{} {
		level, _ := strconv.Atoi(params["level"].(string))
		d.PwmWrite(byte(level))
		return nil
	})
	d.AddCommand("ServoWrite", func(params map[string]interface{}) interface{} {
		level, _ := strconv.Atoi(params["level"].(string))
		d.ServoWrite(byte(level))
		return nil
	})

	return d
}

func (d *DirectPinDriver) adaptor() DirectPin {
	return d.Adaptor().(DirectPin)
}

// Starts the DirectPinDriver. Returns true on successful start of the driver
func (d *DirectPinDriver) Start() error { return nil }

// Halts the DirectPinDriver. Returns true on successful halt of the driver
func (d *DirectPinDriver) Halt() error { return nil }

// DigitalRead returns the current digital state of the pin
func (d *DirectPinDriver) DigitalRead() int {
	return d.adaptor().DigitalRead(d.Pin())
}

// DigitalWrite writes to the pin
func (d *DirectPinDriver) DigitalWrite(level byte) {
	d.adaptor().DigitalWrite(d.Pin(), level)
}

// AnalogRead reads the current analog reading of the pin
func (d *DirectPinDriver) AnalogRead() int {
	return d.adaptor().AnalogRead(d.Pin())
}

// AnalogWrite writes to the pin
func (d *DirectPinDriver) AnalogWrite(level byte) {
	d.adaptor().AnalogWrite(d.Pin(), level)
}

// PwmWrite writes to the pin
func (d *DirectPinDriver) PwmWrite(level byte) {
	d.adaptor().PwmWrite(d.Pin(), level)
}

// ServoWrite writes to the pin
func (d *DirectPinDriver) ServoWrite(level byte) {
	d.adaptor().ServoWrite(d.Pin(), level)
}
