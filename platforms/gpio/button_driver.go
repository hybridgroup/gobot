package gpio

import (
	"time"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*ButtonDriver)(nil)

// Represents a digital Button
type ButtonDriver struct {
	Active     bool
	pin        string
	name       string
	interval   time.Duration
	connection gobot.Connection
	gobot.Eventer
}

// NewButtonDriver return a new ButtonDriver given a DigitalReader, name and pin
func NewButtonDriver(a DigitalReader, name string, pin string, v ...time.Duration) *ButtonDriver {
	b := &ButtonDriver{
		name:       name,
		connection: a.(gobot.Connection),
		pin:        pin,
		Active:     false,
		Eventer:    gobot.NewEventer(),
		interval:   10 * time.Millisecond,
	}

	if len(v) > 0 {
		b.interval = v[0]
	}

	b.AddEvent("push")
	b.AddEvent("release")
	b.AddEvent("error")

	return b
}

func (b *ButtonDriver) adaptor() DigitalReader {
	return b.Connection().(DigitalReader)
}

// Starts the ButtonDriver and reads the state of the button at the given Driver.Interval().
// Returns true on successful start of the driver.
//
// Emits the Events:
// 	"push"    int - On button push
//	"release" int - On button release
//	"error" error - On button error
func (b *ButtonDriver) Start() (errs []error) {
	state := 0
	go func() {
		for {
			newValue, err := b.readState()
			if err != nil {
				gobot.Publish(b.Event("error"), err)
			} else if newValue != state && newValue != -1 {
				state = newValue
				b.update(newValue)
			}
			<-time.After(b.interval)
		}
	}()
	return
}

// Halt returns true on a successful halt of the driver
func (b *ButtonDriver) Halt() (errs []error) { return }

func (b *ButtonDriver) Name() string                 { return b.name }
func (b *ButtonDriver) Pin() string                  { return b.pin }
func (b *ButtonDriver) Connection() gobot.Connection { return b.connection }

func (b *ButtonDriver) readState() (val int, err error) {
	return b.adaptor().DigitalRead(b.Pin())
}

func (b *ButtonDriver) update(newValue int) {
	if newValue == 1 {
		b.Active = true
		gobot.Publish(b.Event("push"), newValue)
	} else {
		b.Active = false
		gobot.Publish(b.Event("release"), newValue)
	}
}
