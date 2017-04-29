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
	gobottest.Refute(t, d.adaptor(), nil)

	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}

func TestNatsDriverName(t *testing.T) {
	d := NewDriver(initTestNatsAdaptor(), "/test/topic")
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "NATS"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestNatsDriverTopic(t *testing.T) {
	d := NewDriver(initTestNatsAdaptor(), "/test/topic")
	d.SetTopic("interestingtopic")
	gobottest.Assert(t, d.Topic(), "interestingtopic")
}
