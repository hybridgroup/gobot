package neurosky

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var n *NeuroskyDriver

func init() {
	n = NewNeuroskyDriver(NewNeuroskyAdaptor("bot", "/dev/null"), "bot")
}

func TestStart(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, n.Start(), true)
}

func TestHalt(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, n.Halt(), true)
}

func TestInit(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, n.Init(), true)
}
