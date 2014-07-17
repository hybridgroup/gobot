package i2c

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestHMC6352Driver() *HMC6352Driver {
	return NewHMC6352Driver(newI2cTestAdaptor("adaptor"), "bot")
}

func TestHMC6352DriverStart(t *testing.T) {
	t.SkipNow()
	d := initTestHMC6352Driver()
	gobot.Assert(t, d.Start(), true)
}
