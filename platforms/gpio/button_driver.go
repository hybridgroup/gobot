package gpio

import (
	"github.com/hybridgroup/gobot"
	"time"
)

var _ gobot.DriverInterface = (*ButtonDriver)(nil)

// Represents a digital Button
type ButtonDriver struct {
	gobot.Driver
	Active bool
}

// NewButtonDriver return a new ButtonDriver given a DigitalReader, name and pin
func NewButtonDriver(a DigitalReader, name string, pin string) *ButtonDriver {
	b := &ButtonDriver{
		Driver: *gobot.NewDriver(
			name,
			"ButtonDriver",
			a.(gobot.AdaptorInterface),
			pin,
		),
		Active: false,
	}

	b.AddEvent("push")
	b.AddEvent("release")
	b.AddEvent("error")

	return b
}

func (b *ButtonDriver) adaptor() DigitalReader {
	return b.Adaptor().(DigitalReader)
}

// Starts the ButtonDriver and reads the state of the button at the given Driver.Interval().
// Returns true on successful start of the driver.
//
// Emits the Events:
// 	"push"    int - On button push
//	"release" int - On button release
//	"error" error - On button error
func (b *ButtonDriver) Start() error {
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
			<-time.After(b.Interval())
		}
	}()
	return nil
}

// Halt returns true on a successful halt of the driver
func (b *ButtonDriver) Halt() error { return nil }

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
