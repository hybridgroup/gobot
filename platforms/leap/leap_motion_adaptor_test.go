package leap

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestLeapMotionAdaptor() *LeapMotionAdaptor {
	a := NewLeapMotionAdaptor("bot", "")
	a.connect = func(l *LeapMotionAdaptor) {}
	return a
}

func TestLeapMotionAdaptorConnect(t *testing.T) {
	a := initTestLeapMotionAdaptor()
	gobot.Assert(t, a.Connect(), nil)
}

func TestLeapMotionAdaptorFinalize(t *testing.T) {
	a := initTestLeapMotionAdaptor()
	gobot.Assert(t, a.Finalize(), nil)
}
