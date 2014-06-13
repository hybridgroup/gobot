package opencv

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var w *WindowDriver

func init() {
	w = NewWindowDriver("bot")
}

func TestWindowStart(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, w.Start(), true)
}

func TestWindowHalt(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, w.Halt(), true)
}

func TestWindowInit(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, w.Init(), true)
}
