package opencv

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var c *CameraDriver

func init() {
	c = NewCameraDriver("bot", 0)
}

func TestCameraStart(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, c.Start(), true)
}

func TestCameraHalt(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, c.Halt(), true)
}

func TestCameraInit(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, c.Init(), true)
}
