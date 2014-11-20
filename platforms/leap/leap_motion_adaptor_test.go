package leap

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestLeapMotionAdaptor() *LeapMotionAdaptor {
	a := NewLeapMotionAdaptor("bot", "")
	a.connect = func(l *LeapMotionAdaptor) (err error) { return nil }
	return a
}

func TestLeapMotionAdaptorConnect(t *testing.T) {
	a := initTestLeapMotionAdaptor()
	gobot.Assert(t, len(a.Connect()), 0)
}

func TestLeapMotionAdaptorFinalize(t *testing.T) {
	a := initTestLeapMotionAdaptor()
	gobot.Assert(t, len(a.Finalize()), 0)
}
