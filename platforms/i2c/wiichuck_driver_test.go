package i2c

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var w *WiichuckDriver

func init() {
	w = NewWiichuckDriver(TestAdaptor{}, "bot")
}

func TestWiichuckStart(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, w.Start(), true)
}
