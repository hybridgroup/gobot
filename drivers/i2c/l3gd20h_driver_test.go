package i2c

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/gobottest"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*HMC6352Driver)(nil)

func initL3GD20HDriver() (driver *L3GD20HDriver) {
	driver, _ = initL3GD20HWithStubbedAdaptor()
	return
}

func initL3GD20HWithStubbedAdaptor() (*L3GD20HDriver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewL3GD20HDriver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewL3GD20HDriver(t *testing.T) {
	var di interface{} = NewL3GD20HDriver(newI2cTestAdaptor())
	d, ok := di.(*L3GD20HDriver)
	if !ok {
		t.Errorf("NewL3GD20HDriver() should have returned a *L3GD20HDriver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "L3GD20H"), true)
	gobottest.Assert(t, d.defaultAddress, 0x6b)
	gobottest.Assert(t, d.Scale(), L3GD20HScale250dps)
}

func TestL3GD20HOptions(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option.
	// Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewL3GD20HDriver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}

func TestL3GD20HWithL3GD20HFullScaleRange(t *testing.T) {
	var tests = map[string]struct {
		scale L3GD20HScale
		want  uint8
	}{
		"250dps": {
			scale: L3GD20HScale250dps,
			want:  0x00,
		},
		"500dps": {
			scale: L3GD20HScale500dps,
			want:  0x10,
		},
		"2001dps": {
			scale: L3GD20HScale2001dps,
			want:  0x20,
		},
		"2000dps": {
			scale: L3GD20HScale2000dps,
			want:  0x30,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d := initL3GD20HDriver()
			// act
			WithL3GD20HFullScaleRange(tc.scale)(d)
			// assert
			gobottest.Assert(t, d.scale, L3GD20HScale(tc.want))
		})
	}
}

func TestL3GD20HScale(t *testing.T) {
	var tests = map[string]struct {
		scale L3GD20HScale
		want  uint8
	}{
		"250dps": {
			scale: L3GD20HScale250dps,
			want:  0x00,
		},
		"500dps": {
			scale: L3GD20HScale500dps,
			want:  0x10,
		},
		"2001dps": {
			scale: L3GD20HScale2001dps,
			want:  0x20,
		},
		"2000dps": {
			scale: L3GD20HScale2000dps,
			want:  0x30,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d := initL3GD20HDriver()
			// act
			d.SetScale(tc.scale)
			// assert
			gobottest.Assert(t, d.scale, L3GD20HScale(tc.want))
		})
	}
}

func TestL3GD20HFullScaleRange(t *testing.T) {
	// sequence to read full scale range
	// * write control register 4 (0x23)
	// * read content and filter FS bits (bit 4, bit 5)
	//
	// arrange
	d, a := initL3GD20HWithStubbedAdaptor()
	a.written = []byte{} // reset values from Start() and previous tests
	readValue := uint8(0x10)
	a.i2cReadImpl = func(b []byte) (int, error) {
		b[0] = readValue
		return len(b), nil
	}
	// act
	got, err := d.FullScaleRange()
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.written), 1)
	gobottest.Assert(t, a.written[0], uint8(0x23))
	gobottest.Assert(t, got, readValue)
}

func TestL3GD20HMeasurement(t *testing.T) {
	// sequence for measurement
	// note: big endian transfer is configured (LSB in lower address, transferred first)
	//
	// * write X-axis angular rate data LSB register (0x28) with auto increment bit set (0xA8)
	// * read 3 x 2 bytes X, Y, Z data, big-endian (LSB, MSB)
	// * scale values by configured range (sensitivity = 1/gain)
	//
	// data table according to data sheet AN4506 example in table 7, supplemented with FS limit values
	sensitivity := float32(0.00875) // FS=245 dps
	var tests = map[string]struct {
		gyroData []byte
		wantX    float32
		wantY    float32
		wantZ    float32
	}{
		"245_200_100dps": {
			gyroData: []byte{0x60, 0x6D, 0x49, 0x59, 0xA4, 0x2C},
			wantX:    245,
			wantY:    199.99875,
			wantZ:    99.995,
		},
		"-100_-200_-245dps": {
			gyroData: []byte{0x5C, 0xD3, 0xB7, 0xA6, 0xA0, 0x92},
			wantX:    -99.995,
			wantY:    -199.99875,
			wantZ:    -245,
		},
		"1_0_-1": {
			gyroData: []byte{0x72, 0x00, 0x00, 0x00, 0x8D, 0xFF},
			wantX:    0.9975,
			wantY:    0,
			wantZ:    -1.00625,
		},
		"raw_range_int16_-32768_0_+32767": {
			gyroData: []byte{0x00, 0x80, 0x00, 0x00, 0xFF, 0x7F},
			wantX:    -286.72,
			wantY:    0,
			wantZ:    286.71124,
		},
		"raw_8_5_-3": {
			gyroData: []byte{0x08, 0x00, 0x05, 0x00, 0xFD, 0xFF},
			wantX:    float32(8) * sensitivity,
			wantY:    float32(5) * sensitivity,
			wantZ:    float32(-3) * sensitivity,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initL3GD20HWithStubbedAdaptor()
			a.written = []byte{} // reset values from Start() and previous tests
			a.i2cReadImpl = func(b []byte) (int, error) {
				copy(b, tc.gyroData)
				return len(b), nil
			}
			// act
			x, y, z, err := d.XYZ()
			// assert
			gobottest.Assert(t, err, nil)
			gobottest.Assert(t, len(a.written), 1)
			gobottest.Assert(t, a.written[0], uint8(0xA8))
			gobottest.Assert(t, x, tc.wantX)
			gobottest.Assert(t, y, tc.wantY)
			gobottest.Assert(t, z, tc.wantZ)
		})
	}
}

func TestL3GD20HMeasurementError(t *testing.T) {
	d, a := initL3GD20HWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}

	d.Start()
	_, _, _, err := d.XYZ()
	gobottest.Assert(t, err, errors.New("read error"))
}

func TestL3GD20HMeasurementWriteError(t *testing.T) {
	d, a := initL3GD20HWithStubbedAdaptor()
	a.i2cWriteImpl = func(b []byte) (int, error) {
		return 0, errors.New("write error")
	}
	_, _, _, err := d.XYZ()
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestL3GD20H_initialize(t *testing.T) {
	// sequence for initialization the device on Start()
	// * write control register 1 (0x20)
	// * write reset (0x00)
	// * write control register 1 (0x20)
	// * prepare register content:
	//  * output data rate for no Low_ODR=100Hz (DR=0x00)
	//  * bandwidth for no Low_ODR=12.5Hz (BW=0x00)
	//  * normal mode and enable all axes (PD=1, X/Y/Z=1)
	// * write register content (0x0F)
	// * write control register 4 (0x23)
	// * prepare register content
	//  * continuous block data update (BDU=0x00)
	//  * use big endian transfer (LSB in lower address, transferred first) (BLE=0x00)
	//  * set full scale selection to configured scale (default=245dps, FS=0x00)
	//  * normal self test (ST=0x00)
	//  * SPI mode to 4 wire (SIM=0x00)
	// * write register content (0x00)
	//
	// all other registers currently untouched
	//
	// arrange, act - initialize() must be called on Start()
	_, a := initL3GD20HWithStubbedAdaptor()
	// assert
	gobottest.Assert(t, len(a.written), 6)
	gobottest.Assert(t, a.written[0], uint8(0x20))
	gobottest.Assert(t, a.written[1], uint8(0x00))
	gobottest.Assert(t, a.written[2], uint8(0x20))
	gobottest.Assert(t, a.written[3], uint8(0x0F))
	gobottest.Assert(t, a.written[4], uint8(0x23))
	gobottest.Assert(t, a.written[5], uint8(0x00))
}
