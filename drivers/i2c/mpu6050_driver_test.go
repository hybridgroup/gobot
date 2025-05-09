package i2c

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*MPU6050Driver)(nil)

func initTestMPU6050WithStubbedAdaptor() (*MPU6050Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewMPU6050Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewMPU6050Driver(t *testing.T) {
	var di interface{} = NewMPU6050Driver(newI2cTestAdaptor())
	d, ok := di.(*MPU6050Driver)
	if !ok {
		require.Fail(t, "NewMPU6050Driver() should have returned a *MPU6050Driver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.name, "MPU6050"))
	assert.Equal(t, 0x68, d.defaultAddress)
	assert.Equal(t, MPU6050DlpfConfig(0x00), d.dlpf)
	assert.Equal(t, MPU6050FrameSyncConfig(0x00), d.frameSync)
	assert.Equal(t, MPU6050AccelFsConfig(0x00), d.accelFs)
	assert.Equal(t, MPU6050GyroFsConfig(0x00), d.gyroFs)
	assert.Equal(t, MPU6050Pwr1ClockConfig(0x01), d.clock)
	assert.InDelta(t, 9.80665, d.gravity, 0.0)
}

func TestMPU6050Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewMPU6050Driver(newI2cTestAdaptor(), WithBus(2), WithMPU6050DigitalFilter(0x06))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
	assert.Equal(t, MPU6050DlpfConfig(0x06), d.dlpf)
}

func TestWithMPU6050FrameSync(t *testing.T) {
	d := NewMPU6050Driver(newI2cTestAdaptor(), WithMPU6050FrameSync(0x07))
	assert.Equal(t, MPU6050FrameSyncConfig(0x07), d.frameSync)
}

func TestWithMPU6050AccelFullScaleRange(t *testing.T) {
	d := NewMPU6050Driver(newI2cTestAdaptor(), WithMPU6050AccelFullScaleRange(0x02))
	assert.Equal(t, MPU6050AccelFsConfig(0x02), d.accelFs)
}

func TestWithMPU6050GyroFullScaleRange(t *testing.T) {
	d := NewMPU6050Driver(newI2cTestAdaptor(), WithMPU6050GyroFullScaleRange(0x03))
	assert.Equal(t, MPU6050GyroFsConfig(0x03), d.gyroFs)
}

func TestWithMPU6050ClockSource(t *testing.T) {
	d := NewMPU6050Driver(newI2cTestAdaptor(), WithMPU6050ClockSource(0x05))
	assert.Equal(t, MPU6050Pwr1ClockConfig(0x05), d.clock)
}

func TestWithMPU6050Gravity(t *testing.T) {
	d := NewMPU6050Driver(newI2cTestAdaptor(), WithMPU6050Gravity(1.0))
	assert.InDelta(t, 1.0, d.gravity, 0.0)
}

func TestMPU6050GetData(t *testing.T) {
	// sequence to read values
	// * reset device and prepare config registers, see test for Start()
	// * write first data register address (0x3B)
	// * read 3 x 2 bytes acceleration.X, Y, Z data, little-endian (MSB, LSB)
	// * read 2 bytes temperature data, little-endian (MSB, LSB)
	// * read 3 x 2 bytes gyroscope.X, Y, Z data, little-endian (MSB, LSB)
	// * scale
	//    Acceleration: raw value / gain * standard gravity [m/s²]
	//    Temperature:  raw value / 340 + 36.53 [°C]
	//    Gyroscope:    raw value / gain [°/s]

	// arrange
	d, adaptor := initTestMPU6050WithStubbedAdaptor()
	_ = d.Start()

	accData := []byte{0x00, 0x01, 0x02, 0x04, 0x08, 0x16}
	tempData := []byte{0x32, 0x64}
	gyroData := []byte{0x16, 0x08, 0x04, 0x02, 0x01, 0x00}

	wantAccel := MPU6050ThreeDData{
		X: 0x0001 / 16384.0 * d.gravity,
		Y: 0x0204 / 16384.0 * d.gravity,
		Z: 0x0816 / 16384.0 * d.gravity,
	}
	wantGyro := MPU6050ThreeDData{
		X: 0x1608 / 131.0,
		Y: 0x0402 / 131.0,
		Z: 0x0100 / 131.0,
	}
	wantTemp := float64(0x3264)/340 + 36.53

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, append(append(accData, tempData...), gyroData...))
		return len(b), nil
	}
	// act
	_ = d.GetData()
	// assert
	assert.Equal(t, wantAccel, d.Accelerometer)
	assert.Equal(t, wantGyro, d.Gyroscope)
	assert.InDelta(t, wantTemp, d.Temperature, 0.0)
}

func TestMPU6050GetDataReadError(t *testing.T) {
	d, adaptor := initTestMPU6050WithStubbedAdaptor()
	_ = d.Start()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}

	require.ErrorContains(t, d.GetData(), "read error")
}

func TestMPU6050GetDataWriteError(t *testing.T) {
	d, adaptor := initTestMPU6050WithStubbedAdaptor()
	_ = d.Start()

	adaptor.i2cWriteImpl = func(b []byte) (int, error) {
		return 0, errors.New("write error")
	}

	require.ErrorContains(t, d.GetData(), "write error")
}

func TestMPU6050_initialize(t *testing.T) {
	// sequence for initialization the device on Start()
	//          reset (according to data sheet)
	// * write power management 1 register address (0x6B)
	// * set device reset bit - write 0x80
	// * read device reset bit until it becomes false (timeout 100 ms)
	//   * write power management 1 register address (0x6B)
	//   * read value
	//   * wait some time before retry
	// * write signal path reset register address (0x68)
	// * reset all sensors - write 0x07
	// * wait 100 ms
	//               config
	// * write general config register address (0x1A)
	// * disable external sync and set filter bandwidth to 260 HZ - write 0x00
	// * write gyroscope config register address (0x1B)
	// * set full scale to 250 °/s - write 0x00
	// * write accelerometer config register address (0x1C)
	// * set full scale to 2 g - write 0x00
	// * write power management 1 register address (0x6B)
	// * set clock source to PLL with X and switch off sleep bit - write 0x01
	// arrange
	a := newI2cTestAdaptor()
	d := NewMPU6050Driver(a)
	readCallCounter := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		readCallCounter++
		// emulate ready
		b[0] = 0x00
		return len(b), nil
	}
	// act, assert - initialize() must be called on Start()
	err := d.Start()
	// assert
	require.NoError(t, err)
	assert.Equal(t, 1, readCallCounter)
	assert.Len(t, a.written, 13)
	assert.Equal(t, uint8(0x6B), a.written[0])
	assert.Equal(t, uint8(0x80), a.written[1])
	assert.Equal(t, uint8(0x6B), a.written[2])
	assert.Equal(t, uint8(0x68), a.written[3])
	assert.Equal(t, uint8(0x07), a.written[4])
	assert.Equal(t, uint8(0x1A), a.written[5])
	assert.Equal(t, uint8(0x00), a.written[6])
	assert.Equal(t, uint8(0x1B), a.written[7])
	assert.Equal(t, uint8(0x00), a.written[8])
	assert.Equal(t, uint8(0x1C), a.written[9])
	assert.Equal(t, uint8(0x00), a.written[10])
	assert.Equal(t, uint8(0x6B), a.written[11])
	assert.Equal(t, uint8(0x01), a.written[12])
}
