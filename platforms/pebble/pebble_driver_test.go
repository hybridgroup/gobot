package pebble

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestPebbleDriver() *PebbleDriver {
	return NewPebbleDriver(NewPebbleAdaptor("adaptor"), "pebble")
}

func TestPebbleDriverStart(t *testing.T) {
	d := initTestPebbleDriver()
	gobot.Assert(t, d.Start(), true)
}

func TestPebbleDriverHalt(t *testing.T) {
	d := initTestPebbleDriver()
	gobot.Assert(t, d.Halt(), true)
}

func TestPebbleDriverNotification(t *testing.T) {
	d := initTestPebbleDriver()
	d.SendNotification("Hello")
	d.SendNotification("World")

	gobot.Assert(t, d.Messages[0], "Hello")
	gobot.Assert(t, d.PendingMessage(), "Hello")
	gobot.Assert(t, d.PendingMessage(), "World")
	gobot.Assert(t, d.PendingMessage(), "")
}
