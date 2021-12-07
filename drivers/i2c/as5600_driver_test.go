package i2c

import (
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*AS5600Driver)(nil)

func initTestAS5600Driver() *AS5600Driver {

	driver, _ := initTestAS5600DriverWithStubbedAdaptor()

	return driver
}

func initTestAS5600DriverWithStubbedAdaptor() (*AS5600Driver, *i2cTestAdaptor) {

	adaptor := newI2cTestAdaptor()

	return NewAS5600Driver(adaptor), adaptor
}

func TestAS5600Driver(t *testing.T) {

	as := initTestAS5600Driver()

	gobottest.Refute(t, as.Connection(), nil)
}

func TestAS5600DriverStart(t *testing.T) {
	var as *AS5600Driver

	as, _ = initTestAS5600DriverWithStubbedAdaptor()
	gobottest.Assert(t, as.Start(), nil)
}

func TestAS5600DriverHalt(t *testing.T) {
	as := initTestAS5600Driver()

	gobottest.Assert(t, as.Halt(), nil)
}

func TestAS5600DriverSetName(t *testing.T) {

	// Does it change the name of the driver
	as := initTestAS5600Driver()
	as.SetName("TESTME")
	gobottest.Assert(t, as.Name(), "TESTME")
}

func TestAS5600DriverOptions(t *testing.T) {

	as := NewAS5600Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, as.GetBusOrDefault(1), 2)
}
