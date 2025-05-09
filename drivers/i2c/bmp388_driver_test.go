package i2c

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*BMP388Driver)(nil)

func initTestBMP388WithStubbedAdaptor() (*BMP388Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()

	readCallCounter := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		readCallCounter++
		if readCallCounter == 1 {
			buf := new(bytes.Buffer)
			// Simulate returning of 0x50 for the
			// ReadByteData(bmp388RegChipID) call in initialisation()
			_ = binary.Write(buf, binary.LittleEndian, uint8(0x50))
			copy(b, buf.Bytes())
			return buf.Len(), nil
		}
		if readCallCounter == 2 {
			// Simulate returning 24 bytes for the coefficients (register bmp388RegCalib00)
			return 24, nil
		}
		return 0, nil
	}
	return NewBMP388Driver(a), a
}

func TestNewBMP388Driver(t *testing.T) {
	var di interface{} = NewBMP388Driver(newI2cTestAdaptor())
	d, ok := di.(*BMP388Driver)
	if !ok {
		require.Fail(t, "NewBMP388Driver() should have returned a *BMP388Driver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "BMP388"))
	assert.Equal(t, 0x77, d.defaultAddress)
	assert.Equal(t, uint8(0x01), d.ctrlPwrMode)          // forced mode
	assert.Equal(t, BMP388IIRFilter(0x00), d.confFilter) // filter off
	assert.NotNil(t, d.calCoeffs)
}

func TestBMP388Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewBMP388Driver(newI2cTestAdaptor(), WithBus(2), WithBMP388IIRFilter(BMP388IIRFilter(0x03)))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
	assert.Equal(t, BMP388IIRFilter(0x03), d.confFilter)
}

func TestBMP388Measurements(t *testing.T) {
	d, a := initTestBMP388WithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		lastWritten := a.written[len(a.written)-1]
		switch lastWritten {
		case bmp388RegChipID:
			// Simulate returning of 0x50 for the
			// ReadByteData(bmp388RegChipID) call in initialisation()
			_ = binary.Write(buf, binary.LittleEndian, uint8(0x50))
		case bmp388RegCalib00:
			// Values produced by dumping data from actual sensor
			buf.Write([]byte{
				36, 107, 156, 73, 246, 104, 255, 189, 245, 35, 0, 151, 101, 184, 122, 243, 246, 211, 64, 14, 196, 0, 0, 0,
			})
		case bmp388RegTempData:
			buf.Write([]byte{0, 28, 127})
		case bmp388RegPressureData:
			buf.Write([]byte{0, 66, 113})
		}

		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	_ = d.Start()
	temp, err := d.Temperature(2)
	require.NoError(t, err)
	assert.InDelta(t, float32(22.906143), temp, 0.0)
	pressure, err := d.Pressure(2)
	require.NoError(t, err)
	assert.InDelta(t, float32(98874.85), pressure, 0.0)
	alt, err := d.Altitude(2)
	require.NoError(t, err)
	assert.InDelta(t, float32(205.89395), alt, 0.0)
}

func TestBMP388TemperatureWriteError(t *testing.T) {
	d, a := initTestBMP388WithStubbedAdaptor()
	_ = d.Start()

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	temp, err := d.Temperature(2)
	require.ErrorContains(t, err, "write error")
	assert.InDelta(t, float32(0.0), temp, 0.0)
}

func TestBMP388TemperatureReadError(t *testing.T) {
	d, a := initTestBMP388WithStubbedAdaptor()
	_ = d.Start()

	a.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}
	temp, err := d.Temperature(2)
	require.ErrorContains(t, err, "read error")
	assert.InDelta(t, float32(0.0), temp, 0.0)
}

func TestBMP388PressureWriteError(t *testing.T) {
	d, a := initTestBMP388WithStubbedAdaptor()
	_ = d.Start()

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	press, err := d.Pressure(2)
	require.ErrorContains(t, err, "write error")
	assert.InDelta(t, float32(0.0), press, 0.0)
}

func TestBMP388PressureReadError(t *testing.T) {
	d, a := initTestBMP388WithStubbedAdaptor()
	_ = d.Start()

	a.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}
	press, err := d.Pressure(2)
	require.ErrorContains(t, err, "read error")
	assert.InDelta(t, float32(0.0), press, 0.0)
}

func TestBMP388_initialization(t *testing.T) {
	// sequence to read and write in initialization():
	// * read chip ID register (0x00) and compare
	// * read 24 bytes (12 x 16 bit calibration data), starting from TC1 register (0x31)
	// * fill calibration struct with data (LSByte read first)
	// * perform a soft reset by command register (0x7E)
	// * prepare the content of config register
	// * write the config register (0x1F)
	// arrange
	d, a := initTestBMP388WithStubbedAdaptor()
	a.written = []byte{} // reset writes of former test
	const (
		wantChipIDReg     = uint8(0x00)
		wantCalibReg      = uint8(0x31)
		wantCommandReg    = uint8(0x7E)
		wantCommandRegVal = uint8(0xB6) // soft reset
		wantConfReg       = uint8(0x1F)
		wantConfRegVal    = uint8(0x00) // no filter
	)
	// Values produced by dumping data from actual sensor
	returnRead := []byte{
		36, 107, 156, 73, 246, 104, 255, 189, 245, 35, 0, 151, 101, 184, 122, 243, 246, 211, 64, 14, 196, 0, 0, 0,
	}
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		if numCallsRead == 1 {
			b[0] = 0x50
		} else {
			copy(b, returnRead)
		}
		return len(b), nil
	}
	// act, assert - initialization() must be called on Start()
	err := d.Start()
	// assert
	require.NoError(t, err)
	assert.Equal(t, 2, numCallsRead)
	assert.Len(t, a.written, 6)
	assert.Equal(t, wantChipIDReg, a.written[0])
	assert.Equal(t, wantCalibReg, a.written[1])
	assert.Equal(t, wantCommandReg, a.written[2])
	assert.Equal(t, wantCommandRegVal, a.written[3])
	assert.Equal(t, wantConfReg, a.written[4])
	assert.Equal(t, wantConfRegVal, a.written[5])
	assert.InDelta(t, float32(7.021568e+06), d.calCoeffs.t1, 0.0)
	assert.InDelta(t, float32(1.7549843e-05), d.calCoeffs.t2, 0.0)
	assert.InDelta(t, float32(-3.5527137e-14), d.calCoeffs.t3, 0.0)
	assert.InDelta(t, float32(-0.015769958), d.calCoeffs.p1, 0.0)
	assert.InDelta(t, float32(-3.5410747e-05), d.calCoeffs.p2, 0.0)
	assert.InDelta(t, float32(8.1490725e-09), d.calCoeffs.p3, 0.0)
	assert.InDelta(t, float32(0), d.calCoeffs.p4, 0.0)
	assert.InDelta(t, float32(208056), d.calCoeffs.p5, 0.0)
	assert.InDelta(t, float32(490.875), d.calCoeffs.p6, 0.0)
	assert.InDelta(t, float32(-0.05078125), d.calCoeffs.p7, 0.0)
	assert.InDelta(t, float32(-0.00030517578), d.calCoeffs.p8, 0.0)
	assert.InDelta(t, float32(5.8957283e-11), d.calCoeffs.p9, 0.0)
}
