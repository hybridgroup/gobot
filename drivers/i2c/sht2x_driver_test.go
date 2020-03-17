package i2c

import (
	"bytes"
	"errors"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*SHT2xDriver)(nil)

// --------- HELPERS
func initTestSHT2xDriver() (driver *SHT2xDriver) {
	driver, _ = initTestSHT2xDriverWithStubbedAdaptor()
	return
}

func initTestSHT2xDriverWithStubbedAdaptor() (*SHT2xDriver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewSHT2xDriver(adaptor), adaptor
}

// --------- TESTS

func TestNewSHT2xDriver(t *testing.T) {
	// Does it return a pointer to an instance of SHT2xDriver?
	var sht2x interface{} = NewSHT2xDriver(newI2cTestAdaptor())
	_, ok := sht2x.(*SHT2xDriver)
	if !ok {
		t.Errorf("NewSHT2xDriver() should have returned a *SHT2xDriver")
	}
}

func TestSHT2xDriver(t *testing.T) {
	sht2x := initTestSHT2xDriver()
	gobottest.Refute(t, sht2x.Connection(), nil)
}

func TestSHT2xDriverStart(t *testing.T) {
	sht2x, _ := initTestSHT2xDriverWithStubbedAdaptor()

	gobottest.Assert(t, sht2x.Start(), nil)
}

func TestSHT2xStartConnectError(t *testing.T) {
	d, adaptor := initTestSHT2xDriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestSHT2xDriverHalt(t *testing.T) {
	sht2x := initTestSHT2xDriver()

	gobottest.Assert(t, sht2x.Halt(), nil)
}

func TestSHT2xDriverReset(t *testing.T) {
	sht2x, adaptor := initTestSHT2xDriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return 0, nil
	}
	sht2x.Start()
	err := sht2x.Reset()
	gobottest.Assert(t, err, nil)
}

func TestSHT2xDriverMeasurements(t *testing.T) {
	sht2x, adaptor := initTestSHT2xDriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values produced by dumping data from actual sensor
		if adaptor.written[len(adaptor.written)-1] == SHT2xTriggerTempMeasureNohold {
			buf.Write([]byte{95, 168, 59})
		} else if adaptor.written[len(adaptor.written)-1] == SHT2xTriggerHumdMeasureNohold {
			buf.Write([]byte{94, 202, 22})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	sht2x.Start()
	temp, err := sht2x.Temperature()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, temp, float32(18.809052))
	hum, err := sht2x.Humidity()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, hum, float32(40.279907))
}

func TestSHT2xDriverAccuracy(t *testing.T) {
	sht2x, adaptor := initTestSHT2xDriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		if adaptor.written[len(adaptor.written)-1] == SHT2xReadUserReg {
			buf.Write([]byte{0x3a})
		} else if adaptor.written[len(adaptor.written)-2] == SHT2xWriteUserReg {
			buf.Write([]byte{adaptor.written[len(adaptor.written)-1]})
		} else {
			return 0, nil
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	sht2x.Start()
	sht2x.SetAccuracy(SHT2xAccuracyLow)
	gobottest.Assert(t, sht2x.Accuracy(), SHT2xAccuracyLow)
	err := sht2x.sendAccuracy()
	gobottest.Assert(t, err, nil)
}

func TestSHT2xDriverTemperatureCrcError(t *testing.T) {
	sht2x, adaptor := initTestSHT2xDriverWithStubbedAdaptor()
	sht2x.Start()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		if adaptor.written[len(adaptor.written)-1] == SHT2xTriggerTempMeasureNohold {
			buf.Write([]byte{95, 168, 0})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	temp, err := sht2x.Temperature()
	gobottest.Assert(t, err, errors.New("Invalid crc"))
	gobottest.Assert(t, temp, float32(0.0))
}

func TestSHT2xDriverHumidityCrcError(t *testing.T) {
	sht2x, adaptor := initTestSHT2xDriverWithStubbedAdaptor()
	sht2x.Start()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		if adaptor.written[len(adaptor.written)-1] == SHT2xTriggerHumdMeasureNohold {
			buf.Write([]byte{94, 202, 0})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	hum, err := sht2x.Humidity()
	gobottest.Assert(t, err, errors.New("Invalid crc"))
	gobottest.Assert(t, hum, float32(0.0))
}

func TestSHT2xDriverTemperatureLengthError(t *testing.T) {
	sht2x, adaptor := initTestSHT2xDriverWithStubbedAdaptor()
	sht2x.Start()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		if adaptor.written[len(adaptor.written)-1] == SHT2xTriggerTempMeasureNohold {
			buf.Write([]byte{95, 168})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	temp, err := sht2x.Temperature()
	gobottest.Assert(t, err, ErrNotEnoughBytes)
	gobottest.Assert(t, temp, float32(0.0))
}

func TestSHT2xDriverHumidityLengthError(t *testing.T) {
	sht2x, adaptor := initTestSHT2xDriverWithStubbedAdaptor()
	sht2x.Start()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		if adaptor.written[len(adaptor.written)-1] == SHT2xTriggerHumdMeasureNohold {
			buf.Write([]byte{94, 202})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	hum, err := sht2x.Humidity()
	gobottest.Assert(t, err, ErrNotEnoughBytes)
	gobottest.Assert(t, hum, float32(0.0))
}

func TestSHT2xDriverSetName(t *testing.T) {
	b := initTestSHT2xDriver()
	b.SetName("TESTME")
	gobottest.Assert(t, b.Name(), "TESTME")
}

func TestSHT2xDriverOptions(t *testing.T) {
	b := NewSHT2xDriver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, b.GetBusOrDefault(1), 2)
}
