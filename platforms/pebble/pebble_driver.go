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

func NewPebbleDriver(adaptor *PebbleAdaptor, name string) *PebbleDriver {
  return &PebbleDriver{
    Driver: gobot.Driver{
      Name: name,
      Events: map[string]chan interface{}{
        "button": make(chan interface{}),
      },
      Commands: []string{
        "PublishEventC",
      },
    },
    Adaptor: adaptor,
  }
}

func (d *PebbleDriver) Start() bool { return true }

func (d *PebbleDriver) Halt() bool { return true }

func (d *PebbleDriver) PublishEvent(name string, data string) {
  gobot.Publish(d.Events[name], data)
}
