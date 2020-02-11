package gpio

import (
	"time"

	"gobot.io/x/gobot"
)

// MakeyButtonDriver Represents a Makey Button
type MakeyButtonDriver struct {
	name                  string
	pin                   string
	halt                  chan bool
	connection            DigitalReader
	connectionInputPullup DigitalReaderInputPullup
	Active                bool
	inputPullup           bool
	interval              time.Duration
	gobot.Eventer
}

// NewMakeyButtonDriver returns a new MakeyButtonDriver with a polling interval of
// 10 Milliseconds given a DigitalReader and pin.
//
// Optionally accepts:
//  time.Duration: Interval at which the ButtonDriver is polled for new information
func NewMakeyButtonDriver(a DigitalReader, pin string, v ...time.Duration) *MakeyButtonDriver {
	m := &MakeyButtonDriver{
		name:        gobot.DefaultName("MakeyButton"),
		connection:  a,
		pin:         pin,
		Active:      false,
		Eventer:     gobot.NewEventer(),
		interval:    10 * time.Millisecond,
		halt:        make(chan bool),
		inputPullup: false,
	}

	if len(v) > 0 {
		m.interval = v[0]
	}

	m.AddEvent(Error)
	m.AddEvent(ButtonPush)
	m.AddEvent(ButtonRelease)

	return m
}

// SetInputPullup permit to put pin mode as INPUT_PULLUP
// Your pletaform must be support it
func (b *MakeyButtonDriver) SetInputPullup() (err error) {
	if reader, ok := b.Connection().(DigitalReaderInputPullup); ok {
		b.connectionInputPullup = reader
		b.inputPullup = true

		return
	}

	err = ErrDigitalReadInputPullupUnsupported

	return
}

// IsInputPullup return if pin is setting as INPUT_PULLUP
func (b *MakeyButtonDriver) IsInputPullup() bool { return b.inputPullup }

// Name returns the MakeyButtonDrivers name
func (b *MakeyButtonDriver) Name() string { return b.name }

// SetName sets the MakeyButtonDrivers name
func (b *MakeyButtonDriver) SetName(n string) { b.name = n }

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
func (b *MakeyButtonDriver) Start() (err error) {
	state := 1
	go func() {
		timer := time.NewTimer(b.interval)
		timer.Stop()
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
				if newValue == 0 {
					b.Active = true
					b.Publish(ButtonPush, newValue)
				} else {
					b.Active = false
					b.Publish(ButtonRelease, newValue)
				}
			}
			timer.Reset(b.interval)
			select {
			case <-timer.C:
			case <-b.halt:
				return
			}
		}
	}()
	return
}

// Halt stops polling the makey button for new information
func (b *MakeyButtonDriver) Halt() (err error) {
	b.halt <- true
	return
}
