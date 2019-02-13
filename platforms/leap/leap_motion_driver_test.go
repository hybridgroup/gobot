package leap

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"sync"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*Driver)(nil)

type NullReadWriteCloser struct {
	mtx        sync.Mutex
	writeError error
}

func (n *NullReadWriteCloser) WriteError(e error) {
	n.mtx.Lock()
	defer n.mtx.Unlock()
	n.writeError = e
}

func (n *NullReadWriteCloser) Write(p []byte) (int, error) {
	n.mtx.Lock()
	defer n.mtx.Unlock()
	return len(p), n.writeError
}
func (n *NullReadWriteCloser) Read(b []byte) (int, error) {
	return len(b), nil
}
func (n *NullReadWriteCloser) Close() error {
	return nil
}

func initTestLeapMotionDriver() (*Driver, *NullReadWriteCloser) {
	a := NewAdaptor("")
	rwc := &NullReadWriteCloser{}
	a.connect = func(port string) (io.ReadWriteCloser, error) {
		return rwc, nil
	}
	a.Connect()

	d := NewDriver(a)
	d.receive = func(ws io.ReadWriteCloser, buf *[]byte) {
		file, _ := ioutil.ReadFile("./test/support/example_frame.json")
		copy(*buf, file)
	}
	return d, rwc
}

func TestLeapMotionDriver(t *testing.T) {
	d, _ := initTestLeapMotionDriver()
	gobottest.Refute(t, d.Connection(), nil)
}

func TestLeapMotionDriverName(t *testing.T) {
	d, _ := initTestLeapMotionDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Leap"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestLeapMotionDriverStart(t *testing.T) {
	d, _ := initTestLeapMotionDriver()
	gobottest.Assert(t, d.Start(), nil)

	d2, rwc := initTestLeapMotionDriver()
	e := errors.New("write error")
	rwc.WriteError(e)
	gobottest.Assert(t, d2.Start(), e)
}

func TestLeapMotionDriverHalt(t *testing.T) {
	d, _ := initTestLeapMotionDriver()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestLeapMotionDriverParser(t *testing.T) {
	d, _ := initTestLeapMotionDriver()
	file, _ := ioutil.ReadFile("./test/support/example_frame.json")
	parsedFrame := d.ParseFrame(file)

	if parsedFrame.Hands == nil || parsedFrame.Pointables == nil || parsedFrame.Gestures == nil {
		t.Errorf("ParseFrame incorrectly parsed frame")
	}

	gobottest.Assert(t, parsedFrame.Timestamp, uint64(134211791358))
	gobottest.Assert(t, parsedFrame.Hands[0].X(), 247.410)
	gobottest.Assert(t, parsedFrame.Hands[0].Y(), 275.868)
	gobottest.Assert(t, parsedFrame.Hands[0].Z(), 132.843)

	gobottest.Assert(t, parsedFrame.Pointables[0].BTipPosition[0], 214.293)
	gobottest.Assert(t, parsedFrame.Pointables[0].BTipPosition[1], 213.865)
	gobottest.Assert(t, parsedFrame.Pointables[0].BTipPosition[2], 95.0224)

	gobottest.Assert(t, parsedFrame.Pointables[0].Bases[0][0][0], -0.468069)
	gobottest.Assert(t, parsedFrame.Pointables[0].Bases[0][0][1], 0.807844)
	gobottest.Assert(t, parsedFrame.Pointables[0].Bases[0][0][2], -0.358190)

	gobottest.Assert(t, parsedFrame.Pointables[0].Width, 19.7871)
}
