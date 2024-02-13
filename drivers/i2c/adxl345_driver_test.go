package i2c

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*ADXL345Driver)(nil)

func initTestADXL345WithStubbedAdaptor() (*ADXL345Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewADXL345Driver(a)
	return d, a
}

func TestNewADXL345Driver(t *testing.T) {
	var di interface{} = NewADXL345Driver(newI2cTestAdaptor())
	d, ok := di.(*ADXL345Driver)
	if !ok {
		require.Fail(t, "NewADXL345Driver() should have returned a *ADXL345Driver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "ADXL345"))
	assert.Equal(t, 0x53, d.defaultAddress)
	assert.Equal(t, uint8(1), d.powerCtl.measure)
	assert.Equal(t, ADXL345FsRangeConfig(0x00), d.dataFormat.fullScaleRange)
	assert.Equal(t, ADXL345RateConfig(0x0A), d.bwRate.rate)
	assert.True(t, d.bwRate.lowPower)
}

func TestADXL345Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewADXL345Driver(newI2cTestAdaptor(), WithBus(2), WithADXL345LowPowerMode(false))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
	assert.False(t, d.bwRate.lowPower)
}

func TestADXL345WithADXL345DataOutputRate(t *testing.T) {
	// arrange
	d, a := initTestADXL345WithStubbedAdaptor()
	a.written = []byte{} // reset writes of former test
	const (
		setVal = ADXL345RateConfig(0x0E) // 1.6kHz
	)
	// act
	WithADXL345DataOutputRate(setVal)(d)
	// assert
	assert.Equal(t, setVal, d.bwRate.rate)
	assert.Empty(t, a.written)
}

func TestADXL345WithADXL345FullScaleRange(t *testing.T) {
	// arrange
	d, a := initTestADXL345WithStubbedAdaptor()
	a.written = []byte{} // reset writes of former test
	const (
		setVal = ADXL345FsRangeConfig(0x02) // +-8 g
	)
	// act
	WithADXL345FullScaleRange(setVal)(d)
	// assert
	assert.Equal(t, setVal, d.dataFormat.fullScaleRange)
	assert.Empty(t, a.written)
}

func TestADXL345UseLowPower(t *testing.T) {
	// sequence to set low power:
	// * set value in data rate structure
	// * write the data rate register (0x2C)
	d, a := initTestADXL345WithStubbedAdaptor()
	_ = d.Start()
	a.written = []byte{} // reset writes of former test
	setVal := !d.bwRate.lowPower
	const (
		wantReg = uint8(0x2C)
		wantVal = uint8(0x0A) // only 100 Hz left over
	)
	// act
	err := d.UseLowPower(setVal)
	// assert
	require.NoError(t, err)
	assert.Equal(t, setVal, d.bwRate.lowPower)
	assert.Len(t, a.written, 2)
	assert.Equal(t, wantReg, a.written[0])
	assert.Equal(t, wantVal, a.written[1])
}

func TestADXL345SetRate(t *testing.T) {
	// sequence to set rate:
	// * set value in data rate structure
	// * write the data rate register (0x2C)
	d, a := initTestADXL345WithStubbedAdaptor()
	_ = d.Start()
	a.written = []byte{} // reset writes of former test
	const (
		setVal  = ADXL345RateConfig(0x0F) // 3.2kHz
		wantReg = uint8(0x2C)
		wantVal = uint8(0x1F) // also low power bit
	)
	// act
	err := d.SetRate(setVal)
	// assert
	require.NoError(t, err)
	assert.Equal(t, setVal, d.bwRate.rate)
	assert.Len(t, a.written, 2)
	assert.Equal(t, wantReg, a.written[0])
	assert.Equal(t, wantVal, a.written[1])
}

func TestADXL345SetRange(t *testing.T) {
	// sequence to set range:
	// * set value in data format structure
	// * write the data format register (0x31)
	d, a := initTestADXL345WithStubbedAdaptor()
	_ = d.Start()
	a.written = []byte{} // reset writes of former test
	const (
		setVal  = ADXL345FsRangeConfig(0x03) // +/- 16 g
		wantReg = uint8(0x31)
		wantVal = uint8(0x03)
	)
	// act
	err := d.SetRange(setVal)
	// assert
	require.NoError(t, err)
	assert.Equal(t, setVal, d.dataFormat.fullScaleRange)
	assert.Len(t, a.written, 2)
	assert.Equal(t, wantReg, a.written[0])
	assert.Equal(t, wantVal, a.written[1])
}

func TestADXL345RawXYZ(t *testing.T) {
	// sequence to read:
	// * prepare read, see test of initialize()
	// * read data output registers (0x32, 3 x 16 bit, LSByte first)
	// * apply two's complement converter
	//
	// arrange
	tests := map[string]struct {
		inputX []uint8
		inputY []uint8
		inputZ []uint8
		wantX  int16
		wantY  int16
		wantZ  int16
	}{
		"+FS_0_-FS": {
			inputX: []uint8{0xFF, 0x07},
			inputY: []uint8{0x00, 0x00},
			inputZ: []uint8{0x00, 0xF8},
			wantX:  (1<<11 - 1),
			wantY:  0,
			wantZ:  -(1 << 11),
		},
		"-4096_-1_+1": {
			inputX: []uint8{0x00, 0xF0},
			inputY: []uint8{0xFF, 0xFF},
			inputZ: []uint8{0x01, 0x00},
			wantX:  -4096,
			wantY:  -1,
			wantZ:  1,
		},
	}
	d, a := initTestADXL345WithStubbedAdaptor()
	_ = d.Start()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			a.written = []byte{} // reset writes of former test and start
			// arrange reads
			returnRead := append(append(tc.inputX, tc.inputY...), tc.inputZ...)
			numCallsRead := 0
			a.i2cReadImpl = func(b []byte) (int, error) {
				numCallsRead++
				copy(b, returnRead)
				return len(b), nil
			}
			// act
			gotX, gotY, gotZ, err := d.RawXYZ()
			// assert
			require.NoError(t, err)
			assert.Equal(t, tc.wantX, gotX)
			assert.Equal(t, tc.wantY, gotY)
			assert.Equal(t, tc.wantZ, gotZ)
			assert.Equal(t, 1, numCallsRead)
			assert.Len(t, a.written, 1)
			assert.Equal(t, uint8(0x32), a.written[0])
		})
	}
}

func TestADXL345RawXYZError(t *testing.T) {
	// arrange
	d, a := initTestADXL345WithStubbedAdaptor()
	_ = d.Start()
	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}
	// act
	x, y, z, err := d.RawXYZ()
	// assert
	require.ErrorContains(t, err, "read error")
	assert.Equal(t, int16(0), x)
	assert.Equal(t, int16(0), y)
	assert.Equal(t, int16(0), z)
}

func TestADXL345XYZ(t *testing.T) {
	// arrange
	tests := map[string]struct {
		inputX []uint8
		inputY []uint8
		inputZ []uint8
		wantX  float64
		wantY  float64
		wantZ  float64
	}{
		"null_value": {
			inputX: []uint8{0, 0},
			inputY: []uint8{0, 0},
			inputZ: []uint8{0, 0},
			wantX:  0,
			wantY:  0,
			wantZ:  0,
		},
		"some_value": {
			inputX: []uint8{218, 0},
			inputY: []uint8{251, 255},
			inputZ: []uint8{100, 0},
			wantX:  0.8515625,
			wantY:  -0.01953125,
			wantZ:  0.390625,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestADXL345WithStubbedAdaptor()
			_ = d.Start()
			a.written = []byte{} // reset writes of former test and start
			// arrange reads
			returnRead := append(append(tc.inputX, tc.inputY...), tc.inputZ...)
			numCallsRead := 0
			a.i2cReadImpl = func(b []byte) (int, error) {
				numCallsRead++
				copy(b, returnRead)
				return len(b), nil
			}
			// act
			x, y, z, _ := d.XYZ()
			// assert
			assert.InDelta(t, tc.wantX, x, 0.0)
			assert.InDelta(t, tc.wantY, y, 0.0)
			assert.InDelta(t, tc.wantZ, z, 0.0)
		})
	}
}

func TestADXL345XYZError(t *testing.T) {
	// arrange
	d, a := initTestADXL345WithStubbedAdaptor()
	_ = d.Start()
	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}
	// act
	x, y, z, err := d.XYZ()
	// assert
	require.ErrorContains(t, err, "read error")
	assert.InDelta(t, 0.0, x, 0.0)
	assert.InDelta(t, 0.0, y, 0.0)
	assert.InDelta(t, 0.0, z, 0.0)
}

func TestADXL345_initialize(t *testing.T) {
	// sequence to prepare read in initialize():
	// * prepare rate register content (data output rate, low power mode)
	// * prepare power control register content (wake up, sleep, measure, auto sleep, link)
	// * prepare data format register (fullScaleRange, justify, fullRes, intInvert, spi, selfTest)
	// * write 3 registers
	// arrange
	d, a := initTestADXL345WithStubbedAdaptor()
	a.written = []byte{} // reset writes of former test
	const (
		wantRateReg      = uint8(0x2C)
		wantRateRegVal   = uint8(0x1A) // 100HZ and low power
		wantPwrReg       = uint8(0x2D)
		wantPwrRegVal    = uint8(0x08) // measurement on
		wantFormatReg    = uint8(0x31)
		wantFormatRegVal = uint8(0x00) // FS to +/-2 g
	)
	// act, assert - initialize() must be called on Start()
	err := d.Start()
	// assert
	require.NoError(t, err)
	assert.Len(t, a.written, 6)
	assert.Equal(t, wantRateReg, a.written[0])
	assert.Equal(t, wantRateRegVal, a.written[1])
	assert.Equal(t, wantPwrReg, a.written[2])
	assert.Equal(t, wantPwrRegVal, a.written[3])
	assert.Equal(t, wantFormatReg, a.written[4])
	assert.Equal(t, wantFormatRegVal, a.written[5])
}

func TestADXL345_shutdown(t *testing.T) {
	// sequence to prepare read in shutdown():
	// * reset the measurement bit in structure
	// * write the power control register (0x2D)
	d, a := initTestADXL345WithStubbedAdaptor()
	_ = d.Start()
	a.written = []byte{} // reset writes of former test
	const (
		wantReg = uint8(0x2D)
		wantVal = uint8(0x00)
	)
	// act, assert - shutdown() must be called on Halt()
	err := d.Halt()
	// assert
	require.NoError(t, err)
	assert.Len(t, a.written, 2)
	assert.Equal(t, wantReg, a.written[0])
	assert.Equal(t, wantVal, a.written[1])
}
