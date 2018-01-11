package i2c

import (
	"errors"
	"strings"
	"testing"

	"bytes"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*BH1750Driver)(nil)

// --------- HELPERS
func initTestBH1750Driver() (driver *BH1750Driver) {
	driver, _ = initTestBH1750DriverWithStubbedAdaptor()
	return
}

func initTestBH1750DriverWithStubbedAdaptor() (*BH1750Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewBH1750Driver(adaptor), adaptor
}

// --------- TESTS

func TestNewBH1750Driver(t *testing.T) {
	// Does it return a pointer to an instance of BH1750Driver?
	var mma interface{} = NewBH1750Driver(newI2cTestAdaptor())
	_, ok := mma.(*BH1750Driver)
	if !ok {
		t.Errorf("NewBH1750Driver() should have returned a *BH1750Driver")
	}
}

// Methods
func TestBH1750Driver(t *testing.T) {
	mma := initTestBH1750Driver()

	gobottest.Refute(t, mma.Connection(), nil)
	gobottest.Assert(t, strings.HasPrefix(mma.Name(), "BH1750"), true)
}

func TestBH1750DriverSetName(t *testing.T) {
	d := initTestBH1750Driver()
	d.SetName("TESTME")
	gobottest.Assert(t, d.Name(), "TESTME")
}

func TestBH1750DriverOptions(t *testing.T) {
	d := NewBH1750Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}

func TestBH1750DriverStart(t *testing.T) {
	d := initTestBH1750Driver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestBH1750StartConnectError(t *testing.T) {
	d, adaptor := initTestBH1750DriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestBH1750DriverStartWriteError(t *testing.T) {
	mma, adaptor := initTestBH1750DriverWithStubbedAdaptor()
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, mma.Start(), errors.New("write error"))
}

func TestBH1750DriverHalt(t *testing.T) {
	d := initTestBH1750Driver()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestBH1750DriverNullLux(t *testing.T) {
	d, _ := initTestBH1750DriverWithStubbedAdaptor()
	d.Start()
	lux, _ := d.Lux()
	gobottest.Assert(t, lux, 0)
}

func TestBH1750DriverLux(t *testing.T) {
	d, adaptor := initTestBH1750DriverWithStubbedAdaptor()
	d.Start()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{0x05, 0xb0})
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}

	lux, _ := d.Lux()
	gobottest.Assert(t, lux, 1213)
}

func TestBH1750DriverNullRawSensorData(t *testing.T) {
	d, _ := initTestBH1750DriverWithStubbedAdaptor()
	d.Start()
	level, _ := d.RawSensorData()
	gobottest.Assert(t, level, 0)
}

func TestBH1750DriverRawSensorData(t *testing.T) {
	d, adaptor := initTestBH1750DriverWithStubbedAdaptor()
	d.Start()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{0x05, 0xb0})
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}

	level, _ := d.RawSensorData()
	gobottest.Assert(t, level, 1456)
}

func TestBH1750DriverLuxError(t *testing.T) {
	d, adaptor := initTestBH1750DriverWithStubbedAdaptor()
	d.Start()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("wrong number of bytes read")
	}

	_, err := d.Lux()
	gobottest.Assert(t, err, errors.New("wrong number of bytes read"))
}

func TestBH1750DriverRawSensorDataError(t *testing.T) {
	d, adaptor := initTestBH1750DriverWithStubbedAdaptor()
	d.Start()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("wrong number of bytes read")
	}

	_, err := d.RawSensorData()
	gobottest.Assert(t, err, errors.New("wrong number of bytes read"))
}

