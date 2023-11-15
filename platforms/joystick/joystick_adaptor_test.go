package joystick

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

func initTestAdaptor() *Adaptor {
	a := NewAdaptor("6")
	a.connect = func(j *Adaptor) error {
		j.joystick = &testJoystick{}
		return nil
	}
	return a
}

func TestJoystickAdaptorName(t *testing.T) {
	a := initTestAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "Joystick"))
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestAdaptorConnect(t *testing.T) {
	a := initTestAdaptor()
	require.NoError(t, a.Connect())

	a = NewAdaptor("6")
	err := a.Connect()
	require.ErrorContains(t, err, "no joystick available")
}

func TestAdaptorFinalize(t *testing.T) {
	a := initTestAdaptor()
	_ = a.Connect()
	require.NoError(t, a.Finalize())
}
