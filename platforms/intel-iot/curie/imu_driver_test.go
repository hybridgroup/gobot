package curie

import (
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"

	"gobot.io/x/gobot/platforms/firmata"
)

var _ gobot.Driver = (*IMUDriver)(nil)

func initTestIMUDriver() *IMUDriver {
	return NewIMUDriver(firmata.NewAdaptor("/dev/null"))
}

func TestIMUDriverHalt(t *testing.T) {
	d := initTestIMUDriver()
	gobottest.Assert(t, d.Halt(), nil)
}
