package nats

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

func connStub(options ...nats.Option) func() (*nats.Conn, error) {
	return func() (*nats.Conn, error) {
		opts := nats.GetDefaultOptions()
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
	a := NewAdaptor("localhost:4222", 19999)
	a.connect = func() (*nats.Conn, error) {
		c := &nats.Conn{}
		return c, nil
	}
	return a
}

func initTestNatsAdaptorWithAuth() *Adaptor {
	a := NewAdaptorWithAuth("localhost:4222", 29999, "user", "pass")
	a.connect = func() (*nats.Conn, error) {
		c := &nats.Conn{}
		return c, nil
	}
	return a
}

func initTestNatsAdaptorTLS(options ...nats.Option) *Adaptor {
	a := NewAdaptor("tls://localhost:4242", 39999, options...)
	a.connect = connStub(options...)
	return a
}

func TestNatsAdaptorName(t *testing.T) {
	a := initTestNatsAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "NATS"))
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestNatsAdaptorReturnsHost(t *testing.T) {
	a := initTestNatsAdaptor()
	assert.Equal(t, "nats://localhost:4222", a.Host)
}

func TestNatsAdaptorWithAuth(t *testing.T) {
	a := initTestNatsAdaptorWithAuth()
	assert.Equal(t, "user", a.username)
	assert.Equal(t, "pass", a.password)
}

func TestNatsAdapterSetsRootCAs(t *testing.T) {
	a := initTestNatsAdaptorTLS(nats.RootCAs("test_certs/catest.pem"))
	assert.Equal(t, "tls://localhost:4242", a.Host)
	_ = a.Connect()
	o := a.client.Opts
	casPool, err := o.RootCAsCB()
	require.NoError(t, err)
	assert.NotNil(t, casPool)
	assert.True(t, o.Secure)
}

func TestNatsAdapterSetsClientCerts(t *testing.T) {
	a := initTestNatsAdaptorTLS(nats.ClientCert("test_certs/client-cert.pem", "test_certs/client-key.pem"))
	assert.Equal(t, "tls://localhost:4242", a.Host)
	_ = a.Connect()
	cert, err := a.client.Opts.TLSCertCB()
	require.NoError(t, err)
	assert.NotNil(t, cert)
	assert.NotNil(t, cert.Leaf)
	assert.True(t, a.client.Opts.Secure)
}

func TestNatsAdapterSetsClientCertsWithUserInfo(t *testing.T) {
	a := initTestNatsAdaptorTLS(nats.ClientCert("test_certs/client-cert.pem", "test_certs/client-key.pem"),
		nats.UserInfo("test", "testwd"))
	assert.Equal(t, "tls://localhost:4242", a.Host)
	_ = a.Connect()
	cert, err := a.client.Opts.TLSCertCB()
	require.NoError(t, err)
	assert.NotNil(t, cert)
	assert.NotNil(t, cert.Leaf)
	assert.True(t, a.client.Opts.Secure)
	assert.Equal(t, "test", a.client.Opts.User)
	assert.Equal(t, "testwd", a.client.Opts.Password)
}

// TODO: implement this test without requiring actual server connection
func TestNatsAdaptorPublishWhenConnected(t *testing.T) {
	t.Skip("TODO: implement this test without requiring actual server connection")
	a := initTestNatsAdaptor()
	_ = a.Connect()
	data := []byte("o")
	assert.True(t, a.Publish("test", data))
}

// TODO: implement this test without requiring actual server connection
func TestNatsAdaptorOnWhenConnected(t *testing.T) {
	t.Skip("TODO: implement this test without requiring actual server connection")
	a := initTestNatsAdaptor()
	_ = a.Connect()
	assert.True(t, a.On("hola", func(msg Message) {
		fmt.Println("hola")
	}))
}

// TODO: implement this test without requiring actual server connection
func TestNatsAdaptorPublishWhenConnectedWithAuth(t *testing.T) {
	t.Skip("TODO: implement this test without requiring actual server connection")
	a := NewAdaptorWithAuth("localhost:4222", 49999, "test", "testwd")
	_ = a.Connect()
	data := []byte("o")
	assert.True(t, a.Publish("test", data))
}

// TODO: implement this test without requiring actual server connection
func TestNatsAdaptorOnWhenConnectedWithAuth(t *testing.T) {
	t.Skip("TODO: implement this test without requiring actual server connection")
	log.Println("###not skipped###")
	a := NewAdaptorWithAuth("localhost:4222", 59999, "test", "testwd")
	_ = a.Connect()
	assert.True(t, a.On("hola", func(msg Message) {
		fmt.Println("hola")
	}))
}

func TestNatsAdaptorFailedConnect(t *testing.T) {
	a := NewAdaptor("localhost:9999", 69999)
	err := a.Connect()
	if err != nil && strings.Contains(err.Error(), "cannot assign requested address") {
		t.Skip("FLAKY: Can not test, because IP or port is in use.")
	}
	require.ErrorContains(t, err, "nats: no servers available for connection")
}

func TestNatsAdaptorFinalize(t *testing.T) {
	a := NewAdaptor("localhost:9999", 79999)
	require.NoError(t, a.Finalize())
}

func TestNatsAdaptorCannotPublishUnlessConnected(t *testing.T) {
	a := NewAdaptor("localhost:9999", 89999)
	data := []byte("o")
	assert.False(t, a.Publish("test", data))
}

func TestNatsAdaptorCannotOnUnlessConnected(t *testing.T) {
	a := NewAdaptor("localhost:9999", 99999)
	assert.False(t, a.On("hola", func(msg Message) {
		fmt.Println("hola")
	}))
}
