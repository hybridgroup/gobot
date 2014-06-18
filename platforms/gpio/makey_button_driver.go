package gpio

import (
	"github.com/hybridgroup/gobot"
)

type MakeyButtonDriver struct {
	gobot.Driver
	Adaptor DigitalReader
	Active  bool
	data    []int
}

func NewMakeyButtonDriver(a DigitalReader, name string, pin string) *MakeyButtonDriver {
	return &MakeyButtonDriver{
		Driver: gobot.Driver{
			Name: name,
			Pin:  pin,
			Events: map[string]*gobot.Event{
				"push":    gobot.NewEvent(),
				"release": gobot.NewEvent(),
			},
		},
		Active:  false,
		Adaptor: a,
	}
}

func (m *MakeyButtonDriver) Start() bool {
	state := 0
	gobot.Every(m.Interval, func() {
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
	return m.Adaptor.DigitalRead(m.Pin)
}

func (m *MakeyButtonDriver) update(newVal int) {
	if newVal == 0 {
		m.Active = true
		gobot.Publish(m.Events["push"], newVal)
	} else {
		m.Active = false
		gobot.Publish(m.Events["release"], newVal)
	}
}
