package gpio

import (
	"time"

	"gobot.io/x/gobot"
)

// ButtonDriver Represents a digital Button
type ButtonDriver struct {
	Active       bool
	DefaultState int
	pin          string
	name         string
	halt         chan bool
	interval     time.Duration
	connection   DigitalReader
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
			newValue, err := b.connection.DigitalRead(b.Pin())
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
