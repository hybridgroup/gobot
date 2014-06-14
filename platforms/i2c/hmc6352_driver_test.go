package i2c

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestHMC6352Driver() *HMC6352Driver {
	return NewHMC6352Driver(TestAdaptor{}, "bot")
}

func TestHMC6352DriverStart(t *testing.T) {
	t.SkipNow()
	d := initTestHMC6352Driver()
	gobot.Expect(t, d.Start(), true)
}
