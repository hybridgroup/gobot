package nats

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*Driver)(nil)

func TestNatsDriver(t *testing.T) {
	d := NewDriver(initTestNatsAdaptor(), "/test/topic")

	gobottest.Assert(t, strings.HasPrefix(d.Name(), "NATS"), true)
	gobottest.Assert(t, strings.HasPrefix(d.Connection().Name(), "NATS"), true)

	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}
