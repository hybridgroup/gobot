package nats

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/nats-io/nats"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

func connStub(options ...nats.Option) func() (*nats.Conn, error) {
	return func() (*nats.Conn, error) {
		opts := nats.DefaultOptions
		for _, opt := range options {
			if err := opt(&opts); err != nil {
				return nil, err
			}
		}
		c := &nats.Conn{Opts: opts}
		return c, nil
	}
}

func initTestNatsAdaptor() *Adaptor {
	a := NewAdaptor("localhost:4222", 9999)
	a.connect = func() (*nats.Conn, error) {
		c := &nats.Conn{}
		return c, nil
	}
	return a
}

func initTestNatsAdaptorWithAuth() *Adaptor {
	a := NewAdaptorWithAuth("localhost:4222", 9999, "user", "pass")
	a.connect = func() (*nats.Conn, error) {
		c := &nats.Conn{}
		return c, nil
	}
	return a
}

func initTestNatsAdaptorTLS(options ...nats.Option) *Adaptor {
	a := NewAdaptor("tls://localhost:4242", 49999, options...)
	a.connect = connStub(options...)
	return a
}

func TestNatsAdaptorName(t *testing.T) {
	a := initTestNatsAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "NATS"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestNatsAdaptorReturnsHost(t *testing.T) {
	a := initTestNatsAdaptor()
	gobottest.Assert(t, a.Host, "nats://localhost:4222")
}

func TestNatsAdaptorWithAuth(t *testing.T) {
	a := initTestNatsAdaptorWithAuth()
	gobottest.Assert(t, a.username, "user")
	gobottest.Assert(t, a.password, "pass")
}

func TestNatsAdapterSetsRootCAs(t *testing.T) {
	a := initTestNatsAdaptorTLS(nats.RootCAs("test_certs/catest.pem"))
	gobottest.Assert(t, a.Host, "tls://localhost:4242")
	a.Connect()
	o := a.client.Opts
	gobottest.Assert(t, len(o.TLSConfig.RootCAs.Subjects()), 1)
	gobottest.Assert(t, o.Secure, true)
}

func TestNatsAdapterSetsClientCerts(t *testing.T) {
	a := initTestNatsAdaptorTLS(nats.ClientCert("test_certs/client-cert.pem", "test_certs/client-key.pem"))
	gobottest.Assert(t, a.Host, "tls://localhost:4242")
	a.Connect()
	certs := a.client.Opts.TLSConfig.Certificates
	gobottest.Assert(t, len(certs), 1)
	gobottest.Assert(t, a.client.Opts.Secure, true)
}

func TestNatsAdapterSetsClientCertsWithUserInfo(t *testing.T) {
	a := initTestNatsAdaptorTLS(nats.ClientCert("test_certs/client-cert.pem", "test_certs/client-key.pem"), nats.UserInfo("test", "testwd"))
	gobottest.Assert(t, a.Host, "tls://localhost:4242")
	a.Connect()
	certs := a.client.Opts.TLSConfig.Certificates
	gobottest.Assert(t, len(certs), 1)
	gobottest.Assert(t, a.client.Opts.Secure, true)
	gobottest.Assert(t, a.client.Opts.User, "test")
	gobottest.Assert(t, a.client.Opts.Password, "testwd")
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
	gobottest.Assert(t, a.On("hola", func(msg Message) {
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
	gobottest.Assert(t, a.On("hola", func(msg Message) {
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
	gobottest.Assert(t, a.On("hola", func(msg Message) {
		fmt.Println("hola")
	}), false)
}
