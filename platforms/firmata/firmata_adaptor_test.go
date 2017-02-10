package firmata

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/platforms/firmata/client"
)

// make sure that this Adaptor fullfills all the required interfaces
var _ gobot.Adaptor = (*Adaptor)(nil)
var _ gpio.DigitalReader = (*Adaptor)(nil)
var _ gpio.DigitalWriter = (*Adaptor)(nil)
var _ aio.AnalogReader = (*Adaptor)(nil)
var _ gpio.PwmWriter = (*Adaptor)(nil)
var _ gpio.ServoWriter = (*Adaptor)(nil)
var _ i2c.Connector = (*Adaptor)(nil)

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

func initTestAdaptor() *Adaptor {
	a := NewAdaptor("/dev/null")
	a.board = newMockFirmataBoard()
	a.PortOpener = func(port string) (io.ReadWriteCloser, error) {
		return &readWriteCloser{}, nil
	}
	a.Connect()
	return a
}

func TestAdaptor(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Assert(t, a.Port(), "/dev/null")
}

func TestAdaptorFinalize(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Assert(t, a.Finalize(), nil)

	a = initTestAdaptor()
	a.board.(*mockFirmataBoard).disconnectError = errors.New("close error")
	gobottest.Assert(t, a.Finalize(), errors.New("close error"))
}

func TestAdaptorConnect(t *testing.T) {
	var openSP = func(port string) (io.ReadWriteCloser, error) {
		return &readWriteCloser{}, nil
	}
	a := NewAdaptor("/dev/null")
	a.PortOpener = openSP
	a.board = newMockFirmataBoard()
	gobottest.Assert(t, a.Connect(), nil)

	a = NewAdaptor("/dev/null")
	a.board = newMockFirmataBoard()
	a.PortOpener = func(port string) (io.ReadWriteCloser, error) {
		return nil, errors.New("connect error")
	}
	gobottest.Assert(t, a.Connect(), errors.New("connect error"))

	a = NewAdaptor(&readWriteCloser{})
	a.board = newMockFirmataBoard()
	gobottest.Assert(t, a.Connect(), nil)

}

func TestAdaptorServoWrite(t *testing.T) {
	a := initTestAdaptor()
	a.ServoWrite("1", 50)
}

func TestAdaptorPwmWrite(t *testing.T) {
	a := initTestAdaptor()
	a.PwmWrite("1", 50)
}

func TestAdaptorDigitalWrite(t *testing.T) {
	a := initTestAdaptor()
	a.DigitalWrite("1", 1)
}

func TestAdaptorDigitalRead(t *testing.T) {
	a := initTestAdaptor()
	val, err := a.DigitalRead("1")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, val, 1)
}

func TestAdaptorAnalogRead(t *testing.T) {
	a := initTestAdaptor()
	val, err := a.AnalogRead("1")
	gobottest.Assert(t, val, 133)
	gobottest.Assert(t, err, nil)
}

func TestAdaptorI2cStart(t *testing.T) {
	a := initTestAdaptor()
	a.GetConnection(0, 0)
}

func TestAdaptorI2cRead(t *testing.T) {
	a := initTestAdaptor()
	i := []byte{100}
	i2cReply := client.I2cReply{Data: i}
	go func() {
		<-time.After(10 * time.Millisecond)
		a.board.Publish(a.board.Event("I2cReply"), i2cReply)
	}()

	con, err := a.GetConnection(0, 0)
	gobottest.Assert(t, err, nil)

	response := []byte{12}
	_, err = con.Read(response)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, response, i)
}

func TestAdaptorI2cWrite(t *testing.T) {
	a := initTestAdaptor()
	con, err := a.GetConnection(0, 0)
	gobottest.Assert(t, err, nil)
	con.Write([]byte{0x00, 0x01})
}

func TestServoConfig(t *testing.T) {
	a := initTestAdaptor()
	err := a.ServoConfig("9", 0, 0)
	gobottest.Assert(t, err, nil)

	// test atoi error
	err = a.ServoConfig("a", 0, 0)
	gobottest.Assert(t, true, strings.Contains(fmt.Sprintf("%v", err), "invalid syntax"))
}
