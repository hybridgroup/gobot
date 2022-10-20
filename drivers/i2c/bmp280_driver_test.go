package i2c

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
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
		t.Errorf("NewBMP280Driver() should have returned a *BMP280Driver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "BMP280"), true)
	gobottest.Assert(t, d.defaultAddress, 0x77)
	gobottest.Assert(t, d.ctrlPwrMode, uint8(0x03))
	gobottest.Assert(t, d.ctrlPressOversamp, BMP280PressureOversampling(0x05))
	gobottest.Assert(t, d.ctrlTempOversamp, BMP280TemperatureOversampling(0x01))
	gobottest.Assert(t, d.confFilter, BMP280IIRFilter(0x00))
	gobottest.Refute(t, d.calCoeffs, nil)
}

func TestBMP280Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewBMP280Driver(newI2cTestAdaptor(), WithBus(2), WithBMP280PressureOversampling(0x04))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
	gobottest.Assert(t, d.ctrlPressOversamp, BMP280PressureOversampling(0x04))
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
	gobottest.Assert(t, d.ctrlTempOversamp, setVal)
	gobottest.Assert(t, len(a.written), 0)
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
	gobottest.Assert(t, d.confFilter, setVal)
	gobottest.Assert(t, len(a.written), 0)
}

func TestBMP280Measurements(t *testing.T) {
	d, adaptor := initTestBMP280WithStubbedAdaptor()
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		// Values produced by dumping data from actual sensor
		if adaptor.written[len(adaptor.written)-1] == bmp280RegCalib00 {
			buf.Write([]byte{126, 109, 214, 102, 50, 0, 54, 149, 220, 213, 208, 11, 64, 30, 166, 255, 249, 255, 172, 38, 10, 216, 189, 16})
		} else if adaptor.written[len(adaptor.written)-1] == bmp280RegTempData {
			buf.Write([]byte{128, 243, 0})
		} else if adaptor.written[len(adaptor.written)-1] == bmp280RegPressureData {
			buf.Write([]byte{77, 23, 48})
		}
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	d.Start()
	temp, err := d.Temperature()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, temp, float32(25.014637))
	pressure, err := d.Pressure()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, pressure, float32(99545.414))
	alt, err := d.Altitude()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, alt, float32(149.22713))
}

func TestBMP280TemperatureWriteError(t *testing.T) {
	d, adaptor := initTestBMP280WithStubbedAdaptor()
	d.Start()

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	temp, err := d.Temperature()
	gobottest.Assert(t, err, errors.New("write error"))
	gobottest.Assert(t, temp, float32(0.0))
}

func TestBMP280TemperatureReadError(t *testing.T) {
	d, adaptor := initTestBMP280WithStubbedAdaptor()
	d.Start()

	adaptor.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}
	temp, err := d.Temperature()
	gobottest.Assert(t, err, errors.New("read error"))
	gobottest.Assert(t, temp, float32(0.0))
}

func TestBMP280PressureWriteError(t *testing.T) {
	d, adaptor := initTestBMP280WithStubbedAdaptor()
	d.Start()

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	press, err := d.Pressure()
	gobottest.Assert(t, err, errors.New("write error"))
	gobottest.Assert(t, press, float32(0.0))
}

func TestBMP280PressureReadError(t *testing.T) {
	d, adaptor := initTestBMP280WithStubbedAdaptor()
	d.Start()

	adaptor.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}
	press, err := d.Pressure()
	gobottest.Assert(t, err, errors.New("read error"))
	gobottest.Assert(t, press, float32(0.0))
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
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, numCallsRead, 1)
	gobottest.Assert(t, len(a.written), 5)
	gobottest.Assert(t, a.written[0], wantCalibReg)
	gobottest.Assert(t, a.written[1], wantCtrlReg)
	gobottest.Assert(t, a.written[2], wantCtrlRegVal)
	gobottest.Assert(t, a.written[3], wantConfReg)
	gobottest.Assert(t, a.written[4], wantConfRegVal)
	gobottest.Assert(t, d.calCoeffs.t1, uint16(27504))
	gobottest.Assert(t, d.calCoeffs.t2, int16(26435))
	gobottest.Assert(t, d.calCoeffs.t3, int16(-1000))
	gobottest.Assert(t, d.calCoeffs.p1, uint16(36477))
	gobottest.Assert(t, d.calCoeffs.p2, int16(-10685))
	gobottest.Assert(t, d.calCoeffs.p3, int16(3024))
	gobottest.Assert(t, d.calCoeffs.p4, int16(2855))
	gobottest.Assert(t, d.calCoeffs.p5, int16(140))
	gobottest.Assert(t, d.calCoeffs.p6, int16(-7))
	gobottest.Assert(t, d.calCoeffs.p7, int16(15500))
	gobottest.Assert(t, d.calCoeffs.p8, int16(-14600))
	gobottest.Assert(t, d.calCoeffs.p9, int16(6000))
}
