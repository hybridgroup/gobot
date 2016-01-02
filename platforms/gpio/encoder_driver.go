package gpio

import (
	"strconv"
	"time"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*EncoderDriver)(nil)

// EncoderDriver represents an encoder driver
// note that the encoder driver uses 2 pins.
type EncoderDriver struct {
	pin        string
	pin2       string
	name       string
	connection DigitalReader
	// Direction can be "", "forward" or "backward"
	Direction string
	// Position when 0 is the initial state of the encoder
	Position int
	interval time.Duration
	halt     chan bool
	gobot.Eventer
}

// NewEncoderDriver return a new EncoderDriver given a DigitalREader, name and pin.
//
// Adds the following API Commands:
func NewEncoderDriver(a DigitalReader, name string, pin string) *EncoderDriver {
	pinInt, err := strconv.Atoi(pin)
	if err != nil {
		panic(err)
	}
	e := &EncoderDriver{
		name:       name,
		pin:        pin,
		pin2:       strconv.Itoa(pinInt + 1),
		connection: a,
		Eventer:    gobot.NewEventer(),
		interval:   10 * time.Millisecond,
		halt:       make(chan bool),
	}

	e.AddEvent(Error)
	e.AddEvent(Turn)

	return e
}

// Start implements the Driver interface
func (e *EncoderDriver) Start() (errs []error) {
	var isReady bool
	go func() {
		for {
			pin1, err := e.connection.DigitalRead(e.Pin())
			pin2, err2 := e.connection.DigitalRead(e.Pin2())
			if err != nil {
				gobot.Publish(e.Event(Error), err)
			}
			if err2 != nil {
				gobot.Publish(e.Event(Error), err2)
			}
			if !isReady && pin1+pin2 == 2 {
				isReady = true
			}

			if isReady {
				if pin1 > pin2 {
					e.Direction = "backward"
					e.Position--
					gobot.Publish(e.Event(Turn), pin1)
				} else if pin1 < pin2 {
					e.Direction = "forward"
					e.Position++
					gobot.Publish(e.Event(Turn), pin2)
				}
			}

			select {
			case <-time.After(e.interval):
			case <-e.halt:
				return
			}
		}

	}()
	return
}

// Halt implements the Driver interface
func (e *EncoderDriver) Halt() (errs []error) {
	e.halt <- true
	return
}

// Name returns the EncoderDriver's name
func (e *EncoderDriver) Name() string { return e.name }

// Pin returns the EncoderDriver's pin
func (e *EncoderDriver) Pin() string { return e.pin }

// Pin2 returns the pin of the secondary pin used to read value from (implicit)
func (e *EncoderDriver) Pin2() string { return e.pin2 }

// Connection returns the EncoderDrivers Connection
func (e *EncoderDriver) Connection() gobot.Connection {
	return e.connection.(gobot.Connection)
}
