package microbit

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"

	"gobot.io/x/gobot/platforms/ble"
)

var _ gobot.Driver = (*AccelerometerDriver)(nil)

func initTestAccelerometerDriver() *AccelerometerDriver {
	d := NewAccelerometerDriver(ble.NewClientAdaptor("D7:99:5A:26:EC:38"))
	return d
}

func TestAccelerometerDriver(t *testing.T) {
	d := initTestAccelerometerDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Microbit Accelerometer"), true)
}
