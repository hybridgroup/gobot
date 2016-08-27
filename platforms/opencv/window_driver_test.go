package opencv

import (
	"path"
	"runtime"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gobottest"
	cv "github.com/lazywei/go-opencv/opencv"
)

var _ gobot.Driver = (*WindowDriver)(nil)

func initTestWindowDriver() *WindowDriver {
	d := NewWindowDriver("bot")
	return d
}

func TestWindowDriver(t *testing.T) {
	d := initTestWindowDriver()
	gobottest.Assert(t, d.Name(), "bot")
	gobottest.Assert(t, d.Connection(), (gobot.Connection)(nil))
}
func TestWindowDriverStart(t *testing.T) {
	d := initTestWindowDriver()
	gobottest.Assert(t, len(d.Start()), 0)
}

func TestWindowDriverHalt(t *testing.T) {
	d := initTestWindowDriver()
	gobottest.Assert(t, len(d.Halt()), 0)
}

func TestWindowDriverShowImage(t *testing.T) {
	d := initTestWindowDriver()
	_, currentfile, _, _ := runtime.Caller(0)
	image := cv.LoadImage(path.Join(path.Dir(currentfile), "lena-256x256.jpg"))
	d.Start()
	d.ShowImage(image)
}
