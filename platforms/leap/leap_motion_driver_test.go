package leap

import (
	"errors"
	"io"
	"io/ioutil"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*Driver)(nil)

type NullReadWriteCloser struct{}

var writeError error

func (NullReadWriteCloser) Write(p []byte) (int, error) {
	return len(p), writeError
}
func (NullReadWriteCloser) Read(b []byte) (int, error) {
	return len(b), nil
}
func (NullReadWriteCloser) Close() error {
	return nil
}

func initTestLeapMotionDriver() *Driver {
	a := NewAdaptor("")
	a.connect = func(port string) (io.ReadWriteCloser, error) {
		return &NullReadWriteCloser{}, nil
	}
	a.Connect()
	receive = func(ws io.ReadWriteCloser, buf *[]byte) {
		file, _ := ioutil.ReadFile("./test/support/example_frame.json")
		copy(*buf, file)
	}
	return NewDriver(a)
}

func TestLeapMotionDriver(t *testing.T) {
	d := initTestLeapMotionDriver()
	gobottest.Refute(t, d.Connection(), nil)
}

func TestLeapMotionDriverStart(t *testing.T) {
	d := initTestLeapMotionDriver()
	gobottest.Assert(t, d.Start(), nil)

	d = initTestLeapMotionDriver()
	writeError = errors.New("write error")
	gobottest.Assert(t, d.Start(), errors.New("write error"))
}

func TestLeapMotionDriverHalt(t *testing.T) {
	d := initTestLeapMotionDriver()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestLeapMotionDriverParser(t *testing.T) {
	d := initTestLeapMotionDriver()
	file, _ := ioutil.ReadFile("./test/support/example_frame.json")
	parsedFrame := d.ParseFrame(file)

	if parsedFrame.Hands == nil || parsedFrame.Pointables == nil || parsedFrame.Gestures == nil {
		t.Errorf("ParseFrame incorrectly parsed frame")
	}

	gobottest.Assert(t, parsedFrame.Timestamp, 4729292670)
	gobottest.Assert(t, parsedFrame.Hands[0].X(), 117.546)
	gobottest.Assert(t, parsedFrame.Hands[0].Y(), 236.007)
	gobottest.Assert(t, parsedFrame.Hands[0].Z(), 76.3394)
}
