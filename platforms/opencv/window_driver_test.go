package opencv

import (
	"path"
	"runtime"
	"testing"

	cv "github.com/hybridgroup/go-opencv/opencv"
	"github.com/hybridgroup/gobot"
)

func initTestWindowDriver() *WindowDriver {
	d := NewWindowDriver("bot")
	d.start = func(w *WindowDriver) {
		w.window = &testWindow{}
	}
	return d
}

func TestWindowDriverStart(t *testing.T) {
	d := initTestWindowDriver()
	gobot.Assert(t, len(d.Start()), 0)
}

func TestWindowDriverHalt(t *testing.T) {
	d := initTestWindowDriver()
	gobot.Assert(t, len(d.Halt()), 0)
}

func TestWindowDriverShowImage(t *testing.T) {
	d := initTestWindowDriver()
	_, currentfile, _, _ := runtime.Caller(0)
	image := cv.LoadImage(path.Join(path.Dir(currentfile), "lena-256x256.jpg"))
	d.Start()
	d.ShowImage(image)
}
