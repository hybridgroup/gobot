package microbit

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"

	"gobot.io/x/gobot/platforms/ble"
)

var _ gobot.Driver = (*LEDDriver)(nil)

func initTestLEDDriver() *LEDDriver {
	d := NewLEDDriver(ble.NewClientAdaptor("D7:99:5A:26:EC:38"))
	return d
}

func TestLEDDriver(t *testing.T) {
	d := initTestLEDDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Microbit LED"), true)
}
