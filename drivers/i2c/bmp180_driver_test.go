package i2c

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*BMP180Driver)(nil)

func initTestBMP180WithStubbedAdaptor() (*BMP180Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewBMP180Driver(adaptor), adaptor
}

func TestNewBMP180Driver(t *testing.T) {
	// Does it return a pointer to an instance of BMP180Driver?
	var di interface{} = NewBMP180Driver(newI2cTestAdaptor())
	d, ok := di.(*BMP180Driver)
	if !ok {
		t.Errorf("NewBMP180Driver() should have returned a *BMP180Driver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "BMP180"), true)
	gobottest.Assert(t, d.defaultAddress, 0x77)
	gobottest.Assert(t, d.oversampling, BMP180OversamplingMode(0x00))
	gobottest.Refute(t, d.calCoeffs, nil)
}

func TestBMP180Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewBMP180Driver(newI2cTestAdaptor(), WithBus(2), WithBMP180OversamplingMode(0x01))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
	gobottest.Assert(t, d.oversampling, BMP180OversamplingMode(0x01))
}

func TestBMP180Measurements(t *testing.T) {
	bmp180, adaptor := initTestBMP180WithStubbedAdaptor()
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
		} else if adaptor.written[len(adaptor.written)-2] == bmp180CtlTemp && adaptor.written[len(adaptor.written)-1] == bmp180RegisterDataMSB {
			binary.Write(buf, binary.BigEndian, int16(27898))
		} else if adaptor.written[len(adaptor.written)-2] == bmp180CtlPressure && adaptor.written[len(adaptor.written)-1] == bmp180RegisterDataMSB {
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

func TestBMP180TemperatureError(t *testing.T) {
	bmp180, adaptor := initTestBMP180WithStubbedAdaptor()
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
		} else if adaptor.written[len(adaptor.written)-2] == bmp180CtlTemp && adaptor.written[len(adaptor.written)-1] == bmp180RegisterDataMSB {
			return 0, errors.New("temp error")
		} else if adaptor.written[len(adaptor.written)-2] == bmp180CtlPressure && adaptor.written[len(adaptor.written)-1] == bmp180RegisterDataMSB {
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

func TestBMP180PressureError(t *testing.T) {
	bmp180, adaptor := initTestBMP180WithStubbedAdaptor()
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
		} else if adaptor.written[len(adaptor.written)-2] == bmp180CtlTemp && adaptor.written[len(adaptor.written)-1] == bmp180RegisterDataMSB {
			binary.Write(buf, binary.BigEndian, int16(27898))
		} else if adaptor.written[len(adaptor.written)-2] == bmp180CtlPressure && adaptor.written[len(adaptor.written)-1] == bmp180RegisterDataMSB {
			return 0, errors.New("press error")
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	bmp180.Start()
	_, err := bmp180.Pressure()
	gobottest.Assert(t, err, errors.New("press error"))
}

func TestBMP180PressureWriteError(t *testing.T) {
	bmp180, adaptor := initTestBMP180WithStubbedAdaptor()
	bmp180.Start()

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	_, err := bmp180.Pressure()
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestBMP180_initialization(t *testing.T) {
	// sequence to read in initialization():
	// * read 22 bytes (11 x 16 bit calibration data), starting from AC1 register (0xAA)
	// * fill calibration struct with data (MSByte read first)
	// arrange
	d, a := initTestBMP180WithStubbedAdaptor()
	a.written = []byte{} // reset writes of former test
	// Values from the datasheet example.
	ac1 := []uint8{0x01, 0x98}
	ac2 := []uint8{0xFF, 0xB8}
	ac3 := []uint8{0xC7, 0xD1}
	ac4 := []uint8{0x7F, 0xE5}
	ac5 := []uint8{0x7F, 0xF5}
	ac6 := []uint8{0x5A, 0x71}
	b1 := []uint8{0x18, 0x2E}
	b2 := []uint8{0x00, 0x04}
	mb := []uint8{0x80, 0x00}
	mc := []uint8{0xDD, 0xF9}
	md := []uint8{0x0B, 0x34}
	returnRead := append(append(append(append(append(ac1, ac2...), ac3...), ac4...), ac5...), ac6...)
	returnRead = append(append(append(append(append(returnRead, b1...), b2...), mb...), mc...), md...)
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		copy(b, returnRead)
		return len(b), nil
	}
	// act, assert - initialization() must be called on Start()
	err := d.Start()
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, numCallsRead, 1)
	gobottest.Assert(t, len(a.written), 1)
	gobottest.Assert(t, a.written[0], uint8(0xAA))
	gobottest.Assert(t, d.calCoeffs.ac1, int16(408))
	gobottest.Assert(t, d.calCoeffs.ac2, int16(-72))
	gobottest.Assert(t, d.calCoeffs.ac3, int16(-14383))
	gobottest.Assert(t, d.calCoeffs.ac4, uint16(32741))
	gobottest.Assert(t, d.calCoeffs.ac5, uint16(32757))
	gobottest.Assert(t, d.calCoeffs.ac6, uint16(23153))
	gobottest.Assert(t, d.calCoeffs.b1, int16(6190))
	gobottest.Assert(t, d.calCoeffs.b2, int16(4))
	gobottest.Assert(t, d.calCoeffs.mb, int16(-32768))
	gobottest.Assert(t, d.calCoeffs.mc, int16(-8711))
	gobottest.Assert(t, d.calCoeffs.md, int16(2868))
}

func TestBMP180_bmp180PauseForReading(t *testing.T) {
	gobottest.Assert(t, bmp180PauseForReading(BMP180UltraLowPower), time.Duration(5*time.Millisecond))
	gobottest.Assert(t, bmp180PauseForReading(BMP180Standard), time.Duration(8*time.Millisecond))
	gobottest.Assert(t, bmp180PauseForReading(BMP180HighResolution), time.Duration(14*time.Millisecond))
	gobottest.Assert(t, bmp180PauseForReading(BMP180UltraHighResolution), time.Duration(26*time.Millisecond))
}
