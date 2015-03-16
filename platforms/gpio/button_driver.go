package gpio

import (
	"time"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*ButtonDriver)(nil)

// ButtonDriver Represents a digital Button
type ButtonDriver struct {
	Active     bool
	pin        string
	name       string
	halt       chan bool
	interval   time.Duration
	connection DigitalReader
	gobot.Eventer
}

// NewButtonDriver returns a new ButtonDriver with a polling interval of
// 10 Milliseconds given a DigitalReader, name and pin.
//
// Optinally accepts:
//  time.Duration: Interval at which the ButtonDriver is polled for new information
func NewButtonDriver(a DigitalReader, name string, pin string, v ...time.Duration) *ButtonDriver {
	b := &ButtonDriver{
		name:       name,
		connection: a,
		pin:        pin,
		Active:     false,
		Eventer:    gobot.NewEventer(),
		interval:   10 * time.Millisecond,
		halt:       make(chan bool),
	}

	if len(v) > 0 {
		b.interval = v[0]
	}

	b.AddEvent(Push)
	b.AddEvent(Release)
	b.AddEvent(Error)

	return b
}

// Start starts the ButtonDriver and polls the state of the button at the given interval.
//
// Emits the Events:
// 	Push int - On button push
//	Release int - On button release
//	Error error - On button error
func (b *ButtonDriver) Start() (errs []error) {
	state := 0
	go func() {
		for {
			newValue, err := b.connection.DigitalRead(b.Pin())
			if err != nil {
				gobot.Publish(b.Event(Error), err)
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
func (b *ButtonDriver) Halt() (errs []error) {
	b.halt <- true
	return
}

// Name returns the ButtonDrivers name
func (b *ButtonDriver) Name() string { return b.name }

// Pin returns the ButtonDrivers pin
func (b *ButtonDriver) Pin() string { return b.pin }

// Connection returns the ButtonDrivers Connection
func (b *ButtonDriver) Connection() gobot.Connection { return b.connection.(gobot.Connection) }

func (b *ButtonDriver) update(newValue int) {
	if newValue == 1 {
		b.Active = true
		gobot.Publish(b.Event(Push), newValue)
	} else {
		b.Active = false
		gobot.Publish(b.Event(Release), newValue)
	}
}
