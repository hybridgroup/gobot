package opencv

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestWindowDriver() *WindowDriver {
	return NewWindowDriver("bot")
}

func TestWindowDriverStart(t *testing.T) {
	t.SkipNow()
	d := initTestWindowDriver()
	gobot.Assert(t, d.Start(), true)
}

func TestWindowDriverHalt(t *testing.T) {
	t.SkipNow()
	d := initTestWindowDriver()
	gobot.Assert(t, d.Halt(), true)
}

func TestWindowDriverInit(t *testing.T) {
	t.SkipNow()
	d := initTestWindowDriver()
	gobot.Assert(t, d.Init(), true)
}
