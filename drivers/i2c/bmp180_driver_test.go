package i2c

import (
	"bytes"
	"encoding/binary"
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
	return NewBMP180Driver(adaptor, BMP180UltraLowPower), adaptor
}

// --------- TESTS

func TestNewBMP180Driver(t *testing.T) {
	// Does it return a pointer to an instance of BMP180Driver?
	var bmp180 interface{} = NewBMP180Driver(newI2cTestAdaptor(), BMP180UltraLowPower)
	_, ok := bmp180.(*BMP180Driver)
	if !ok {
		t.Errorf("NewBMP180Driver() should have returned a *BMP180Driver")
	}
}

// Methods
func TestBMP180Driver(t *testing.T) {
	bmp180 := initTestBMP180Driver()

	gobottest.Refute(t, bmp180.Connection(), nil)
	gobottest.Assert(t, bmp180.interval, 10*time.Millisecond)

	bmp180 = NewBMP180Driver(newI2cTestAdaptor(), BMP180UltraLowPower, 100*time.Millisecond)
	gobottest.Assert(t, bmp180.interval, 100*time.Millisecond)

	bmp180 = NewBMP180Driver(newI2cTestAdaptor(), BMP180Standard)
	gobottest.Assert(t, bmp180.Mode(), BMP180Standard)
}

func TestBMP180DriverStart(t *testing.T) {
	bmp180, adaptor := initTestBMP180DriverWithStubbedAdaptor()

	adaptor.i2cReadImpl = func() ([]byte, error) {
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
		return buf.Bytes(), nil
	}
	gobottest.Assert(t, bmp180.Start(), nil)
	time.Sleep(100 * time.Millisecond)
	gobottest.Assert(t, bmp180.calibrationCoefficients.ac1, int16(408))
	gobottest.Assert(t,	bmp180.Pressure, float32(69964))
	gobottest.Assert(t,	bmp180.Temperature, float32(15.0))
}

func TestBMP180DriverHalt(t *testing.T) {
	bmp180 := initTestBMP180Driver()

	gobottest.Assert(t, bmp180.Halt(), nil)
}
