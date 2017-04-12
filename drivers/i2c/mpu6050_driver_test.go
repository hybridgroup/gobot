package i2c

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

// ensure that MPU6050Driver fulfills Gobot Driver interface
var _ gobot.Driver = (*MPU6050Driver)(nil)

// --------- HELPERS
func initTestMPU6050Driver() (driver *MPU6050Driver) {
	driver, _ = initTestMPU6050DriverWithStubbedAdaptor()
	return
}

func initTestMPU6050DriverWithStubbedAdaptor() (*MPU6050Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewMPU6050Driver(adaptor), adaptor
}

// --------- TESTS

func TestNewMPU6050Driver(t *testing.T) {
	// Does it return a pointer to an instance of MPU6050Driver?
	var bm interface{} = NewMPU6050Driver(newI2cTestAdaptor())
	_, ok := bm.(*MPU6050Driver)
	if !ok {
		t.Errorf("NewMPU6050Driver() should have returned a *MPU6050Driver")
	}
}

func TestMPU6050DriverName(t *testing.T) {
	mpu := initTestMPU6050Driver()
	gobottest.Refute(t, mpu.Connection(), nil)
	gobottest.Assert(t, strings.HasPrefix(mpu.Name(), "MPU6050"), true)
}

func TestMPU6050DriverOptions(t *testing.T) {
	mpu := NewMPU6050Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, mpu.GetBusOrDefault(1), 2)
}

// Methods
func TestMPU6050DriverStart(t *testing.T) {
	mpu := initTestMPU6050Driver()

	gobottest.Assert(t, mpu.Start(), nil)
}

func TestMPU6050StartConnectError(t *testing.T) {
	d, adaptor := initTestMPU6050DriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestMPU6050DriverStartWriteError(t *testing.T) {
	mpu, adaptor := initTestMPU6050DriverWithStubbedAdaptor()
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, mpu.Start(), errors.New("write error"))
}

func TestMPU6050DriverHalt(t *testing.T) {
	mpu := initTestMPU6050Driver()

	gobottest.Assert(t, mpu.Halt(), nil)
}

func TestMPU6050DriverReadData(t *testing.T) {
	mpu, adaptor := initTestMPU6050DriverWithStubbedAdaptor()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x00, 0x01, 0x02, 0x04})
		return 4, nil
	}
	mpu.Start()
	mpu.GetData()
	gobottest.Assert(t, mpu.Temperature, int16(36))
}

func TestMPU6050DriverGetDataReadError(t *testing.T) {
	mpu, adaptor := initTestMPU6050DriverWithStubbedAdaptor()
	mpu.Start()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}

	gobottest.Assert(t, mpu.GetData(), errors.New("read error"))
}

func TestMPU6050DriverGetDataWriteError(t *testing.T) {
	mpu, adaptor := initTestMPU6050DriverWithStubbedAdaptor()
	mpu.Start()

	adaptor.i2cWriteImpl = func(b []byte) (int, error) {
		return 0, errors.New("write error")
	}

	gobottest.Assert(t, mpu.GetData(), errors.New("write error"))
}

func TestMPU6050DriverSetName(t *testing.T) {
	mpu := initTestMPU6050Driver()
	mpu.SetName("TESTME")
	gobottest.Assert(t, mpu.Name(), "TESTME")
}
