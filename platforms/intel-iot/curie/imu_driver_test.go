package curie

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"

	"gobot.io/x/gobot/platforms/firmata"
	"gobot.io/x/gobot/platforms/firmata/client"
)

var _ gobot.Driver = (*IMUDriver)(nil)

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

func initTestIMUDriver() *IMUDriver {
	a := firmata.NewAdaptor("/dev/null")
	a.Board = newMockFirmataBoard()
	a.PortOpener = func(port string) (io.ReadWriteCloser, error) {
		return &readWriteCloser{}, nil
	}
	return NewIMUDriver(a)
}

func TestIMUDriverStart(t *testing.T) {
	d := initTestIMUDriver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestIMUDriverHalt(t *testing.T) {
	d := initTestIMUDriver()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestIMUDriverDefaultName(t *testing.T) {
	d := initTestIMUDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "CurieIMU"), true)
}

func TestIMUDriverSetName(t *testing.T) {
	d := initTestIMUDriver()
	d.SetName("mybot")
	gobottest.Assert(t, d.Name(), "mybot")
}

func TestIMUDriverConnection(t *testing.T) {
	d := initTestIMUDriver()
	gobottest.Refute(t, d.Connection(), nil)
}

func TestIMUDriverReadAccelerometer(t *testing.T) {
	d := initTestIMUDriver()
	d.Start()
	gobottest.Assert(t, d.ReadAccelerometer(), nil)
}

func TestIMUDriverReadAccelerometerData(t *testing.T) {
	result, err := parseAccelerometerData([]byte{})
	gobottest.Assert(t, err, errors.New("Invalid data"))

	result, err = parseAccelerometerData([]byte{0xF0, 0x11, 0x00, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0xf7})
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, result, &AccelerometerData{X: 1920, Y: 1920, Z: 1920})
}

func TestIMUDriverReadGyroscope(t *testing.T) {
	d := initTestIMUDriver()
	d.Start()
	gobottest.Assert(t, d.ReadGyroscope(), nil)
}

func TestIMUDriverReadGyroscopeData(t *testing.T) {
	result, err := parseGyroscopeData([]byte{})
	gobottest.Assert(t, err, errors.New("Invalid data"))

	result, err = parseGyroscopeData([]byte{0xF0, 0x11, 0x01, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0xf7})
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, result, &GyroscopeData{X: 1920, Y: 1920, Z: 1920})
}

func TestIMUDriverReadTemperature(t *testing.T) {
	d := initTestIMUDriver()
	d.Start()
	gobottest.Assert(t, d.ReadTemperature(), nil)
}

func TestIMUDriverReadTemperatureData(t *testing.T) {
	result, err := parseTemperatureData([]byte{})
	gobottest.Assert(t, err, errors.New("Invalid data"))

	result, err = parseTemperatureData([]byte{0xF0, 0x11, 0x02, 0x00, 0x02, 0x03, 0x04, 0xf7})
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, result, float32(31.546875))
}

func TestIMUDriverEnableShockDetection(t *testing.T) {
	d := initTestIMUDriver()
	d.Start()
	gobottest.Assert(t, d.EnableShockDetection(true), nil)
}

func TestIMUDriverShockDetectData(t *testing.T) {
	result, err := parseShockData([]byte{})
	gobottest.Assert(t, err, errors.New("Invalid data"))

	result, err = parseShockData([]byte{0xF0, 0x11, 0x03, 0x00, 0x02, 0xf7})
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, result, &ShockData{Axis: 0, Direction: 2})
}

func TestIMUDriverEnableStepCounter(t *testing.T) {
	d := initTestIMUDriver()
	d.Start()
	gobottest.Assert(t, d.EnableStepCounter(true), nil)
}

func TestIMUDriverStepCountData(t *testing.T) {
	result, err := parseStepData([]byte{})
	gobottest.Assert(t, err, errors.New("Invalid data"))

	result, err = parseStepData([]byte{0xF0, 0x11, 0x04, 0x00, 0x02, 0xf7})
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, result, int16(256))
}

func TestIMUDriverEnableTapDetection(t *testing.T) {
	d := initTestIMUDriver()
	d.Start()
	gobottest.Assert(t, d.EnableTapDetection(true), nil)
}

func TestIMUDriverTapDetectData(t *testing.T) {
	result, err := parseTapData([]byte{})
	gobottest.Assert(t, err, errors.New("Invalid data"))

	result, err = parseTapData([]byte{0xF0, 0x11, 0x05, 0x00, 0x02, 0xf7})
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, result, &TapData{Axis: 0, Direction: 2})
}

func TestIMUDriverEnableReadMotion(t *testing.T) {
	d := initTestIMUDriver()
	d.Start()
	gobottest.Assert(t, d.ReadMotion(), nil)
}

func TestIMUDriverReadMotionData(t *testing.T) {
	result, err := parseMotionData([]byte{})
	gobottest.Assert(t, err, errors.New("Invalid data"))

	result, err = parseMotionData([]byte{0xF0, 0x11, 0x06, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0xf7})
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, result, &MotionData{AX: 1920, AY: 1920, AZ: 1920, GX: 1920, GY: 1920, GZ: 1920})
}

func TestIMUDriverHandleEvents(t *testing.T) {
	d := initTestIMUDriver()
	d.Start()

	gobottest.Assert(t, d.handleEvent([]byte{0xF0, 0x11, 0x00, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0xf7}), nil)
	gobottest.Assert(t, d.handleEvent([]byte{0xF0, 0x11, 0x01, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0xf7}), nil)
	gobottest.Assert(t, d.handleEvent([]byte{0xF0, 0x11, 0x02, 0x00, 0x02, 0x03, 0x04, 0xf7}), nil)
	gobottest.Assert(t, d.handleEvent([]byte{0xF0, 0x11, 0x03, 0x00, 0x02, 0xf7}), nil)
	gobottest.Assert(t, d.handleEvent([]byte{0xF0, 0x11, 0x04, 0x00, 0x02, 0xf7}), nil)
	gobottest.Assert(t, d.handleEvent([]byte{0xF0, 0x11, 0x05, 0x00, 0x02, 0xf7}), nil)
	gobottest.Assert(t, d.handleEvent([]byte{0xF0, 0x11, 0x06, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0xf7}), nil)
}
