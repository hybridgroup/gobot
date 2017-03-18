package microbit

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"

	"gobot.io/x/gobot/platforms/ble"
)

var _ gobot.Driver = (*MagnetometerDriver)(nil)

func initTestMagnetometerDriver() *MagnetometerDriver {
	d := NewMagnetometerDriver(ble.NewClientAdaptor("D7:99:5A:26:EC:38"))
	return d
}

func TestMagnetometerDriver(t *testing.T) {
	d := initTestMagnetometerDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Microbit Magnetometer"), true)
}
