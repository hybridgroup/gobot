package joystick

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var j *JoystickAdaptor

func init() {
	j = NewJoystickAdaptor("bot")
}

func TestFinalize(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, j.Finalize(), true)
}
func TestConnect(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, j.Connect(), true)
}
