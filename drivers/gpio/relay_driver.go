package gpio

import "gobot.io/x/gobot"

// RelayDriver represents a digital relay
type RelayDriver struct {
	pin        string
	name       string
	connection DigitalWriter
	high       bool
	Inverted   bool
	gobot.Commander
}

// NewRelayDriver return a new RelayDriver given a DigitalWriter and pin.
//
// Adds the following API Commands:
//	"Toggle" - See RelayDriver.Toggle
//	"On" - See RelayDriver.On
//	"Off" - See RelayDriver.Off
func NewRelayDriver(a DigitalWriter, pin string) *RelayDriver {
	l := &RelayDriver{
		name:       gobot.DefaultName("Relay"),
		pin:        pin,
		connection: a,
		high:       false,
		Inverted:   false,
		Commander:  gobot.NewCommander(),
	}

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
func (l *RelayDriver) Start() (err error) { return }

// Halt implements the Driver interface
func (l *RelayDriver) Halt() (err error) { return }

// Name returns the RelayDrivers name
func (l *RelayDriver) Name() string { return l.name }

// SetName sets the RelayDrivers name
func (l *RelayDriver) SetName(n string) { l.name = n }

// Pin returns the RelayDrivers name
func (l *RelayDriver) Pin() string { return l.pin }

// Connection returns the RelayDrivers Connection
func (l *RelayDriver) Connection() gobot.Connection {
	return l.connection.(gobot.Connection)
}

// State return true if the relay is On and false if the relay is Off
func (l *RelayDriver) State() bool {
	if l.Inverted {
		return !l.high
	}
	return l.high
}

// On sets the relay to a high state.
func (l *RelayDriver) On() (err error) {
	newValue := byte(1)
	if l.Inverted {
		newValue = 0
	}
	if err = l.connection.DigitalWrite(l.Pin(), newValue); err != nil {
		return
	}

	if l.Inverted {
		l.high = false
	} else {
		l.high = true
	}

	return
}

// Off sets the relay to a low state.
func (l *RelayDriver) Off() (err error) {
	newValue := byte(0)
	if l.Inverted {
		newValue = 1
	}
	if err = l.connection.DigitalWrite(l.Pin(), newValue); err != nil {
		return
	}

	if l.Inverted {
		l.high = true
	} else {
		l.high = false
	}

	return
}

// Toggle sets the relay to the opposite of it's current state
func (l *RelayDriver) Toggle() (err error) {
	if l.State() {
		err = l.Off()
	} else {
		err = l.On()
	}
	return
}
