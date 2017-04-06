package sphero

import (
	"errors"
	"io"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
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
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "Sphero"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestSpheroAdaptor(t *testing.T) {
	a, _ := initTestSpheroAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "Sphero"), true)
	gobottest.Assert(t, a.Port(), "/dev/null")
}

func TestSpheroAdaptorReconnect(t *testing.T) {
	a, _ := initTestSpheroAdaptor()
	a.Connect()
	gobottest.Assert(t, a.connected, true)
	a.Reconnect()
	gobottest.Assert(t, a.connected, true)
	a.Disconnect()
	gobottest.Assert(t, a.connected, false)
	a.Reconnect()
	gobottest.Assert(t, a.connected, true)
}

func TestSpheroAdaptorFinalize(t *testing.T) {
	a, rwc := initTestSpheroAdaptor()
	a.Connect()
	gobottest.Assert(t, a.Finalize(), nil)

	rwc.testAdaptorClose = func() error {
		return errors.New("close error")
	}

	a.connected = true
	gobottest.Assert(t, a.Finalize(), errors.New("close error"))
}

func TestSpheroAdaptorConnect(t *testing.T) {
	a, _ := initTestSpheroAdaptor()
	gobottest.Assert(t, a.Connect(), nil)

	a.connect = func(string) (io.ReadWriteCloser, error) {
		return nil, errors.New("connect error")
	}

	gobottest.Assert(t, a.Connect(), errors.New("connect error"))
}
