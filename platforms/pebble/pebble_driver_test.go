package pebble

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var driver *PebbleDriver

func init() {
	driver = NewPebbleDriver(NewPebbleAdaptor("adaptor"), "pebble")
}

func TestStart(t *testing.T) {
	gobot.Expect(t, driver.Start(), true)
}

func TestHalt(t *testing.T) {
	gobot.Expect(t, driver.Halt(), true)
}

func TestNotification(t *testing.T) {
	driver.SendNotification("Hello")
	driver.SendNotification("World")

	gobot.Expect(t, driver.Messages[0], "Hello")
	gobot.Expect(t, driver.PendingMessage(), "Hello")
	gobot.Expect(t, driver.PendingMessage(), "World")
	gobot.Expect(t, driver.PendingMessage(), "")
}
