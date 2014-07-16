package opencv

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestCameraDriver() *CameraDriver {
	return NewCameraDriver("bot", "")
}

func TestCameraDriverStart(t *testing.T) {
	d := initTestCameraDriver()
	gobot.Expect(t, d.Start(), true)
}

func TestCameraDriverHalt(t *testing.T) {
	d := initTestCameraDriver()
	gobot.Expect(t, d.Halt(), true)
}

func TestCameraDriverInit(t *testing.T) {
	d := initTestCameraDriver()
	gobot.Expect(t, d.Init(), true)
}
