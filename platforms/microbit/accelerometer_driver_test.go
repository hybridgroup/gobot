package microbit

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*AccelerometerDriver)(nil)

func initTestAccelerometerDriver() *AccelerometerDriver {
	d := NewAccelerometerDriver(NewBleTestAdaptor())
	return d
}

func TestAccelerometerDriver(t *testing.T) {
	d := initTestAccelerometerDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Microbit Accelerometer"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestAccelerometerDriverStartAndHalt(t *testing.T) {
	d := initTestAccelerometerDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}
