package i2c

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*MMA7660Driver)(nil)

// --------- HELPERS
func initTestMMA7660Driver() (driver *MMA7660Driver) {
	driver, _ = initTestMMA7660DriverWithStubbedAdaptor()
	return
}

func initTestMMA7660DriverWithStubbedAdaptor() (*MMA7660Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewMMA7660Driver(adaptor), adaptor
}

// --------- TESTS

func TestNewMMA7660Driver(t *testing.T) {
	// Does it return a pointer to an instance of MMA7660Driver?
	var mma interface{} = NewMMA7660Driver(newI2cTestAdaptor())
	_, ok := mma.(*MMA7660Driver)
	if !ok {
		t.Errorf("NewMMA7660Driver() should have returned a *MMA7660Driver")
	}
}

// Methods
func TestMMA7660Driver(t *testing.T) {
	mma := initTestMMA7660Driver()

	gobottest.Refute(t, mma.Connection(), nil)
	gobottest.Assert(t, strings.HasPrefix(mma.Name(), "MMA7660"), true)
}
