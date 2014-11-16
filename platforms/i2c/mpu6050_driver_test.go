package i2c

import (
	"testing"

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

func TestMPU6050Driver(t *testing.T) {
	// Does it implement gobot.DriverInterface?
	var _ gobot.DriverInterface = (*MPU6050Driver)(nil)

	// Does its adaptor implements the I2cInterface?
	driver := initTestMPU6050Driver()
	var _ I2cInterface = driver.adaptor()
}

func TestNewMPU6050Driver(t *testing.T) {
	// Does it return a pointer to an instance of MPU6050Driver?
	var bm interface{} = NewMPU6050Driver(newI2cTestAdaptor("adaptor"), "bot")
	_, ok := bm.(*MPU6050Driver)
	if !ok {
		t.Errorf("NewMPU6050Driver() should have returned a *MPU6050Driver")
	}
}

// Methods
func TestMPU6050DriverStart(t *testing.T) {
	mpu := initTestMPU6050Driver()

	gobot.Assert(t, mpu.Start(), nil)
}

func TestMPU6050DriverHalt(t *testing.T) {
	mpu := initTestMPU6050Driver()

	gobot.Assert(t, mpu.Halt(), nil)
}
