package i2c

import (
	"bytes"
	"errors"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*BME280Driver)(nil)

// --------- HELPERS
func initTestBME280Driver() (driver *BME280Driver) {
	driver, _ = initTestBME280DriverWithStubbedAdaptor()
	return
}

func initTestBME280DriverWithStubbedAdaptor() (*BME280Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewBME280Driver(adaptor), adaptor
}

// --------- TESTS

func TestNewBME280Driver(t *testing.T) {
	// Does it return a pointer to an instance of BME280Driver?
	var bme280 interface{} = NewBME280Driver(newI2cTestAdaptor())
	_, ok := bme280.(*BME280Driver)
	if !ok {
		t.Errorf("NewBME280Driver() should have returned a *BME280Driver")
	}
}

func TestBME280Driver(t *testing.T) {
	bme280 := initTestBME280Driver()
	gobottest.Refute(t, bme280.Connection(), nil)
}

func TestBME280DriverStart(t *testing.T) {
	bme280, adaptor := initTestBME280DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		// Simulate returning a single byte for the
		// ReadByteData(bmp280RegisterControl) call in Start()
		return 1, nil
	}
	gobottest.Assert(t, bme280.Start(), nil)
}

func TestBME280StartConnectError(t *testing.T) {
	d, adaptor := initTestBME280DriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestBME280DriverStartWriteError(t *testing.T) {
	bme280, adaptor := initTestBME280DriverWithStubbedAdaptor()
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, bme280.Start(), errors.New("write error"))
}

func TestBME280DriverStartReadError(t *testing.T) {
	bme280, adaptor := initTestBME280DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}
	gobottest.Assert(t, bme280.Start(), errors.New("read error"))
}

func TestBME280DriverHalt(t *testing.T) {
	bme280 := initTestBME280Driver()

	gobottest.Assert(t, bme280.Halt(), nil)
}

func TestBME280DriverMeasurements(t *testing.T) {
	bme280, adaptor := initTestBME280DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values produced by dumping data from actual sensor
		if adaptor.written[len(adaptor.written)-1] == bmp280RegisterCalib00 {
			buf.Write([]byte{126, 109, 214, 102, 50, 0, 54, 149, 220, 213, 208, 11, 64, 30, 166, 255, 249, 255, 172, 38, 10, 216, 189, 16})
		} else if adaptor.written[len(adaptor.written)-1] == bme280RegisterCalibDigH1 {
			buf.Write([]byte{75})
		} else if adaptor.written[len(adaptor.written)-1] == bmp280RegisterTempData {
			buf.Write([]byte{129, 0, 0})
		} else if adaptor.written[len(adaptor.written)-1] == bme280RegisterCalibDigH2LSB {
			buf.Write([]byte{112, 1, 0, 19, 1, 0, 30})
		} else if adaptor.written[len(adaptor.written)-1] == bme280RegisterHumidityMSB {
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

func TestBME280DriverInitH1Error(t *testing.T) {
	bme280, adaptor := initTestBME280DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values produced by dumping data from actual sensor
		if adaptor.written[len(adaptor.written)-1] == bmp280RegisterCalib00 {
			buf.Write([]byte{126, 109, 214, 102, 50, 0, 54, 149, 220, 213, 208, 11, 64, 30, 166, 255, 249, 255, 172, 38, 10, 216, 189, 16})
		} else if adaptor.written[len(adaptor.written)-1] == bme280RegisterCalibDigH1 {
			return 0, errors.New("h1 read error")
		} else if adaptor.written[len(adaptor.written)-1] == bme280RegisterCalibDigH2LSB {
			buf.Write([]byte{112, 1, 0, 19, 1, 0, 30})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}

	gobottest.Assert(t, bme280.Start(), errors.New("h1 read error"))
}

func TestBME280DriverInitH2Error(t *testing.T) {
	bme280, adaptor := initTestBME280DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values produced by dumping data from actual sensor
		if adaptor.written[len(adaptor.written)-1] == bmp280RegisterCalib00 {
			buf.Write([]byte{126, 109, 214, 102, 50, 0, 54, 149, 220, 213, 208, 11, 64, 30, 166, 255, 249, 255, 172, 38, 10, 216, 189, 16})
		} else if adaptor.written[len(adaptor.written)-1] == bme280RegisterCalibDigH1 {
			buf.Write([]byte{75})
		} else if adaptor.written[len(adaptor.written)-1] == bme280RegisterCalibDigH2LSB {
			return 0, errors.New("h2 read error")
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}

	gobottest.Assert(t, bme280.Start(), errors.New("h2 read error"))
}

func TestBME280DriverHumidityWriteError(t *testing.T) {
	bme280, adaptor := initTestBME280DriverWithStubbedAdaptor()
	bme280.Start()

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	hum, err := bme280.Humidity()
	gobottest.Assert(t, err, errors.New("write error"))
	gobottest.Assert(t, hum, float32(0.0))
}

func TestBME280DriverHumidityReadError(t *testing.T) {
	bme280, adaptor := initTestBME280DriverWithStubbedAdaptor()
	bme280.Start()

	adaptor.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}
	hum, err := bme280.Humidity()
	gobottest.Assert(t, err, errors.New("read error"))
	gobottest.Assert(t, hum, float32(0.0))
}

func TestBME280DriverHumidityNotEnabled(t *testing.T) {
	bme280, adaptor := initTestBME280DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values produced by dumping data from actual sensor
		if adaptor.written[len(adaptor.written)-1] == bmp280RegisterCalib00 {
			buf.Write([]byte{126, 109, 214, 102, 50, 0, 54, 149, 220, 213, 208, 11, 64, 30, 166, 255, 249, 255, 172, 38, 10, 216, 189, 16})
		} else if adaptor.written[len(adaptor.written)-1] == bme280RegisterCalibDigH1 {
			buf.Write([]byte{75})
		} else if adaptor.written[len(adaptor.written)-1] == bmp280RegisterTempData {
			buf.Write([]byte{129, 0, 0})
		} else if adaptor.written[len(adaptor.written)-1] == bme280RegisterCalibDigH2LSB {
			buf.Write([]byte{112, 1, 0, 19, 1, 0, 30})
		} else if adaptor.written[len(adaptor.written)-1] == bme280RegisterHumidityMSB {
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

func TestBME280DriverSetName(t *testing.T) {
	b := initTestBME280Driver()
	b.SetName("TESTME")
	gobottest.Assert(t, b.Name(), "TESTME")
}

func TestBME280DriverOptions(t *testing.T) {
	b := NewBME280Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, b.GetBusOrDefault(1), 2)
}
