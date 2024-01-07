package curie

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/firmata"
	"gobot.io/x/gobot/v2/platforms/firmata/client"
)

var _ gobot.Driver = (*IMUDriver)(nil)

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
	require.NoError(t, d.Start())
}

func TestIMUDriverHalt(t *testing.T) {
	d := initTestIMUDriver()
	require.NoError(t, d.Halt())
}

func TestIMUDriverDefaultName(t *testing.T) {
	d := initTestIMUDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "CurieIMU"))
}

func TestIMUDriverSetName(t *testing.T) {
	d := initTestIMUDriver()
	d.SetName("mybot")
	assert.Equal(t, "mybot", d.Name())
}

func TestIMUDriverConnection(t *testing.T) {
	d := initTestIMUDriver()
	assert.NotNil(t, d.Connection())
}

func TestIMUDriverReadAccelerometer(t *testing.T) {
	d := initTestIMUDriver()
	_ = d.Start()
	require.NoError(t, d.ReadAccelerometer())
}

func TestIMUDriverReadAccelerometerData(t *testing.T) {
	_, err := parseAccelerometerData([]byte{})
	require.ErrorContains(t, err, "Invalid data")

	result, err := parseAccelerometerData([]byte{0xF0, 0x11, 0x00, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0xf7})
	require.NoError(t, err)
	assert.Equal(t, &AccelerometerData{X: 1920, Y: 1920, Z: 1920}, result)
}

func TestIMUDriverReadGyroscope(t *testing.T) {
	d := initTestIMUDriver()
	_ = d.Start()
	require.NoError(t, d.ReadGyroscope())
}

func TestIMUDriverReadGyroscopeData(t *testing.T) {
	_, err := parseGyroscopeData([]byte{})
	require.ErrorContains(t, err, "Invalid data")

	result, err := parseGyroscopeData([]byte{0xF0, 0x11, 0x01, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0xf7})
	require.NoError(t, err)
	assert.Equal(t, &GyroscopeData{X: 1920, Y: 1920, Z: 1920}, result)
}

func TestIMUDriverReadTemperature(t *testing.T) {
	d := initTestIMUDriver()
	_ = d.Start()
	require.NoError(t, d.ReadTemperature())
}

func TestIMUDriverReadTemperatureData(t *testing.T) {
	_, err := parseTemperatureData([]byte{})
	require.ErrorContains(t, err, "Invalid data")

	result, err := parseTemperatureData([]byte{0xF0, 0x11, 0x02, 0x00, 0x02, 0x03, 0x04, 0xf7})
	require.NoError(t, err)
	assert.InDelta(t, float32(31.546875), result, 0.0)
}

func TestIMUDriverEnableShockDetection(t *testing.T) {
	d := initTestIMUDriver()
	_ = d.Start()
	require.NoError(t, d.EnableShockDetection(true))
}

func TestIMUDriverShockDetectData(t *testing.T) {
	_, err := parseShockData([]byte{})
	require.ErrorContains(t, err, "Invalid data")

	result, err := parseShockData([]byte{0xF0, 0x11, 0x03, 0x00, 0x02, 0xf7})
	require.NoError(t, err)
	assert.Equal(t, &ShockData{Axis: 0, Direction: 2}, result)
}

func TestIMUDriverEnableStepCounter(t *testing.T) {
	d := initTestIMUDriver()
	_ = d.Start()
	require.NoError(t, d.EnableStepCounter(true))
}

func TestIMUDriverStepCountData(t *testing.T) {
	_, err := parseStepData([]byte{})
	require.ErrorContains(t, err, "Invalid data")

	result, err := parseStepData([]byte{0xF0, 0x11, 0x04, 0x00, 0x02, 0xf7})
	require.NoError(t, err)
	assert.Equal(t, int16(256), result)
}

func TestIMUDriverEnableTapDetection(t *testing.T) {
	d := initTestIMUDriver()
	_ = d.Start()
	require.NoError(t, d.EnableTapDetection(true))
}

func TestIMUDriverTapDetectData(t *testing.T) {
	_, err := parseTapData([]byte{})
	require.ErrorContains(t, err, "Invalid data")

	result, err := parseTapData([]byte{0xF0, 0x11, 0x05, 0x00, 0x02, 0xf7})
	require.NoError(t, err)
	assert.Equal(t, &TapData{Axis: 0, Direction: 2}, result)
}

func TestIMUDriverEnableReadMotion(t *testing.T) {
	d := initTestIMUDriver()
	_ = d.Start()
	require.NoError(t, d.ReadMotion())
}

func TestIMUDriverReadMotionData(t *testing.T) {
	_, err := parseMotionData([]byte{})
	require.ErrorContains(t, err, "Invalid data")

	result, err := parseMotionData([]byte{
		0xF0, 0x11, 0x06, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0xf7,
	})
	require.NoError(t, err)
	assert.Equal(t, &MotionData{AX: 1920, AY: 1920, AZ: 1920, GX: 1920, GY: 1920, GZ: 1920}, result)
}

func TestIMUDriverHandleEvents(t *testing.T) {
	d := initTestIMUDriver()
	_ = d.Start()

	require.NoError(t, d.handleEvent([]byte{0xF0, 0x11, 0x00, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0xf7}))
	require.NoError(t, d.handleEvent([]byte{0xF0, 0x11, 0x01, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0xf7}))
	require.NoError(t, d.handleEvent([]byte{0xF0, 0x11, 0x02, 0x00, 0x02, 0x03, 0x04, 0xf7}))
	require.NoError(t, d.handleEvent([]byte{0xF0, 0x11, 0x03, 0x00, 0x02, 0xf7}))
	require.NoError(t, d.handleEvent([]byte{0xF0, 0x11, 0x04, 0x00, 0x02, 0xf7}))
	require.NoError(t, d.handleEvent([]byte{0xF0, 0x11, 0x05, 0x00, 0x02, 0xf7}))
	require.NoError(t, d.handleEvent([]byte{
		0xF0, 0x11, 0x06, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0x00, 0x0f, 0xf7,
	}))
}
