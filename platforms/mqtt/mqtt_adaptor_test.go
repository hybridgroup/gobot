package mqtt

import (
	"fmt"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gobottest"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

func initTestMqttAdaptor() *Adaptor {
	return NewAdaptor("localhost:1883", "client")
}

func TestMqttAdaptorConnect(t *testing.T) {
	a := initTestMqttAdaptor()
	gobottest.Assert(t, a.Connect().Error(), "1 error(s) occurred:\n\n* Network Error : Unknown protocol")
}

func TestMqttAdaptorFinalize(t *testing.T) {
	a := initTestMqttAdaptor()
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestMqttAdaptorCannotPublishUnlessConnected(t *testing.T) {
	a := initTestMqttAdaptor()
	data := []byte("o")
	gobottest.Assert(t, a.Publish("test", data), false)
}

func TestMqttAdaptorPublishWhenConnected(t *testing.T) {
	a := initTestMqttAdaptor()
	a.Connect()
	data := []byte("o")
	gobottest.Assert(t, a.Publish("test", data), true)
}

func TestMqttAdaptorCannotOnUnlessConnected(t *testing.T) {
	a := initTestMqttAdaptor()
	gobottest.Assert(t, a.On("hola", func(data []byte) {
		fmt.Println("hola")
	}), false)
}

func TestMqttAdaptorOnWhenConnected(t *testing.T) {
	a := initTestMqttAdaptor()
	a.Connect()
	gobottest.Assert(t, a.On("hola", func(data []byte) {
		fmt.Println("hola")
	}), true)
}
