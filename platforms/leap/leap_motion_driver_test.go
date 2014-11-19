package leap

import (
	"io"
	"io/ioutil"
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestLeapMotionDriver() *LeapMotionDriver {
	a := NewLeapMotionAdaptor("bot", "")
	a.connect = func(l *LeapMotionAdaptor) (err error) {
		l.ws = new(gobot.NullReadWriteCloser)
		return nil
	}
	a.Connect()
	receive = func(ws io.ReadWriteCloser) []byte {
		file, _ := ioutil.ReadFile("./test/support/example_frame.json")
		return file
	}
	return NewLeapMotionDriver(a, "bot")
}

func TestLeapMotionDriverStart(t *testing.T) {
	//t.SkipNow()
	d := initTestLeapMotionDriver()
	gobot.Assert(t, d.Start(), nil)
}

func TestLeapMotionDriverHalt(t *testing.T) {
	d := initTestLeapMotionDriver()
	gobot.Assert(t, d.Halt(), nil)
}

func TestLeapMotionDriverParser(t *testing.T) {
	d := initTestLeapMotionDriver()
	file, _ := ioutil.ReadFile("./test/support/example_frame.json")
	parsedFrame := d.ParseFrame(file)

	if parsedFrame.Hands == nil || parsedFrame.Pointables == nil || parsedFrame.Gestures == nil {
		t.Errorf("ParseFrame incorrectly parsed frame")
	}

	gobot.Assert(t, parsedFrame.Timestamp, 4729292670)
	gobot.Assert(t, parsedFrame.Hands[0].X(), 117.546)
	gobot.Assert(t, parsedFrame.Hands[0].Y(), 236.007)
	gobot.Assert(t, parsedFrame.Hands[0].Z(), 76.3394)
}
