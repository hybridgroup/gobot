package joystick

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestJoystickAdaptor() *JoystickAdaptor {
	return NewJoystickAdaptor("bot")
}

func TestJoystickAdaptorConnect(t *testing.T) {
	t.SkipNow()
	a := initTestJoystickAdaptor()
	gobot.Assert(t, a.Connect(), true)
}

func TestJoystickAdaptorFinalize(t *testing.T) {
	t.SkipNow()
	a := initTestJoystickAdaptor()
	gobot.Assert(t, a.Finalize(), true)
}
