package opencv

import (
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gobottest"
)

var _ gobot.Driver = (*CameraDriver)(nil)

func initTestCameraDriver() *CameraDriver {
	d := NewCameraDriver("bot", "")
	d.start = func(c *CameraDriver) (err error) {
		d.camera = &testCapture{}
		return nil
	}
	return d
}

func TestCameraDriver(t *testing.T) {
	d := initTestCameraDriver()
	gobottest.Assert(t, d.Name(), "bot")
	gobottest.Assert(t, d.Connection(), (gobot.Connection)(nil))
}
func TestCameraDriverStart(t *testing.T) {
	sem := make(chan bool)
	d := initTestCameraDriver()
	gobottest.Assert(t, len(d.Start()), 0)
	d.On(d.Event("frame"), func(data interface{}) {
		sem <- true
	})
	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Event \"frame\" was not published")
	}

	d = NewCameraDriver("bot", "")
	gobottest.Assert(t, len(d.Start()), 0)

	d = NewCameraDriver("bot", true)
	gobottest.Refute(t, len(d.Start()), 0)

}

func TestCameraDriverHalt(t *testing.T) {
	d := initTestCameraDriver()
	gobottest.Assert(t, len(d.Halt()), 0)
}
