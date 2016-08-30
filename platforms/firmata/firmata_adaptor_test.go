package firmata

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gobottest"
	"github.com/hybridgroup/gobot/platforms/firmata/client"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/i2c"
)

var _ gobot.Adaptor = (*FirmataAdaptor)(nil)

var _ gpio.DigitalReader = (*FirmataAdaptor)(nil)
var _ gpio.DigitalWriter = (*FirmataAdaptor)(nil)
var _ gpio.AnalogReader = (*FirmataAdaptor)(nil)
var _ gpio.PwmWriter = (*FirmataAdaptor)(nil)
var _ gpio.ServoWriter = (*FirmataAdaptor)(nil)

var _ i2c.I2c = (*FirmataAdaptor)(nil)

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

	m.AddEvent("I2cReply")
	return m
}

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
func (mockFirmataBoard) I2cRead(int, int) error          { return nil }
func (mockFirmataBoard) I2cWrite(int, []byte) error      { return nil }
func (mockFirmataBoard) I2cConfig(int) error             { return nil }
func (mockFirmataBoard) ServoConfig(int, int, int) error { return nil }

func initTestFirmataAdaptor() *FirmataAdaptor {
	a := NewFirmataAdaptor("board", "/dev/null")
	a.board = newMockFirmataBoard()
	a.openSP = func(port string) (io.ReadWriteCloser, error) {
		return &readWriteCloser{}, nil
	}
	a.Connect()
	return a
}

func TestFirmataAdaptor(t *testing.T) {
	a := initTestFirmataAdaptor()
	gobottest.Assert(t, a.Name(), "board")
	gobottest.Assert(t, a.Port(), "/dev/null")
}

func TestFirmataAdaptorFinalize(t *testing.T) {
	a := initTestFirmataAdaptor()
	gobottest.Assert(t, len(a.Finalize()), 0)

	a = initTestFirmataAdaptor()
	a.board.(*mockFirmataBoard).disconnectError = errors.New("close error")
	gobottest.Assert(t, a.Finalize()[0], errors.New("close error"))
}

func TestFirmataAdaptorConnect(t *testing.T) {
	var openSP = func(port string) (io.ReadWriteCloser, error) {
		return &readWriteCloser{}, nil
	}
	a := NewFirmataAdaptor("board", "/dev/null")
	a.openSP = openSP
	a.board = newMockFirmataBoard()
	gobottest.Assert(t, len(a.Connect()), 0)

	a = NewFirmataAdaptor("board", "/dev/null")
	a.board = newMockFirmataBoard()
	a.openSP = func(port string) (io.ReadWriteCloser, error) {
		return nil, errors.New("connect error")
	}
	gobottest.Assert(t, a.Connect()[0], errors.New("connect error"))

	a = NewFirmataAdaptor("board", &readWriteCloser{})
	a.board = newMockFirmataBoard()
	gobottest.Assert(t, len(a.Connect()), 0)

}

func TestFirmataAdaptorServoWrite(t *testing.T) {
	a := initTestFirmataAdaptor()
	a.ServoWrite("1", 50)
}

func TestFirmataAdaptorPwmWrite(t *testing.T) {
	a := initTestFirmataAdaptor()
	a.PwmWrite("1", 50)
}

func TestFirmataAdaptorDigitalWrite(t *testing.T) {
	a := initTestFirmataAdaptor()
	a.DigitalWrite("1", 1)
}

func TestFirmataAdaptorDigitalRead(t *testing.T) {
	a := initTestFirmataAdaptor()
	val, err := a.DigitalRead("1")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, val, 1)
}

func TestFirmataAdaptorAnalogRead(t *testing.T) {
	a := initTestFirmataAdaptor()
	val, err := a.AnalogRead("1")
	gobottest.Assert(t, val, 133)
	gobottest.Assert(t, err, nil)
}

func TestFirmataAdaptorI2cStart(t *testing.T) {
	a := initTestFirmataAdaptor()
	a.I2cStart(0x00)
}
func TestFirmataAdaptorI2cRead(t *testing.T) {
	a := initTestFirmataAdaptor()
	i := []byte{100}
	i2cReply := client.I2cReply{Data: i}
	go func() {
		<-time.After(10 * time.Millisecond)
		a.Publish(a.board.Event("I2cReply"), i2cReply)
	}()
	data, err := a.I2cRead(0x00, 1)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, data, i)
}
func TestFirmataAdaptorI2cWrite(t *testing.T) {
	a := initTestFirmataAdaptor()
	a.I2cWrite(0x00, []byte{0x00, 0x01})
}

func TestServoConfig(t *testing.T) {
	a := initTestFirmataAdaptor()
	err := a.ServoConfig("9", 0, 0)
	gobottest.Assert(t, err, nil)

	// test atoi error
	err = a.ServoConfig("a", 0, 0)
	gobottest.Assert(t, true, strings.Contains(fmt.Sprintf("%v", err), "invalid syntax"))
}
