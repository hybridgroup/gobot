package i2c

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*SHT3xDriver)(nil)

func initTestSHT3xDriverWithStubbedAdaptor() (*SHT3xDriver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewSHT3xDriver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewSHT3xDriver(t *testing.T) {
	var di interface{} = NewSHT3xDriver(newI2cTestAdaptor())
	d, ok := di.(*SHT3xDriver)
	if !ok {
		t.Errorf("NewSHT3xDriver() should have returned a *SHT3xDriver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "SHT3x"), true)
	gobottest.Assert(t, d.defaultAddress, 0x44)
}

func TestSHT3xOptions(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewSHT3xDriver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}

func TestSHT3xStart(t *testing.T) {
	d := NewSHT3xDriver(newI2cTestAdaptor())
	gobottest.Assert(t, d.Start(), nil)
}

func TestSHT3xHalt(t *testing.T) {
	d, _ := initTestSHT3xDriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestSHT3xSampleNormal(t *testing.T) {
	d, a := initTestSHT3xDriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0xbe, 0xef, 0x92, 0xbe, 0xef, 0x92})
		return 6, nil
	}

	temp, rh, _ := d.Sample()
	gobottest.Assert(t, temp, float32(85.523003))
	gobottest.Assert(t, rh, float32(74.5845))

	// check the temp with the units in F
	d.Units = "F"
	temp, _, _ = d.Sample()
	gobottest.Assert(t, temp, float32(185.9414))
}

func TestSHT3xSampleBadCrc(t *testing.T) {
	d, a := initTestSHT3xDriverWithStubbedAdaptor()
	// Check that the 1st crc failure is caught
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0xbe, 0xef, 0x00, 0xbe, 0xef, 0x92})
		return 6, nil
	}

	_, _, err := d.Sample()
	gobottest.Assert(t, err, ErrInvalidCrc)

	// Check that the 2nd crc failure is caught
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0xbe, 0xef, 0x92, 0xbe, 0xef, 0x00})
		return 6, nil
	}

	_, _, err = d.Sample()
	gobottest.Assert(t, err, ErrInvalidCrc)
}

func TestSHT3xSampleBadRead(t *testing.T) {
	d, a := initTestSHT3xDriverWithStubbedAdaptor()
	// Check that the 1st crc failure is caught
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0xbe, 0xef, 0x00, 0xbe, 0xef})
		return 5, nil
	}

	_, _, err := d.Sample()
	gobottest.Assert(t, err, ErrNotEnoughBytes)
}

func TestSHT3xSampleUnits(t *testing.T) {
	d, a := initTestSHT3xDriverWithStubbedAdaptor()
	// Check that the 1st crc failure is caught
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0xbe, 0xef, 0x92, 0xbe, 0xef, 0x92})
		return 6, nil
	}

	d.Units = "K"
	_, _, err := d.Sample()
	gobottest.Assert(t, err, ErrInvalidTemp)
}

// Test internal sendCommandDelayGetResponse
func TestSHT3xSCDGRIoFailures(t *testing.T) {
	d, a := initTestSHT3xDriverWithStubbedAdaptor()
	invalidRead := errors.New("Read error")
	invalidWrite := errors.New("Write error")

	// Only send 5 bytes
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0xbe, 0xef, 0x92, 0xbe, 0xef})
		return 5, nil
	}

	_, err := d.sendCommandDelayGetResponse(nil, nil, 2)
	gobottest.Assert(t, err, ErrNotEnoughBytes)

	// Don't read any bytes and return an error
	a.i2cReadImpl = func([]byte) (int, error) {
		return 0, invalidRead
	}

	_, err = d.sendCommandDelayGetResponse(nil, nil, 1)
	gobottest.Assert(t, err, invalidRead)

	// Don't write any bytes and return an error
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 42, invalidWrite
	}

	_, err = d.sendCommandDelayGetResponse(nil, nil, 1)
	gobottest.Assert(t, err, invalidWrite)
}

// Test Heater and getStatusRegister
func TestSHT3xHeater(t *testing.T) {
	d, a := initTestSHT3xDriverWithStubbedAdaptor()
	// heater enabled
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x20, 0x00, 0x5d})
		return 3, nil
	}

	status, err := d.Heater()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, status, true)

	// heater disabled
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x00, 0x00, 0x81})
		return 3, nil
	}

	status, err = d.Heater()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, status, false)

	// heater crc failed
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x00, 0x00, 0x00})
		return 3, nil
	}

	_, err = d.Heater()
	gobottest.Assert(t, err, ErrInvalidCrc)

	// heater read failed
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x00, 0x00})
		return 2, nil
	}

	_, err = d.Heater()
	gobottest.Refute(t, err, nil)
}

func TestSHT3xSetHeater(t *testing.T) {
	d, _ := initTestSHT3xDriverWithStubbedAdaptor()
	d.SetHeater(false)
	d.SetHeater(true)
}

func TestSHT3xSetAccuracy(t *testing.T) {
	d, _ := initTestSHT3xDriverWithStubbedAdaptor()

	gobottest.Assert(t, d.Accuracy(), byte(SHT3xAccuracyHigh))

	err := d.SetAccuracy(SHT3xAccuracyMedium)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, d.Accuracy(), byte(SHT3xAccuracyMedium))

	err = d.SetAccuracy(SHT3xAccuracyLow)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, d.Accuracy(), byte(SHT3xAccuracyLow))

	err = d.SetAccuracy(0xff)
	gobottest.Assert(t, err, ErrInvalidAccuracy)
}

func TestSHT3xSerialNumber(t *testing.T) {
	d, a := initTestSHT3xDriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x20, 0x00, 0x5d, 0xbe, 0xef, 0x92})
		return 6, nil
	}

	sn, err := d.SerialNumber()

	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, sn, uint32(0x2000beef))
}
