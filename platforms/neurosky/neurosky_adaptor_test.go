package neurosky

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestNeuroskyAdaptor() *NeuroskyAdaptor {
	a := NewNeuroskyAdaptor("bot", "/dev/null")
	a.connect = func(n *NeuroskyAdaptor) {
		n.sp = gobot.NullReadWriteCloser{}
	}
	return a
}

func TestNeuroskyAdaptorConnect(t *testing.T) {
	a := initTestNeuroskyAdaptor()
	gobot.Assert(t, a.Connect(), true)
}

func TestNeuroskyAdaptorFinalize(t *testing.T) {
	a := initTestNeuroskyAdaptor()
	a.Connect()
	gobot.Assert(t, a.Finalize(), true)
}
