package pebble

import (
	"github.com/hybridgroup/gobot"
)

type PebbleDriver struct {
	gobot.Driver
	Adaptor *PebbleAdaptor
}

type PebbleInterface interface {
}

func NewPebble(adaptor *PebbleAdaptor) *PebbleDriver {
	d := new(PebbleDriver)
	d.Events = make(map[string]chan interface{})
	d.Adaptor = adaptor
	d.Commands = []string{
      "PublishEventC",
  }
	return d
}

func (me *PebbleDriver) Init() bool {
  me.Events["button"] = make(chan interface{})

  return true
}

func (me *PebbleDriver) Start() bool { return true }

func (me *PebbleDriver) Halt() bool { return true }

func (sd *PebbleDriver) PublishEvent(name string, data string) {
  gobot.Publish(sd.Events[name], data)
}
