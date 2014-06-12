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
	p := &PebbleDriver{
		Driver: gobot.Driver{
			Name: name,
			Events: map[string]*gobot.Event{
				"button": gobot.NewEvent(),
				"accel":  gobot.NewEvent(),
				"tap":    gobot.NewEvent(),
			},
			Commands: make(map[string]func(map[string]interface{}) interface{}),
		},
		Messages: []string{},
		Adaptor:  adaptor,
	}

	p.Driver.AddCommand("PublishEvent", func(params map[string]interface{}) interface{} {
		p.PublishEvent(params["name"].(string), params["data"].(string))
		return nil
	})

	p.Driver.AddCommand("SendNotification", func(params map[string]interface{}) interface{} {
		p.SendNotification(params["message"].(string))
		return nil
	})

	p.Driver.AddCommand("PendingMessage", func(params map[string]interface{}) interface{} {
		m := make(map[string]string)
		m["result"] = p.PendingMessage()
		return m
	})

	return p
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
