package pebble

import (
	"github.com/hybridgroup/gobot"
)

type PebbleDriver struct {
	gobot.Driver
	Messages []string
}

// NewPebbleDriver creates a new pebble driver with specified name
// Adds following events:
//		button - Sent when a pebble button is pressed
//		accel - Pebble watch acceleromenter data
//		tab - When a pebble watch tap event is detected
// And the following API commands:
//		"publish_event"
//		"send_notification"
//		"pending_message"
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

	p.AddCommand("send_notification", func(params map[string]interface{}) interface{} {
		p.SendNotification(params["message"].(string))
		return nil
	})

	p.AddCommand("pending_message", func(params map[string]interface{}) interface{} {
		return p.PendingMessage()
	})

	return p
}

// Start returns true if driver is initialized correctly
func (d *PebbleDriver) Start() bool { return true }

// Halt returns true if driver is halted succesfully
func (d *PebbleDriver) Halt() bool { return true }

// PublishEvent publishes event with specified name and data in gobot
func (d *PebbleDriver) PublishEvent(name string, data string) {
	gobot.Publish(d.Event(name), data)
}

// SendNotification appends message to list of notifications to be sent to watch
func (d *PebbleDriver) SendNotification(message string) string {
	d.Messages = append(d.Messages, message)
	return message
}

// PendingMessages returns messages to be sent as notifications to pebble
// (Not intented to be used directly)
func (d *PebbleDriver) PendingMessage() string {
	if len(d.Messages) < 1 {
		return ""
	}
	m := d.Messages[0]
	d.Messages = d.Messages[1:]

	return m
}
