package gpio

import "github.com/hybridgroup/gobot"

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
		return l.Brightness(level)
	})

	l.AddCommand("Toggle", func(params map[string]interface{}) interface{} {
		return l.Toggle()
	})

	l.AddCommand("On", func(params map[string]interface{}) interface{} {
		return l.On()
	})

	l.AddCommand("Off", func(params map[string]interface{}) interface{} {
		return l.Off()
	})

	return l
}

func (l *LedDriver) adaptor() PwmDigitalWriter {
	return l.Adaptor().(PwmDigitalWriter)
}

// Start starts the LedDriver. Returns true on successful start of the driver
func (l *LedDriver) Start() (errs []error) { return }

// Halt halts the LedDriver. Returns true on successful halt of the driver
func (l *LedDriver) Halt() (errs []error) { return }

// State return true if the led is On and false if the led is Off
func (l *LedDriver) State() bool {
	return l.high
}

// On sets the led to a high state. Returns true on success
func (l *LedDriver) On() (err error) {
	err = l.changeState(1)
	if err != nil {
		return
	}
	l.high = true
	return
}

// Off sets the led to a low state. Returns true on success
func (l *LedDriver) Off() (err error) {
	err = l.changeState(0)
	if err != nil {
		return
	}
	l.high = false
	return
}

// Toggle sets the led to the opposite of it's current state
func (l *LedDriver) Toggle() (err error) {
	if l.State() {
		err = l.Off()
	} else {
		err = l.On()
	}
	return
}

// Brightness sets the led to the specified level of brightness
func (l *LedDriver) Brightness(level byte) (err error) {
	return l.adaptor().PwmWrite(l.Pin(), level)
}

func (l *LedDriver) changeState(level byte) (err error) {
	return l.adaptor().DigitalWrite(l.Pin(), level)
}
