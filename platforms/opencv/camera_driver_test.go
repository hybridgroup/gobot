package opencv

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestCameraDriver() *CameraDriver {
	return NewCameraDriver("bot", "")
}

func TestCameraDriverStart(t *testing.T) {
  t.SkipNow()
	d := initTestCameraDriver()
	gobot.Expect(t, d.Start(), true)
}

func TestCameraDriverStartPanic(t *testing.T) {
  recovered := false
  defer func() {
    if r := recover(); r != nil {
      recovered = true
    }
  }()

  NewCameraDriver("bot", false).Start()
	gobot.Expect(t, recovered, true)
}

func TestCameraDriverHalt(t *testing.T) {
	d := initTestCameraDriver()
	gobot.Expect(t, d.Halt(), true)
}

func TestCameraDriverInit(t *testing.T) {
	d := initTestCameraDriver()
	gobot.Expect(t, d.Init(), true)
}
