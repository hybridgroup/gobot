package microbit

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*LEDDriver)(nil)

func initTestLEDDriver() *LEDDriver {
	d := NewLEDDriver(NewBleTestAdaptor())
	return d
}

func TestLEDDriver(t *testing.T) {
	d := initTestLEDDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Microbit LED"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestLEDDriverStartAndHalt(t *testing.T) {
	d := initTestLEDDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}
