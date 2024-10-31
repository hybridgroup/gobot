package mavlink

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

type nullReadWriteCloser struct{}

var payload = []byte{
	0xFE, 0x09, 0x4E, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x03, 0x51, 0x04, 0x03, 0x1C, 0x7F,
}

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

func initTestMavlinkAdaptor() *Adaptor {
	m := NewAdaptor("/dev/null")
	m.sp = nullReadWriteCloser{}
	m.connect = func(port string) (io.ReadWriteCloser, error) { return nil, nil } //nolint:nilnil // ok for tests
	return m
}

func TestMavlinkAdaptor(t *testing.T) {
	a := initTestMavlinkAdaptor()
	assert.Equal(t, "/dev/null", a.Port())
}

func TestMavlinkAdaptorName(t *testing.T) {
	a := initTestMavlinkAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "Mavlink"))
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestMavlinkAdaptorConnect(t *testing.T) {
	a := initTestMavlinkAdaptor()
	require.NoError(t, a.Connect())

	a.connect = func(port string) (io.ReadWriteCloser, error) { return nil, errors.New("connect error") }
	require.ErrorContains(t, a.Connect(), "connect error")
}

func TestMavlinkAdaptorFinalize(t *testing.T) {
	a := initTestMavlinkAdaptor()
	require.NoError(t, a.Finalize())

	testAdaptorClose = func() error {
		return errors.New("close error")
	}
	require.ErrorContains(t, a.Finalize(), "close error")
}
