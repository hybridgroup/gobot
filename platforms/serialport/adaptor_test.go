package serialport

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func initTestAdaptor() (*Adaptor, *nullReadWriteCloser) {
	a := NewAdaptor("/dev/null")
	rwc := NewNullReadWriteCloser()

	a.connect = func(string) (io.ReadWriteCloser, error) {
		return rwc, nil
	}
	return a, rwc
}

func TestNewAdaptor(t *testing.T) {
	a := NewAdaptor("/dev/null")
	assert.True(t, strings.HasPrefix(a.Name(), "Serial"))
	assert.Equal(t, "/dev/null", a.Port())
}

func TestName(t *testing.T) {
	a, _ := initTestAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "Serial"))
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestReconnect(t *testing.T) {
	a, _ := initTestAdaptor()
	require.NoError(t, a.Connect())
	assert.True(t, a.connected)
	require.NoError(t, a.Reconnect())
	assert.True(t, a.connected)
	require.NoError(t, a.Disconnect())
	assert.False(t, a.connected)
	require.NoError(t, a.Reconnect())
	assert.True(t, a.connected)
}

func TestFinalize(t *testing.T) {
	a, rwc := initTestAdaptor()
	require.NoError(t, a.Connect())
	require.NoError(t, a.Finalize())

	rwc.testAdaptorClose = func() error {
		return errors.New("close error")
	}

	a.connected = true
	require.ErrorContains(t, a.Finalize(), "close error")
}

func TestConnect(t *testing.T) {
	a, _ := initTestAdaptor()
	require.NoError(t, a.Connect())

	a.connect = func(string) (io.ReadWriteCloser, error) {
		return nil, errors.New("connect error")
	}

	require.ErrorContains(t, a.Connect(), "connect error")
}
