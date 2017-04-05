package microbit

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*TemperatureDriver)(nil)

func initTestTemperatureDriver() *TemperatureDriver {
	d := NewTemperatureDriver(NewBleTestAdaptor())
	return d
}

func TestTemperatureDriver(t *testing.T) {
	d := initTestTemperatureDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Microbit Temperature"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestTemperatureDriverStartAndHalt(t *testing.T) {
	d := initTestTemperatureDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}
