package joystick

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestJoystickAdaptor() *JoystickAdaptor {
	a := NewJoystickAdaptor("bot")
	a.connect = func(j *JoystickAdaptor) {
		j.joystick = &testJoystick{}
	}
	return a
}

func TestJoystickAdaptorConnect(t *testing.T) {
	a := initTestJoystickAdaptor()
	gobot.Assert(t, a.Connect(), true)
}

func TestJoystickAdaptorFinalize(t *testing.T) {
	a := initTestJoystickAdaptor()
	a.Connect()
	gobot.Assert(t, a.Finalize(), true)
}
