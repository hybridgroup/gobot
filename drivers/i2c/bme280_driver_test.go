package i2c

import (
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*BME280Driver)(nil)

// --------- HELPERS
func initTestBME280Driver() (driver *BME280Driver) {
	driver, _ = initTestBME280DriverWithStubbedAdaptor()
	return
}

func initTestBME280DriverWithStubbedAdaptor() (*BME280Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewBME280Driver(adaptor), adaptor
}

// --------- TESTS

func TestNewBME280Driver(t *testing.T) {
	// Does it return a pointer to an instance of BME280Driver?
	var bme280 interface{} = NewBME280Driver(newI2cTestAdaptor())
	_, ok := bme280.(*BME280Driver)
	if !ok {
		t.Errorf("NewBME280Driver() should have returned a *BME280Driver")
	}
}

func TestBME280Driver(t *testing.T) {
	bme280 := initTestBME280Driver()
	gobottest.Refute(t, bme280.Connection(), nil)
}

func TestBME280DriverStart(t *testing.T) {
	bme280, _ := initTestBME280DriverWithStubbedAdaptor()
	gobottest.Assert(t, bme280.Start(), nil)
}

func TestBME280DriverHalt(t *testing.T) {
	bme280 := initTestBME280Driver()

	gobottest.Assert(t, bme280.Halt(), nil)
}

// TODO: implement test
func TestBME280DriverMeasurements(t *testing.T) {

}

func TestBME280DriverSetName(t *testing.T) {
	b := initTestBME280Driver()
	b.SetName("TESTME")
	gobottest.Assert(t, b.Name(), "TESTME")
}

func TestBME280DriverOptions(t *testing.T) {
	b := NewBME280Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, b.GetBusOrDefault(1), 2)
}
