package i2c

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestWiichuckDriver() *WiichuckDriver {
	return NewWiichuckDriver(TestAdaptor{}, "bot")
}

func TestWiichuckDriverStart(t *testing.T) {
	t.SkipNow()
	d := initTestWiichuckDriver()
	gobot.Expect(t, d.Start(), true)
}
