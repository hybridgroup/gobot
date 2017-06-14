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
func (mockFirmataBoard) WriteSysex(data []byte) error    { return nil }

func initTestAdaptor() *Adaptor {
	a := NewAdaptor("/dev/null")
	a.Board = newMockFirmataBoard()
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

func TestAdaptorI2cStart(t *testing.T) {
	a := initTestAdaptor()
	i2c, err := a.GetConnection(0, 0)
	gobottest.Assert(t, err, nil)
	gobottest.Refute(t, i2c, nil)
	gobottest.Assert(t, i2c.Close(), nil)
}

func TestAdaptorI2cRead(t *testing.T) {
	a := initTestAdaptor()
	i := []byte{100}
	i2cReply := client.I2cReply{Data: i}
	go func() {
		<-time.After(10 * time.Millisecond)
		a.Board.Publish(a.Board.Event("I2cReply"), i2cReply)
	}()

	con, err := a.GetConnection(0, 0)
	gobottest.Assert(t, err, nil)

	response := []byte{12}
	_, err = con.Read(response)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, response, i)
}

func TestAdaptorI2cReadByte(t *testing.T) {
	a := initTestAdaptor()
	i := []byte{100}
	i2cReply := client.I2cReply{Data: i}
	go func() {
		<-time.After(10 * time.Millisecond)
		a.Board.Publish(a.Board.Event("I2cReply"), i2cReply)
	}()

	con, err := a.GetConnection(0, 0)
	gobottest.Assert(t, err, nil)

	var val byte
	val, err = con.ReadByte()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, val, uint8(100))
}

func TestAdaptorI2cReadByteData(t *testing.T) {
	a := initTestAdaptor()
	i := []byte{100}
	i2cReply := client.I2cReply{Data: i}
	go func() {
		<-time.After(10 * time.Millisecond)
		a.Board.Publish(a.Board.Event("I2cReply"), i2cReply)
	}()

	con, err := a.GetConnection(0, 0)
	gobottest.Assert(t, err, nil)

	var val byte
	val, err = con.ReadByteData(0x01)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, val, uint8(100))
}

func TestAdaptorI2cReadWordData(t *testing.T) {
	a := initTestAdaptor()
	i := []byte{100}
	i2cReply := client.I2cReply{Data: i}
	go func() {
		<-time.After(10 * time.Millisecond)
		a.Board.Publish(a.Board.Event("I2cReply"), i2cReply)
	}()

	con, err := a.GetConnection(0, 0)
	gobottest.Assert(t, err, nil)

	var val uint16
	val, err = con.ReadWordData(0x01)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, val, uint16(100))
}

func TestAdaptorI2cWrite(t *testing.T) {
	a := initTestAdaptor()
	con, err := a.GetConnection(0, 0)
	gobottest.Assert(t, err, nil)
	written, _ := con.Write([]byte{0x00, 0x01})
	gobottest.Assert(t, written, 2)
}

func TestAdaptorI2cWrite20bytes(t *testing.T) {
	a := initTestAdaptor()
	con, err := a.GetConnection(0, 0)
	gobottest.Assert(t, err, nil)
	written, _ := con.Write([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19})
	gobottest.Assert(t, written, 20)
}

func TestAdaptorI2cWriteByte(t *testing.T) {
	a := initTestAdaptor()
	con, err := a.GetConnection(0, 0)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, con.WriteByte(0x00), nil)
}

func TestAdaptorI2cWriteByteData(t *testing.T) {
	a := initTestAdaptor()
	con, err := a.GetConnection(0, 0)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, con.WriteByteData(0x00, 0x02), nil)
}

func TestAdaptorI2cWriteWordData(t *testing.T) {
	a := initTestAdaptor()
	con, err := a.GetConnection(0, 0)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, con.WriteWordData(0x00, 0x02), nil)
}

func TestAdaptorI2cWriteBlockData(t *testing.T) {
	a := initTestAdaptor()
	con, err := a.GetConnection(0, 0)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, con.WriteBlockData(0x00, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}), nil)
}

func TestServoConfig(t *testing.T) {
	a := initTestAdaptor()
	err := a.ServoConfig("9", 0, 0)
	gobottest.Assert(t, err, nil)

	// test atoi error
	err = a.ServoConfig("a", 0, 0)
	gobottest.Assert(t, true, strings.Contains(fmt.Sprintf("%v", err), "invalid syntax"))
}

func TestDefaultBus(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Assert(t, a.GetDefaultBus(), 0)
}

func TestGetConnectionInvalidBus(t *testing.T) {
	a := initTestAdaptor()
	_, err := a.GetConnection(0x01, 99)
	gobottest.Assert(t, err, errors.New("Invalid bus number 99, only 0 is supported"))
}
