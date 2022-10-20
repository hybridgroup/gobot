package i2c

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*BME280Driver)(nil)

func initTestBME280WithStubbedAdaptor() (*BME280Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewBME280Driver(adaptor), adaptor
}

func TestNewBME280Driver(t *testing.T) {
	var di interface{} = NewBME280Driver(newI2cTestAdaptor())
	d, ok := di.(*BME280Driver)
	if !ok {
		t.Errorf("NewBME280Driver() should have returned a *BME280Driver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "BMP280"), true)
	gobottest.Assert(t, d.defaultAddress, 0x77)
	gobottest.Assert(t, d.ctrlPwrMode, uint8(0x03))
	gobottest.Assert(t, d.ctrlPressOversamp, BMP280PressureOversampling(0x05))
	gobottest.Assert(t, d.ctrlTempOversamp, BMP280TemperatureOversampling(0x01))
	gobottest.Assert(t, d.ctrlHumOversamp, BME280HumidityOversampling(0x05))
	gobottest.Assert(t, d.confFilter, BMP280IIRFilter(0x00))
	gobottest.Refute(t, d.calCoeffs, nil)
}

func TestBME280Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewBME280Driver(newI2cTestAdaptor(), WithBus(2),
		WithBME280PressureOversampling(0x01),
		WithBME280TemperatureOversampling(0x02),
		WithBME280IIRFilter(0x03),
		WithBME280HumidityOversampling(0x04))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
	gobottest.Assert(t, d.ctrlPressOversamp, BMP280PressureOversampling(0x01))
	gobottest.Assert(t, d.ctrlTempOversamp, BMP280TemperatureOversampling(0x02))
	gobottest.Assert(t, d.confFilter, BMP280IIRFilter(0x03))
	gobottest.Assert(t, d.ctrlHumOversamp, BME280HumidityOversampling(0x04))
}

func TestBME280Measurements(t *testing.T) {
	bme280, adaptor := initTestBME280WithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values produced by dumping data from actual sensor
		if adaptor.written[len(adaptor.written)-1] == bmp280RegCalib00 {
			buf.Write([]byte{126, 109, 214, 102, 50, 0, 54, 149, 220, 213, 208, 11, 64, 30, 166, 255, 249, 255, 172, 38, 10, 216, 189, 16})
		} else if adaptor.written[len(adaptor.written)-1] == bme280RegCalibDigH1 {
			buf.Write([]byte{75})
		} else if adaptor.written[len(adaptor.written)-1] == bmp280RegTempData {
			buf.Write([]byte{129, 0, 0})
		} else if adaptor.written[len(adaptor.written)-1] == bme280RegCalibDigH2LSB {
			buf.Write([]byte{112, 1, 0, 19, 1, 0, 30})
		} else if adaptor.written[len(adaptor.written)-1] == bme280RegHumidityMSB {
			buf.Write([]byte{111, 83})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	bme280.Start()
	hum, err := bme280.Humidity()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, hum, float32(51.20179))
}

func TestBME280InitH1Error(t *testing.T) {
	bme280, adaptor := initTestBME280WithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values produced by dumping data from actual sensor
		if adaptor.written[len(adaptor.written)-1] == bmp280RegCalib00 {
			buf.Write([]byte{126, 109, 214, 102, 50, 0, 54, 149, 220, 213, 208, 11, 64, 30, 166, 255, 249, 255, 172, 38, 10, 216, 189, 16})
		} else if adaptor.written[len(adaptor.written)-1] == bme280RegCalibDigH1 {
			return 0, errors.New("h1 read error")
		} else if adaptor.written[len(adaptor.written)-1] == bme280RegCalibDigH2LSB {
			buf.Write([]byte{112, 1, 0, 19, 1, 0, 30})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}

	gobottest.Assert(t, bme280.Start(), errors.New("h1 read error"))
}

func TestBME280InitH2Error(t *testing.T) {
	bme280, adaptor := initTestBME280WithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values produced by dumping data from actual sensor
		if adaptor.written[len(adaptor.written)-1] == bmp280RegCalib00 {
			buf.Write([]byte{126, 109, 214, 102, 50, 0, 54, 149, 220, 213, 208, 11, 64, 30, 166, 255, 249, 255, 172, 38, 10, 216, 189, 16})
		} else if adaptor.written[len(adaptor.written)-1] == bme280RegCalibDigH1 {
			buf.Write([]byte{75})
		} else if adaptor.written[len(adaptor.written)-1] == bme280RegCalibDigH2LSB {
			return 0, errors.New("h2 read error")
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}

	gobottest.Assert(t, bme280.Start(), errors.New("h2 read error"))
}

func TestBME280HumidityWriteError(t *testing.T) {
	bme280, adaptor := initTestBME280WithStubbedAdaptor()
	bme280.Start()

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	hum, err := bme280.Humidity()
	gobottest.Assert(t, err, errors.New("write error"))
	gobottest.Assert(t, hum, float32(0.0))
}

func TestBME280HumidityReadError(t *testing.T) {
	bme280, adaptor := initTestBME280WithStubbedAdaptor()
	bme280.Start()

	adaptor.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}
	hum, err := bme280.Humidity()
	gobottest.Assert(t, err, errors.New("read error"))
	gobottest.Assert(t, hum, float32(0.0))
}

func TestBME280HumidityNotEnabled(t *testing.T) {
	bme280, adaptor := initTestBME280WithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values produced by dumping data from actual sensor
		if adaptor.written[len(adaptor.written)-1] == bmp280RegCalib00 {
			buf.Write([]byte{126, 109, 214, 102, 50, 0, 54, 149, 220, 213, 208, 11, 64, 30, 166, 255, 249, 255, 172, 38, 10, 216, 189, 16})
		} else if adaptor.written[len(adaptor.written)-1] == bme280RegCalibDigH1 {
			buf.Write([]byte{75})
		} else if adaptor.written[len(adaptor.written)-1] == bmp280RegTempData {
			buf.Write([]byte{129, 0, 0})
		} else if adaptor.written[len(adaptor.written)-1] == bme280RegCalibDigH2LSB {
			buf.Write([]byte{112, 1, 0, 19, 1, 0, 30})
		} else if adaptor.written[len(adaptor.written)-1] == bme280RegHumidityMSB {
			buf.Write([]byte{0x80, 0x00})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	bme280.Start()
	hum, err := bme280.Humidity()
	gobottest.Assert(t, err, errors.New("Humidity disabled"))
	gobottest.Assert(t, hum, float32(0.0))
}

func TestBME280_initializationBME280(t *testing.T) {
	bme280, adaptor := initTestBME280WithStubbedAdaptor()
	readCallCounter := 0
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		readCallCounter++
		if readCallCounter == 1 {
			// Simulate returning 24 bytes for the coefficients (register bmp280RegCalib00)
			return 24, nil
		}
		if readCallCounter == 2 {
			// Simulate returning a single byte for the hc.h1 value (register bme280RegCalibDigH1)
			return 1, nil
		}
		if readCallCounter == 3 {
			// Simulate returning 7 bytes for the coefficients (register bme280RegCalibDigH2LSB)
			return 7, nil
		}
		if readCallCounter == 4 {
			// Simulate returning 1 byte for the cmr (register bmp280RegControl)
			return 1, nil
		}
		return 0, nil
	}
	gobottest.Assert(t, bme280.Start(), nil)
}
