package leap

import (
	"errors"
	"io"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gobottest"
)

var _ gobot.Adaptor = (*LeapMotionAdaptor)(nil)

func initTestLeapMotionAdaptor() *LeapMotionAdaptor {
	a := NewLeapMotionAdaptor("bot", "")
	a.connect = func(port string) (io.ReadWriteCloser, error) { return nil, nil }
	return a
}

func TestLeapMotionAdaptor(t *testing.T) {
	a := NewLeapMotionAdaptor("bot", "127.0.0.1")
	gobottest.Assert(t, a.Name(), "bot")
	gobottest.Assert(t, a.Port(), "127.0.0.1")
}
func TestLeapMotionAdaptorConnect(t *testing.T) {
	a := initTestLeapMotionAdaptor()
	gobottest.Assert(t, len(a.Connect()), 0)

	a.connect = func(port string) (io.ReadWriteCloser, error) {
		return nil, errors.New("connection error")
	}
	gobottest.Assert(t, a.Connect()[0], errors.New("connection error"))
}

func TestLeapMotionAdaptorFinalize(t *testing.T) {
	a := initTestLeapMotionAdaptor()
	gobottest.Assert(t, len(a.Finalize()), 0)
}
