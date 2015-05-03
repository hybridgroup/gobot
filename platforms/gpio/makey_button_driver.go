package gpio

import (
	"time"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*MakeyButtonDriver)(nil)

// MakeyButtonDriver Represents a Makey Button
type MakeyButtonDriver struct {
	name       string
	pin        string
	halt       chan bool
	connection DigitalReader
	Active     bool
	interval   time.Duration
	gobot.Eventer
}

// NewMakeyButtonDriver returns a new MakeyButtonDriver with a polling interval of
// 10 Milliseconds given a DigitalReader, name and pin.
//
// Optinally accepts:
//  time.Duration: Interval at which the ButtonDriver is polled for new information
func NewMakeyButtonDriver(a DigitalReader, name string, pin string, v ...time.Duration) *MakeyButtonDriver {
	m := &MakeyButtonDriver{
		name:       name,
		connection: a,
		pin:        pin,
		Active:     false,
		Eventer:    gobot.NewEventer(),
		interval:   10 * time.Millisecond,
		halt:       make(chan bool),
	}

	if len(v) > 0 {
		m.interval = v[0]
	}

	m.AddEvent(Error)
	m.AddEvent(Push)
	m.AddEvent(Release)

	return m
}

// Name returns the MakeyButtonDrivers name
func (b *MakeyButtonDriver) Name() string { return b.name }

// Pin returns the MakeyButtonDrivers pin
func (b *MakeyButtonDriver) Pin() string { return b.pin }

// Connection returns the MakeyButtonDrivers Connection
func (b *MakeyButtonDriver) Connection() gobot.Connection { return b.connection.(gobot.Connection) }

// Start starts the MakeyButtonDriver and polls the state of the button at the given interval.
//
// Emits the Events:
// 	Push int - On button push
//	Release int - On button release
//	Error error - On button error
func (b *MakeyButtonDriver) Start() (errs []error) {
	state := 1
	go func() {
		for {
			newValue, err := b.connection.DigitalRead(b.Pin())
			if err != nil {
				gobot.Publish(b.Event(Error), err)
			} else if newValue != state && newValue != -1 {
				state = newValue
				if newValue == 0 {
					b.Active = true
					gobot.Publish(b.Event(Push), newValue)
				} else {
					b.Active = false
					gobot.Publish(b.Event(Release), newValue)
				}
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

// Halt stops polling the makey button for new information
func (b *MakeyButtonDriver) Halt() (errs []error) {
	b.halt <- true
	return
}
