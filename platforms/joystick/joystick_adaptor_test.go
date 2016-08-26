package joystick

import (
	"errors"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gobottest"
)

var _ gobot.Adaptor = (*JoystickAdaptor)(nil)

func initTestJoystickAdaptor() *JoystickAdaptor {
	a := NewJoystickAdaptor("bot")
	a.connect = func(j *JoystickAdaptor) (err error) {
		j.joystick = &testJoystick{}
		return nil
	}
	return a
}

func TestJoystickAdaptorConnect(t *testing.T) {
	a := initTestJoystickAdaptor()
	gobottest.Assert(t, len(a.Connect()), 0)

	a = NewJoystickAdaptor("bot")
	gobottest.Assert(t, a.Connect()[0], errors.New("No joystick available"))
}

func TestJoystickAdaptorFinalize(t *testing.T) {
	a := initTestJoystickAdaptor()
	a.Connect()
	gobottest.Assert(t, len(a.Finalize()), 0)
}
