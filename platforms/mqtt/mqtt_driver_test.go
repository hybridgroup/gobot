package mqtt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*Driver)(nil)

func TestMqttDriver(t *testing.T) {
	d := NewDriver(initTestMqttAdaptor(), "/test/topic")

	assert.True(t, strings.HasPrefix(d.Name(), "MQTT"))
	assert.True(t, strings.HasPrefix(d.Connection().Name(), "MQTT"))

	require.NoError(t, d.Start())
	require.NoError(t, d.Halt())
}

func TestMqttDriverName(t *testing.T) {
	d := NewDriver(initTestMqttAdaptor(), "/test/topic")
	assert.True(t, strings.HasPrefix(d.Name(), "MQTT"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestMqttDriverTopic(t *testing.T) {
	d := NewDriver(initTestMqttAdaptor(), "/test/topic")
	assert.Equal(t, "/test/topic", d.Topic())
	d.SetTopic("/test/newtopic")
	assert.Equal(t, "/test/newtopic", d.Topic())
}

func TestMqttDriverPublish(t *testing.T) {
	a := initTestMqttAdaptor()
	d := NewDriver(a, "/test/topic")
	_ = a.Connect()
	_ = d.Start()
	defer func() { _ = d.Halt() }()
	assert.True(t, d.Publish([]byte{0x01, 0x02, 0x03}))
}

func TestMqttDriverPublishError(t *testing.T) {
	a := initTestMqttAdaptor()
	d := NewDriver(a, "/test/topic")
	_ = d.Start()
	defer func() { _ = d.Halt() }()
	assert.False(t, d.Publish([]byte{0x01, 0x02, 0x03}))
}
