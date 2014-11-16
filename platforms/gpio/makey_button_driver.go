package gpio

import (
	"github.com/hybridgroup/gobot"
	"time"
)

var _ gobot.DriverInterface = (*MakeyButtonDriver)(nil)

// Represents a Makey Button
type MakeyButtonDriver struct {
	gobot.Driver
	Active bool
	data   []int
}

// NewMakeyButtonDriver returns a new MakeyButtonDriver given a DigitalRead, name and pin.
func NewMakeyButtonDriver(a DigitalReader, name string, pin string) *MakeyButtonDriver {
	m := &MakeyButtonDriver{
		Driver: *gobot.NewDriver(
			name,
			"MakeyButtonDriver",
			a.(gobot.AdaptorInterface),
			pin,
		),
		Active: false,
	}

	m.AddEvent("error")
	m.AddEvent("push")
	m.AddEvent("release")

	return m
}

func (b *MakeyButtonDriver) adaptor() DigitalReader {
	return b.Adaptor().(DigitalReader)
}

// Starts the MakeyButtonDriver and reads the state of the button at the given Driver.Interval().
// Returns true on successful start of the driver.
//
// Emits the Events:
// 	"push"    int - On button push
//	"release" int - On button release
func (m *MakeyButtonDriver) Start() error {
	state := 0
	go func() {
		for {
			newValue, err := m.readState()
			if err != nil {
				gobot.Publish(m.Event("error"), err)
			} else if newValue != state && newValue != -1 {
				state = newValue
				if newValue == 0 {
					m.Active = true
					gobot.Publish(m.Event("push"), newValue)
				} else {
					m.Active = false
					gobot.Publish(m.Event("release"), newValue)
				}
			}
		}
		<-time.After(m.Interval())
	}()
	return nil
}

// Halt returns true on a successful halt of the driver
func (m *MakeyButtonDriver) Halt() error { return nil }

func (m *MakeyButtonDriver) readState() (val int, err error) {
	return m.adaptor().DigitalRead(m.Pin())
}
