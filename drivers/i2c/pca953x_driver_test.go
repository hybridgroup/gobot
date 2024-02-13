package i2c

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*PCA953xDriver)(nil)

func initPCA953xTestDriverWithStubbedAdaptor() (*PCA953xDriver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewPCA953xDriver(a)
	_ = d.Start()
	return d, a
}

func TestNewPCA953xDriver(t *testing.T) {
	// arrange, act
	var di interface{} = NewPCA953xDriver(newI2cTestAdaptor())
	// assert
	d, ok := di.(*PCA953xDriver)
	if !ok {
		require.Fail(t, "NewPCA953xDriver() should have returned a *PCA953xDriver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "PCA953x"))
	assert.Equal(t, 0x63, d.defaultAddress)
}

func TestPCA953xWriteGPIO(t *testing.T) {
	// sequence to write:
	// * choose LED select register according to the given GPIO index (0x05 for 0..3, 0x06 for 4..7)
	// * read current state of LED select register (write reg, read val)
	// * modify 2 bits according to given index of GPIO
	// * write the new state to the LED select register (write reg, write val)
	tests := map[string]struct {
		idx         uint8
		ls0State    uint8
		ls1State    uint8
		val         uint8
		wantWritten []uint8
		wantErr     error
	}{
		"out_0_0": {
			idx:         0,
			ls0State:    0xFE,
			ls1State:    0xAF,
			val:         0,
			wantWritten: []byte{0x05, 0x05, 0xFD}, // set lowest bits to "01" for ls0
		},
		"out_0_1": {
			idx:         0,
			ls0State:    0xFF,
			ls1State:    0xAF,
			val:         1,
			wantWritten: []byte{0x05, 0x05, 0xFC}, // set lowest bits to "00" for ls0
		},
		"out_5_0": {
			idx:         5,
			ls0State:    0xAF,
			ls1State:    0xFB,
			val:         0,
			wantWritten: []byte{0x06, 0x06, 0xF7}, // set bit 2,3 to "01" for ls1
		},
		"out_5_1": {
			idx:         5,
			ls0State:    0xAF,
			ls1State:    0xFF,
			val:         1,
			wantWritten: []byte{0x06, 0x06, 0xF3}, // set bit 2,3 to "00" for ls1
		},
		"read_error": {
			idx:         3,
			wantWritten: []byte{0x05},
			wantErr:     fmt.Errorf("a read error"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initPCA953xTestDriverWithStubbedAdaptor()
			a.written = []byte{} // reset writes of Start() and former test
			a.i2cReadImpl = func(b []byte) (int, error) {
				if a.written[0] == 0x05 {
					b[0] = tc.ls0State
				}
				if a.written[0] == 0x06 {
					b[0] = tc.ls1State
				}
				return 1, tc.wantErr
			}
			// act
			err := d.WriteGPIO(tc.idx, tc.val)
			// assert
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantWritten, a.written)
		})
	}
}

func TestPCA953xReadGPIO(t *testing.T) {
	// sequence to read:
	// * read current state of INPUT register (write reg 0x00, read val)
	// * convert bit position to output value
	tests := map[string]struct {
		idx     uint8
		want    uint8
		wantErr error
	}{
		"in_0_0": {
			idx:  0,
			want: 0,
		},
		"in_0_1": {
			idx:  0,
			want: 1,
		},
		"in_2_0": {
			idx:  2,
			want: 0,
		},
		"in_2_1": {
			idx:  2,
			want: 1,
		},
		"in_7_0": {
			idx:  7,
			want: 0,
		},
		"in_7_1": {
			idx:  7,
			want: 1,
		},
		"read_error": {
			idx:     2,
			want:    0,
			wantErr: fmt.Errorf("a read error"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			const wantReg = uint8(0x00) // input register
			d, a := initPCA953xTestDriverWithStubbedAdaptor()
			a.written = []byte{} // reset writes of Start() and former test
			bits := tc.want << tc.idx
			a.i2cReadImpl = func(b []byte) (int, error) {
				b[0] = bits
				return 1, tc.wantErr
			}
			// act
			got, err := d.ReadGPIO(tc.idx)
			// assert
			assert.Equal(t, tc.wantErr, err)
			assert.Len(t, a.written, 1)
			assert.Equal(t, wantReg, a.written[0])
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestPCA953xWritePeriod(t *testing.T) {
	// sequence to write:
	// * calculate PSC value (0..255) from given value in seconds, valid values are 0.00658 ... 1.68 [s]
	// * choose PSC0 (0x01) or PSC1 (0x03) frequency prescaler register by the given index
	// * write the value to the register (write reg, write val)
	tests := map[string]struct {
		idx         uint8
		val         float32
		wantWritten []uint8
	}{
		"write_ok_psc0": {
			idx:         0,
			val:         1,
			wantWritten: []byte{0x01, 151},
		},
		"write_ok_psc1": {
			idx:         2,
			val:         0.5,
			wantWritten: []byte{0x03, 75},
		},
		"write_limited_noerror": {
			idx:         0,
			val:         2,
			wantWritten: []byte{0x01, 255},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initPCA953xTestDriverWithStubbedAdaptor()
			a.written = []byte{} // reset writes of Start() and former test
			// act
			err := d.WritePeriod(tc.idx, tc.val)
			// assert
			require.NoError(t, err)
			assert.Equal(t, tc.wantWritten, a.written)
		})
	}
}

func TestPCA953xReadPeriod(t *testing.T) {
	// sequence to write:
	// * choose PSC0 (0x01) or PSC1 (0x03) frequency prescaler register by the given index
	// * read the value from the register (write reg, write val)
	// * calculate value in seconds from PSC value
	tests := map[string]struct {
		idx         uint8
		val         uint8
		want        float32
		wantWritten []uint8
		wantErr     error
	}{
		"read_ok_psc0": {
			idx:         0,
			val:         151,
			want:        1,
			wantWritten: []byte{0x01},
		},
		"read_ok_psc1": {
			idx:         1,
			val:         75,
			want:        0.5,
			wantWritten: []byte{0x03},
		},
		"read_error": {
			idx:         5,
			val:         75,
			want:        -1,
			wantWritten: []byte{0x03},
			wantErr:     fmt.Errorf("read psc error"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initPCA953xTestDriverWithStubbedAdaptor()
			a.written = []byte{} // reset writes of Start() and former test
			a.i2cReadImpl = func(b []byte) (int, error) {
				b[0] = tc.val
				return 1, tc.wantErr
			}
			// act
			got, err := d.ReadPeriod(tc.idx)
			// assert
			assert.Equal(t, tc.wantErr, err)
			assert.InDelta(t, tc.want, got, 0.0)
			assert.Equal(t, tc.wantWritten, a.written)
		})
	}
}

func TestPCA953xWriteFrequency(t *testing.T) {
	// sequence to write:
	// * calculate PSC value (0..255) from given value in Hz, valid values are 0.6 ... 152 [Hz]
	// * choose PSC0 (0x01) or PSC1 (0x03) frequency prescaler register by the given index
	// * write the value to the register (write reg, write val)
	tests := map[string]struct {
		idx         uint8
		val         float32
		wantWritten []uint8
	}{
		"write_ok_psc0": {
			idx:         0,
			val:         1,
			wantWritten: []byte{0x01, 151},
		},
		"write_ok_psc1": {
			idx:         5,
			val:         2,
			wantWritten: []byte{0x03, 75},
		},
		"write_limited_noerror": {
			idx:         0,
			val:         153,
			wantWritten: []byte{0x01, 0},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initPCA953xTestDriverWithStubbedAdaptor()
			a.written = []byte{} // reset writes of Start() and former test
			// act
			err := d.WriteFrequency(tc.idx, tc.val)
			// assert
			require.NoError(t, err)
			assert.Equal(t, tc.wantWritten, a.written)
		})
	}
}

func TestPCA953xReadFrequency(t *testing.T) {
	// sequence to write:
	// * choose PSC0 (0x01) or PSC1 (0x03) frequency prescaler register by the given index
	// * read the value from the register (write reg, write val)
	// * calculate value in Hz from PSC value
	tests := map[string]struct {
		idx         uint8
		val         uint8
		want        float32
		wantWritten []uint8
		wantErr     error
	}{
		"read_ok_psc0": {
			idx:         0,
			val:         75,
			want:        2,
			wantWritten: []byte{0x01},
		},
		"read_ok_psc1": {
			idx:         1,
			val:         151,
			want:        1,
			wantWritten: []byte{0x03},
		},
		"read_error": {
			idx:         3,
			val:         75,
			want:        -1,
			wantWritten: []byte{0x03},
			wantErr:     fmt.Errorf("read psc error"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initPCA953xTestDriverWithStubbedAdaptor()
			a.written = []byte{} // reset writes of Start() and former test
			a.i2cReadImpl = func(b []byte) (int, error) {
				b[0] = tc.val
				return 1, tc.wantErr
			}
			// act
			got, err := d.ReadFrequency(tc.idx)
			// assert
			assert.Equal(t, tc.wantErr, err)
			assert.InDelta(t, tc.want, got, 0.0)
			assert.Equal(t, tc.wantWritten, a.written)
		})
	}
}

func TestPCA953xWriteDutyCyclePercent(t *testing.T) {
	// sequence to write:
	// * calculate PWM value (0..255) from given value in percent, valid values are 0 ... 100 [%]
	// * choose PWM0 (0x02) or PWM1 (0x04) pwm register by the given index
	// * write the value to the register (write reg, write val)
	tests := map[string]struct {
		idx         uint8
		val         float32
		wantWritten []uint8
	}{
		"write_ok_pwm0": {
			idx:         0,
			val:         10,
			wantWritten: []byte{0x02, 26},
		},
		"write_ok_pwm1": {
			idx:         5,
			val:         50,
			wantWritten: []byte{0x04, 128},
		},
		"write_limited_noerror": {
			idx:         1,
			val:         101,
			wantWritten: []byte{0x04, 255},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initPCA953xTestDriverWithStubbedAdaptor()
			a.written = []byte{} // reset writes of Start() and former test
			// act
			err := d.WriteDutyCyclePercent(tc.idx, tc.val)
			// assert
			require.NoError(t, err)
			assert.Equal(t, tc.wantWritten, a.written)
		})
	}
}

func TestPCA953xReadDutyCyclePercent(t *testing.T) {
	// sequence to write:
	// * choose PWM0 (0x02) or PWM1 (0x04) pwm register by the given index
	// * read the value from the register (write reg, write val)
	// * calculate value percent from PWM value
	tests := map[string]struct {
		idx         uint8
		val         uint8
		want        float32
		wantWritten []uint8
		wantErr     error
	}{
		"read_ok_psc0": {
			idx:         0,
			val:         128,
			want:        50.19608,
			wantWritten: []byte{0x02},
		},
		"read_ok_psc1": {
			idx:         1,
			val:         26,
			want:        10.196078,
			wantWritten: []byte{0x04},
		},
		"read_error": {
			idx:         0,
			val:         75,
			want:        -1,
			wantWritten: []byte{0x02},
			wantErr:     fmt.Errorf("read psc error"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initPCA953xTestDriverWithStubbedAdaptor()
			a.written = []byte{} // reset writes of Start() and former test
			a.i2cReadImpl = func(b []byte) (int, error) {
				b[0] = tc.val
				return 1, tc.wantErr
			}
			// act
			got, err := d.ReadDutyCyclePercent(tc.idx)
			// assert
			assert.Equal(t, tc.wantErr, err)
			assert.InDelta(t, tc.want, got, 0.0)
			assert.Equal(t, tc.wantWritten, a.written)
		})
	}
}

func TestPCA953x_readRegister(t *testing.T) {
	// arrange
	const (
		wantRegAddress    = pca953xRegister(0x03)
		wantReadByteCount = 1
		wantRegVal        = uint8(0x04)
	)
	readByteCount := 0
	d, a := initPCA953xTestDriverWithStubbedAdaptor()
	// prepare all writes
	numCallsWrite := 0
	a.i2cWriteImpl = func([]byte) (int, error) {
		numCallsWrite++
		return 0, nil
	}
	// prepare all reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		readByteCount = len(b)
		b[0] = wantRegVal
		return readByteCount, nil
	}
	// act
	val, err := d.readRegister(wantRegAddress)
	// assert
	require.NoError(t, err)
	assert.Equal(t, 1, numCallsRead)
	assert.Equal(t, 1, numCallsWrite)
	assert.Equal(t, wantRegVal, val)
	assert.Equal(t, wantReadByteCount, readByteCount)
	assert.Len(t, a.written, 1)
	assert.Equal(t, uint8(wantRegAddress), a.written[0])
}

func TestPCA953x_writeRegister(t *testing.T) {
	// arrange
	const (
		wantRegAddress = pca953xRegister(0x03)
		wantRegVal     = uint8(0x97)
		wantByteCount  = 2
	)
	d, a := initPCA953xTestDriverWithStubbedAdaptor()
	// prepare all writes
	numCallsWrite := 0
	a.i2cWriteImpl = func(b []byte) (int, error) {
		numCallsWrite++
		return 0, nil
	}
	// act
	err := d.writeRegister(wantRegAddress, wantRegVal)
	// assert
	require.NoError(t, err)
	assert.Equal(t, 1, numCallsWrite)
	assert.Equal(t, 1, numCallsWrite)
	assert.Len(t, a.written, wantByteCount)
	assert.Equal(t, uint8(wantRegAddress), a.written[0])
	assert.Equal(t, wantRegVal, a.written[1])
}

func TestPCA953x_pca953xCalcPsc(t *testing.T) {
	// arrange
	tests := map[string]struct {
		period  float32
		want    uint8
		wantErr error
	}{
		"error_to_small": {period: 0.0065, want: 0, wantErr: errToSmallPeriod},
		"minimum":        {period: 0.0066, want: 0, wantErr: nil},
		"one":            {period: 1, want: 151, wantErr: nil},
		"maximum":        {period: 1.684, want: 255, wantErr: nil},
		"error_to_big5":  {period: 1.685, want: 255, wantErr: errToBigPeriod},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// act
			val, err := pca953xCalcPsc(tc.period)
			// assert
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, val)
		})
	}
}

func TestPCA953x_pca953xCalcPeriod(t *testing.T) {
	// arrange
	tests := map[string]struct {
		psc  uint8
		want float32
	}{
		"minimum":  {psc: 0, want: 0.0066},
		"one":      {psc: 1, want: 0.0132},
		"one_want": {psc: 151, want: 1},
		"maximum":  {psc: 255, want: 1.6842},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// act
			val := pca953xCalcPeriod(tc.psc)
			// assert
			assert.InDelta(t, tc.want, float32(math.Round(float64(val)*10000)/10000), 0.0)
		})
	}
}

func TestPCA953x_pca953xCalcPwm(t *testing.T) {
	// arrange
	tests := map[string]struct {
		percent float32
		want    uint8
		wantErr error
	}{
		"error_to_small": {percent: -0.1, want: 0, wantErr: errToSmallDutyCycle},
		"zero":           {percent: 0, want: 0, wantErr: nil},
		"below_medium":   {percent: 49.9, want: 127, wantErr: nil},
		"medium":         {percent: 50, want: 128, wantErr: nil},
		"maximum":        {percent: 100, want: 255, wantErr: nil},
		"error_to_big":   {percent: 100.1, want: 255, wantErr: errToBigDutyCycle},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// act
			val, err := pca953xCalcPwm(tc.percent)
			// assert
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, val)
		})
	}
}

func TestPCA953x_pca953xCalcDutyCyclePercent(t *testing.T) {
	// arrange
	tests := map[string]struct {
		pwm  uint8
		want float32
	}{
		"minimum":      {pwm: 0, want: 0},
		"below_medium": {pwm: 127, want: 49.8},
		"medium":       {pwm: 128, want: 50.2},
		"maximum":      {pwm: 255, want: 100},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// act
			val := pca953xCalcDutyCyclePercent(tc.pwm)
			// assert
			assert.InDelta(t, tc.want, float32(math.Round(float64(val)*10)/10), 0.0)
		})
	}
}
