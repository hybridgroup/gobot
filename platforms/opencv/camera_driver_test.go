package opencv

import (
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*CameraDriver)(nil)

func initTestCameraDriver() *CameraDriver {
	d := NewCameraDriver("")
	d.start = func(c *CameraDriver) (err error) {
		d.camera = &testCapture{}
		return nil
	}
	return d
}

func TestCameraDriver(t *testing.T) {
	d := initTestCameraDriver()
	gobottest.Assert(t, d.Name(), "Camera")
	gobottest.Assert(t, d.Connection(), (gobot.Connection)(nil))
}

func TestCameraDriverName(t *testing.T) {
	d := initTestCameraDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Camera"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestCameraDriverStart(t *testing.T) {
	sem := make(chan bool)
	d := initTestCameraDriver()
	gobottest.Assert(t, d.Start(), nil)
	d.On(d.Event("frame"), func(data interface{}) {
		sem <- true
	})
	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Event \"frame\" was not published")
	}

	d = NewCameraDriver("")
	gobottest.Assert(t, d.Start(), nil)

	d = NewCameraDriver(true)
	gobottest.Refute(t, d.Start(), nil)
}

func TestCameraDriverHalt(t *testing.T) {
	d := initTestCameraDriver()
	gobottest.Assert(t, d.Halt(), nil)
}
