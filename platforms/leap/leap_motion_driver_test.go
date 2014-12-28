package leap

import (
	"errors"
	"io"
	"io/ioutil"
	"testing"

	"github.com/hybridgroup/gobot"
)

type NullReadWriteCloser struct{}

var writeError error = nil

func (NullReadWriteCloser) Write(p []byte) (int, error) {
	return len(p), writeError
}
func (NullReadWriteCloser) Read(b []byte) (int, error) {
	return len(b), nil
}
func (NullReadWriteCloser) Close() error {
	return nil
}

func initTestLeapMotionDriver() *LeapMotionDriver {
	a := NewLeapMotionAdaptor("bot", "")
	a.connect = func(port string) (io.ReadWriteCloser, error) {
		return &NullReadWriteCloser{}, nil
	}
	a.Connect()
	receive = func(ws io.ReadWriteCloser, buf *[]byte) {
		file, _ := ioutil.ReadFile("./test/support/example_frame.json")
		copy(*buf, file)
	}
	return NewLeapMotionDriver(a, "bot")
}

func TestLeapMotionDriver(t *testing.T) {
	d := initTestLeapMotionDriver()
	gobot.Assert(t, d.Name(), "bot")
	gobot.Assert(t, d.Connection().Name(), "bot")
}
func TestLeapMotionDriverStart(t *testing.T) {
	d := initTestLeapMotionDriver()
	gobot.Assert(t, len(d.Start()), 0)

	d = initTestLeapMotionDriver()
	writeError = errors.New("write error")
	gobot.Assert(t, d.Start()[0], errors.New("write error"))

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
