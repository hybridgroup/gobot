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
	//nolint:forcetypeassert // no error return value, so there is no better way
	d := &ButtonDriver{
		Driver:       NewDriver(a.(gobot.Connection), "Button"),
		Eventer:      gobot.NewEventer(),
		pin:          pin,
		active:       false,
		defaultState: 0,
		interval:     10 * time.Millisecond,
		halt:         make(chan bool),
	}
	d.afterStart = d.initialize
	d.beforeHalt = d.shutdown

	if len(v) > 0 {
		d.interval = v[0]
	}

	d.AddEvent(ButtonPush)
	d.AddEvent(ButtonRelease)
	d.AddEvent(Error)

	return d
}

// Pin returns the ButtonDrivers pin
func (d *ButtonDriver) Pin() string { return d.pin }

// Active gets the current state
func (d *ButtonDriver) Active() bool {
	// ensure that read and write can not interfere
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.active
}

// SetDefaultState for the next start.
func (d *ButtonDriver) SetDefaultState(s int) {
	// ensure that read and write can not interfere
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.defaultState = s
}

// initialize the ButtonDriver and polls the state of the button at the given interval.
//
// Emits the Events:
//
//	Push int - On button push
//	Release int - On button release
//	Error error - On button error
func (d *ButtonDriver) initialize() error {
	state := d.defaultState
	go func() {
		for {
			newValue, err := d.connection.(DigitalReader).DigitalRead(d.Pin())
			if err != nil {
				d.Publish(Error, err)
			} else if newValue != state && newValue != -1 {
				state = newValue
				d.update(newValue)
			}
			select {
			case <-time.After(d.interval):
			case <-d.halt:
				return
			}
		}
	}()
	return nil
}

func (d *ButtonDriver) shutdown() error {
	d.halt <- true
	return nil
}

func (d *ButtonDriver) update(newValue int) {
	// ensure that read and write can not interfere
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if newValue != d.defaultState {
		d.active = true
		d.Publish(ButtonPush, newValue)
	} else {
		d.active = false
		d.Publish(ButtonRelease, newValue)
	}
}
