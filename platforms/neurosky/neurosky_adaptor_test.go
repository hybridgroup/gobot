package neurosky

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

type NullReadWriteCloser struct{}

func (NullReadWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}
func (NullReadWriteCloser) Read(b []byte) (int, error) {
	return len(b), nil
}
func (NullReadWriteCloser) Close() error {
	return nil
}

func initTestNeuroskyAdaptor() *NeuroskyAdaptor {
	a := NewNeuroskyAdaptor("bot", "/dev/null")
	a.connect = func(n *NeuroskyAdaptor) (err error) {
		n.sp = NullReadWriteCloser{}
		return nil
	}
	return a
}

func TestNeuroskyAdaptorConnect(t *testing.T) {
	a := initTestNeuroskyAdaptor()
	gobot.Assert(t, len(a.Connect()), 0)
}

func TestNeuroskyAdaptorFinalize(t *testing.T) {
	a := initTestNeuroskyAdaptor()
	a.Connect()
	gobot.Assert(t, len(a.Finalize()), 0)
}
