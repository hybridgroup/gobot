package neurosky

import (
	"errors"
	"io"
	"strings"
	"sync"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

type NullReadWriteCloser struct {
	mtx        sync.Mutex
	readError  error
	closeError error
}

// func NewNullReadWriteCloser() *NullReadWriteCloser {
// 	return NullReadWriteCloser{
//
// 	}
// }

func (n *NullReadWriteCloser) ReadError(e error) {
	n.mtx.Lock()
	defer n.mtx.Unlock()
	n.readError = e
}

func (n *NullReadWriteCloser) CloseError(e error) {
	n.mtx.Lock()
	defer n.mtx.Unlock()
	n.closeError = e
}

func (n *NullReadWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}

func (n *NullReadWriteCloser) Read(b []byte) (int, error) {
	n.mtx.Lock()
	defer n.mtx.Unlock()
	return len(b), n.readError
}

func (n *NullReadWriteCloser) Close() error {
	n.mtx.Lock()
	defer n.mtx.Unlock()
	return n.closeError
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

func TestNeuroskyAdaptorName(t *testing.T) {
	a := NewAdaptor("/dev/null")
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "Neurosky"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
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
	rwc := &NullReadWriteCloser{}
	a := NewAdaptor("/dev/null")
	a.connect = func(n *Adaptor) (io.ReadWriteCloser, error) {
		return rwc, nil
	}
	a.Connect()
	gobottest.Assert(t, a.Finalize(), nil)

	rwc.CloseError(errors.New("close error"))
	a.Connect()
	gobottest.Assert(t, a.Finalize(), errors.New("close error"))
}
