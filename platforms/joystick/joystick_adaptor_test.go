package joystick

import (
	"errors"
	"testing"

	"github.com/hybridgroup/gobot"
)

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
	gobot.Assert(t, len(a.Connect()), 0)

	a = NewJoystickAdaptor("bot")
	gobot.Assert(t, a.Connect()[0], errors.New("No joystick available"))
}

func TestJoystickAdaptorFinalize(t *testing.T) {
	a := initTestJoystickAdaptor()
	a.Connect()
	gobot.Assert(t, len(a.Finalize()), 0)
}
