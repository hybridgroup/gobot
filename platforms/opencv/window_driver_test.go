package opencv

import (
	cv "github.com/hybridgroup/go-opencv/opencv"
	"github.com/hybridgroup/gobot"
	"testing"
  "path"
  "runtime"
)

func initTestWindowDriver() *WindowDriver {
	return NewWindowDriver("bot")
}

func TestWindowDriverStart(t *testing.T) {
  t.SkipNow()
	d := initTestWindowDriver()
	gobot.Expect(t, d.Start(), true)
}

func TestWindowDriverHalt(t *testing.T) {
	d := initTestWindowDriver()
	gobot.Expect(t, d.Halt(), true)
}

func TestWindowDriverInit(t *testing.T) {
	d := initTestWindowDriver()
	gobot.Expect(t, d.Init(), true)
}

func TestWindowDriverShowImage(t *testing.T) {
  t.SkipNow()
	d := initTestWindowDriver()
  _, currentfile, _, _ := runtime.Caller(0)
  image := cv.LoadImage(path.Join(path.Dir(currentfile), "test.png"))

  d.Start()
  d.ShowImage(image)
}
