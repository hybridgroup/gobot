package i2c

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
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
		require.Fail(t, "NewBMP180Driver() should have returned a *BMP180Driver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "BMP180"))
	assert.Equal(t, 0x77, d.defaultAddress)
	assert.Equal(t, BMP180OversamplingMode(0x00), d.oversampling)
	assert.NotNil(t, d.calCoeffs)
}

func TestBMP180Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewBMP180Driver(newI2cTestAdaptor(), WithBus(2), WithBMP180OversamplingMode(0x01))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
	assert.Equal(t, BMP180OversamplingMode(0x01), d.oversampling)
}

func TestBMP180Measurements(t *testing.T) {
	bmp180, adaptor := initTestBMP180WithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values from the datasheet example.
		switch {
		case adaptor.written[len(adaptor.written)-1] == bmp180RegisterAC1MSB:
			_ = binary.Write(buf, binary.BigEndian, int16(408))
			_ = binary.Write(buf, binary.BigEndian, int16(-72))
			_ = binary.Write(buf, binary.BigEndian, int16(-14383))
			_ = binary.Write(buf, binary.BigEndian, uint16(32741))
			_ = binary.Write(buf, binary.BigEndian, uint16(32757))
			_ = binary.Write(buf, binary.BigEndian, uint16(23153))
			_ = binary.Write(buf, binary.BigEndian, int16(6190))
			_ = binary.Write(buf, binary.BigEndian, int16(4))
			_ = binary.Write(buf, binary.BigEndian, int16(-32768))
			_ = binary.Write(buf, binary.BigEndian, int16(-8711))
			_ = binary.Write(buf, binary.BigEndian, int16(2868))
		case adaptor.written[len(adaptor.written)-2] == bmp180CtlTemp &&
			adaptor.written[len(adaptor.written)-1] == bmp180RegisterDataMSB:
			_ = binary.Write(buf, binary.BigEndian, int16(27898))
		case adaptor.written[len(adaptor.written)-2] == bmp180CtlPressure &&
			adaptor.written[len(adaptor.written)-1] == bmp180RegisterDataMSB:
			_ = binary.Write(buf, binary.BigEndian, int16(23843))
			// XLSB, not used in this test.
			buf.WriteByte(0)
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	_ = bmp180.Start()
	temp, err := bmp180.Temperature()
	require.NoError(t, err)
	assert.InDelta(t, float32(15.0), temp, 0.0)
	pressure, err := bmp180.Pressure()
	require.NoError(t, err)
	assert.InDelta(t, float32(69964), pressure, 0.0)
}

func TestBMP180TemperatureError(t *testing.T) {
	bmp180, adaptor := initTestBMP180WithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values from the datasheet example.
		switch {
		case adaptor.written[len(adaptor.written)-1] == bmp180RegisterAC1MSB:
			_ = binary.Write(buf, binary.BigEndian, int16(408))
			_ = binary.Write(buf, binary.BigEndian, int16(-72))
			_ = binary.Write(buf, binary.BigEndian, int16(-14383))
			_ = binary.Write(buf, binary.BigEndian, uint16(32741))
			_ = binary.Write(buf, binary.BigEndian, uint16(32757))
			_ = binary.Write(buf, binary.BigEndian, uint16(23153))
			_ = binary.Write(buf, binary.BigEndian, int16(6190))
			_ = binary.Write(buf, binary.BigEndian, int16(4))
			_ = binary.Write(buf, binary.BigEndian, int16(-32768))
			_ = binary.Write(buf, binary.BigEndian, int16(-8711))
			_ = binary.Write(buf, binary.BigEndian, int16(2868))
		case adaptor.written[len(adaptor.written)-2] == bmp180CtlTemp &&
			adaptor.written[len(adaptor.written)-1] == bmp180RegisterDataMSB:
			return 0, errors.New("temp error")
		case adaptor.written[len(adaptor.written)-2] == bmp180CtlPressure &&
			adaptor.written[len(adaptor.written)-1] == bmp180RegisterDataMSB:
			_ = binary.Write(buf, binary.BigEndian, int16(23843))
			// XLSB, not used in this test.
			buf.WriteByte(0)
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	_ = bmp180.Start()
	_, err := bmp180.Temperature()
	require.ErrorContains(t, err, "temp error")
}

func TestBMP180PressureError(t *testing.T) {
	bmp180, adaptor := initTestBMP180WithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values from the datasheet example.
		switch {
		case adaptor.written[len(adaptor.written)-1] == bmp180RegisterAC1MSB:
			_ = binary.Write(buf, binary.BigEndian, int16(408))
			_ = binary.Write(buf, binary.BigEndian, int16(-72))
			_ = binary.Write(buf, binary.BigEndian, int16(-14383))
			_ = binary.Write(buf, binary.BigEndian, uint16(32741))
			_ = binary.Write(buf, binary.BigEndian, uint16(32757))
			_ = binary.Write(buf, binary.BigEndian, uint16(23153))
			_ = binary.Write(buf, binary.BigEndian, int16(6190))
			_ = binary.Write(buf, binary.BigEndian, int16(4))
			_ = binary.Write(buf, binary.BigEndian, int16(-32768))
			_ = binary.Write(buf, binary.BigEndian, int16(-8711))
			_ = binary.Write(buf, binary.BigEndian, int16(2868))
		case adaptor.written[len(adaptor.written)-2] == bmp180CtlTemp &&
			adaptor.written[len(adaptor.written)-1] == bmp180RegisterDataMSB:
			_ = binary.Write(buf, binary.BigEndian, int16(27898))
		case adaptor.written[len(adaptor.written)-2] == bmp180CtlPressure &&
			adaptor.written[len(adaptor.written)-1] == bmp180RegisterDataMSB:
			return 0, errors.New("press error")
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	_ = bmp180.Start()
	_, err := bmp180.Pressure()
	require.ErrorContains(t, err, "press error")
}

func TestBMP180PressureWriteError(t *testing.T) {
	bmp180, adaptor := initTestBMP180WithStubbedAdaptor()
	_ = bmp180.Start()

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	_, err := bmp180.Pressure()
	require.ErrorContains(t, err, "write error")
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
	require.NoError(t, err)
	assert.Equal(t, 1, numCallsRead)
	assert.Len(t, a.written, 1)
	assert.Equal(t, uint8(0xAA), a.written[0])
	assert.Equal(t, int16(408), d.calCoeffs.ac1)
	assert.Equal(t, int16(-72), d.calCoeffs.ac2)
	assert.Equal(t, int16(-14383), d.calCoeffs.ac3)
	assert.Equal(t, uint16(32741), d.calCoeffs.ac4)
	assert.Equal(t, uint16(32757), d.calCoeffs.ac5)
	assert.Equal(t, uint16(23153), d.calCoeffs.ac6)
	assert.Equal(t, int16(6190), d.calCoeffs.b1)
	assert.Equal(t, int16(4), d.calCoeffs.b2)
	assert.Equal(t, int16(-32768), d.calCoeffs.mb)
	assert.Equal(t, int16(-8711), d.calCoeffs.mc)
	assert.Equal(t, int16(2868), d.calCoeffs.md)
}

func TestBMP180_bmp180PauseForReading(t *testing.T) {
	assert.Equal(t, 5*time.Millisecond, bmp180PauseForReading(BMP180UltraLowPower))
	assert.Equal(t, 8*time.Millisecond, bmp180PauseForReading(BMP180Standard))
	assert.Equal(t, 14*time.Millisecond, bmp180PauseForReading(BMP180HighResolution))
	assert.Equal(t, 26*time.Millisecond, bmp180PauseForReading(BMP180UltraHighResolution))
}
