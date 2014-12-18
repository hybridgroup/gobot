package leap

import (
	"errors"
	"io"
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestLeapMotionAdaptor() *LeapMotionAdaptor {
	a := NewLeapMotionAdaptor("bot", "")
	a.connect = func(port string) (io.ReadWriteCloser, error) { return nil, nil }
	return a
}

func TestLeapMotionAdaptor(t *testing.T) {
	a := NewLeapMotionAdaptor("bot", "127.0.0.1")
	gobot.Assert(t, a.Name(), "bot")
	gobot.Assert(t, a.Port(), "127.0.0.1")
}
func TestLeapMotionAdaptorConnect(t *testing.T) {
	a := initTestLeapMotionAdaptor()
	gobot.Assert(t, len(a.Connect()), 0)

	a.connect = func(port string) (io.ReadWriteCloser, error) {
		return nil, errors.New("connection error")
	}
	gobot.Assert(t, a.Connect()[0], errors.New("connection error"))
}

func TestLeapMotionAdaptorFinalize(t *testing.T) {
	a := initTestLeapMotionAdaptor()
	gobot.Assert(t, len(a.Finalize()), 0)
}
