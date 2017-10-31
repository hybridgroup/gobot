package sprkplus

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*SPRKPlusDriver)(nil)

func initTestSPRKPlusDriver() *SPRKPlusDriver {
	d := NewDriver(NewBleTestAdaptor())
	return d
}

func TestSPRKPlusDriver(t *testing.T) {
	d := initTestSPRKPlusDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "SPRK"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestSPRKPlusDriverStartAndHalt(t *testing.T) {
	d := initTestSPRKPlusDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}
