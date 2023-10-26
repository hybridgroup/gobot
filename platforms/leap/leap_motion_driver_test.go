package leap

import (
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
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
	_ = a.Connect()

	d := NewDriver(a)
	d.receive = func(ws io.ReadWriteCloser, buf *[]byte) {
		file, _ := os.ReadFile("./test/support/example_frame.json")
		copy(*buf, file)
	}
	return d, rwc
}

func TestLeapMotionDriver(t *testing.T) {
	d, _ := initTestLeapMotionDriver()
	assert.NotNil(t, d.Connection())
}

func TestLeapMotionDriverName(t *testing.T) {
	d, _ := initTestLeapMotionDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "Leap"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestLeapMotionDriverStart(t *testing.T) {
	d, _ := initTestLeapMotionDriver()
	assert.NoError(t, d.Start())

	d2, rwc := initTestLeapMotionDriver()
	e := errors.New("write error")
	rwc.WriteError(e)
	assert.Equal(t, e, d2.Start())
}

func TestLeapMotionDriverHalt(t *testing.T) {
	d, _ := initTestLeapMotionDriver()
	assert.NoError(t, d.Halt())
}

func TestLeapMotionDriverParser(t *testing.T) {
	d, _ := initTestLeapMotionDriver()
	file, _ := os.ReadFile("./test/support/example_frame.json")
	parsedFrame, _ := d.ParseFrame(file)

	if parsedFrame.Hands == nil || parsedFrame.Pointables == nil || parsedFrame.Gestures == nil {
		t.Errorf("ParseFrame incorrectly parsed frame")
	}

	assert.Equal(t, uint64(134211791358), parsedFrame.Timestamp)
	assert.Equal(t, 247.410, parsedFrame.Hands[0].X())
	assert.Equal(t, 275.868, parsedFrame.Hands[0].Y())
	assert.Equal(t, 132.843, parsedFrame.Hands[0].Z())

	assert.Equal(t, 214.293, parsedFrame.Pointables[0].BTipPosition[0])
	assert.Equal(t, 213.865, parsedFrame.Pointables[0].BTipPosition[1])
	assert.Equal(t, 95.0224, parsedFrame.Pointables[0].BTipPosition[2])

	assert.Equal(t, -0.468069, parsedFrame.Pointables[0].Bases[0][0][0])
	assert.Equal(t, 0.807844, parsedFrame.Pointables[0].Bases[0][0][1])
	assert.Equal(t, -0.358190, parsedFrame.Pointables[0].Bases[0][0][2])

	assert.Equal(t, 19.7871, parsedFrame.Pointables[0].Width)
}
