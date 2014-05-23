package gpio

import (
	"github.com/hybridgroup/gobot"
)

type ButtonDriver struct {
	gobot.Driver
	Adaptor DigitalReader
	Active  bool
}

func NewButtonDriver(a DigitalReader, name string, pin string) *ButtonDriver {
	return &ButtonDriver{
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

func (b *ButtonDriver) Start() bool {
	state := 0
	go func() {
		for {
			new_value := b.readState()
			if new_value != state && new_value != -1 {
				state = new_value
				b.update(new_value)
			}
		}
	}()
	return true
}
func (b *ButtonDriver) Halt() bool { return true }
func (b *ButtonDriver) Init() bool { return true }

func (b *ButtonDriver) readState() int {
	return b.Adaptor.DigitalRead(b.Pin)
}

func (b *ButtonDriver) update(new_val int) {
	if new_val == 1 {
		b.Active = true
		gobot.Publish(b.Events["push"], new_val)
	} else {
		b.Active = false
		gobot.Publish(b.Events["release"], new_val)
	}
}
