//go:build !windows
// +build !windows

//nolint:forcetypeassert // ok here
package firmata

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/firmata/client"
)

// make sure that this Adaptor fulfills all required analog and digital interfaces
var (
	_ gobot.Adaptor      = (*Adaptor)(nil)
	_ gpio.DigitalReader = (*Adaptor)(nil)
	_ gpio.DigitalWriter = (*Adaptor)(nil)
	_ aio.AnalogReader   = (*Adaptor)(nil)
	_ gpio.PwmWriter     = (*Adaptor)(nil)
	_ gpio.ServoWriter   = (*Adaptor)(nil)
	_ FirmataAdaptor     = (*Adaptor)(nil)
)

type readWriteCloser struct{}

func (readWriteCloser) Write(p []byte) (int, error) {
	return testWriteData.Write(p)
}

var (
	testReadData  = []byte{}
	testWriteData = bytes.Buffer{}
)

func (readWriteCloser) Read(b []byte) (int, error) {
	size := len(b)
	if len(testReadData) < size {
		size = len(testReadData)
	}
	copy(b, testReadData[:size])
	testReadData = testReadData[size:]

	return size, nil
}

func (readWriteCloser) Close() error {
	return nil
}

type mockFirmataBoard struct {
	disconnectError error
	gobot.Eventer
	pins []client.Pin
}

func newMockFirmataBoard() *mockFirmataBoard {
	m := &mockFirmataBoard{
		Eventer:         gobot.NewEventer(),
		disconnectError: nil,
		pins:            make([]client.Pin, 100),
	}

	m.pins[1].Value = 1
	m.pins[15].Value = 133

	return m
}

// setup mock for GPIO, PWM and servo tests
func (mockFirmataBoard) Connect(io.ReadWriteCloser) error { return nil }

func (m mockFirmataBoard) Disconnect() error {
	return m.disconnectError
}

func (m mockFirmataBoard) Pins() []client.Pin {
	return m.pins
}
func (mockFirmataBoard) AnalogWrite(int, int) error      { return nil }
func (mockFirmataBoard) SetPinMode(int, int) error       { return nil }
func (mockFirmataBoard) ReportAnalog(int, int) error     { return nil }
func (mockFirmataBoard) ReportDigital(int, int) error    { return nil }
func (mockFirmataBoard) DigitalWrite(int, int) error     { return nil }
func (mockFirmataBoard) ServoConfig(int, int, int) error { return nil }
func (mockFirmataBoard) WriteSysex([]byte) error         { return nil }

// i2c functions unused in this test scenarios
func (mockFirmataBoard) I2cRead(int, int) error     { return nil }
func (mockFirmataBoard) I2cWrite(int, []byte) error { return nil }
func (mockFirmataBoard) I2cConfig(int) error        { return nil }

func initTestAdaptor() *Adaptor {
	a := NewAdaptor("/dev/null")
	a.Board = newMockFirmataBoard()
	a.PortOpener = func(port string) (io.ReadWriteCloser, error) {
		return &readWriteCloser{}, nil
	}
	_ = a.Connect()
	return a
}

func TestNewAdaptor(t *testing.T) {
	a := NewAdaptor("/dev/null")
	assert.Equal(t, "/dev/null", a.Port())
}

func TestAdaptorFinalize(t *testing.T) {
	a := initTestAdaptor()
	require.NoError(t, a.Finalize())

	a = initTestAdaptor()
	a.Board.(*mockFirmataBoard).disconnectError = errors.New("close error")
	require.ErrorContains(t, a.Finalize(), "close error")
}

func TestAdaptorConnect(t *testing.T) {
	openSP := func(port string) (io.ReadWriteCloser, error) {
		return &readWriteCloser{}, nil
	}
	a := NewAdaptor("/dev/null")
	a.PortOpener = openSP
	a.Board = newMockFirmataBoard()
	require.NoError(t, a.Connect())

	a = NewAdaptor("/dev/null")
	a.Board = newMockFirmataBoard()
	a.PortOpener = func(port string) (io.ReadWriteCloser, error) {
		return nil, errors.New("connect error")
	}
	require.ErrorContains(t, a.Connect(), "connect error")

	a = NewAdaptor(&readWriteCloser{})
	a.Board = newMockFirmataBoard()
	require.NoError(t, a.Connect())

	a = NewAdaptor("/dev/null")
	a.Board = nil
	require.NoError(t, a.Disconnect())
}

func TestAdaptorServoWrite(t *testing.T) {
	a := initTestAdaptor()
	require.NoError(t, a.ServoWrite("1", 50))
}

func TestAdaptorServoWriteBadPin(t *testing.T) {
	a := initTestAdaptor()
	require.Error(t, a.ServoWrite("xyz", 50))
}

func TestAdaptorPwmWrite(t *testing.T) {
	a := initTestAdaptor()
	require.NoError(t, a.PwmWrite("1", 50))
}

func TestAdaptorPwmWriteBadPin(t *testing.T) {
	a := initTestAdaptor()
	require.Error(t, a.PwmWrite("xyz", 50))
}

func TestAdaptorDigitalWrite(t *testing.T) {
	a := initTestAdaptor()
	require.NoError(t, a.DigitalWrite("1", 1))
}

func TestAdaptorDigitalWriteBadPin(t *testing.T) {
	a := initTestAdaptor()
	require.Error(t, a.DigitalWrite("xyz", 50))
}

func TestAdaptorDigitalRead(t *testing.T) {
	a := initTestAdaptor()
	val, err := a.DigitalRead("1")
	require.NoError(t, err)
	assert.Equal(t, 1, val)

	val, err = a.DigitalRead("0")
	require.NoError(t, err)
	assert.Equal(t, 0, val)
}

func TestAdaptorDigitalReadBadPin(t *testing.T) {
	a := initTestAdaptor()
	_, err := a.DigitalRead("xyz")
	require.Error(t, err)
}

func TestAdaptorAnalogRead(t *testing.T) {
	a := initTestAdaptor()
	val, err := a.AnalogRead("1")
	assert.Equal(t, 133, val)
	require.NoError(t, err)

	val, err = a.AnalogRead("0")
	require.NoError(t, err)
	assert.Equal(t, 0, val)
}

func TestAdaptorAnalogReadBadPin(t *testing.T) {
	a := initTestAdaptor()
	_, err := a.AnalogRead("xyz")
	require.Error(t, err)
}

func TestServoConfig(t *testing.T) {
	a := initTestAdaptor()
	err := a.ServoConfig("9", 0, 0)
	require.NoError(t, err)

	// test atoi error
	err = a.ServoConfig("a", 0, 0)
	require.ErrorContains(t, err, "invalid syntax")
}
