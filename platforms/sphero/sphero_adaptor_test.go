package sphero

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestSpheroAdaptor() *SpheroAdaptor {
	a := NewSpheroAdaptor("bot", "/dev/null")
	a.sp = gobot.NullReadWriteCloser{}
	a.connect = func(a *SpheroAdaptor) {}
	return a
}

func TestSpheroAdaptorFinalize(t *testing.T) {
	a := initTestSpheroAdaptor()
	gobot.Assert(t, a.Finalize(), true)
}
func TestSpheroAdaptorConnect(t *testing.T) {
	a := initTestSpheroAdaptor()
	gobot.Assert(t, a.Connect(), true)
}
