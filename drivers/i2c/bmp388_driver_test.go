package i2c

import (
	"bytes"
	"encoding/binary"
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
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Simulate returning of 0x50 for the
		// ReadByteData(bmp388RegisterChipID) call in initialisation()
		binary.Write(buf, binary.LittleEndian, uint8(0x50))
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
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
	bmp388, adaptor := initTestBMP388DriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, bmp388.Start(), errors.New("Invalid i2c connection"))
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
		if len(adaptor.written) == 0 {
			// Simulate returning of 0x50 for the
			// ReadByteData(bmp388RegisterChipID) call in initialisation()
			binary.Write(buf, binary.LittleEndian, uint8(0x50))
		} else if adaptor.written[len(adaptor.written)-1] == bmp388RegisterCalib00 {
			// Values produced by dumping data from actual sensor
			buf.Write([]byte{36, 107, 156, 73, 246, 104, 255, 189, 245, 35, 0, 151, 101, 184, 122, 243, 246, 211, 64, 14, 196, 0, 0, 0})
		} else if adaptor.written[len(adaptor.written)-1] == bmp388RegisterTempData {
			buf.Write([]byte{0, 28, 127})
		} else if adaptor.written[len(adaptor.written)-1] == bmp388RegisterPressureData {
			buf.Write([]byte{0, 66, 113})
		}

		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	bmp388.Start()
	temp, err := bmp388.Temperature(2)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, temp, float32(22.906143))
	pressure, err := bmp388.Pressure(2)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, pressure, float32(98874.85))
	alt, err := bmp388.Altitude(2)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, alt, float32(205.89395))
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
