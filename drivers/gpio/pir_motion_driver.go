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
	b := &PIRMotionDriver{
		Driver:   NewDriver(a.(gobot.Connection), "PIRMotion"),
		Eventer:  gobot.NewEventer(),
		pin:      pin,
		active:   false,
		interval: 10 * time.Millisecond,
		halt:     make(chan bool),
	}
	b.afterStart = b.initialize
	b.beforeHalt = b.shutdown

	if len(v) > 0 {
		b.interval = v[0]
	}

	b.AddEvent(MotionDetected)
	b.AddEvent(MotionStopped)
	b.AddEvent(Error)

	return b
}

// Pin returns the PIRMotionDriver pin
func (p *PIRMotionDriver) Pin() string { return p.pin }

// Active gets the current state
func (p *PIRMotionDriver) Active() bool {
	// ensure that read and write can not interfere
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.active
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
func (p *PIRMotionDriver) initialize() error {
	go func() {
		for {
			newValue, err := p.connection.(DigitalReader).DigitalRead(p.Pin())
			if err != nil {
				p.Publish(Error, err)
			}
			p.update(newValue)
			select {
			case <-time.After(p.interval):
			case <-p.halt:
				return
			}
		}
	}()
	return nil
}

// shutdown stops polling
func (p *PIRMotionDriver) shutdown() error {
	p.halt <- true
	return nil
}

func (p *PIRMotionDriver) update(newValue int) {
	// ensure that read and write can not interfere
	p.mutex.Lock()
	defer p.mutex.Unlock()

	switch newValue {
	case 1:
		if !p.active {
			p.active = true
			p.Publish(MotionDetected, newValue)
		}
	case 0:
		if p.active {
			p.active = false
			p.Publish(MotionStopped, newValue)
		}
	}
}
