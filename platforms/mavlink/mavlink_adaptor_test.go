package mavlink

import (
	"errors"
	"io"
	"testing"

	"github.com/hybridgroup/gobot"
)

type nullReadWriteCloser struct{}

var payload = []byte{0xFE, 0x09, 0x4E, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x03, 0x51, 0x04, 0x03, 0x1C, 0x7F}

var testAdaptorRead = func(p []byte) (int, error) {
	return len(p), nil
}

func (nullReadWriteCloser) Write(p []byte) (int, error) {
	return testAdaptorRead(p)
}

var testAdaptorWrite = func(b []byte) (int, error) {
	if len(payload) > 0 {
		copy(b, payload[:len(b)])
		payload = payload[len(b):]
		return len(b), nil
	}
	return 0, errors.New("out of bytes")
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

func initTestMavlinkAdaptor() *MavlinkAdaptor {
	m := NewMavlinkAdaptor("myAdaptor", "/dev/null")
	m.sp = nullReadWriteCloser{}
	m.connect = func(port string) (io.ReadWriteCloser, error) { return nil, nil }
	return m
}

func TestMavlinkAdaptor(t *testing.T) {
	a := initTestMavlinkAdaptor()
	gobot.Assert(t, a.Name(), "myAdaptor")
	gobot.Assert(t, a.Port(), "/dev/null")
}
func TestMavlinkAdaptorConnect(t *testing.T) {
	a := initTestMavlinkAdaptor()
	gobot.Assert(t, len(a.Connect()), 0)

	a.connect = func(port string) (io.ReadWriteCloser, error) { return nil, errors.New("connect error") }
	gobot.Assert(t, a.Connect()[0], errors.New("connect error"))
}

func TestMavlinkAdaptorFinalize(t *testing.T) {
	a := initTestMavlinkAdaptor()
	gobot.Assert(t, len(a.Finalize()), 0)

	testAdaptorClose = func() error {
		return errors.New("close error")
	}
	gobot.Assert(t, a.Finalize()[0], errors.New("close error"))
}
