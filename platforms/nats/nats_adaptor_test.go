package nats

import (
	"fmt"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gobottest"
)

var _ gobot.Adaptor = (*NatsAdaptor)(nil)

func initTestNatsAdaptor() *NatsAdaptor {
	return NewNatsAdaptor("Nats", "localhost:4222", 9999)
}

// These tests only succeed when there is a nats server available
func TestNatsAdaptorPublishWhenConnected(t *testing.T) {
	a := initTestNatsAdaptor()
	a.Connect()
	data := []byte("o")
	gobottest.Assert(t, a.Publish("test", data), true)
}

func TestNatsAdaptorOnWhenConnected(t *testing.T) {
	a := initTestNatsAdaptor()
	a.Connect()
	gobottest.Assert(t, a.On("hola", func(data []byte) {
		fmt.Println("hola")
	}), true)
}

// These tests only succeed when there is no nats server available
func TestNatsAdaptorConnect(t *testing.T) {
	a := initTestNatsAdaptor()
	gobottest.Assert(t, a.Connect()[0].Error(), "nats: no servers available for connection")
}

func TestNatsAdaptorFinalize(t *testing.T) {
	a := initTestNatsAdaptor()
	gobottest.Assert(t, len(a.Finalize()), 0)
}

func TestNatsAdaptorCannotPublishUnlessConnected(t *testing.T) {
	a := initTestNatsAdaptor()
	data := []byte("o")
	gobottest.Assert(t, a.Publish("test", data), false)
}

func TestNatsAdaptorCannotOnUnlessConnected(t *testing.T) {
	a := initTestNatsAdaptor()
	gobottest.Assert(t, a.On("hola", func(data []byte) {
		fmt.Println("hola")
	}), false)
}
