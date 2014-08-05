package pebble

import (
	"github.com/hybridgroup/gobot"
)

type PebbleDriver struct {
	gobot.Driver
	Messages []string
}

type PebbleInterface interface {
}

func NewPebbleDriver(adaptor *PebbleAdaptor, name string) *PebbleDriver {
	p := &PebbleDriver{
		Driver: *gobot.NewDriver(
			name,
			"PebbleDriver",
			adaptor,
		),
		Messages: []string{},
	}

	p.AddEvent("button")
	p.AddEvent("accel")
	p.AddEvent("tap")

	p.AddCommand("publish_event", func(params map[string]interface{}) interface{} {
		p.PublishEvent(params["name"].(string), params["data"].(string))
		return nil
	})

	p.AddCommand("SendNotification", func(params map[string]interface{}) interface{} {
		p.SendNotification(params["message"].(string))
		return nil
	})

	p.AddCommand("pending_message", func(params map[string]interface{}) interface{} {
		return p.PendingMessage()
	})

	return p
}
func (d *PebbleDriver) adaptor() *PebbleAdaptor {
	return d.Adaptor().(*PebbleAdaptor)
}

func (d *PebbleDriver) Start() bool { return true }

func (d *PebbleDriver) Halt() bool { return true }

func (d *PebbleDriver) PublishEvent(name string, data string) {
	gobot.Publish(d.Event(name), data)
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
