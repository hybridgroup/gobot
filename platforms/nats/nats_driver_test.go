package nats

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*Driver)(nil)

func TestNatsDriver(t *testing.T) {
	d := NewDriver(initTestNatsAdaptor(), "/test/topic")

	assert.True(t, strings.HasPrefix(d.Name(), "NATS"))
	assert.True(t, strings.HasPrefix(d.Connection().Name(), "NATS"))
	assert.NotNil(t, d.adaptor())

	assert.NoError(t, d.Start())
	assert.NoError(t, d.Halt())
}

func TestNatsDriverName(t *testing.T) {
	d := NewDriver(initTestNatsAdaptor(), "/test/topic")
	assert.True(t, strings.HasPrefix(d.Name(), "NATS"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestNatsDriverTopic(t *testing.T) {
	d := NewDriver(initTestNatsAdaptor(), "/test/topic")
	d.SetTopic("interestingtopic")
	assert.Equal(t, "interestingtopic", d.Topic())
}
