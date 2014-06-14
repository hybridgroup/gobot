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
	gobot.Expect(t, d.Start(), true)
}

func TestPebbleDriverHalt(t *testing.T) {
	d := initTestPebbleDriver()
	gobot.Expect(t, d.Halt(), true)
}

func TestPebbleDriverNotification(t *testing.T) {
	d := initTestPebbleDriver()
	d.SendNotification("Hello")
	d.SendNotification("World")

	gobot.Expect(t, d.Messages[0], "Hello")
	gobot.Expect(t, d.PendingMessage(), "Hello")
	gobot.Expect(t, d.PendingMessage(), "World")
	gobot.Expect(t, d.PendingMessage(), "")
}
