package sphero

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestSpheroDriver() *SpheroDriver {
	a := NewSpheroAdaptor("bot", "/dev/null")
	a.sp = gobot.NullReadWriteCloser{}
	return NewSpheroDriver(a, "bot")
}

func TestSpheroDriverStart(t *testing.T) {
	d := initTestSpheroDriver()
	gobot.Assert(t, d.Start(), true)
}

func TestSpheroDriverHalt(t *testing.T) {
	d := initTestSpheroDriver()
	gobot.Assert(t, d.Halt(), true)
}
