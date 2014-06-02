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
	i := len(d.Messages) - 1
	if i < 0 {
		return ""
	}
	m := d.Messages[i]
	d.Messages = append(d.Messages[i+1:], d.Messages[:i]...)

	return m
}
