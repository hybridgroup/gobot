package pebble

import (
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*Driver)(nil)

func initTestDriver() *Driver {
	return NewDriver(NewAdaptor())
}

func TestDriverStart(t *testing.T) {
	d := initTestDriver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestDriverHalt(t *testing.T) {
	d := initTestDriver()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestDriver(t *testing.T) {
	d := initTestDriver()

	gobottest.Assert(t, d.Name(), "Pebble")
	gobottest.Assert(t, d.Connection().Name(), "Pebble")

	sem := make(chan bool)
	d.SendNotification("Hello")
	d.SendNotification("World")

	gobottest.Assert(t, d.Messages[0], "Hello")
	gobottest.Assert(t, d.PendingMessage(), "Hello")
	gobottest.Assert(t, d.PendingMessage(), "World")
	gobottest.Assert(t, d.PendingMessage(), "")

	d.On(d.Event("button"), func(data interface{}) {
		sem <- true
	})

	d.PublishEvent("button", "")

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Button Event was not published")
	}

	d.On(d.Event("accel"), func(data interface{}) {
		sem <- true
	})

	d.Command("publish_event")(map[string]interface{}{"name": "accel", "data": "100"})

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Accel Event was not published")
	}

	d.Command("send_notification")(map[string]interface{}{"message": "Hey buddy!"})
	gobottest.Assert(t, d.Messages[0], "Hey buddy!")

	message := d.Command("pending_message")(map[string]interface{}{})
	gobottest.Assert(t, message, "Hey buddy!")
}
