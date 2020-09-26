package i2c

import (
	"bytes"
	"errors"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*BMP388Driver)(nil)

// --------- HELPERS
func initTestBMP388Driver() (driver *BMP388Driver) {
	driver, _ = initTestBMP388DriverWithStubbedAdaptor()
	return
}

func initTestBMP388DriverWithStubbedAdaptor() (*BMP388Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewBMP388Driver(adaptor), adaptor
}

// --------- TESTS

func TestNewBMP388Driver(t *testing.T) {
	// Does it return a pointer to an instance of BMP388Driver?
	var bmp388 interface{} = NewBMP388Driver(newI2cTestAdaptor())
	_, ok := bmp388.(*BMP388Driver)
	if !ok {
		t.Errorf("NewBMP388Driver() should have returned a *BMP388Driver")
	}
}

func TestBMP388Driver(t *testing.T) {
	bmp388 := initTestBMP388Driver()
	gobottest.Refute(t, bmp388.Connection(), nil)
}

func TestBMP388DriverStart(t *testing.T) {
	bmp388, _ := initTestBMP388DriverWithStubbedAdaptor()
	gobottest.Assert(t, bmp388.Start(), nil)
}

func TestBMP388StartConnectError(t *testing.T) {
	d, adaptor := initTestBMP388DriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestBMP388DriverStartWriteError(t *testing.T) {
	bmp388, adaptor := initTestBMP388DriverWithStubbedAdaptor()
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, bmp388.Start(), errors.New("write error"))
}

func TestBMP388DriverStartReadError(t *testing.T) {
	bmp388, adaptor := initTestBMP388DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}
	gobottest.Assert(t, bmp388.Start(), errors.New("read error"))
}

func TestBMP388DriverHalt(t *testing.T) {
	bmp388 := initTestBMP388Driver()

	gobottest.Assert(t, bmp388.Halt(), nil)
}

func TestBMP388DriverMeasurements(t *testing.T) {
	bmp388, adaptor := initTestBMP388DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values produced by dumping data from actual sensor
		if adaptor.written[len(adaptor.written)-1] == bmp388RegisterCalib00 {
			buf.Write([]byte{126, 109, 214, 102, 50, 0, 54, 149, 220, 213, 208, 11, 64, 30, 166, 255, 249, 255, 172, 38, 10, 216, 189, 16})
		} else if adaptor.written[len(adaptor.written)-1] == bmp388RegisterTempData {
			buf.Write([]byte{128, 243, 0})
		} else if adaptor.written[len(adaptor.written)-1] == bmp388RegisterPressureData {
			buf.Write([]byte{77, 23, 48})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	bmp388.Start()
	temp, err := bmp388.Temperature(2)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, temp, float32(25.014637))
	pressure, err := bmp388.Pressure(2)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, pressure, float32(99545.414))
	alt, err := bmp388.Altitude(2)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, alt, float32(149.22713))
}

func TestBMP388DriverTemperatureWriteError(t *testing.T) {
	bmp388, adaptor := initTestBMP388DriverWithStubbedAdaptor()
	bmp388.Start()

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	temp, err := bmp388.Temperature(2)
	gobottest.Assert(t, err, errors.New("write error"))
	gobottest.Assert(t, temp, float32(0.0))
}

func TestBMP388DriverTemperatureReadError(t *testing.T) {
	bmp388, adaptor := initTestBMP388DriverWithStubbedAdaptor()
	bmp388.Start()

	adaptor.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}
	temp, err := bmp388.Temperature(2)
	gobottest.Assert(t, err, errors.New("read error"))
	gobottest.Assert(t, temp, float32(0.0))
}

func TestBMP388DriverPressureWriteError(t *testing.T) {
	bmp388, adaptor := initTestBMP388DriverWithStubbedAdaptor()
	bmp388.Start()

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	press, err := bmp388.Pressure(2)
	gobottest.Assert(t, err, errors.New("write error"))
	gobottest.Assert(t, press, float32(0.0))
}

func TestBMP388DriverPressureReadError(t *testing.T) {
	bmp388, adaptor := initTestBMP388DriverWithStubbedAdaptor()
	bmp388.Start()

	adaptor.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}
	press, err := bmp388.Pressure(2)
	gobottest.Assert(t, err, errors.New("read error"))
	gobottest.Assert(t, press, float32(0.0))
}

func TestBMP388DriverSetName(t *testing.T) {
	b := initTestBMP388Driver()
	b.SetName("TESTME")
	gobottest.Assert(t, b.Name(), "TESTME")
}

func TestBMP388DriverOptions(t *testing.T) {
	b := NewBMP388Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, b.GetBusOrDefault(1), 2)
}
