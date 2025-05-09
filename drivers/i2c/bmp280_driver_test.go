package i2c

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*BMP280Driver)(nil)

func initTestBMP280WithStubbedAdaptor() (*BMP280Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewBMP280Driver(adaptor), adaptor
}

func TestNewBMP280Driver(t *testing.T) {
	var di interface{} = NewBMP280Driver(newI2cTestAdaptor())
	d, ok := di.(*BMP280Driver)
	if !ok {
		require.Fail(t, "NewBMP280Driver() should have returned a *BMP280Driver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "BMP280"))
	assert.Equal(t, 0x77, d.defaultAddress)
	assert.Equal(t, uint8(0x03), d.ctrlPwrMode)
	assert.Equal(t, BMP280PressureOversampling(0x05), d.ctrlPressOversamp)
	assert.Equal(t, BMP280TemperatureOversampling(0x01), d.ctrlTempOversamp)
	assert.Equal(t, BMP280IIRFilter(0x00), d.confFilter)
	assert.NotNil(t, d.calCoeffs)
}

func TestBMP280Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewBMP280Driver(newI2cTestAdaptor(), WithBus(2), WithBMP280PressureOversampling(0x04))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
	assert.Equal(t, BMP280PressureOversampling(0x04), d.ctrlPressOversamp)
}

func TestWithBMP280TemperatureOversampling(t *testing.T) {
	// arrange
	d, a := initTestBMP280WithStubbedAdaptor()
	a.written = []byte{} // reset writes of former test
	const (
		setVal = BMP280TemperatureOversampling(0x04) // 8 x
	)
	// act
	WithBMP280TemperatureOversampling(setVal)(d)
	// assert
	assert.Equal(t, setVal, d.ctrlTempOversamp)
	assert.Empty(t, a.written)
}

func TestWithBMP280IIRFilter(t *testing.T) {
	// arrange
	d, a := initTestBMP280WithStubbedAdaptor()
	a.written = []byte{} // reset writes of former test
	const (
		setVal = BMP280IIRFilter(0x02) // 4 x
	)
	// act
	WithBMP280IIRFilter(setVal)(d)
	// assert
	assert.Equal(t, setVal, d.confFilter)
	assert.Empty(t, a.written)
}

func TestBMP280Measurements(t *testing.T) {
	d, adaptor := initTestBMP280WithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values produced by dumping data from actual sensor
		switch {
		case adaptor.written[len(adaptor.written)-1] == bmp280RegCalib00:
			buf.Write([]byte{
				126, 109, 214, 102, 50, 0, 54, 149, 220, 213, 208, 11, 64, 30, 166, 255, 249, 255, 172, 38, 10, 216, 189, 16,
			})
		case adaptor.written[len(adaptor.written)-1] == bmp280RegTempData:
			buf.Write([]byte{128, 243, 0})
		case adaptor.written[len(adaptor.written)-1] == bmp280RegPressureData:
			buf.Write([]byte{77, 23, 48})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	_ = d.Start()
	temp, err := d.Temperature()
	require.NoError(t, err)
	assert.InDelta(t, float32(25.014637), temp, 0.0)
	pressure, err := d.Pressure()
	require.NoError(t, err)
	assert.InDelta(t, float32(99545.414), pressure, 0.0)
	alt, err := d.Altitude()
	require.NoError(t, err)
	assert.InDelta(t, float32(149.22713), alt, 0.0)
}

func TestBMP280TemperatureWriteError(t *testing.T) {
	d, adaptor := initTestBMP280WithStubbedAdaptor()
	_ = d.Start()

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	temp, err := d.Temperature()
	require.ErrorContains(t, err, "write error")
	assert.InDelta(t, float32(0.0), temp, 0.0)
}

func TestBMP280TemperatureReadError(t *testing.T) {
	d, adaptor := initTestBMP280WithStubbedAdaptor()
	_ = d.Start()

	adaptor.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}
	temp, err := d.Temperature()
	require.ErrorContains(t, err, "read error")
	assert.InDelta(t, float32(0.0), temp, 0.0)
}

func TestBMP280PressureWriteError(t *testing.T) {
	d, adaptor := initTestBMP280WithStubbedAdaptor()
	_ = d.Start()

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	press, err := d.Pressure()
	require.ErrorContains(t, err, "write error")
	assert.InDelta(t, float32(0.0), press, 0.0)
}

func TestBMP280PressureReadError(t *testing.T) {
	d, adaptor := initTestBMP280WithStubbedAdaptor()
	_ = d.Start()

	adaptor.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}
	press, err := d.Pressure()
	require.ErrorContains(t, err, "read error")
	assert.InDelta(t, float32(0.0), press, 0.0)
}

func TestBMP280_initialization(t *testing.T) {
	// sequence to read and write in initialization():
	// * read 24 bytes (12 x 16 bit calibration data), starting from TC1 register (0x88)
	// * fill calibration struct with data (LSByte read first)
	// * prepare the content of control register
	// * write the control register (0xF4)
	// * prepare the content of config register
	// * write the config register (0xF5)
	// arrange
	d, a := initTestBMP280WithStubbedAdaptor()
	a.written = []byte{} // reset writes of former test
	const (
		wantCalibReg   = uint8(0x88)
		wantCtrlReg    = uint8(0xF4)
		wantCtrlRegVal = uint8(0x37) // normal power mode, 16 x pressure and 1 x temperature oversampling
		wantConfReg    = uint8(0xF5)
		wantConfRegVal = uint8(0x00) // no SPI, no filter, smallest standby (unused, because normal power mode)
	)
	// Values from the datasheet example.
	t1 := []uint8{0x70, 0x6B}
	t2 := []uint8{0x43, 0x67}
	t3 := []uint8{0x18, 0xFC}
	p1 := []uint8{0x7D, 0x8E}
	p2 := []uint8{0x43, 0xD6}
	p3 := []uint8{0xD0, 0x0B}
	p4 := []uint8{0x27, 0x0B}
	p5 := []uint8{0x8C, 0x00}
	p6 := []uint8{0xF9, 0xFF}
	p7 := []uint8{0x8C, 0x3C}
	p8 := []uint8{0xF8, 0xC6}
	p9 := []uint8{0x70, 0x17}
	returnRead := append(append(append(append(append(append(t1, t2...), t3...), p1...), p2...), p3...), p4...)
	returnRead = append(append(append(append(append(returnRead, p5...), p6...), p7...), p8...), p9...)
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
	assert.Len(t, a.written, 5)
	assert.Equal(t, wantCalibReg, a.written[0])
	assert.Equal(t, wantCtrlReg, a.written[1])
	assert.Equal(t, wantCtrlRegVal, a.written[2])
	assert.Equal(t, wantConfReg, a.written[3])
	assert.Equal(t, wantConfRegVal, a.written[4])
	assert.Equal(t, uint16(27504), d.calCoeffs.t1)
	assert.Equal(t, int16(26435), d.calCoeffs.t2)
	assert.Equal(t, int16(-1000), d.calCoeffs.t3)
	assert.Equal(t, uint16(36477), d.calCoeffs.p1)
	assert.Equal(t, int16(-10685), d.calCoeffs.p2)
	assert.Equal(t, int16(3024), d.calCoeffs.p3)
	assert.Equal(t, int16(2855), d.calCoeffs.p4)
	assert.Equal(t, int16(140), d.calCoeffs.p5)
	assert.Equal(t, int16(-7), d.calCoeffs.p6)
	assert.Equal(t, int16(15500), d.calCoeffs.p7)
	assert.Equal(t, int16(-14600), d.calCoeffs.p8)
	assert.Equal(t, int16(6000), d.calCoeffs.p9)
}
