package i2c

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/gobottest"
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
		t.Errorf("NewMPU6050Driver() should have returned a *MPU6050Driver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.name, "MPU6050"), true)
	gobottest.Assert(t, d.defaultAddress, 0x68)
	gobottest.Assert(t, d.dlpf, MPU6050DlpfConfig(0x00))
	gobottest.Assert(t, d.frameSync, MPU6050FrameSyncConfig(0x00))
	gobottest.Assert(t, d.accelFs, MPU6050AccelFsConfig(0x00))
	gobottest.Assert(t, d.gyroFs, MPU6050GyroFsConfig(0x00))
	gobottest.Assert(t, d.clock, MPU6050Pwr1ClockConfig(0x01))
	gobottest.Assert(t, d.gravity, 9.80665)
}

func TestMPU6050Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewMPU6050Driver(newI2cTestAdaptor(), WithBus(2), WithMPU6050DigitalFilter(0x06))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
	gobottest.Assert(t, d.dlpf, MPU6050DlpfConfig(0x06))
}

func TestWithMPU6050FrameSync(t *testing.T) {
	d := NewMPU6050Driver(newI2cTestAdaptor(), WithMPU6050FrameSync(0x07))
	gobottest.Assert(t, d.frameSync, MPU6050FrameSyncConfig(0x07))
}

func TestWithMPU6050AccelFullScaleRange(t *testing.T) {
	d := NewMPU6050Driver(newI2cTestAdaptor(), WithMPU6050AccelFullScaleRange(0x02))
	gobottest.Assert(t, d.accelFs, MPU6050AccelFsConfig(0x02))
}

func TestWithMPU6050GyroFullScaleRange(t *testing.T) {
	d := NewMPU6050Driver(newI2cTestAdaptor(), WithMPU6050GyroFullScaleRange(0x03))
	gobottest.Assert(t, d.gyroFs, MPU6050GyroFsConfig(0x03))
}

func TestWithMPU6050ClockSource(t *testing.T) {
	d := NewMPU6050Driver(newI2cTestAdaptor(), WithMPU6050ClockSource(0x05))
	gobottest.Assert(t, d.clock, MPU6050Pwr1ClockConfig(0x05))
}

func TestWithMPU6050Gravity(t *testing.T) {
	d := NewMPU6050Driver(newI2cTestAdaptor(), WithMPU6050Gravity(1.0))
	gobottest.Assert(t, d.gravity, 1.0)
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
	d.Start()

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
	d.GetData()
	// assert
	gobottest.Assert(t, d.Accelerometer, wantAccel)
	gobottest.Assert(t, d.Gyroscope, wantGyro)
	gobottest.Assert(t, d.Temperature, wantTemp)
}

func TestMPU6050GetDataReadError(t *testing.T) {
	d, adaptor := initTestMPU6050WithStubbedAdaptor()
	d.Start()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}

	gobottest.Assert(t, d.GetData(), errors.New("read error"))
}

func TestMPU6050GetDataWriteError(t *testing.T) {
	d, adaptor := initTestMPU6050WithStubbedAdaptor()
	d.Start()

	adaptor.i2cWriteImpl = func(b []byte) (int, error) {
		return 0, errors.New("write error")
	}

	gobottest.Assert(t, d.GetData(), errors.New("write error"))
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
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, readCallCounter, 1)
	gobottest.Assert(t, len(a.written), 13)
	gobottest.Assert(t, a.written[0], uint8(0x6B))
	gobottest.Assert(t, a.written[1], uint8(0x80))
	gobottest.Assert(t, a.written[2], uint8(0x6B))
	gobottest.Assert(t, a.written[3], uint8(0x68))
	gobottest.Assert(t, a.written[4], uint8(0x07))
	gobottest.Assert(t, a.written[5], uint8(0x1A))
	gobottest.Assert(t, a.written[6], uint8(0x00))
	gobottest.Assert(t, a.written[7], uint8(0x1B))
	gobottest.Assert(t, a.written[8], uint8(0x00))
	gobottest.Assert(t, a.written[9], uint8(0x1C))
	gobottest.Assert(t, a.written[10], uint8(0x00))
	gobottest.Assert(t, a.written[11], uint8(0x6B))
	gobottest.Assert(t, a.written[12], uint8(0x01))
}
