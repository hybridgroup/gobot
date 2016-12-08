package nats

import (
	"errors"
	"fmt"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
	"github.com/nats-io/nats"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

func initTestNatsAdaptor() *Adaptor {
	a := NewAdaptor("localhost:4222", 9999)
	a.connect = func() (*nats.Conn, error) {
		c := &nats.Conn{}
		return c, nil
	}
	return a
}

func TestNatsAdaptorReturnsHost(t *testing.T) {
	a := initTestNatsAdaptor()
	gobottest.Assert(t, a.Host, "localhost:4222")
}

// TODO: implement this test without requiring actual server connection
func TestNatsAdaptorPublishWhenConnected(t *testing.T) {
	t.Skip("TODO: implement this test without requiring actual server connection")
	a := initTestNatsAdaptor()
	a.Connect()
	data := []byte("o")
	gobottest.Assert(t, a.Publish("test", data), true)
}

// TODO: implement this test without requiring actual server connection
func TestNatsAdaptorOnWhenConnected(t *testing.T) {
	t.Skip("TODO: implement this test without requiring actual server connection")
	a := initTestNatsAdaptor()
	a.Connect()
	gobottest.Assert(t, a.On("hola", func(data []byte) {
		fmt.Println("hola")
	}), true)
}

// TODO: implement this test without requiring actual server connection
func TestNatsAdaptorPublishWhenConnectedWithAuth(t *testing.T) {
	t.Skip("TODO: implement this test without requiring actual server connection")
	a := NewAdaptorWithAuth("localhost:4222", 9999, "test", "testwd")
	a.Connect()
	data := []byte("o")
	gobottest.Assert(t, a.Publish("test", data), true)
}

// TODO: implement this test without requiring actual server connection
func TestNatsAdaptorOnWhenConnectedWithAuth(t *testing.T) {
	t.Skip("TODO: implement this test without requiring actual server connection")
	a := NewAdaptorWithAuth("localhost:4222", 9999, "test", "testwd")
	a.Connect()
	gobottest.Assert(t, a.On("hola", func(data []byte) {
		fmt.Println("hola")
	}), true)
}

func TestNatsAdaptorFailedConnect(t *testing.T) {
	a := NewAdaptor("localhost:9999", 9999)
	gobottest.Assert(t, a.Connect(), errors.New("nats: no servers available for connection"))
}

func TestNatsAdaptorFinalize(t *testing.T) {
	a := NewAdaptor("localhost:9999", 9999)
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestNatsAdaptorCannotPublishUnlessConnected(t *testing.T) {
	a := NewAdaptor("localhost:9999", 9999)
	data := []byte("o")
	gobottest.Assert(t, a.Publish("test", data), false)
}

func TestNatsAdaptorCannotOnUnlessConnected(t *testing.T) {
	a := NewAdaptor("localhost:9999", 9999)
	gobottest.Assert(t, a.On("hola", func(data []byte) {
		fmt.Println("hola")
	}), false)
}
