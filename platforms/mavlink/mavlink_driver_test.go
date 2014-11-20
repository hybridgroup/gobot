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
	gobot.Assert(t, len(d.Start()), 0)
}

func TestMavlinkDriverHalt(t *testing.T) {
	d := initTestMavlinkDriver()
	gobot.Assert(t, len(d.Halt()), 0)
}
