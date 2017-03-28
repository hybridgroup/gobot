package i2c

import (
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*BMP280Driver)(nil)

// --------- HELPERS
func initTestBMP280Driver() (driver *BMP280Driver) {
	driver, _ = initTestBMP280DriverWithStubbedAdaptor()
	return
}

func initTestBMP280DriverWithStubbedAdaptor() (*BMP280Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewBMP280Driver(adaptor), adaptor
}

// --------- TESTS

func TestNewBMP280Driver(t *testing.T) {
	// Does it return a pointer to an instance of BME280Driver?
	var bmp280 interface{} = NewBMP280Driver(newI2cTestAdaptor())
	_, ok := bmp280.(*BMP280Driver)
	if !ok {
		t.Errorf("NewBMP280Driver() should have returned a *BMP280Driver")
	}
}

func TestBMP280Driver(t *testing.T) {
	bmp280 := initTestBMP280Driver()
	gobottest.Refute(t, bmp280.Connection(), nil)
}

func TestBMP280DriverStart(t *testing.T) {
	bmp280, _ := initTestBMP280DriverWithStubbedAdaptor()
	gobottest.Assert(t, bmp280.Start(), nil)
}

func TestBMP280DriverHalt(t *testing.T) {
	bmp280 := initTestBMP280Driver()

	gobottest.Assert(t, bmp280.Halt(), nil)
}

// TODO: implement test
func TestBMP280DriverMeasurements(t *testing.T) {

}

func TestBMP280DriverSetName(t *testing.T) {
	b := initTestBMP280Driver()
	b.SetName("TESTME")
	gobottest.Assert(t, b.Name(), "TESTME")
}

func TestBMP280DriverOptions(t *testing.T) {
	b := NewBMP280Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, b.GetBusOrDefault(1), 2)
}
