package opencv

import (
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
)

func initTestCameraDriver() *CameraDriver {
	d := NewCameraDriver("bot", "")
	d.start = func(c *CameraDriver) {
		d.camera = &testCapture{}
	}
	return d
}

func TestCameraDriverStart(t *testing.T) {
	sem := make(chan bool)
	d := initTestCameraDriver()
	gobot.Assert(t, d.Start(), true)
	gobot.On(d.Event("frame"), func(data interface{}) {
		sem <- true
	})
	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Event \"frame\" was not published")
	}

}
func TestCameraDriver(t *testing.T) {
	d := NewCameraDriver("bot", "")
	d.Start()
	gobot.Refute(t, d.camera, nil)

	defer func() {
		r := recover()
		if r != nil {
			gobot.Assert(t, "unknown camera source", r)
		} else {
			t.Errorf("Did not return Unknown camera error")
		}
	}()
	d = NewCameraDriver("bot", true)
	d.Start()
}

func TestCameraDriverHalt(t *testing.T) {
	d := initTestCameraDriver()
	gobot.Assert(t, d.Halt(), true)
}
