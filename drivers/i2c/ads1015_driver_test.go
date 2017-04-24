package i2c

import (
	"errors"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/gobottest"
)

// the ADS1015Driver is a Driver
var _ gobot.Driver = (*ADS1015Driver)(nil)

// that supports the AnalogReader interface
var _ aio.AnalogReader = (*ADS1015Driver)(nil)

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
	d := NewADS1015Driver(newI2cTestAdaptor(), WithBus(2), WithADS1015Gain(ADS1015RegConfigPga2048V))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
	gobottest.Assert(t, d.gain, uint16(ADS1015RegConfigPga2048V))
}

func TestADS1015StartConnectError(t *testing.T) {
	d, adaptor := initTestADS1015DriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

// --------- DRIVER SPECIFIC TESTS

func TestADS1015DriverAnalogRead(t *testing.T) {
	d, adaptor := initTestADS1015DriverWithStubbedAdaptor()
	d.Start()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99, 1})
		return 2, nil
	}

	val, err := d.AnalogRead("0")
	gobottest.Assert(t, val, 1584)
	gobottest.Assert(t, err, nil)

	val, err = d.AnalogRead("0-1")
	gobottest.Assert(t, val, 1584)
	gobottest.Assert(t, err, nil)

	val, err = d.AnalogRead("2-3")
	gobottest.Assert(t, val, 1584)
	gobottest.Assert(t, err, nil)
}

func TestADS1015DriverAnalogReadError(t *testing.T) {
	d, a := initTestADS1015DriverWithStubbedAdaptor()
	d.Start()

	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}

	_, err := d.AnalogRead("0")
	gobottest.Assert(t, err, errors.New("read error"))
}

func TestADS1015DriverAnalogReadInvalidPin(t *testing.T) {
	d, _ := initTestADS1015DriverWithStubbedAdaptor()

	_, err := d.AnalogRead("99")
	gobottest.Assert(t, err, errors.New("Invalid channel."))
}
