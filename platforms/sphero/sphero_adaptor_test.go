package sphero

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestSpheroAdaptor() *SpheroAdaptor {
	a := NewSpheroAdaptor("bot", "/dev/null")
	a.sp = gobot.NullReadWriteCloser{}
	a.connect = func(a *SpheroAdaptor) (err error) { return nil }
	return a
}

func TestSpheroAdaptorFinalize(t *testing.T) {
	a := initTestSpheroAdaptor()
	gobot.Assert(t, len(a.Finalize()), 0)
}
func TestSpheroAdaptorConnect(t *testing.T) {
	a := initTestSpheroAdaptor()
	gobot.Assert(t, len(a.Connect()), 0)
}
