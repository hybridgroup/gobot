package leap

import (
	"github.com/hybridgroup/gobot"
	"io/ioutil"
	"testing"
)

var d *LeapMotionDriver

func init() {
	d = NewLeapMotionDriver(NewLeapMotionAdaptor("bot", "/dev/null"), "bot")
}

func TestStart(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, d.Start(), true)
}

func TestHalt(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, d.Halt(), true)
}

func TestInit(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, d.Init(), true)
}

func TestParser(t *testing.T) {
	file, _ := ioutil.ReadFile("./test/support/example_frame.json")
	parsedFrame := d.ParseFrame(file)

	if parsedFrame.Hands == nil || parsedFrame.Pointables == nil || parsedFrame.Gestures == nil {
		t.Errorf("ParseFrame incorrectly parsed frame")
	}

	parsedFrame = d.ParseFrame([]byte{})
	gobot.Expect(t, parsedFrame.Timestamp, 0)
}
