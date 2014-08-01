package mavlink

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestMavlinkDriver() *MavlinkDriver {
	return NewMavlinkDriver(NewMavlinkAdaptor("myAdaptor", "/dev/null"), "myDriver")
}

func TestMavlinkDriverStart(t *testing.T) {
	t.SkipNow()
	d := initTestMavlinkDriver()
	gobot.Assert(t, d.Start(), true)
}

func TestMavlinkDriverHalt(t *testing.T) {
	t.SkipNow()
	d := initTestMavlinkDriver()
	gobot.Assert(t, d.Halt(), true)
}
