package sphero

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

type nullReadWriteCloser struct {
	testAdaptorRead  func(p []byte) (int, error)
	testAdaptorWrite func(b []byte) (int, error)
	testAdaptorClose func() error
}

func (n *nullReadWriteCloser) Write(p []byte) (int, error) {
	return n.testAdaptorWrite(p)
}

func (n *nullReadWriteCloser) Read(b []byte) (int, error) {
	return n.testAdaptorRead(b)
}

func (n *nullReadWriteCloser) Close() error {
	return n.testAdaptorClose()
}

func NewNullReadWriteCloser() *nullReadWriteCloser {
	return &nullReadWriteCloser{
		testAdaptorRead: func(p []byte) (int, error) {
			return len(p), nil
		},
		testAdaptorWrite: func(b []byte) (int, error) {
			return len(b), nil
		},
		testAdaptorClose: func() error {
			return nil
		},
	}
}

func initTestSpheroAdaptor() (*Adaptor, *nullReadWriteCloser) {
	a := NewAdaptor("/dev/null")
	rwc := NewNullReadWriteCloser()

	a.connect = func(string) (io.ReadWriteCloser, error) {
		return rwc, nil
	}
	return a, rwc
}

func TestSpheroAdaptorName(t *testing.T) {
	a, _ := initTestSpheroAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "Sphero"))
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestSpheroAdaptor(t *testing.T) {
	a, _ := initTestSpheroAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "Sphero"))
	assert.Equal(t, "/dev/null", a.Port())
}

func TestSpheroAdaptorReconnect(t *testing.T) {
	a, _ := initTestSpheroAdaptor()
	_ = a.Connect()
	assert.True(t, a.connected)
	_ = a.Reconnect()
	assert.True(t, a.connected)
	_ = a.Disconnect()
	assert.False(t, a.connected)
	_ = a.Reconnect()
	assert.True(t, a.connected)
}

func TestSpheroAdaptorFinalize(t *testing.T) {
	a, rwc := initTestSpheroAdaptor()
	_ = a.Connect()
	assert.NoError(t, a.Finalize())

	rwc.testAdaptorClose = func() error {
		return errors.New("close error")
	}

	a.connected = true
	assert.ErrorContains(t, a.Finalize(), "close error")
}

func TestSpheroAdaptorConnect(t *testing.T) {
	a, _ := initTestSpheroAdaptor()
	assert.NoError(t, a.Connect())

	a.connect = func(string) (io.ReadWriteCloser, error) {
		return nil, errors.New("connect error")
	}

	assert.ErrorContains(t, a.Connect(), "connect error")
}
