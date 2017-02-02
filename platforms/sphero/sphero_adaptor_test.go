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

type nullReadWriteCloser struct{}

var testAdaptorRead = func(p []byte) (int, error) {
	return len(p), nil
}

func (nullReadWriteCloser) Write(p []byte) (int, error) {
	return testAdaptorRead(p)
}

var testAdaptorWrite = func(b []byte) (int, error) {
	return len(b), nil
}

func (nullReadWriteCloser) Read(b []byte) (int, error) {
	return testAdaptorWrite(b)
}

var testAdaptorClose = func() error {
	return nil
}

func (nullReadWriteCloser) Close() error {
	return testAdaptorClose()
}

func initTestSpheroAdaptor() *Adaptor {
	a := NewAdaptor("/dev/null")
	a.connect = func(string) (io.ReadWriteCloser, error) {
		return &nullReadWriteCloser{}, nil
	}
	return a
}

func TestSpheroAdaptor(t *testing.T) {
	a := initTestSpheroAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "Sphero"), true)
	gobottest.Assert(t, a.Port(), "/dev/null")
}

func TestSpheroAdaptorReconnect(t *testing.T) {
	a := initTestSpheroAdaptor()
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
	a := initTestSpheroAdaptor()
	a.Connect()
	gobottest.Assert(t, a.Finalize(), nil)

	testAdaptorClose = func() error {
		return errors.New("close error")
	}

	a.connected = true
	gobottest.Assert(t, a.Finalize(), errors.New("close error"))
}

func TestSpheroAdaptorConnect(t *testing.T) {
	a := initTestSpheroAdaptor()
	gobottest.Assert(t, a.Connect(), nil)

	a.connect = func(string) (io.ReadWriteCloser, error) {
		return nil, errors.New("connect error")
	}

	gobottest.Assert(t, a.Connect(), errors.New("connect error"))
}
