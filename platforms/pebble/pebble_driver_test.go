package pebble

import (
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
)

func initTestPebbleDriver() *PebbleDriver {
	return NewPebbleDriver(NewPebbleAdaptor("adaptor"), "pebble")
}

func TestPebbleDriverStart(t *testing.T) {
	d := initTestPebbleDriver()
	gobot.Assert(t, len(d.Start()), 0)
}

func TestPebbleDriverHalt(t *testing.T) {
	d := initTestPebbleDriver()
	gobot.Assert(t, len(d.Halt()), 0)
}

func TestPebbleDriver(t *testing.T) {
	d := initTestPebbleDriver()

	gobot.Assert(t, d.Name(), "pebble")
	gobot.Assert(t, d.Connection().Name(), "adaptor")

	sem := make(chan bool)
	d.SendNotification("Hello")
	d.SendNotification("World")

	gobot.Assert(t, d.Messages[0], "Hello")
	gobot.Assert(t, d.PendingMessage(), "Hello")
	gobot.Assert(t, d.PendingMessage(), "World")
	gobot.Assert(t, d.PendingMessage(), "")

	gobot.On(d.Event("button"), func(data interface{}) {
		sem <- true
	})

	d.PublishEvent("button", "")

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Button Event was not published")
	}

	gobot.On(d.Event("accel"), func(data interface{}) {
		sem <- true
	})

	d.Command("publish_event")(map[string]interface{}{"name": "accel", "data": "100"})

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Accel Event was not published")
	}

	d.Command("send_notification")(map[string]interface{}{"message": "Hey buddy!"})
	gobot.Assert(t, d.Messages[0], "Hey buddy!")

	message := d.Command("pending_message")(map[string]interface{}{})
	gobot.Assert(t, message, "Hey buddy!")

}
