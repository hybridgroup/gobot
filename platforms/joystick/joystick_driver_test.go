package joystick

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestJoystickDriver() *JoystickDriver {
	return NewJoystickDriver(NewJoystickAdaptor("bot"), "bot", "/dev/null")
}

func TestJoystickDriverStart(t *testing.T) {
	t.SkipNow()
	d := initTestJoystickDriver()
	gobot.Expect(t, d.Start(), true)
}

func TestJoystickDriverHalt(t *testing.T) {
	t.SkipNow()
	d := initTestJoystickDriver()
	gobot.Expect(t, d.Halt(), true)
}

func TestJoystickDriverInit(t *testing.T) {
	t.SkipNow()
	d := initTestJoystickDriver()
	gobot.Expect(t, d.Init(), true)
}
