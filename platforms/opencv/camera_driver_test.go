package opencv

import (
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
)

func initTestCameraDriver() *CameraDriver {
	d := NewCameraDriver("bot", "")
	d.start = func(c *CameraDriver) (err error) {
		d.camera = &testCapture{}
		return nil
	}
	return d
}

func TestCameraDriverStart(t *testing.T) {
	sem := make(chan bool)
	d := initTestCameraDriver()
	gobot.Assert(t, d.Start(), nil)
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
	err := d.Start()
	gobot.Assert(t, err, nil)

	d = NewCameraDriver("bot", true)
	err = d.Start()
	gobot.Refute(t, err, nil)
}

func TestCameraDriverHalt(t *testing.T) {
	d := initTestCameraDriver()
	gobot.Assert(t, d.Halt(), nil)
}
