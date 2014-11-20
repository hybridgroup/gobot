package sphero

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestSpheroAdaptor() *SpheroAdaptor {
	a := NewSpheroAdaptor("bot", "/dev/null")
	a.sp = gobot.NullReadWriteCloser{}
	a.connect = func(a *SpheroAdaptor) (err error) { return nil }
	return a
}

func TestSpheroAdaptorFinalize(t *testing.T) {
	a := initTestSpheroAdaptor()
	gobot.Assert(t, a.Finalize(), nil)
}
func TestSpheroAdaptorConnect(t *testing.T) {
	a := initTestSpheroAdaptor()
	gobot.Assert(t, a.Connect(), nil)
}
