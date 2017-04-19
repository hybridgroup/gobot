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

func TestMqttDriverName(t *testing.T) {
	d := NewDriver(initTestMqttAdaptor(), "/test/topic")
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "MQTT"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestMqttDriverTopic(t *testing.T) {
	d := NewDriver(initTestMqttAdaptor(), "/test/topic")
	gobottest.Assert(t, d.Topic(), "/test/topic")
	d.SetTopic("/test/newtopic")
	gobottest.Assert(t, d.Topic(), "/test/newtopic")
}

func TestMqttDriverPublish(t *testing.T) {
	a := initTestMqttAdaptor()
	d := NewDriver(a, "/test/topic")
	a.Connect()
	d.Start()
	defer d.Halt()
	gobottest.Assert(t, d.Publish([]byte{0x01, 0x02, 0x03}), true)
}

func TestMqttDriverPublishError(t *testing.T) {
	a := initTestMqttAdaptor()
	d := NewDriver(a, "/test/topic")
	d.Start()
	defer d.Halt()
	gobottest.Assert(t, d.Publish([]byte{0x01, 0x02, 0x03}), false)
}
