package i2c

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*SHT3xDriver)(nil)

// --------- HELPERS
func initTestSHT3xDriver() (driver *SHT3xDriver) {
	driver, _ = initTestSHT3xDriverWithStubbedAdaptor()
	return
}

func initTestSHT3xDriverWithStubbedAdaptor() (*SHT3xDriver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewSHT3xDriver(adaptor), adaptor
}

// --------- TESTS

func TestNewSHT3xDriver(t *testing.T) {
	// Does it return a pointer to an instance of SHT3xDriver?
	var bm interface{} = NewSHT3xDriver(newI2cTestAdaptor())
	_, ok := bm.(*SHT3xDriver)
	if !ok {
		t.Errorf("NewSHT3xDriver() should have returned a *SHT3xDriver")
	}

	b := NewSHT3xDriver(newI2cTestAdaptor())
	gobottest.Refute(t, b.Connection(), nil)
}

// Methods

func TestSHT3xDriverStart(t *testing.T) {
	sht3x, _ := initTestSHT3xDriverWithStubbedAdaptor()

	gobottest.Assert(t, sht3x.Start(), nil)
}

func TestSHT3xStartConnectError(t *testing.T) {
	d, adaptor := initTestSHT3xDriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestSHT3xDriverHalt(t *testing.T) {
	sht3x := initTestSHT3xDriver()

	gobottest.Assert(t, sht3x.Halt(), nil)
}

// Test Name & SetName
func TestSHT3xDriverName(t *testing.T) {
	sht3x := initTestSHT3xDriver()

	gobottest.Assert(t, strings.HasPrefix(sht3x.Name(), "SHT3x"), true)
	sht3x.SetName("Sensor")
	gobottest.Assert(t, sht3x.Name(), "Sensor")
}

func TestSHT3xDriverOptions(t *testing.T) {
	d := NewSHT3xDriver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}

// Test Accuracy & SetAccuracy
func TestSHT3xDriverSetAccuracy(t *testing.T) {
	sht3x := initTestSHT3xDriver()

	gobottest.Assert(t, sht3x.Accuracy(), byte(SHT3xAccuracyHigh))

	err := sht3x.SetAccuracy(SHT3xAccuracyMedium)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, sht3x.Accuracy(), byte(SHT3xAccuracyMedium))

	err = sht3x.SetAccuracy(SHT3xAccuracyLow)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, sht3x.Accuracy(), byte(SHT3xAccuracyLow))

	err = sht3x.SetAccuracy(0xff)
	gobottest.Assert(t, err, ErrInvalidAccuracy)
}

// Test Sample
func TestSHT3xDriverSampleNormal(t *testing.T) {
	sht3x, adaptor := initTestSHT3xDriverWithStubbedAdaptor()

	gobottest.Assert(t, sht3x.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0xbe, 0xef, 0x92, 0xbe, 0xef, 0x92})
		return 6, nil
	}

	temp, rh, _ := sht3x.Sample()
	gobottest.Assert(t, temp, float32(85.523003))
	gobottest.Assert(t, rh, float32(74.5845))

	// check the temp with the units in F
	sht3x.Units = "F"
	temp, _, _ = sht3x.Sample()
	gobottest.Assert(t, temp, float32(185.9414))
}

func TestSHT3xDriverSampleBadCrc(t *testing.T) {
	sht3x, adaptor := initTestSHT3xDriverWithStubbedAdaptor()

	gobottest.Assert(t, sht3x.Start(), nil)

	// Check that the 1st crc failure is caught
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0xbe, 0xef, 0x00, 0xbe, 0xef, 0x92})
		return 6, nil
	}

	_, _, err := sht3x.Sample()
	gobottest.Assert(t, err, ErrInvalidCrc)

	// Check that the 2nd crc failure is caught
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0xbe, 0xef, 0x92, 0xbe, 0xef, 0x00})
		return 6, nil
	}

	_, _, err = sht3x.Sample()
	gobottest.Assert(t, err, ErrInvalidCrc)
}

func TestSHT3xDriverSampleBadRead(t *testing.T) {
	sht3x, adaptor := initTestSHT3xDriverWithStubbedAdaptor()

	gobottest.Assert(t, sht3x.Start(), nil)

	// Check that the 1st crc failure is caught
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0xbe, 0xef, 0x00, 0xbe, 0xef})
		return 5, nil
	}

	_, _, err := sht3x.Sample()
	gobottest.Assert(t, err, ErrNotEnoughBytes)
}

func TestSHT3xDriverSampleUnits(t *testing.T) {
	sht3x, adaptor := initTestSHT3xDriverWithStubbedAdaptor()

	gobottest.Assert(t, sht3x.Start(), nil)

	// Check that the 1st crc failure is caught
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0xbe, 0xef, 0x92, 0xbe, 0xef, 0x92})
		return 6, nil
	}

	sht3x.Units = "K"
	_, _, err := sht3x.Sample()
	gobottest.Assert(t, err, ErrInvalidTemp)
}

// Test internal sendCommandDelayGetResponse
func TestSHT3xDriverSCDGRIoFailures(t *testing.T) {
	sht3x, adaptor := initTestSHT3xDriverWithStubbedAdaptor()

	gobottest.Assert(t, sht3x.Start(), nil)

	invalidRead := errors.New("Read error")
	invalidWrite := errors.New("Write error")

	// Only send 5 bytes
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0xbe, 0xef, 0x92, 0xbe, 0xef})
		return 5, nil
	}

	_, err := sht3x.sendCommandDelayGetResponse(nil, nil, 2)
	gobottest.Assert(t, err, ErrNotEnoughBytes)

	// Don't read any bytes and return an error
	adaptor.i2cReadImpl = func([]byte) (int, error) {
		return 0, invalidRead
	}

	_, err = sht3x.sendCommandDelayGetResponse(nil, nil, 1)
	gobottest.Assert(t, err, invalidRead)

	// Don't write any bytes and return an error
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 42, invalidWrite
	}

	_, err = sht3x.sendCommandDelayGetResponse(nil, nil, 1)
	gobottest.Assert(t, err, invalidWrite)
}

// Test Heater and getStatusRegister
func TestSHT3xDriverHeater(t *testing.T) {
	sht3x, adaptor := initTestSHT3xDriverWithStubbedAdaptor()

	gobottest.Assert(t, sht3x.Start(), nil)

	// heater enabled
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x20, 0x00, 0x5d})
		return 3, nil
	}

	status, err := sht3x.Heater()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, status, true)

	// heater disabled
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x00, 0x00, 0x81})
		return 3, nil
	}

	status, err = sht3x.Heater()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, status, false)

	// heater crc failed
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x00, 0x00, 0x00})
		return 3, nil
	}

	status, err = sht3x.Heater()
	gobottest.Assert(t, err, ErrInvalidCrc)

	// heater read failed
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x00, 0x00})
		return 2, nil
	}

	status, err = sht3x.Heater()
	gobottest.Refute(t, err, nil)
}

// Test SetHeater
func TestSHT3xDriverSetHeater(t *testing.T) {
	sht3x, _ := initTestSHT3xDriverWithStubbedAdaptor()

	gobottest.Assert(t, sht3x.Start(), nil)

	sht3x.SetHeater(false)
	sht3x.SetHeater(true)
}

// Test SerialNumber
func TestSHT3xDriverSerialNumber(t *testing.T) {
	sht3x, adaptor := initTestSHT3xDriverWithStubbedAdaptor()

	gobottest.Assert(t, sht3x.Start(), nil)

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x20, 0x00, 0x5d, 0xbe, 0xef, 0x92})
		return 6, nil
	}

	sn, err := sht3x.SerialNumber()

	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, sn, uint32(0x2000beef))
}
