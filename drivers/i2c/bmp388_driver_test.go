package i2c

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/gobottest"
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
			binary.Write(buf, binary.LittleEndian, uint8(0x50))
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
		t.Errorf("NewBMP388Driver() should have returned a *BMP388Driver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "BMP388"), true)
	gobottest.Assert(t, d.defaultAddress, 0x77)
	gobottest.Assert(t, d.ctrlPwrMode, uint8(0x01))          // forced mode
	gobottest.Assert(t, d.confFilter, BMP388IIRFilter(0x00)) // filter off
	gobottest.Refute(t, d.calCoeffs, nil)
}

func TestBMP388Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewBMP388Driver(newI2cTestAdaptor(), WithBus(2), WithBMP388IIRFilter(BMP388IIRFilter(0x03)))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
	gobottest.Assert(t, d.confFilter, BMP388IIRFilter(0x03))
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
			binary.Write(buf, binary.LittleEndian, uint8(0x50))
		case bmp388RegCalib00:
			// Values produced by dumping data from actual sensor
			buf.Write([]byte{36, 107, 156, 73, 246, 104, 255, 189, 245, 35, 0, 151, 101, 184, 122, 243, 246, 211, 64, 14, 196, 0, 0, 0})
		case bmp388RegTempData:
			buf.Write([]byte{0, 28, 127})
		case bmp388RegPressureData:
			buf.Write([]byte{0, 66, 113})
		}

		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	d.Start()
	temp, err := d.Temperature(2)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, temp, float32(22.906143))
	pressure, err := d.Pressure(2)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, pressure, float32(98874.85))
	alt, err := d.Altitude(2)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, alt, float32(205.89395))
}

func TestBMP388TemperatureWriteError(t *testing.T) {
	d, a := initTestBMP388WithStubbedAdaptor()
	d.Start()

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	temp, err := d.Temperature(2)
	gobottest.Assert(t, err, errors.New("write error"))
	gobottest.Assert(t, temp, float32(0.0))
}

func TestBMP388TemperatureReadError(t *testing.T) {
	d, a := initTestBMP388WithStubbedAdaptor()
	d.Start()

	a.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}
	temp, err := d.Temperature(2)
	gobottest.Assert(t, err, errors.New("read error"))
	gobottest.Assert(t, temp, float32(0.0))
}

func TestBMP388PressureWriteError(t *testing.T) {
	d, a := initTestBMP388WithStubbedAdaptor()
	d.Start()

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	press, err := d.Pressure(2)
	gobottest.Assert(t, err, errors.New("write error"))
	gobottest.Assert(t, press, float32(0.0))
}

func TestBMP388PressureReadError(t *testing.T) {
	d, a := initTestBMP388WithStubbedAdaptor()
	d.Start()

	a.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}
	press, err := d.Pressure(2)
	gobottest.Assert(t, err, errors.New("read error"))
	gobottest.Assert(t, press, float32(0.0))
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
	returnRead := []byte{36, 107, 156, 73, 246, 104, 255, 189, 245, 35, 0, 151, 101, 184, 122, 243, 246, 211, 64, 14, 196, 0, 0, 0}
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
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, numCallsRead, 2)
	gobottest.Assert(t, len(a.written), 6)
	gobottest.Assert(t, a.written[0], wantChipIDReg)
	gobottest.Assert(t, a.written[1], wantCalibReg)
	gobottest.Assert(t, a.written[2], wantCommandReg)
	gobottest.Assert(t, a.written[3], wantCommandRegVal)
	gobottest.Assert(t, a.written[4], wantConfReg)
	gobottest.Assert(t, a.written[5], wantConfRegVal)
	gobottest.Assert(t, d.calCoeffs.t1, float32(7.021568e+06))
	gobottest.Assert(t, d.calCoeffs.t2, float32(1.7549843e-05))
	gobottest.Assert(t, d.calCoeffs.t3, float32(-3.5527137e-14))
	gobottest.Assert(t, d.calCoeffs.p1, float32(-0.015769958))
	gobottest.Assert(t, d.calCoeffs.p2, float32(-3.5410747e-05))
	gobottest.Assert(t, d.calCoeffs.p3, float32(8.1490725e-09))
	gobottest.Assert(t, d.calCoeffs.p4, float32(0))
	gobottest.Assert(t, d.calCoeffs.p5, float32(208056))
	gobottest.Assert(t, d.calCoeffs.p6, float32(490.875))
	gobottest.Assert(t, d.calCoeffs.p7, float32(-0.05078125))
	gobottest.Assert(t, d.calCoeffs.p8, float32(-0.00030517578))
	gobottest.Assert(t, d.calCoeffs.p9, float32(5.8957283e-11))
}
