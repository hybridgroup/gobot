package microbit

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"

	"gobot.io/x/gobot/platforms/ble"
)

var _ gobot.Driver = (*TemperatureDriver)(nil)

func initTestTemperatureDriver() *TemperatureDriver {
	d := NewTemperatureDriver(ble.NewClientAdaptor("D7:99:5A:26:EC:38"))
	return d
}

func TestTemperatureDriver(t *testing.T) {
	d := initTestTemperatureDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Microbit Temperature"), true)
}
