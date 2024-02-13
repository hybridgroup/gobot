package i2c

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
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
		require.Fail(t, "NewSHT3xDriver() should have returned a *SHT3xDriver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "SHT3x"))
	assert.Equal(t, 0x44, d.defaultAddress)
}

func TestSHT3xOptions(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewSHT3xDriver(newI2cTestAdaptor(), WithBus(2))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
}

func TestSHT3xStart(t *testing.T) {
	d := NewSHT3xDriver(newI2cTestAdaptor())
	require.NoError(t, d.Start())
}

func TestSHT3xHalt(t *testing.T) {
	d, _ := initTestSHT3xDriverWithStubbedAdaptor()
	require.NoError(t, d.Halt())
}

func TestSHT3xSampleNormal(t *testing.T) {
	d, a := initTestSHT3xDriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0xbe, 0xef, 0x92, 0xbe, 0xef, 0x92})
		return 6, nil
	}

	temp, rh, _ := d.Sample()
	assert.InDelta(t, float32(85.523003), temp, 0.0)
	assert.InDelta(t, float32(74.5845), rh, 0.0)

	// check the temp with the units in F
	d.Units = "F"
	temp, _, _ = d.Sample()
	assert.InDelta(t, float32(185.9414), temp, 0.0)
}

func TestSHT3xSampleBadCrc(t *testing.T) {
	d, a := initTestSHT3xDriverWithStubbedAdaptor()
	// Check that the 1st crc failure is caught
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0xbe, 0xef, 0x00, 0xbe, 0xef, 0x92})
		return 6, nil
	}

	_, _, err := d.Sample()
	assert.Equal(t, ErrInvalidCrc, err)

	// Check that the 2nd crc failure is caught
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0xbe, 0xef, 0x92, 0xbe, 0xef, 0x00})
		return 6, nil
	}

	_, _, err = d.Sample()
	assert.Equal(t, ErrInvalidCrc, err)
}

func TestSHT3xSampleBadRead(t *testing.T) {
	d, a := initTestSHT3xDriverWithStubbedAdaptor()
	// Check that the 1st crc failure is caught
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0xbe, 0xef, 0x00, 0xbe, 0xef})
		return 5, nil
	}

	_, _, err := d.Sample()
	assert.Equal(t, ErrNotEnoughBytes, err)
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
	assert.Equal(t, ErrInvalidTemp, err)
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
	assert.Equal(t, ErrNotEnoughBytes, err)

	// Don't read any bytes and return an error
	a.i2cReadImpl = func([]byte) (int, error) {
		return 0, invalidRead
	}

	_, err = d.sendCommandDelayGetResponse(nil, nil, 1)
	assert.Equal(t, invalidRead, err)

	// Don't write any bytes and return an error
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 42, invalidWrite
	}

	_, err = d.sendCommandDelayGetResponse(nil, nil, 1)
	assert.Equal(t, invalidWrite, err)
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
	require.NoError(t, err)
	assert.True(t, status)

	// heater disabled
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x00, 0x00, 0x81})
		return 3, nil
	}

	status, err = d.Heater()
	require.NoError(t, err)
	assert.False(t, status)

	// heater crc failed
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x00, 0x00, 0x00})
		return 3, nil
	}

	_, err = d.Heater()
	assert.Equal(t, ErrInvalidCrc, err)

	// heater read failed
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x00, 0x00})
		return 2, nil
	}

	_, err = d.Heater()
	require.Error(t, err)
}

func TestSHT3xSetHeater(t *testing.T) {
	d, _ := initTestSHT3xDriverWithStubbedAdaptor()
	_ = d.SetHeater(false)
	_ = d.SetHeater(true)
}

func TestSHT3xSetAccuracy(t *testing.T) {
	d, _ := initTestSHT3xDriverWithStubbedAdaptor()

	assert.Equal(t, byte(SHT3xAccuracyHigh), d.Accuracy())

	err := d.SetAccuracy(SHT3xAccuracyMedium)
	require.NoError(t, err)
	assert.Equal(t, byte(SHT3xAccuracyMedium), d.Accuracy())

	err = d.SetAccuracy(SHT3xAccuracyLow)
	require.NoError(t, err)
	assert.Equal(t, byte(SHT3xAccuracyLow), d.Accuracy())

	err = d.SetAccuracy(0xff)
	assert.Equal(t, ErrInvalidAccuracy, err)
}

func TestSHT3xSerialNumber(t *testing.T) {
	d, a := initTestSHT3xDriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x20, 0x00, 0x5d, 0xbe, 0xef, 0x92})
		return 6, nil
	}

	sn, err := d.SerialNumber()

	require.NoError(t, err)
	assert.Equal(t, uint32(0x2000beef), sn)
}
