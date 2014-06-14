package neurosky

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestNeuroskyDriver() *NeuroskyDriver {
	return NewNeuroskyDriver(NewNeuroskyAdaptor("bot", "/dev/null"), "bot")
}

func TestNeuroskyDriverStart(t *testing.T) {
	t.SkipNow()
	d := initTestNeuroskyDriver()
	gobot.Expect(t, d.Start(), true)
}

func TestNeuroskyDriverHalt(t *testing.T) {
	t.SkipNow()
	d := initTestNeuroskyDriver()
	gobot.Expect(t, d.Halt(), true)
}

func TestNeuroskyDriverInit(t *testing.T) {
	t.SkipNow()
	d := initTestNeuroskyDriver()
	gobot.Expect(t, d.Init(), true)
}
