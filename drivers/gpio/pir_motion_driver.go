package gpio

import (
	"time"

	"gobot.io/x/gobot"
)

// PIRMotionDriver represents a digital Proximity Infra Red (PIR) motion detecter
type PIRMotionDriver struct {
	Active                bool
	pin                   string
	name                  string
	halt                  chan bool
	interval              time.Duration
	inputPullup           bool
	connection            DigitalReader
	connectionInputPullup DigitalReaderInputPullup
	gobot.Eventer
}

// NewPIRMotionDriver returns a new PIRMotionDriver with a polling interval of
// 10 Milliseconds given a DigitalReader and pin.
//
// Optionally accepts:
//  time.Duration: Interval at which the PIRMotionDriver is polled for new information
func NewPIRMotionDriver(a DigitalReader, pin string, v ...time.Duration) *PIRMotionDriver {
	b := &PIRMotionDriver{
		name:        gobot.DefaultName("PIRMotion"),
		connection:  a,
		pin:         pin,
		Active:      false,
		Eventer:     gobot.NewEventer(),
		interval:    10 * time.Millisecond,
		halt:        make(chan bool),
		inputPullup: false,
	}

	if len(v) > 0 {
		b.interval = v[0]
	}

	b.AddEvent(MotionDetected)
	b.AddEvent(MotionStopped)
	b.AddEvent(Error)

	return b
}

// Start starts the PIRMotionDriver and polls the state of the sensor at the given interval.
//
// Emits the Events:
// 	MotionDetected - On motion detected
//	MotionStopped int - On motion stopped
//	Error error - On button error
//
// The PIRMotionDriver will send the MotionDetected event over and over,
// just as long as motion is still being detected.
// It will only send the MotionStopped event once, however, until
// motion starts being detected again
func (p *PIRMotionDriver) Start() (err error) {
	go func() {
		for {
			var newValue int
			var err error
			if p.IsInputPullup() {
				newValue, err = p.connectionInputPullup.DigitalReadInputPullup(p.Pin())
			} else {
				newValue, err = p.connection.DigitalRead(p.Pin())
			}
			if err != nil {
				p.Publish(Error, err)
			}
			switch newValue {
			case 1:
				if !p.Active {
					p.Active = true
					p.Publish(MotionDetected, newValue)
				}
			case 0:
				if p.Active {
					p.Active = false
					p.Publish(MotionStopped, newValue)
				}
			}

			select {
			case <-time.After(p.interval):
			case <-p.halt:
				return
			}
		}
	}()
	return
}

// SetInputPullup permit to put pin mode as INPUT_PULLUP
// Your pletaform must be support it
func (p *PIRMotionDriver) SetInputPullup() (err error) {
	if reader, ok := p.Connection().(DigitalReaderInputPullup); ok {
		p.connectionInputPullup = reader
		p.inputPullup = true

		return
	}

	err = ErrDigitalReadInputPullupUnsupported

	return
}

// IsInputPullup return if pin is setting as INPUT_PULLUP
func (p *PIRMotionDriver) IsInputPullup() bool { return p.inputPullup }

// Halt stops polling the button for new information
func (p *PIRMotionDriver) Halt() (err error) {
	p.halt <- true
	return
}

// Name returns the PIRMotionDriver name
func (p *PIRMotionDriver) Name() string { return p.name }

// SetName sets the PIRMotionDriver name
func (p *PIRMotionDriver) SetName(n string) { p.name = n }

// Pin returns the PIRMotionDriver pin
func (p *PIRMotionDriver) Pin() string { return p.pin }

// Connection returns the PIRMotionDriver Connection
func (p *PIRMotionDriver) Connection() gobot.Connection { return p.connection.(gobot.Connection) }
