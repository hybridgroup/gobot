package gpio

import (
	"time"

	"gobot.io/x/gobot"
)

// ButtonDriver Represents a digital Button
type ButtonDriver struct {
	Active                bool
	DefaultState          int
	pin                   string
	name                  string
	halt                  chan bool
	interval              time.Duration
	connection            DigitalReader
	connectionInputPullup DigitalReaderInputPullup
	inputPullup           bool
	gobot.Eventer
}

// NewButtonDriver returns a new ButtonDriver with a polling interval of
// 10 Milliseconds given a DigitalReader and pin.
//
// Optionally accepts:
//  time.Duration: Interval at which the ButtonDriver is polled for new information
func NewButtonDriver(a DigitalReader, pin string, v ...time.Duration) *ButtonDriver {
	b := &ButtonDriver{
		name:         gobot.DefaultName("Button"),
		connection:   a,
		pin:          pin,
		Active:       false,
		DefaultState: 0,
		Eventer:      gobot.NewEventer(),
		interval:     10 * time.Millisecond,
		halt:         make(chan bool),
		inputPullup:  false,
	}

	if len(v) > 0 {
		b.interval = v[0]
	}

	b.AddEvent(ButtonPush)
	b.AddEvent(ButtonRelease)
	b.AddEvent(Error)

	return b
}

// Start starts the ButtonDriver and polls the state of the button at the given interval.
//
// Emits the Events:
// 	Push int - On button push
//	Release int - On button release
//	Error error - On button error
func (b *ButtonDriver) Start() (err error) {
	state := b.DefaultState
	go func() {
		for {
			var newValue int
			var err error
			if b.IsInputPullup() {
				newValue, err = b.connectionInputPullup.DigitalReadInputPullup(b.Pin())
			} else {
				newValue, err = b.connection.DigitalRead(b.Pin())
			}

			if err != nil {
				b.Publish(Error, err)
			} else if newValue != state && newValue != -1 {
				state = newValue
				b.update(newValue)
			}
			select {
			case <-time.After(b.interval):
			case <-b.halt:
				return
			}
		}
	}()
	return
}

// SetInputPullup permit to put pin mode as INPUT_PULLUP
// Your pletaform must be support it
func (b *ButtonDriver) SetInputPullup() (err error) {
	if reader, ok := b.Connection().(DigitalReaderInputPullup); ok {
		b.connectionInputPullup = reader
		b.inputPullup = true

		return
	}

	err = ErrDigitalReadInputPullupUnsupported

	return
}

// IsInputPullup return if pin is setting as INPUT_PULLUP
func (b *ButtonDriver) IsInputPullup() bool { return b.inputPullup }

// Halt stops polling the button for new information
func (b *ButtonDriver) Halt() (err error) {
	b.halt <- true
	return
}

// Name returns the ButtonDrivers name
func (b *ButtonDriver) Name() string { return b.name }

// SetName sets the ButtonDrivers name
func (b *ButtonDriver) SetName(n string) { b.name = n }

// Pin returns the ButtonDrivers pin
func (b *ButtonDriver) Pin() string { return b.pin }

// Connection returns the ButtonDrivers Connection
func (b *ButtonDriver) Connection() gobot.Connection { return b.connection.(gobot.Connection) }

func (b *ButtonDriver) update(newValue int) {
	if newValue != b.DefaultState {
		b.Active = true
		b.Publish(ButtonPush, newValue)
	} else {
		b.Active = false
		b.Publish(ButtonRelease, newValue)
	}
}
