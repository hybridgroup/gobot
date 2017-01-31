package nats

import (
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*Driver)(nil)

func TestNatsDriver(t *testing.T) {
	d := NewDriver(initTestNatsAdaptor(), "/test/topic")

	gobottest.Assert(t, d.Name(), "NATS")
	gobottest.Assert(t, d.Connection().Name(), "NATS")

	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}
