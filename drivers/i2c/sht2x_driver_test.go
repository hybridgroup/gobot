package i2c

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*SHT2xDriver)(nil)

func initTestSHT2xDriverWithStubbedAdaptor() (*SHT2xDriver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewSHT2xDriver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewSHT2xDriver(t *testing.T) {
	var di interface{} = NewSHT2xDriver(newI2cTestAdaptor())
	d, ok := di.(*SHT2xDriver)
	if !ok {
		require.Fail(t, "NewSHT2xDriver() should have returned a *SHT2xDriver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "SHT2x"))
	assert.Equal(t, 0x40, d.defaultAddress)
}

func TestSHT2xOptions(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	b := NewSHT2xDriver(newI2cTestAdaptor(), WithBus(2))
	assert.Equal(t, 2, b.GetBusOrDefault(1))
}

func TestSHT2xStart(t *testing.T) {
	d := NewSHT2xDriver(newI2cTestAdaptor())
	require.NoError(t, d.Start())
}

func TestSHT2xHalt(t *testing.T) {
	d, _ := initTestSHT2xDriverWithStubbedAdaptor()
	require.NoError(t, d.Halt())
}

func TestSHT2xReset(t *testing.T) {
	d, a := initTestSHT2xDriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, nil
	}
	_ = d.Start()
	err := d.Reset()
	require.NoError(t, err)
}

func TestSHT2xMeasurements(t *testing.T) {
	d, a := initTestSHT2xDriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values produced by dumping data from actual sensor
		if a.written[len(a.written)-1] == SHT2xTriggerTempMeasureNohold {
			buf.Write([]byte{95, 168, 59})
		} else if a.written[len(a.written)-1] == SHT2xTriggerHumdMeasureNohold {
			buf.Write([]byte{94, 202, 22})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	_ = d.Start()
	temp, err := d.Temperature()
	require.NoError(t, err)
	assert.InDelta(t, float32(18.809052), temp, 1.0e-5)
	hum, err := d.Humidity()
	require.NoError(t, err)
	assert.InDelta(t, float32(40.279907), hum, 0.0)
}

func TestSHT2xAccuracy(t *testing.T) {
	d, a := initTestSHT2xDriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		switch {
		case a.written[len(a.written)-1] == SHT2xReadUserReg:
			buf.Write([]byte{0x3a})
		case a.written[len(a.written)-2] == SHT2xWriteUserReg:
			buf.Write([]byte{a.written[len(a.written)-1]})
		default:
			return 0, nil
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	_ = d.Start()
	_ = d.SetAccuracy(SHT2xAccuracyLow)
	assert.Equal(t, SHT2xAccuracyLow, d.Accuracy())
	err := d.sendAccuracy()
	require.NoError(t, err)
}

func TestSHT2xTemperatureCrcError(t *testing.T) {
	d, a := initTestSHT2xDriverWithStubbedAdaptor()
	_ = d.Start()

	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		if a.written[len(a.written)-1] == SHT2xTriggerTempMeasureNohold {
			buf.Write([]byte{95, 168, 0})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	temp, err := d.Temperature()
	require.ErrorContains(t, err, "Invalid crc")
	assert.InDelta(t, float32(0.0), temp, 0.0)
}

func TestSHT2xHumidityCrcError(t *testing.T) {
	d, a := initTestSHT2xDriverWithStubbedAdaptor()
	_ = d.Start()

	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		if a.written[len(a.written)-1] == SHT2xTriggerHumdMeasureNohold {
			buf.Write([]byte{94, 202, 0})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	hum, err := d.Humidity()
	require.ErrorContains(t, err, "Invalid crc")
	assert.InDelta(t, float32(0.0), hum, 0.0)
}

func TestSHT2xTemperatureLengthError(t *testing.T) {
	d, a := initTestSHT2xDriverWithStubbedAdaptor()
	_ = d.Start()

	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		if a.written[len(a.written)-1] == SHT2xTriggerTempMeasureNohold {
			buf.Write([]byte{95, 168})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	temp, err := d.Temperature()
	assert.Equal(t, ErrNotEnoughBytes, err)
	assert.InDelta(t, float32(0.0), temp, 0.0)
}

func TestSHT2xHumidityLengthError(t *testing.T) {
	d, a := initTestSHT2xDriverWithStubbedAdaptor()
	_ = d.Start()

	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		if a.written[len(a.written)-1] == SHT2xTriggerHumdMeasureNohold {
			buf.Write([]byte{94, 202})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	hum, err := d.Humidity()
	assert.Equal(t, ErrNotEnoughBytes, err)
	assert.InDelta(t, float32(0.0), hum, 0.0)
}
