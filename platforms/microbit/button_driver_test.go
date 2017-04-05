package microbit

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*ButtonDriver)(nil)

func initTestButtonDriver() *ButtonDriver {
	d := NewButtonDriver(NewBleTestAdaptor())
	return d
}

func TestButtonDriver(t *testing.T) {
	d := initTestButtonDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Microbit Button"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestButtonDriverStartAndHalt(t *testing.T) {
	d := initTestButtonDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}
