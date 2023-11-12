package gpio

import (
	"time"

	"gobot.io/x/gobot/v2"
)

// PIRMotionDriver represents a digital Proximity Infra Red (PIR) motion detecter
type PIRMotionDriver struct {
	*Driver
	gobot.Eventer
	pin      string
	active   bool
	halt     chan bool
	interval time.Duration
}

// NewPIRMotionDriver returns a new PIRMotionDriver with a polling interval of
// 10 Milliseconds given a DigitalReader and pin.
//
// Optionally accepts:
//
//	time.Duration: Interval at which the PIRMotionDriver is polled for new information
func NewPIRMotionDriver(a DigitalReader, pin string, v ...time.Duration) *PIRMotionDriver {
	//nolint:forcetypeassert // no error return value, so there is no better way
	d := &PIRMotionDriver{
		Driver:   NewDriver(a.(gobot.Connection), "PIRMotion"),
		Eventer:  gobot.NewEventer(),
		pin:      pin,
		active:   false,
		interval: 10 * time.Millisecond,
		halt:     make(chan bool),
	}
	d.afterStart = d.initialize
	d.beforeHalt = d.shutdown

	if len(v) > 0 {
		d.interval = v[0]
	}

	d.AddEvent(MotionDetected)
	d.AddEvent(MotionStopped)
	d.AddEvent(Error)

	return d
}

// Pin returns the PIRMotionDriver pin
func (d *PIRMotionDriver) Pin() string { return d.pin }

// Active gets the current state
func (d *PIRMotionDriver) Active() bool {
	// ensure that read and write can not interfere
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.active
}

// initialize the PIRMotionDriver and polls the state of the sensor at the given interval.
//
// Emits the Events:
//
//	MotionDetected - On motion detected
//	MotionStopped int - On motion stopped
//	Error error - On button error
//
// The PIRMotionDriver will send the MotionDetected event over and over,
// just as long as motion is still being detected.
// It will only send the MotionStopped event once, however, until
// motion starts being detected again
func (d *PIRMotionDriver) initialize() error {
	go func() {
		for {
			newValue, err := d.connection.(DigitalReader).DigitalRead(d.Pin())
			if err != nil {
				d.Publish(Error, err)
			}
			d.update(newValue)
			select {
			case <-time.After(d.interval):
			case <-d.halt:
				return
			}
		}
	}()
	return nil
}

// shutdown stops polling
func (d *PIRMotionDriver) shutdown() error {
	d.halt <- true
	return nil
}

func (d *PIRMotionDriver) update(newValue int) {
	// ensure that read and write can not interfere
	d.mutex.Lock()
	defer d.mutex.Unlock()

	switch newValue {
	case 1:
		if !d.active {
			d.active = true
			d.Publish(MotionDetected, newValue)
		}
	case 0:
		if d.active {
			d.active = false
			d.Publish(MotionStopped, newValue)
		}
	}
}
