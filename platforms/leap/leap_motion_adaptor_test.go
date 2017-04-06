package leap

import (
	"errors"
	"io"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

func initTestLeapMotionAdaptor() *Adaptor {
	a := NewAdaptor("")
	a.connect = func(port string) (io.ReadWriteCloser, error) { return nil, nil }
	return a
}

func TestLeapMotionAdaptor(t *testing.T) {
	a := NewAdaptor("127.0.0.1")
	gobottest.Assert(t, a.Port(), "127.0.0.1")
}

func TestLeapMotionAdaptorName(t *testing.T) {
	a := NewAdaptor("127.0.0.1")
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "Leap"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestLeapMotionAdaptorConnect(t *testing.T) {
	a := initTestLeapMotionAdaptor()
	gobottest.Assert(t, a.Connect(), nil)

	a.connect = func(port string) (io.ReadWriteCloser, error) {
		return nil, errors.New("connection error")
	}
	gobottest.Assert(t, a.Connect(), errors.New("connection error"))
}

func TestLeapMotionAdaptorFinalize(t *testing.T) {
	a := initTestLeapMotionAdaptor()
	gobottest.Assert(t, a.Finalize(), nil)
}
