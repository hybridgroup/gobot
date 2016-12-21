package pebble

import (
	"gobot.io/x/gobot"
)

type Driver struct {
	name       string
	connection gobot.Connection
	gobot.Commander
	gobot.Eventer
	Messages []string
}

// NewDriver creates a new pebble driver
// Adds following events:
//		button - Sent when a pebble button is pressed
//		accel - Pebble watch acceleromenter data
//		tab - When a pebble watch tap event is detected
// And the following API commands:
//		"publish_event"
//		"send_notification"
//		"pending_message"
func NewDriver(adaptor *Adaptor) *Driver {
	p := &Driver{
		name:       "Pebble",
		connection: adaptor,
		Messages:   []string{},
		Eventer:    gobot.NewEventer(),
		Commander:  gobot.NewCommander(),
	}

	p.AddEvent("button")
	p.AddEvent("accel")
	p.AddEvent("tap")

	p.AddCommand("publish_event", func(params map[string]interface{}) interface{} {
		p.PublishEvent(params["name"].(string), params["data"].(string))
		return nil
	})

	p.AddCommand("send_notification", func(params map[string]interface{}) interface{} {
		p.SendNotification(params["message"].(string))
		return nil
	})

	p.AddCommand("pending_message", func(params map[string]interface{}) interface{} {
		return p.PendingMessage()
	})

	return p
}
func (d *Driver) Name() string                 { return d.name }
func (d *Driver) SetName(n string)             { d.name = n }
func (d *Driver) Connection() gobot.Connection { return d.connection }

// Start returns true if driver is initialized correctly
func (d *Driver) Start() (err error) { return }

// Halt returns true if driver is halted successfully
func (d *Driver) Halt() (err error) { return }

// PublishEvent publishes event with specified name and data in gobot
func (d *Driver) PublishEvent(name string, data string) {
	d.Publish(d.Event(name), data)
}

// SendNotification appends message to list of notifications to be sent to watch
func (d *Driver) SendNotification(message string) string {
	d.Messages = append(d.Messages, message)
	return message
}

// PendingMessages returns messages to be sent as notifications to pebble
// (Not intended to be used directly)
func (d *Driver) PendingMessage() string {
	if len(d.Messages) < 1 {
		return ""
	}
	m := d.Messages[0]
	d.Messages = d.Messages[1:]

	return m
}
