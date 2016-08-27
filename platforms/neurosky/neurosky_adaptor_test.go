package neurosky

import (
	"errors"
	"io"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gobottest"
)

var _ gobot.Adaptor = (*NeuroskyAdaptor)(nil)

type NullReadWriteCloser struct{}

func (NullReadWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}

var readError error = nil

func (NullReadWriteCloser) Read(b []byte) (int, error) {
	return len(b), readError
}

var closeError error = nil

func (NullReadWriteCloser) Close() error {
	return closeError
}

func initTestNeuroskyAdaptor() *NeuroskyAdaptor {
	a := NewNeuroskyAdaptor("bot", "/dev/null")
	a.connect = func(n *NeuroskyAdaptor) (io.ReadWriteCloser, error) {
		return &NullReadWriteCloser{}, nil
	}
	return a
}

func TestNeuroskyAdaptor(t *testing.T) {
	a := NewNeuroskyAdaptor("bot", "/dev/null")
	gobottest.Assert(t, a.Name(), "bot")
	gobottest.Assert(t, a.Port(), "/dev/null")
}
func TestNeuroskyAdaptorConnect(t *testing.T) {
	a := initTestNeuroskyAdaptor()
	gobottest.Assert(t, len(a.Connect()), 0)

	a.connect = func(n *NeuroskyAdaptor) (io.ReadWriteCloser, error) {
		return nil, errors.New("connection error")
	}
	gobottest.Assert(t, a.Connect()[0], errors.New("connection error"))
}

func TestNeuroskyAdaptorFinalize(t *testing.T) {
	a := initTestNeuroskyAdaptor()
	a.Connect()
	gobottest.Assert(t, len(a.Finalize()), 0)

	closeError = errors.New("close error")
	a.Connect()
	gobottest.Assert(t, a.Finalize()[0], errors.New("close error"))
}
