package mqtt

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestMqttAdaptor() *MqttAdaptor {
	return NewMqttAdaptor("mqtt", "localhost:1883")
}

func TestMqttAdaptorConnect(t *testing.T) {
	a := initTestMqttAdaptor()
	gobot.Assert(t, a.Connect(), true)
}

func TestMqttAdaptorFinalize(t *testing.T) {
	a := initTestMqttAdaptor()
	gobot.Assert(t, a.Finalize(), true)
}
