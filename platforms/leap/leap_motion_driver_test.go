package leap

import (
	"io"
	"io/ioutil"
	"testing"

	"github.com/hybridgroup/gobot"
)

type NullReadWriteCloser struct{}

func (NullReadWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}
func (NullReadWriteCloser) Read(b []byte) (int, error) {
	return len(b), nil
}
func (NullReadWriteCloser) Close() error {
	return nil
}

func initTestLeapMotionDriver() *LeapMotionDriver {
	a := NewLeapMotionAdaptor("bot", "")
	a.connect = func(l *LeapMotionAdaptor) (err error) {
		l.ws = new(NullReadWriteCloser)
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
	d := initTestLeapMotionDriver()
	gobot.Assert(t, len(d.Start()), 0)
}

func TestLeapMotionDriverHalt(t *testing.T) {
	d := initTestLeapMotionDriver()
	gobot.Assert(t, len(d.Halt()), 0)
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
