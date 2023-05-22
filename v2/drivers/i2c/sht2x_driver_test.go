package i2c

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/gobottest"
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
		t.Errorf("NewSHT2xDriver() should have returned a *SHT2xDriver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "SHT2x"), true)
	gobottest.Assert(t, d.defaultAddress, 0x40)
}

func TestSHT2xOptions(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	b := NewSHT2xDriver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, b.GetBusOrDefault(1), 2)
}

func TestSHT2xStart(t *testing.T) {
	d := NewSHT2xDriver(newI2cTestAdaptor())
	gobottest.Assert(t, d.Start(), nil)
}

func TestSHT2xHalt(t *testing.T) {
	d, _ := initTestSHT2xDriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestSHT2xReset(t *testing.T) {
	d, a := initTestSHT2xDriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, nil
	}
	d.Start()
	err := d.Reset()
	gobottest.Assert(t, err, nil)
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
	d.Start()
	temp, err := d.Temperature()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, temp, float32(18.809052))
	hum, err := d.Humidity()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, hum, float32(40.279907))
}

func TestSHT2xAccuracy(t *testing.T) {
	d, a := initTestSHT2xDriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		if a.written[len(a.written)-1] == SHT2xReadUserReg {
			buf.Write([]byte{0x3a})
		} else if a.written[len(a.written)-2] == SHT2xWriteUserReg {
			buf.Write([]byte{a.written[len(a.written)-1]})
		} else {
			return 0, nil
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	d.Start()
	d.SetAccuracy(SHT2xAccuracyLow)
	gobottest.Assert(t, d.Accuracy(), SHT2xAccuracyLow)
	err := d.sendAccuracy()
	gobottest.Assert(t, err, nil)
}

func TestSHT2xTemperatureCrcError(t *testing.T) {
	d, a := initTestSHT2xDriverWithStubbedAdaptor()
	d.Start()

	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		if a.written[len(a.written)-1] == SHT2xTriggerTempMeasureNohold {
			buf.Write([]byte{95, 168, 0})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	temp, err := d.Temperature()
	gobottest.Assert(t, err, errors.New("Invalid crc"))
	gobottest.Assert(t, temp, float32(0.0))
}

func TestSHT2xHumidityCrcError(t *testing.T) {
	d, a := initTestSHT2xDriverWithStubbedAdaptor()
	d.Start()

	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		if a.written[len(a.written)-1] == SHT2xTriggerHumdMeasureNohold {
			buf.Write([]byte{94, 202, 0})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	hum, err := d.Humidity()
	gobottest.Assert(t, err, errors.New("Invalid crc"))
	gobottest.Assert(t, hum, float32(0.0))
}

func TestSHT2xTemperatureLengthError(t *testing.T) {
	d, a := initTestSHT2xDriverWithStubbedAdaptor()
	d.Start()

	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		if a.written[len(a.written)-1] == SHT2xTriggerTempMeasureNohold {
			buf.Write([]byte{95, 168})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	temp, err := d.Temperature()
	gobottest.Assert(t, err, ErrNotEnoughBytes)
	gobottest.Assert(t, temp, float32(0.0))
}

func TestSHT2xHumidityLengthError(t *testing.T) {
	d, a := initTestSHT2xDriverWithStubbedAdaptor()
	d.Start()

	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		if a.written[len(a.written)-1] == SHT2xTriggerHumdMeasureNohold {
			buf.Write([]byte{94, 202})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	hum, err := d.Humidity()
	gobottest.Assert(t, err, ErrNotEnoughBytes)
	gobottest.Assert(t, hum, float32(0.0))
}
