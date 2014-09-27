package gpio

import (
	"github.com/hybridgroup/gobot"
)

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
func (m *MakeyButtonDriver) Start() bool {
	state := 0
	gobot.Every(m.Interval(), func() {
		newValue := m.readState()
		if newValue != state && newValue != -1 {
			state = newValue
			m.update(newValue)
		}
	})
	return true
}

// Halt returns true on a successful halt of the driver
func (m *MakeyButtonDriver) Halt() bool { return true }

func (m *MakeyButtonDriver) readState() int {
	return m.adaptor().DigitalRead(m.Pin())
}

func (m *MakeyButtonDriver) update(newVal int) {
	if newVal == 0 {
		m.Active = true
		gobot.Publish(m.Event("push"), newVal)
	} else {
		m.Active = false
		gobot.Publish(m.Event("release"), newVal)
	}
}
