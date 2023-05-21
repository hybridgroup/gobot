package firmata

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/gobottest"
	"gobot.io/x/gobot/v2/platforms/firmata/client"
)

// make sure that this Adaptor fulfills all required analog and digital interfaces
var _ gobot.Adaptor = (*Adaptor)(nil)
var _ gpio.DigitalReader = (*Adaptor)(nil)
var _ gpio.DigitalWriter = (*Adaptor)(nil)
var _ aio.AnalogReader = (*Adaptor)(nil)
var _ gpio.PwmWriter = (*Adaptor)(nil)
var _ gpio.ServoWriter = (*Adaptor)(nil)
var _ FirmataAdaptor = (*Adaptor)(nil)

type readWriteCloser struct{}

func (readWriteCloser) Write(p []byte) (int, error) {
	return testWriteData.Write(p)
}

var testReadData = []byte{}
var testWriteData = bytes.Buffer{}

func (readWriteCloser) Read(b []byte) (int, error) {
	size := len(b)
	if len(testReadData) < size {
		size = len(testReadData)
	}
	copy(b, []byte(testReadData)[:size])
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
	a.Connect()
	return a
}

func TestNewAdaptor(t *testing.T) {
	a := NewAdaptor("/dev/null")
	gobottest.Assert(t, a.Port(), "/dev/null")
}

func TestAdaptorFinalize(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Assert(t, a.Finalize(), nil)

	a = initTestAdaptor()
	a.Board.(*mockFirmataBoard).disconnectError = errors.New("close error")
	gobottest.Assert(t, a.Finalize(), errors.New("close error"))
}

func TestAdaptorConnect(t *testing.T) {
	var openSP = func(port string) (io.ReadWriteCloser, error) {
		return &readWriteCloser{}, nil
	}
	a := NewAdaptor("/dev/null")
	a.PortOpener = openSP
	a.Board = newMockFirmataBoard()
	gobottest.Assert(t, a.Connect(), nil)

	a = NewAdaptor("/dev/null")
	a.Board = newMockFirmataBoard()
	a.PortOpener = func(port string) (io.ReadWriteCloser, error) {
		return nil, errors.New("connect error")
	}
	gobottest.Assert(t, a.Connect(), errors.New("connect error"))

	a = NewAdaptor(&readWriteCloser{})
	a.Board = newMockFirmataBoard()
	gobottest.Assert(t, a.Connect(), nil)

	a = NewAdaptor("/dev/null")
	a.Board = nil
	gobottest.Assert(t, a.Disconnect(), nil)
}

func TestAdaptorServoWrite(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Assert(t, a.ServoWrite("1", 50), nil)
}

func TestAdaptorServoWriteBadPin(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Refute(t, a.ServoWrite("xyz", 50), nil)
}

func TestAdaptorPwmWrite(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Assert(t, a.PwmWrite("1", 50), nil)
}

func TestAdaptorPwmWriteBadPin(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Refute(t, a.PwmWrite("xyz", 50), nil)
}

func TestAdaptorDigitalWrite(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Assert(t, a.DigitalWrite("1", 1), nil)
}

func TestAdaptorDigitalWriteBadPin(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Refute(t, a.DigitalWrite("xyz", 50), nil)
}

func TestAdaptorDigitalRead(t *testing.T) {
	a := initTestAdaptor()
	val, err := a.DigitalRead("1")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, val, 1)

	val, err = a.DigitalRead("0")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, val, 0)
}

func TestAdaptorDigitalReadBadPin(t *testing.T) {
	a := initTestAdaptor()
	_, err := a.DigitalRead("xyz")
	gobottest.Refute(t, err, nil)
}

func TestAdaptorAnalogRead(t *testing.T) {
	a := initTestAdaptor()
	val, err := a.AnalogRead("1")
	gobottest.Assert(t, val, 133)
	gobottest.Assert(t, err, nil)

	val, err = a.AnalogRead("0")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, val, 0)
}

func TestAdaptorAnalogReadBadPin(t *testing.T) {
	a := initTestAdaptor()
	_, err := a.AnalogRead("xyz")
	gobottest.Refute(t, err, nil)
}

func TestServoConfig(t *testing.T) {
	a := initTestAdaptor()
	err := a.ServoConfig("9", 0, 0)
	gobottest.Assert(t, err, nil)

	// test atoi error
	err = a.ServoConfig("a", 0, 0)
	gobottest.Assert(t, true, strings.Contains(fmt.Sprintf("%v", err), "invalid syntax"))
}
