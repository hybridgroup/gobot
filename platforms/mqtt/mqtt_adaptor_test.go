package mqtt

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

func initTestMqttAdaptor() *Adaptor {
	return NewAdaptor("tcp://localhost:1883", "client")
}

func TestMqttAdaptorName(t *testing.T) {
	a := initTestMqttAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "MQTT"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestMqttAdaptorPort(t *testing.T) {
	a := initTestMqttAdaptor()
	gobottest.Assert(t, a.Port(), "tcp://localhost:1883")
}

func TestMqttAdaptorAutoReconnect(t *testing.T) {
	a := initTestMqttAdaptor()
	gobottest.Assert(t, a.AutoReconnect(), false)
	a.SetAutoReconnect(true)
	gobottest.Assert(t, a.AutoReconnect(), true)
}

func TestMqttAdaptorCleanSession(t *testing.T) {
	a := initTestMqttAdaptor()
	gobottest.Assert(t, a.CleanSession(), true)
	a.SetCleanSession(false)
	gobottest.Assert(t, a.CleanSession(), false)
}

func TestMqttAdaptorUseSSL(t *testing.T) {
	a := initTestMqttAdaptor()
	gobottest.Assert(t, a.UseSSL(), false)
	a.SetUseSSL(true)
	gobottest.Assert(t, a.UseSSL(), true)
}

func TestMqttAdaptorUseServerCert(t *testing.T) {
	a := initTestMqttAdaptor()
	gobottest.Assert(t, a.ServerCert(), "")
	a.SetServerCert("/path/to/server.cert")
	gobottest.Assert(t, a.ServerCert(), "/path/to/server.cert")
}

func TestMqttAdaptorUseClientCert(t *testing.T) {
	a := initTestMqttAdaptor()
	gobottest.Assert(t, a.ClientCert(), "")
	a.SetClientCert("/path/to/client.cert")
	gobottest.Assert(t, a.ClientCert(), "/path/to/client.cert")
}

func TestMqttAdaptorUseClientKey(t *testing.T) {
	a := initTestMqttAdaptor()
	gobottest.Assert(t, a.ClientKey(), "")
	a.SetClientKey("/path/to/client.key")
	gobottest.Assert(t, a.ClientKey(), "/path/to/client.key")
}

func TestMqttAdaptorConnectError(t *testing.T) {
	a := initTestMqttAdaptor()

	err := a.Connect()
	gobottest.Assert(t, strings.Contains(err.Error(), "connection refused"), true)
}

func TestMqttAdaptorConnectSSLError(t *testing.T) {
	a := initTestMqttAdaptor()
	a.SetUseSSL(true)
	err := a.Connect()
	gobottest.Assert(t, strings.Contains(err.Error(), "connection refused"), true)
}

func TestMqttAdaptorConnectWithAuthError(t *testing.T) {
	a := NewAdaptorWithAuth("localhost:1883", "client", "user", "pass")
	var expected error
	expected = multierror.Append(expected, errors.New("Network Error : Unknown protocol"))

	gobottest.Assert(t, a.Connect(), expected)
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
	gobottest.Assert(t, a.On("hola", func(msg Message) {
		fmt.Println("hola")
	}), false)
}

func TestMqttAdaptorOnWhenConnected(t *testing.T) {
	a := initTestMqttAdaptor()
	a.Connect()
	gobottest.Assert(t, a.On("hola", func(msg Message) {
		fmt.Println("hola")
	}), true)
}
