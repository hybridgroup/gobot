package i2c

import (
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*ADS1015Driver)(nil)

// --------- HELPERS
func initTestADS1015Driver() (driver *ADS1015Driver) {
	driver, _ = initTestADS1015DriverWithStubbedAdaptor()
	return
}

func initTestADS1015DriverWithStubbedAdaptor() (*ADS1015Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewADS1015Driver(adaptor), adaptor
}

// --------- BASE TESTS
func TestNewADS1015Driver(t *testing.T) {
	// Does it return a pointer to an instance of ADS1015Driver?
	var bm interface{} = NewADS1015Driver(newI2cTestAdaptor())
	_, ok := bm.(*ADS1015Driver)
	if !ok {
		t.Errorf("NewADS1015Driver() should have returned a *ADS1015Driver")
	}
}

func TestADS1015DriverSetName(t *testing.T) {
	d := initTestADS1015Driver()
	d.SetName("TESTME")
	gobottest.Assert(t, d.Name(), "TESTME")
}

func TestADS1015DriverOptions(t *testing.T) {
	d := NewADS1015Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}

// --------- DRIVER SPECIFIC TESTS
