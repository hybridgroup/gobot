package mavlink

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestMavlinkDriver() *MavlinkDriver {
	m := NewMavlinkAdaptor("myAdaptor", "/dev/null")
	m.sp = gobot.NullReadWriteCloser{}
	m.connect = func(a *MavlinkAdaptor) (err error) { return nil }
	return NewMavlinkDriver(m, "myDriver")
}

func TestMavlinkDriverStart(t *testing.T) {
	d := initTestMavlinkDriver()
	gobot.Assert(t, d.Start(), nil)
}

func TestMavlinkDriverHalt(t *testing.T) {
	d := initTestMavlinkDriver()
	gobot.Assert(t, d.Halt(), nil)
}
