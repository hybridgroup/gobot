package i2c

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*BMP180Driver)(nil)

// --------- HELPERS
func initTestBMP180Driver() (driver *BMP180Driver) {
	driver, _ = initTestBMP180DriverWithStubbedAdaptor()
	return
}

func initTestBMP180DriverWithStubbedAdaptor() (*BMP180Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewBMP180Driver(adaptor), adaptor
}

// --------- TESTS

func TestNewBMP180Driver(t *testing.T) {
	// Does it return a pointer to an instance of BMP180Driver?
	var bmp180 interface{} = NewBMP180Driver(newI2cTestAdaptor())
	_, ok := bmp180.(*BMP180Driver)
	if !ok {
		t.Errorf("NewBMP180Driver() should have returned a *BMP180Driver")
	}
}

func TestBMP180Driver(t *testing.T) {
	bmp180 := initTestBMP180Driver()
	gobottest.Refute(t, bmp180.Connection(), nil)
}

func TestBMP180DriverStart(t *testing.T) {
	bmp180, _ := initTestBMP180DriverWithStubbedAdaptor()
	gobottest.Assert(t, bmp180.Start(), nil)
}

func TestBMP180StartConnectError(t *testing.T) {
	d, adaptor := initTestBMP180DriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestBMP180DriverStartWriteError(t *testing.T) {
	bmp180, adaptor := initTestBMP180DriverWithStubbedAdaptor()
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, bmp180.Start(), errors.New("write error"))
}

func TestBMP180DriverHalt(t *testing.T) {
	bmp180 := initTestBMP180Driver()

	gobottest.Assert(t, bmp180.Halt(), nil)
}

func TestBMP180DriverMeasurements(t *testing.T) {
	bmp180, adaptor := initTestBMP180DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values from the datasheet example.
		if adaptor.written[len(adaptor.written)-1] == bmp180RegisterAC1MSB {
			binary.Write(buf, binary.BigEndian, int16(408))
			binary.Write(buf, binary.BigEndian, int16(-72))
			binary.Write(buf, binary.BigEndian, int16(-14383))
			binary.Write(buf, binary.BigEndian, uint16(32741))
			binary.Write(buf, binary.BigEndian, uint16(32757))
			binary.Write(buf, binary.BigEndian, uint16(23153))
			binary.Write(buf, binary.BigEndian, int16(6190))
			binary.Write(buf, binary.BigEndian, int16(4))
			binary.Write(buf, binary.BigEndian, int16(-32768))
			binary.Write(buf, binary.BigEndian, int16(-8711))
			binary.Write(buf, binary.BigEndian, int16(2868))
		} else if adaptor.written[len(adaptor.written)-2] == bmp180CmdTemp && adaptor.written[len(adaptor.written)-1] == bmp180RegisterTempMSB {
			binary.Write(buf, binary.BigEndian, int16(27898))
		} else if adaptor.written[len(adaptor.written)-2] == bmp180CmdPressure && adaptor.written[len(adaptor.written)-1] == bmp180RegisterPressureMSB {
			binary.Write(buf, binary.BigEndian, int16(23843))
			// XLSB, not used in this test.
			buf.WriteByte(0)
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	bmp180.Start()
	temp, err := bmp180.Temperature()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, temp, float32(15.0))
	pressure, err := bmp180.Pressure()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, pressure, float32(69964))
}

func TestBMP180DriverTemperatureError(t *testing.T) {
	bmp180, adaptor := initTestBMP180DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values from the datasheet example.
		if adaptor.written[len(adaptor.written)-1] == bmp180RegisterAC1MSB {
			binary.Write(buf, binary.BigEndian, int16(408))
			binary.Write(buf, binary.BigEndian, int16(-72))
			binary.Write(buf, binary.BigEndian, int16(-14383))
			binary.Write(buf, binary.BigEndian, uint16(32741))
			binary.Write(buf, binary.BigEndian, uint16(32757))
			binary.Write(buf, binary.BigEndian, uint16(23153))
			binary.Write(buf, binary.BigEndian, int16(6190))
			binary.Write(buf, binary.BigEndian, int16(4))
			binary.Write(buf, binary.BigEndian, int16(-32768))
			binary.Write(buf, binary.BigEndian, int16(-8711))
			binary.Write(buf, binary.BigEndian, int16(2868))
		} else if adaptor.written[len(adaptor.written)-2] == bmp180CmdTemp && adaptor.written[len(adaptor.written)-1] == bmp180RegisterTempMSB {
			return 0, errors.New("temp error")
		} else if adaptor.written[len(adaptor.written)-2] == bmp180CmdPressure && adaptor.written[len(adaptor.written)-1] == bmp180RegisterPressureMSB {
			binary.Write(buf, binary.BigEndian, int16(23843))
			// XLSB, not used in this test.
			buf.WriteByte(0)
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	bmp180.Start()
	_, err := bmp180.Temperature()
	gobottest.Assert(t, err, errors.New("temp error"))
}

func TestBMP180DriverPressureError(t *testing.T) {
	bmp180, adaptor := initTestBMP180DriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values from the datasheet example.
		if adaptor.written[len(adaptor.written)-1] == bmp180RegisterAC1MSB {
			binary.Write(buf, binary.BigEndian, int16(408))
			binary.Write(buf, binary.BigEndian, int16(-72))
			binary.Write(buf, binary.BigEndian, int16(-14383))
			binary.Write(buf, binary.BigEndian, uint16(32741))
			binary.Write(buf, binary.BigEndian, uint16(32757))
			binary.Write(buf, binary.BigEndian, uint16(23153))
			binary.Write(buf, binary.BigEndian, int16(6190))
			binary.Write(buf, binary.BigEndian, int16(4))
			binary.Write(buf, binary.BigEndian, int16(-32768))
			binary.Write(buf, binary.BigEndian, int16(-8711))
			binary.Write(buf, binary.BigEndian, int16(2868))
		} else if adaptor.written[len(adaptor.written)-2] == bmp180CmdTemp && adaptor.written[len(adaptor.written)-1] == bmp180RegisterTempMSB {
			binary.Write(buf, binary.BigEndian, int16(27898))
		} else if adaptor.written[len(adaptor.written)-2] == bmp180CmdPressure && adaptor.written[len(adaptor.written)-1] == bmp180RegisterPressureMSB {
			return 0, errors.New("press error")
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	bmp180.Start()
	_, err := bmp180.Pressure()
	gobottest.Assert(t, err, errors.New("press error"))
}

func TestBMP180DriverPressureWriteError(t *testing.T) {
	bmp180, adaptor := initTestBMP180DriverWithStubbedAdaptor()
	bmp180.Start()

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	_, err := bmp180.Pressure()
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestBMP180DriverSetName(t *testing.T) {
	b := initTestBMP180Driver()
	b.SetName("TESTME")
	gobottest.Assert(t, b.Name(), "TESTME")
}

func TestBMP180DriverOptions(t *testing.T) {
	b := NewBMP180Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, b.GetBusOrDefault(1), 2)
}

func TestBMP180PauseForReading(t *testing.T) {
	gobottest.Assert(t, pauseForReading(BMP180UltraLowPower), time.Duration(5*time.Millisecond))
	gobottest.Assert(t, pauseForReading(BMP180Standard), time.Duration(8*time.Millisecond))
	gobottest.Assert(t, pauseForReading(BMP180HighResolution), time.Duration(14*time.Millisecond))
	gobottest.Assert(t, pauseForReading(BMP180UltraHighResolution), time.Duration(26*time.Millisecond))
}
