package i2c

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var h *HMC6352Driver

func init() {
	h = NewHMC6352Driver(TestAdaptor{}, "bot")
}

func TestHMC6352Start(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, h.Start(), true)
}
