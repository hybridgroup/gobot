package gobotGPIO

import (
	"github.com/hybridgroup/gobot"
)

type ButtonInterface interface {
	DigitalRead(string) int
}

type Button struct {
	gobot.Driver
	Adaptor ButtonInterface
	Active  bool
}

func NewButton(a ButtonInterface) *Button {
	b := new(Button)
	b.Active = false
	b.Adaptor = a
	b.Events = make(map[string]chan interface{})
	b.Events["push"] = make(chan interface{})
	b.Events["release"] = make(chan interface{})
	b.Commands = []string{}
	return b
}

func (b *Button) Start() bool {
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
func (b *Button) Halt() bool { return true }
func (b *Button) Init() bool { return true }

func (b *Button) readState() int {
	return b.Adaptor.DigitalRead(b.Pin)
}

func (b *Button) update(new_val int) {
	if new_val == 1 {
		b.Active = true
		gobot.Publish(b.Events["push"], new_val)
	} else {
		b.Active = false
		gobot.Publish(b.Events["release"], new_val)
	}
}
