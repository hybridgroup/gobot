package i2c

import (
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

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

func TestMPU6050Driver(t *testing.T) {
	mpu := initTestMPU6050Driver()
	gobottest.Refute(t, mpu.Connection(), nil)
	gobottest.Assert(t, mpu.interval, 10*time.Millisecond)

	mpu = NewMPU6050Driver(newI2cTestAdaptor(), 100*time.Millisecond)
	gobottest.Assert(t, mpu.interval, 100*time.Millisecond)
}

// Methods
func TestMPU6050DriverStart(t *testing.T) {
	mpu := initTestMPU6050Driver()

	gobottest.Assert(t, mpu.Start(), nil)
}

func TestMPU6050DriverHalt(t *testing.T) {
	mpu := initTestMPU6050Driver()

	gobottest.Assert(t, mpu.Halt(), nil)
}
