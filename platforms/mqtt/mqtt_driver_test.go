package mqtt

import (
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*Driver)(nil)

func TestMqttDriver(t *testing.T) {
	d := NewDriver(initTestMqttAdaptor(), "/test/topic")

	gobottest.Assert(t, d.Name(), "MQTT")
	gobottest.Assert(t, d.Connection().Name(), "MQTT")

	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}
