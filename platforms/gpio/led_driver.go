package gpio

import (
	"github.com/hybridgroup/gobot"
)

var _ gobot.DriverInterface = (*LedDriver)(nil)

// Represents a digital Led
type LedDriver struct {
	gobot.Driver
	high bool
}

// NewLedDriver return a new LedDriver  given a PwmDigitalWriter, name and pin.
//
// Adds the following API Commands:
//	"Brightness" - See LedDriver.Brightness
//	"Toggle" - See LedDriver.Toggle
//	"On" - See LedDriver.On
//	"Off" - See LedDriver.Off
func NewLedDriver(a PwmDigitalWriter, name string, pin string) *LedDriver {
	l := &LedDriver{
		Driver: *gobot.NewDriver(
			name,
			"LedDriver",
			pin,
			a.(gobot.AdaptorInterface),
		),
		high: false,
	}

	l.AddCommand("Brightness", func(params map[string]interface{}) interface{} {
		level := byte(params["level"].(float64))
		l.Brightness(level)
		return nil
	})

	l.AddCommand("Toggle", func(params map[string]interface{}) interface{} {
		l.Toggle()
		return nil
	})

	l.AddCommand("On", func(params map[string]interface{}) interface{} {
		l.On()
		return nil
	})

	l.AddCommand("Off", func(params map[string]interface{}) interface{} {
		l.Off()
		return nil
	})

	return l
}

func (l *LedDriver) adaptor() PwmDigitalWriter {
	return l.Adaptor().(PwmDigitalWriter)
}

// Start starts the LedDriver. Returns true on successful start of the driver
func (l *LedDriver) Start() error { return nil }

// Halt halts the LedDriver. Returns true on successful halt of the driver
func (l *LedDriver) Halt() error { return nil }

// State return true if the led is On and false if the led is Off
func (l *LedDriver) State() bool {
	return l.high
}

// On sets the led to a high state. Returns true on success
func (l *LedDriver) On() bool {
	l.changeState(1)
	l.high = true
	return true
}

// Off sets the led to a low state. Returns true on success
func (l *LedDriver) Off() bool {
	l.changeState(0)
	l.high = false
	return true
}

// Toggle sets the led to the opposite of it's current state
func (l *LedDriver) Toggle() {
	if l.State() {
		l.Off()
	} else {
		l.On()
	}
}

// Brightness sets the led to the specified level of brightness
func (l *LedDriver) Brightness(level byte) {
	l.adaptor().PwmWrite(l.Pin(), level)
}

func (l *LedDriver) changeState(level byte) {
	l.adaptor().DigitalWrite(l.Pin(), level)
}
