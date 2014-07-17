package i2c

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestWiichuckDriver() *WiichuckDriver {
	return NewWiichuckDriver(newI2cTestAdaptor("adaptor"), "bot")
}

func TestWiichuckDriverStart(t *testing.T) {
	t.SkipNow()
	d := initTestWiichuckDriver()
	gobot.Assert(t, d.Start(), true)
}
