package gpio

import "gobot.io/x/gobot"

// LedDriver represents a digital Led
type LedDriver struct {
	pin        string
	name       string
	connection DigitalWriter
	high       bool
	gobot.Commander
}

// NewLedDriver return a new LedDriver given a DigitalWriter and pin.
//
// Adds the following API Commands:
//	"Brightness" - See LedDriver.Brightness
//	"Toggle" - See LedDriver.Toggle
//	"On" - See LedDriver.On
//	"Off" - See LedDriver.Off
func NewLedDriver(a DigitalWriter, pin string) *LedDriver {
	l := &LedDriver{
		name:       gobot.DefaultName("LED"),
		pin:        pin,
		connection: a,
		high:       false,
		Commander:  gobot.NewCommander(),
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

// Start implements the Driver interface
func (l *LedDriver) Start() (err error) { return }

// Halt implements the Driver interface
func (l *LedDriver) Halt() (err error) { return }

// Name returns the LedDrivers name
func (l *LedDriver) Name() string { return l.name }

// SetName sets the LedDrivers name
func (l *LedDriver) SetName(n string) { l.name = n }

// Pin returns the LedDrivers name
func (l *LedDriver) Pin() string { return l.pin }

// Connection returns the LedDrivers Connection
func (l *LedDriver) Connection() gobot.Connection {
	return l.connection.(gobot.Connection)
}

// State return true if the led is On and false if the led is Off
func (l *LedDriver) State() bool {
	return l.high
}

// On sets the led to a high state.
func (l *LedDriver) On() (err error) {
	if err = l.connection.DigitalWrite(l.Pin(), 1); err != nil {
		return
	}
	l.high = true
	return
}

// Off sets the led to a low state.
func (l *LedDriver) Off() (err error) {
	if err = l.connection.DigitalWrite(l.Pin(), 0); err != nil {
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
	if writer, ok := l.connection.(PwmWriter); ok {
		return writer.PwmWrite(l.Pin(), level)
	}
	return ErrPwmWriteUnsupported
}
