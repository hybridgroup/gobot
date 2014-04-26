package gobotGPIO

import (
	"github.com/hybridgroup/gobot"
)

type MakeyButtonInterface interface {
	DigitalRead(string) int
}

type MakeyButton struct {
	gobot.Driver
	Adaptor MakeyButtonInterface
	Active  bool
	data    []int
}

func NewMakeyButton(a MakeyButtonInterface) *MakeyButton {
	b := new(MakeyButton)
	b.Active = false
	b.Adaptor = a
	b.Events = make(map[string]chan interface{})
	b.Events["push"] = make(chan interface{})
	b.Events["release"] = make(chan interface{})
	b.Commands = []string{}
	return b
}

func (b *MakeyButton) Start() bool {
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
func (b *MakeyButton) Halt() bool { return true }
func (b *MakeyButton) Init() bool { return true }

func (b *MakeyButton) readState() int {
	return b.Adaptor.DigitalRead(b.Pin)
}

func (b *MakeyButton) update(new_val int) {
	if new_val == 0 {
		b.Active = true
		gobot.Publish(b.Events["push"], new_val)
	} else {
		b.Active = false
		gobot.Publish(b.Events["release"], new_val)
	}
}
