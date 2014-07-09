package gpio

import (
	"github.com/hybridgroup/gobot"
)

type MakeyButtonDriver struct {
	gobot.Driver
	Active bool
	data   []int
}

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

	m.Driver.AddEvent("push")
	m.Driver.AddEvent("release")

	return m
}

func (b *MakeyButtonDriver) adaptor() DigitalReader {
	return b.Driver.Adaptor().(DigitalReader)
}

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
func (m *MakeyButtonDriver) Halt() bool { return true }
func (m *MakeyButtonDriver) Init() bool { return true }

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
