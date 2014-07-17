package leap

import (
	"github.com/hybridgroup/gobot"
	"io/ioutil"
	"testing"
)

func initTestLeapMotionDriver() *LeapMotionDriver {
	return NewLeapMotionDriver(NewLeapMotionAdaptor("bot", "/dev/null"), "bot")
}

func TestLeapMotionDriverStart(t *testing.T) {
	t.SkipNow()
	d := initTestLeapMotionDriver()
	gobot.Assert(t, d.Start(), true)
}

func TestLeapMotionDriverHalt(t *testing.T) {
	t.SkipNow()
	d := initTestLeapMotionDriver()
	gobot.Assert(t, d.Halt(), true)
}

func TestLeapMotionDriverInit(t *testing.T) {
	t.SkipNow()
	d := initTestLeapMotionDriver()
	gobot.Assert(t, d.Init(), true)
}

func TestLeapMotionDriverParser(t *testing.T) {
	d := initTestLeapMotionDriver()
	file, _ := ioutil.ReadFile("./test/support/example_frame.json")
	parsedFrame := d.ParseFrame(file)

	if parsedFrame.Hands == nil || parsedFrame.Pointables == nil || parsedFrame.Gestures == nil {
		t.Errorf("ParseFrame incorrectly parsed frame")
	}

	parsedFrame = d.ParseFrame([]byte{})
	gobot.Assert(t, parsedFrame.Timestamp, 0)
}
