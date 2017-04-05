package microbit

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*MagnetometerDriver)(nil)

func initTestMagnetometerDriver() *MagnetometerDriver {
	d := NewMagnetometerDriver(NewBleTestAdaptor())
	return d
}

func TestMagnetometerDriver(t *testing.T) {
	d := initTestMagnetometerDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Microbit Magnetometer"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestMagnetometerDriverStartAndHalt(t *testing.T) {
	d := initTestMagnetometerDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}
