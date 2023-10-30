package neurosky

import (
	"errors"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

type NullReadWriteCloser struct {
	mtx        sync.Mutex
	readError  error
	closeError error
}

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
	assert.Equal(t, "/dev/null", a.Port())
}

func TestNeuroskyAdaptorName(t *testing.T) {
	a := NewAdaptor("/dev/null")
	assert.True(t, strings.HasPrefix(a.Name(), "Neurosky"))
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestNeuroskyAdaptorConnect(t *testing.T) {
	a := initTestNeuroskyAdaptor()
	assert.NoError(t, a.Connect())

	a.connect = func(n *Adaptor) (io.ReadWriteCloser, error) {
		return nil, errors.New("connection error")
	}
	assert.ErrorContains(t, a.Connect(), "connection error")
}

func TestNeuroskyAdaptorFinalize(t *testing.T) {
	rwc := &NullReadWriteCloser{}
	a := NewAdaptor("/dev/null")
	a.connect = func(n *Adaptor) (io.ReadWriteCloser, error) {
		return rwc, nil
	}
	_ = a.Connect()
	assert.NoError(t, a.Finalize())

	rwc.CloseError(errors.New("close error"))
	_ = a.Connect()
	assert.ErrorContains(t, a.Finalize(), "close error")
}
