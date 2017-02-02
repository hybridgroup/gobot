package mqtt

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*Driver)(nil)

func TestMqttDriver(t *testing.T) {
	d := NewDriver(initTestMqttAdaptor(), "/test/topic")

	gobottest.Assert(t, strings.HasPrefix(d.Name(), "MQTT"), true)
	gobottest.Assert(t, strings.HasPrefix(d.Connection().Name(), "MQTT"), true)

	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}
