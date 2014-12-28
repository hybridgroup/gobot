package sphero

import (
	"errors"
	"io"
	"testing"

	"github.com/hybridgroup/gobot"
)

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

func initTestSpheroAdaptor() *SpheroAdaptor {
	a := NewSpheroAdaptor("bot", "/dev/null")
	a.connect = func(string) (io.ReadWriteCloser, error) {
		return &nullReadWriteCloser{}, nil
	}
	return a
}

func TestSpheroAdaptor(t *testing.T) {
	a := initTestSpheroAdaptor()
	gobot.Assert(t, a.Name(), "bot")
	gobot.Assert(t, a.Port(), "/dev/null")
}

func TestSpheroAdaptorReconnect(t *testing.T) {
	a := initTestSpheroAdaptor()
	a.Connect()
	gobot.Assert(t, a.connected, true)
	a.Reconnect()
	gobot.Assert(t, a.connected, true)
	a.Disconnect()
	gobot.Assert(t, a.connected, false)
	a.Reconnect()
	gobot.Assert(t, a.connected, true)
}

func TestSpheroAdaptorFinalize(t *testing.T) {
	a := initTestSpheroAdaptor()
	a.Connect()
	gobot.Assert(t, len(a.Finalize()), 0)

	testAdaptorClose = func() error {
		return errors.New("close error")
	}

	a.connected = true
	gobot.Assert(t, a.Finalize()[0], errors.New("close error"))
}

func TestSpheroAdaptorConnect(t *testing.T) {
	a := initTestSpheroAdaptor()
	gobot.Assert(t, len(a.Connect()), 0)

	a.connect = func(string) (io.ReadWriteCloser, error) {
		return nil, errors.New("connect error")
	}

	gobot.Assert(t, a.Connect()[0], errors.New("connect error"))
}
