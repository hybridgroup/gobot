package opencv

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestWindowDriver() *WindowDriver {
	return NewWindowDriver("bot")
}

func TestWindowDriverStart(t *testing.T) {
	d := initTestWindowDriver()
	gobot.Expect(t, d.Start(), true)
}

func TestWindowDriverHalt(t *testing.T) {
	d := initTestWindowDriver()
	gobot.Expect(t, d.Halt(), true)
}

func TestWindowDriverInit(t *testing.T) {
	d := initTestWindowDriver()
	gobot.Expect(t, d.Init(), true)
}
