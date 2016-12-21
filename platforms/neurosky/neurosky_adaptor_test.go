package neurosky

import (
	"errors"
	"io"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

type NullReadWriteCloser struct{}

func (NullReadWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}

var readError error = nil

func (NullReadWriteCloser) Read(b []byte) (int, error) {
	return len(b), readError
}

var closeError error

func (NullReadWriteCloser) Close() error {
	return closeError
}

func initTestNeuroskyAdaptor() *Adaptor {
	a := NewAdaptor("/dev/null")
	a.connect = func(n *Adaptor) (io.ReadWriteCloser, error) {
		return &NullReadWriteCloser{}, nil
	}
	return a
}

func TestNeuroskyAdaptor(t *testing.T) {
	a := NewAdaptor("/dev/null")
	gobottest.Assert(t, a.Port(), "/dev/null")
}

func TestNeuroskyAdaptorConnect(t *testing.T) {
	a := initTestNeuroskyAdaptor()
	gobottest.Assert(t, a.Connect(), nil)

	a.connect = func(n *Adaptor) (io.ReadWriteCloser, error) {
		return nil, errors.New("connection error")
	}
	gobottest.Assert(t, a.Connect(), errors.New("connection error"))
}

func TestNeuroskyAdaptorFinalize(t *testing.T) {
	a := initTestNeuroskyAdaptor()
	a.Connect()
	gobottest.Assert(t, a.Finalize(), nil)

	closeError = errors.New("close error")
	a.Connect()
	gobottest.Assert(t, a.Finalize(), errors.New("close error"))
}
