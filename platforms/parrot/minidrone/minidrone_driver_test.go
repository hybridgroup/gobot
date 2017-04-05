package minidrone

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*Driver)(nil)

func initTestMinidroneDriver() *Driver {
	d := NewDriver(NewBleTestAdaptor())
	return d
}

func TestMinidroneDriver(t *testing.T) {
	d := initTestMinidroneDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Minidrone"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestBatteryDriverStartAndHalt(t *testing.T) {
	d := initTestMinidroneDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}
