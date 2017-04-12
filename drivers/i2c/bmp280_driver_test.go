package i2c

import (
	"bytes"
	"errors"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*BMP280Driver)(nil)

// --------- HELPERS
func initTestBMP280Driver() (driver *BMP280Driver) {
	driver, _ = initTestBMP280DriverWithStubbedAdaptor()
	return
}

func initTestBMP280DriverWithStubbedAdaptor() (*BMP280Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewBMP280Driver(adaptor), adaptor
}

// --------- TESTS

func TestNewBMP280Driver(t *testing.T) {
	// Does it return a pointer to an instance of BME280Driver?
	var bmp280 interface{} = NewBMP280Driver(newI2cTestAdaptor())
	_, ok := bmp280.(*BMP280Driver)
	if !ok {
		t.Errorf("NewBMP280Driver() should have returned a *BMP280Driver")
	}
}

func TestBMP280Driver(t *testing.T) {
	bmp280 := initTestBMP280Driver()
	gobottest.Refute(t, bmp280.Connection(), nil)
}

func TestBMP280DriverStart(t *testing.T) {
	bmp280, _ := initTestBMP280DriverWithStubbedAdaptor()
	gobottest.Assert(t, bmp280.Start(), nil)
}

func TestBMP280StartConnectError(t *testing.T) {
	d, adaptor := initTestBMP280DriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestBMP280DriverStartWriteError(t *testing.T) {
	bmp280, adaptor := initTestBMP280DriverWithStubbedAdaptor()
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, bmp280.Start(), errors.New("write error"))
}

func TestBMP280DriverStartReadError(t *testing.T) {
	bmp280, adaptor := initTestBMP280DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}
	gobottest.Assert(t, bmp280.Start(), errors.New("read error"))
}

func TestBMP280DriverHalt(t *testing.T) {
	bmp280 := initTestBMP280Driver()

	gobottest.Assert(t, bmp280.Halt(), nil)
}

func TestBMP280DriverMeasurements(t *testing.T) {
	bmp280, adaptor := initTestBMP280DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values produced by dumping data from actual sensor
		if adaptor.written[len(adaptor.written)-1] == bmp280RegisterCalib00 {
			buf.Write([]byte{126, 109, 214, 102, 50, 0, 54, 149, 220, 213, 208, 11, 64, 30, 166, 255, 249, 255, 172, 38, 10, 216, 189, 16})
		} else if adaptor.written[len(adaptor.written)-1] == bmp280RegisterTempData {
			buf.Write([]byte{128, 243, 0})
		} else if adaptor.written[len(adaptor.written)-1] == bmp280RegisterPressureData {
			buf.Write([]byte{77, 23, 48})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	bmp280.Start()
	temp, err := bmp280.Temperature()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, temp, float32(25.014637))
	pressure, err := bmp280.Pressure()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, pressure, float32(99545.414))
	alt, err := bmp280.Altitude()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, alt, float32(149.22713))
}

func TestBMP280DriverTemperatureWriteError(t *testing.T) {
	bmp280, adaptor := initTestBMP280DriverWithStubbedAdaptor()
	bmp280.Start()

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	temp, err := bmp280.Temperature()
	gobottest.Assert(t, err, errors.New("write error"))
	gobottest.Assert(t, temp, float32(0.0))
}

func TestBMP280DriverTemperatureReadError(t *testing.T) {
	bmp280, adaptor := initTestBMP280DriverWithStubbedAdaptor()
	bmp280.Start()

	adaptor.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}
	temp, err := bmp280.Temperature()
	gobottest.Assert(t, err, errors.New("read error"))
	gobottest.Assert(t, temp, float32(0.0))
}

func TestBMP280DriverPressureWriteError(t *testing.T) {
	bmp280, adaptor := initTestBMP280DriverWithStubbedAdaptor()
	bmp280.Start()

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	press, err := bmp280.Pressure()
	gobottest.Assert(t, err, errors.New("write error"))
	gobottest.Assert(t, press, float32(0.0))
}

func TestBMP280DriverPressureReadError(t *testing.T) {
	bmp280, adaptor := initTestBMP280DriverWithStubbedAdaptor()
	bmp280.Start()

	adaptor.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}
	press, err := bmp280.Pressure()
	gobottest.Assert(t, err, errors.New("read error"))
	gobottest.Assert(t, press, float32(0.0))
}

func TestBMP280DriverSetName(t *testing.T) {
	b := initTestBMP280Driver()
	b.SetName("TESTME")
	gobottest.Assert(t, b.Name(), "TESTME")
}

func TestBMP280DriverOptions(t *testing.T) {
	b := NewBMP280Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, b.GetBusOrDefault(1), 2)
}
