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
	gobot.Assert(t, a.Connect(), nil)

	a = NewJoystickAdaptor("bot")
	gobot.Assert(t, a.Connect(), errors.New("No joystick available"))
}

func TestJoystickAdaptorFinalize(t *testing.T) {
	a := initTestJoystickAdaptor()
	a.Connect()
	gobot.Assert(t, a.Finalize(), nil)
}
