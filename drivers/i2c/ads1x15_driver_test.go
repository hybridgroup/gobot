package i2c

import (
	"errors"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/gobottest"
)

// the ADS1x15Driver is a Driver
var _ gobot.Driver = (*ADS1x15Driver)(nil)

// that supports the AnalogReader interface
var _ aio.AnalogReader = (*ADS1x15Driver)(nil)

// --------- HELPERS
func initTestADS1x15Driver() (driver *ADS1x15Driver) {
	driver, _ = initTestADS1x15DriverWithStubbedAdaptor()
	return
}

func initTestADS1x15DriverWithStubbedAdaptor() (*ADS1x15Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewADS1015Driver(adaptor), adaptor
}

// --------- BASE TESTS
func TestNewADS1x15Driver(t *testing.T) {
	// Does it return a pointer to an instance of ADS1x15Driver?
	var bm interface{} = NewADS1015Driver(newI2cTestAdaptor())
	_, ok := bm.(*ADS1x15Driver)
	if !ok {
		t.Errorf("NewADS1x15Driver() should have returned a *ADS1x15Driver")
	}
}

func TestADS1x15DriverSetName(t *testing.T) {
	d := initTestADS1x15Driver()
	d.SetName("TESTME")
	gobottest.Assert(t, d.Name(), "TESTME")
}

func TestADS1x15DriverOptions(t *testing.T) {
	d := NewADS1015Driver(newI2cTestAdaptor(), WithBus(2), WithADS1x15Gain(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
	gobottest.Assert(t, d.DefaultGain, 2)
}

func TestADS1x15StartAndHalt(t *testing.T) {
	d, _ := initTestADS1x15DriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}

func TestADS1x15StartConnectError(t *testing.T) {
	d, adaptor := initTestADS1x15DriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

// --------- DRIVER SPECIFIC TESTS

func TestADS1x15DriverAnalogRead(t *testing.T) {
	d, adaptor := initTestADS1x15DriverWithStubbedAdaptor()
	d.Start()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99, 1})
		return 2, nil
	}

	val, err := d.AnalogRead("0")
	gobottest.Assert(t, val, 25345)
	gobottest.Assert(t, err, nil)

	val, err = d.AnalogRead("1")
	gobottest.Assert(t, val, 25345)
	gobottest.Assert(t, err, nil)

	val, err = d.AnalogRead("2")
	gobottest.Assert(t, val, 25345)
	gobottest.Assert(t, err, nil)

	val, err = d.AnalogRead("3")
	gobottest.Assert(t, val, 25345)
	gobottest.Assert(t, err, nil)

	val, err = d.AnalogRead("0-1")
	gobottest.Assert(t, val, 25345)
	gobottest.Assert(t, err, nil)

	val, err = d.AnalogRead("2-3")
	gobottest.Assert(t, val, 25345)
	gobottest.Assert(t, err, nil)
}

func TestADS1x15DriverAnalogReadError(t *testing.T) {
	d, a := initTestADS1x15DriverWithStubbedAdaptor()
	d.Start()

	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}

	_, err := d.AnalogRead("0")
	gobottest.Assert(t, err, errors.New("read error"))
}

func TestADS1x15DriverAnalogReadInvalidPin(t *testing.T) {
	d, _ := initTestADS1x15DriverWithStubbedAdaptor()

	_, err := d.AnalogRead("99")
	gobottest.Assert(t, err, errors.New("Invalid channel, must be between 0 and 3"))
}

func TestADS1x15DriverAnalogReadWriteError(t *testing.T) {
	d, a := initTestADS1x15DriverWithStubbedAdaptor()
	d.Start()

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	_, err := d.AnalogRead("0")
	gobottest.Assert(t, err, errors.New("write error"))

	_, err = d.AnalogRead("0-1")
	gobottest.Assert(t, err, errors.New("write error"))

	_, err = d.AnalogRead("2-3")
	gobottest.Assert(t, err, errors.New("write error"))
}
