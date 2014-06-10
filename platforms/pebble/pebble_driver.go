package pebble

import (
	"github.com/hybridgroup/gobot"
)

type PebbleDriver struct {
	gobot.Driver
	Messages []string
	Adaptor  *PebbleAdaptor
}

type PebbleInterface interface {
}

func NewPebbleDriver(adaptor *PebbleAdaptor, name string) *PebbleDriver {
	return &PebbleDriver{
		Driver: gobot.Driver{
			Name: name,
			Events: map[string]chan interface{}{
				"button": make(chan interface{}),
				"accel":  make(chan interface{}),
				"tap":    make(chan interface{}),
			},
			Commands: []string{
				"PublishEventC",
				"SendNotificationC",
				"PendingMessageC",
			},
		},
		Messages: []string{},
		Adaptor:  adaptor,
	}
}

func (d *PebbleDriver) Start() bool { return true }

func (d *PebbleDriver) Halt() bool { return true }

func (d *PebbleDriver) PublishEvent(name string, data string) {
	gobot.Publish(d.Events[name], data)
}

func (d *PebbleDriver) SendNotification(message string) string {
	d.Messages = append(d.Messages, message)
	return message
}

func (d *PebbleDriver) PendingMessage() string {
	if len(d.Messages) < 1 {
		return ""
	}
	m := d.Messages[0]
	d.Messages = d.Messages[1:]

	return m
}
