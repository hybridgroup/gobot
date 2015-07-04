package firmata

import (
	"errors"
	"io"
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata/client"
)

type readWriteCloser struct{}

func (readWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}

var testReadData = []byte{}

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
func (mockFirmataBoard) AnalogWrite(int, int) error   { return nil }
func (mockFirmataBoard) SetPinMode(int, int) error    { return nil }
func (mockFirmataBoard) ReportAnalog(int, int) error  { return nil }
func (mockFirmataBoard) ReportDigital(int, int) error { return nil }
func (mockFirmataBoard) DigitalWrite(int, int) error  { return nil }
func (mockFirmataBoard) I2cRead(int, int) error       { return nil }
func (mockFirmataBoard) I2cWrite(int, []byte) error   { return nil }
func (mockFirmataBoard) I2cConfig(int) error          { return nil }

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
	gobot.Assert(t, a.Name(), "board")
	gobot.Assert(t, a.Port(), "/dev/null")
}

func TestFirmataAdaptorFinalize(t *testing.T) {
	a := initTestFirmataAdaptor()
	gobot.Assert(t, len(a.Finalize()), 0)

	a = initTestFirmataAdaptor()
	a.board.(*mockFirmataBoard).disconnectError = errors.New("close error")
	gobot.Assert(t, a.Finalize()[0], errors.New("close error"))
}

func TestFirmataAdaptorConnect(t *testing.T) {
	var openSP = func(port string) (io.ReadWriteCloser, error) {
		return &readWriteCloser{}, nil
	}
	a := NewFirmataAdaptor("board", "/dev/null")
	a.openSP = openSP
	a.board = newMockFirmataBoard()
	gobot.Assert(t, len(a.Connect()), 0)

	a = NewFirmataAdaptor("board", "/dev/null")
	a.board = newMockFirmataBoard()
	a.openSP = func(port string) (io.ReadWriteCloser, error) {
		return nil, errors.New("connect error")
	}
	gobot.Assert(t, a.Connect()[0], errors.New("connect error"))

	a = NewFirmataAdaptor("board", &readWriteCloser{})
	a.board = newMockFirmataBoard()
	gobot.Assert(t, len(a.Connect()), 0)

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
	gobot.Assert(t, err, nil)
	gobot.Assert(t, val, 1)
}

func TestFirmataAdaptorAnalogRead(t *testing.T) {
	a := initTestFirmataAdaptor()
	val, err := a.AnalogRead("1")
	gobot.Assert(t, val, 133)
	gobot.Assert(t, err, nil)
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
		gobot.Publish(a.board.Event("I2cReply"), i2cReply)
	}()
	data, err := a.I2cRead(0x00, 1)
	gobot.Assert(t, err, nil)
	gobot.Assert(t, data, i)
}
func TestFirmataAdaptorI2cWrite(t *testing.T) {
	a := initTestFirmataAdaptor()
	a.I2cWrite(0x00, []byte{0x00, 0x01})
}
