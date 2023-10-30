package joystick

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

func initTestAdaptor() *Adaptor {
	a := NewAdaptor("6")
	a.connect = func(j *Adaptor) (err error) {
		j.joystick = &testJoystick{}
		return nil
	}
	return a
}

func TestJoystickAdaptorName(t *testing.T) {
	a := initTestAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "Joystick"))
	a.SetName("NewName")
	assert.Equal(t, a.Name(), "NewName")
}

func TestAdaptorConnect(t *testing.T) {
	a := initTestAdaptor()
	assert.NoError(t, a.Connect())

	a = NewAdaptor("6")
	err := a.Connect()
	assert.True(t, strings.HasPrefix(err.Error(), "No joystick available"))
}

func TestAdaptorFinalize(t *testing.T) {
	a := initTestAdaptor()
	_ = a.Connect()
	assert.NoError(t, a.Finalize())
}
