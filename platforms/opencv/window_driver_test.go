package opencv

import (
	"path"
	"runtime"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
	cv "github.com/lazywei/go-opencv/opencv"
)

var _ gobot.Driver = (*WindowDriver)(nil)

func initTestWindowDriver() *WindowDriver {
	d := NewWindowDriver()
	return d
}

func TestWindowDriver(t *testing.T) {
	d := initTestWindowDriver()
	gobottest.Assert(t, d.Name(), "Window")
	gobottest.Assert(t, d.Connection(), (gobot.Connection)(nil))
}

func TestWindowDriverStart(t *testing.T) {
	d := initTestWindowDriver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestWindowDriverHalt(t *testing.T) {
	d := initTestWindowDriver()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestWindowDriverShowImage(t *testing.T) {
	d := initTestWindowDriver()
	_, currentfile, _, _ := runtime.Caller(0)
	image := cv.LoadImage(path.Join(path.Dir(currentfile), "lena-256x256.jpg"))
	d.Start()
	d.ShowImage(image)
}
