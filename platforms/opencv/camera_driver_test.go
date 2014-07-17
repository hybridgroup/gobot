package opencv

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestCameraDriver() *CameraDriver {
	return NewCameraDriver("bot", 0)
}

func TestCameraDriverStart(t *testing.T) {
	t.SkipNow()
	d := initTestCameraDriver()
	gobot.Assert(t, d.Start(), true)
}

func TestCameraDriverHalt(t *testing.T) {
	t.SkipNow()
	d := initTestCameraDriver()
	gobot.Assert(t, d.Halt(), true)
}

func TestCameraDriverInit(t *testing.T) {
	t.SkipNow()
	d := initTestCameraDriver()
	gobot.Assert(t, d.Init(), true)
}
