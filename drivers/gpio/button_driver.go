package gpio

import (
	"time"

	"gobot.io/x/gobot/v2"
)

// ButtonDriver Represents a digital Button
type ButtonDriver struct {
	*Driver
	gobot.Eventer
	pin          string
	active       bool
	defaultState int
	halt         chan bool
	interval     time.Duration
}

// NewButtonDriver returns a new ButtonDriver with a polling interval of
// 10 Milliseconds given a DigitalReader and pin.
//
// Optionally accepts:
//
//	time.Duration: Interval at which the ButtonDriver is polled for new information
func NewButtonDriver(a DigitalReader, pin string, v ...time.Duration) *ButtonDriver {
	b := &ButtonDriver{
		Driver:       NewDriver(a.(gobot.Connection), "Button"),
		Eventer:      gobot.NewEventer(),
		pin:          pin,
		active:       false,
		defaultState: 0,
		interval:     10 * time.Millisecond,
		halt:         make(chan bool),
	}
	b.afterStart = b.initialize
	b.beforeHalt = b.shutdown

	if len(v) > 0 {
		b.interval = v[0]
	}

	b.AddEvent(ButtonPush)
	b.AddEvent(ButtonRelease)
	b.AddEvent(Error)

	return b
}

// Pin returns the ButtonDrivers pin
func (b *ButtonDriver) Pin() string { return b.pin }

// Active gets the current state
func (b *ButtonDriver) Active() bool {
	// ensure that read and write can not interfere
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.active
}

// SetDefaultState for the next start.
func (b *ButtonDriver) SetDefaultState(s int) {
	// ensure that read and write can not interfere
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.defaultState = s
}

// initialize the ButtonDriver and polls the state of the button at the given interval.
//
// Emits the Events:
//
//	Push int - On button push
//	Release int - On button release
//	Error error - On button error
func (b *ButtonDriver) initialize() (err error) {
	state := b.defaultState
	go func() {
		for {
			newValue, err := b.connection.(DigitalReader).DigitalRead(b.Pin())
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

func (b *ButtonDriver) shutdown() (err error) {
	b.halt <- true
	return nil
}

func (b *ButtonDriver) update(newValue int) {
	// ensure that read and write can not interfere
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if newValue != b.defaultState {
		b.active = true
		b.Publish(ButtonPush, newValue)
	} else {
		b.active = false
		b.Publish(ButtonRelease, newValue)
	}
}
