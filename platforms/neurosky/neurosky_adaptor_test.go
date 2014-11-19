package neurosky

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestNeuroskyAdaptor() *NeuroskyAdaptor {
	a := NewNeuroskyAdaptor("bot", "/dev/null")
	a.connect = func(n *NeuroskyAdaptor) (err error) {
		n.sp = gobot.NullReadWriteCloser{}
		return nil
	}
	return a
}

func TestNeuroskyAdaptorConnect(t *testing.T) {
	a := initTestNeuroskyAdaptor()
	gobot.Assert(t, a.Connect(), nil)
}

func TestNeuroskyAdaptorFinalize(t *testing.T) {
	a := initTestNeuroskyAdaptor()
	a.Connect()
	gobot.Assert(t, a.Finalize(), nil)
}
