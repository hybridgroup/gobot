package joystick

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var d *JoystickDriver

func init() {
	d = NewJoystickDriver(NewJoystickAdaptor("bot"), "bot", "/dev/null")
}

func TestStart(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, d.Start(), true)
}

func TestHalt(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, d.Halt(), true)
}

func TestInit(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, d.Init(), true)
}
