package i2c

import (
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
)

// --------- HELPERS
func initTestMPU6050Driver() (driver *MPU6050Driver) {
	driver, _ = initTestMPU6050DriverWithStubbedAdaptor()
	return
}

func initTestMPU6050DriverWithStubbedAdaptor() (*MPU6050Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor("adaptor")
	return NewMPU6050Driver(adaptor, "bot"), adaptor
}

// --------- TESTS

func TestNewMPU6050Driver(t *testing.T) {
	// Does it return a pointer to an instance of MPU6050Driver?
	var bm interface{} = NewMPU6050Driver(newI2cTestAdaptor("adaptor"), "bot")
	_, ok := bm.(*MPU6050Driver)
	if !ok {
		t.Errorf("NewMPU6050Driver() should have returned a *MPU6050Driver")
	}
}

func TestMPU6050Driver(t *testing.T) {
	mpu := initTestMPU6050Driver()
	gobot.Assert(t, mpu.Name(), "bot")
	gobot.Assert(t, mpu.Connection().Name(), "adaptor")
	gobot.Assert(t, mpu.interval, 10*time.Millisecond)

	mpu = NewMPU6050Driver(newI2cTestAdaptor("adaptor"), "bot", 100*time.Millisecond)
	gobot.Assert(t, mpu.interval, 100*time.Millisecond)
}

// Methods
func TestMPU6050DriverStart(t *testing.T) {
	mpu := initTestMPU6050Driver()

	gobot.Assert(t, len(mpu.Start()), 0)
}

func TestMPU6050DriverHalt(t *testing.T) {
	mpu := initTestMPU6050Driver()

	gobot.Assert(t, len(mpu.Halt()), 0)
}
