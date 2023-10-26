package mqtt

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

func initTestMqttAdaptor() *Adaptor {
	return NewAdaptor("tcp://localhost:1883", "client")
}

func TestMqttAdaptorName(t *testing.T) {
	a := initTestMqttAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "MQTT"))
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestMqttAdaptorPort(t *testing.T) {
	a := initTestMqttAdaptor()
	assert.Equal(t, "tcp://localhost:1883", a.Port())
}

func TestMqttAdaptorAutoReconnect(t *testing.T) {
	a := initTestMqttAdaptor()
	assert.False(t, a.AutoReconnect())
	a.SetAutoReconnect(true)
	assert.True(t, a.AutoReconnect())
}

func TestMqttAdaptorCleanSession(t *testing.T) {
	a := initTestMqttAdaptor()
	assert.True(t, a.CleanSession())
	a.SetCleanSession(false)
	assert.False(t, a.CleanSession())
}

func TestMqttAdaptorUseSSL(t *testing.T) {
	a := initTestMqttAdaptor()
	assert.False(t, a.UseSSL())
	a.SetUseSSL(true)
	assert.True(t, a.UseSSL())
}

func TestMqttAdaptorUseServerCert(t *testing.T) {
	a := initTestMqttAdaptor()
	assert.Equal(t, "", a.ServerCert())
	a.SetServerCert("/path/to/server.cert")
	assert.Equal(t, "/path/to/server.cert", a.ServerCert())
}

func TestMqttAdaptorUseClientCert(t *testing.T) {
	a := initTestMqttAdaptor()
	assert.Equal(t, "", a.ClientCert())
	a.SetClientCert("/path/to/client.cert")
	assert.Equal(t, "/path/to/client.cert", a.ClientCert())
}

func TestMqttAdaptorUseClientKey(t *testing.T) {
	a := initTestMqttAdaptor()
	assert.Equal(t, "", a.ClientKey())
	a.SetClientKey("/path/to/client.key")
	assert.Equal(t, "/path/to/client.key", a.ClientKey())
}

func TestMqttAdaptorConnectError(t *testing.T) {
	a := NewAdaptor("tcp://localhost:1884", "client")

	err := a.Connect()
	assert.Contains(t, err.Error(), "connection refused")
}

func TestMqttAdaptorConnectSSLError(t *testing.T) {
	a := NewAdaptor("tcp://localhost:1884", "client")
	a.SetUseSSL(true)
	err := a.Connect()
	assert.Contains(t, err.Error(), "connection refused")
}

func TestMqttAdaptorConnectWithAuthError(t *testing.T) {
	a := NewAdaptorWithAuth("xyz://localhost:1883", "client", "user", "pass")
	assert.ErrorContains(t, a.Connect(), "network Error : unknown protocol")
}

func TestMqttAdaptorFinalize(t *testing.T) {
	a := initTestMqttAdaptor()
	assert.NoError(t, a.Finalize())
}

func TestMqttAdaptorCannotPublishUnlessConnected(t *testing.T) {
	a := initTestMqttAdaptor()
	data := []byte("o")
	assert.False(t, a.Publish("test", data))
}

func TestMqttAdaptorPublishWhenConnected(t *testing.T) {
	a := initTestMqttAdaptor()
	_ = a.Connect()
	data := []byte("o")
	assert.True(t, a.Publish("test", data))
}

func TestMqttAdaptorCannotOnUnlessConnected(t *testing.T) {
	a := initTestMqttAdaptor()
	assert.False(t, a.On("hola", func(msg Message) {
		fmt.Println("hola")
	}))
}

func TestMqttAdaptorOnWhenConnected(t *testing.T) {
	a := initTestMqttAdaptor()
	_ = a.Connect()
	assert.True(t, a.On("hola", func(msg Message) {
		fmt.Println("hola")
	}))
}

func TestMqttAdaptorQoS(t *testing.T) {
	a := initTestMqttAdaptor()
	a.SetQoS(1)
	assert.Equal(t, a.qos, 1)
}
