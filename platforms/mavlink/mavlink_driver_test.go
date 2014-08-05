package mavlink

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestMavlinkDriver() *MavlinkDriver {
	m := NewMavlinkAdaptor("myAdaptor", "/dev/null")
	m.sp = gobot.NullReadWriteCloser{}
	m.connect = func(a *MavlinkAdaptor) {}
	return NewMavlinkDriver(m, "myDriver")
}

func TestMavlinkDriverStart(t *testing.T) {
	d := initTestMavlinkDriver()
	gobot.Assert(t, d.Start(), true)
}

func TestMavlinkDriverHalt(t *testing.T) {
	d := initTestMavlinkDriver()
	gobot.Assert(t, d.Halt(), true)
}
