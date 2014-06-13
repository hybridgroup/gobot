package neurosky

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestNeuroskyAdaptor() *NeuroskyAdaptor {
	return NewNeuroskyAdaptor("bot", "/dev/null")
}

func TestNeuroskyAdaptorFinalize(t *testing.T) {
	t.SkipNow()
	a := initTestNeuroskyAdaptor()
	gobot.Expect(t, a.Finalize(), true)
}
func TestNeuroskyAdaptorConnect(t *testing.T) {
	t.SkipNow()
	a := initTestNeuroskyAdaptor()
	gobot.Expect(t, a.Connect(), true)
}
