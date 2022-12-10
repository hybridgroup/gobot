package i2c

import (
	"errors"
	"strings"
	"testing"

	"bytes"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*BH1750Driver)(nil)

func initTestBH1750DriverWithStubbedAdaptor() (*BH1750Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewBH1750Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewBH1750Driver(t *testing.T) {
	var di interface{} = NewBH1750Driver(newI2cTestAdaptor())
	d, ok := di.(*BH1750Driver)
	if !ok {
		t.Errorf("NewBH1750Driver() should have returned a *BH1750Driver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "BH1750"), true)
	gobottest.Assert(t, d.defaultAddress, 0x23)
}

func TestBH1750Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewBH1750Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}

func TestBH1750Start(t *testing.T) {
	d := NewBH1750Driver(newI2cTestAdaptor())
	gobottest.Assert(t, d.Start(), nil)
}

func TestBH1750Halt(t *testing.T) {
	d, _ := initTestBH1750DriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestBH1750NullLux(t *testing.T) {
	d, _ := initTestBH1750DriverWithStubbedAdaptor()
	lux, _ := d.Lux()
	gobottest.Assert(t, lux, 0)
}

func TestBH1750Lux(t *testing.T) {
	d, a := initTestBH1750DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{0x05, 0xb0})
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}

	lux, _ := d.Lux()
	gobottest.Assert(t, lux, 1213)
}

func TestBH1750NullRawSensorData(t *testing.T) {
	d, _ := initTestBH1750DriverWithStubbedAdaptor()
	level, _ := d.RawSensorData()
	gobottest.Assert(t, level, 0)
}

func TestBH1750RawSensorData(t *testing.T) {
	d, a := initTestBH1750DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{0x05, 0xb0})
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}

	level, _ := d.RawSensorData()
	gobottest.Assert(t, level, 1456)
}

func TestBH1750LuxError(t *testing.T) {
	d, a := initTestBH1750DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("wrong number of bytes read")
	}

	_, err := d.Lux()
	gobottest.Assert(t, err, errors.New("wrong number of bytes read"))
}

func TestBH1750RawSensorDataError(t *testing.T) {
	d, a := initTestBH1750DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("wrong number of bytes read")
	}

	_, err := d.RawSensorData()
	gobottest.Assert(t, err, errors.New("wrong number of bytes read"))
}
