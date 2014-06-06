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
			Events: map[string]chan interface{}{
				"push":    make(chan interface{}),
				"release": make(chan interface{}),
			},
		},
		Active:  false,
		Adaptor: a,
	}
}

func (m *MakeyButtonDriver) Start() bool {
	state := 0
	gobot.Every(m.Interval, func() {
		new_value := m.readState()
		if new_value != state && new_value != -1 {
			state = new_value
			m.update(new_value)
		}
	})
	return true
}
func (m *MakeyButtonDriver) Halt() bool { return true }
func (m *MakeyButtonDriver) Init() bool { return true }

func (m *MakeyButtonDriver) readState() int {
	return m.Adaptor.DigitalRead(m.Pin)
}

func (m *MakeyButtonDriver) update(new_val int) {
	if new_val == 0 {
		m.Active = true
		gobot.Publish(m.Events["push"], new_val)
	} else {
		m.Active = false
		gobot.Publish(m.Events["release"], new_val)
	}
}
