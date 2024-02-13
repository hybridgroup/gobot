package i2c

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*HMC5883LDriver)(nil)

func initTestHMC5883LWithStubbedAdaptor() (*HMC5883LDriver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	return NewHMC5883LDriver(a), a
}

func TestNewHMC5883LDriver(t *testing.T) {
	var di interface{} = NewHMC5883LDriver(newI2cTestAdaptor())
	d, ok := di.(*HMC5883LDriver)
	if !ok {
		require.Fail(t, "NewHMC5883LDriver() should have returned a *HMC5883LDriver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.name, "HMC5883L"))
	assert.Equal(t, 0x1E, d.defaultAddress)
	assert.Equal(t, uint8(8), d.samplesAvg)
	assert.Equal(t, uint32(15000), d.outputRate)
	assert.Equal(t, int8(0), d.applyBias)
	assert.Equal(t, 0, d.measurementMode)
	assert.InDelta(t, 390.0, d.gain, 0.0)
}

func TestHMC5883LOptions(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewHMC5883LDriver(newI2cTestAdaptor(), WithBus(2), WithHMC5883LSamplesAveraged(4))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
	assert.Equal(t, uint8(4), d.samplesAvg)
}

func TestHMC5883LWithHMC5883LDataOutputRate(t *testing.T) {
	d := NewHMC5883LDriver(newI2cTestAdaptor())
	WithHMC5883LDataOutputRate(7500)(d)
	assert.Equal(t, uint32(7500), d.outputRate)
}

func TestHMC5883LWithHMC5883LApplyBias(t *testing.T) {
	d := NewHMC5883LDriver(newI2cTestAdaptor())
	WithHMC5883LApplyBias(-1)(d)
	assert.Equal(t, int8(-1), d.applyBias)
}

func TestHMC5883LWithHMC5883LGain(t *testing.T) {
	d := NewHMC5883LDriver(newI2cTestAdaptor())
	WithHMC5883LGain(230)(d)
	assert.InDelta(t, 230.0, d.gain, 0.0)
}

func TestHMC5883LRead(t *testing.T) {
	// arrange
	tests := map[string]struct {
		inputX []uint8
		inputY []uint8
		inputZ []uint8
		gain   float64
		wantX  float64
		wantY  float64
		wantZ  float64
	}{
		"+FS_0_-FS_resolution_0.73mG": {
			inputX: []uint8{0x07, 0xFF},
			inputY: []uint8{0x00, 0x00},
			inputZ: []uint8{0xF8, 0x00},
			gain:   1370,
			wantX:  2047.0 / 1370,
			wantY:  0,
			wantZ:  -2048.0 / 1370,
		},
		"+1_-4096_-1_resolution_0.73mG": {
			inputX: []uint8{0x00, 0x01},
			inputY: []uint8{0xF0, 0x00},
			inputZ: []uint8{0xFF, 0xFF},
			gain:   1370,
			wantX:  1.0 / 1370,
			wantY:  -4096.0 / 1370,
			wantZ:  -1.0 / 1370,
		},
		"+FS_0_-FS_resolution_4.35mG": {
			inputX: []uint8{0x07, 0xFF},
			inputY: []uint8{0x00, 0x00},
			inputZ: []uint8{0xF8, 0x00},
			gain:   230,
			wantX:  2047.0 / 230,
			wantY:  0,
			wantZ:  -2048.0 / 230,
		},
		"-1_+1_-4096_resolution_4.35mG": {
			inputX: []uint8{0xFF, 0xFF},
			inputY: []uint8{0x00, 0x01},
			inputZ: []uint8{0xF0, 0x00},
			gain:   230,
			wantX:  -1.0 / 230,
			wantY:  1.0 / 230,
			wantZ:  -4096.0 / 230,
		},
	}
	a := newI2cTestAdaptor()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			d := NewHMC5883LDriver(a, WithHMC5883LGain(int(tc.gain)))
			_ = d.Start()
			// arrange reads
			returnRead := append(append(tc.inputX, tc.inputZ...), tc.inputY...)
			a.i2cReadImpl = func(b []byte) (int, error) {
				copy(b, returnRead)
				return len(b), nil
			}
			// act
			gotX, gotY, gotZ, err := d.Read()
			// assert
			require.NoError(t, err)
			assert.InDelta(t, tc.wantX, gotX, 0.0)
			assert.InDelta(t, tc.wantY, gotY, 0.0)
			assert.InDelta(t, tc.wantZ, gotZ, 0.0)
		})
	}
}

func TestHMC5883L_readRawData(t *testing.T) {
	// sequence to read:
	// * prepare read, see test of Start()
	// * read data output registers (3 x 16 bit, MSByte first)
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
			inputX: []uint8{0x07, 0xFF},
			inputY: []uint8{0x00, 0x00},
			inputZ: []uint8{0xF8, 0x00},
			wantX:  (1<<11 - 1),
			wantY:  0,
			wantZ:  -(1 << 11),
		},
		"-4096_-1_+1": {
			inputX: []uint8{0xF0, 0x00},
			inputY: []uint8{0xFF, 0xFF},
			inputZ: []uint8{0x00, 0x01},
			wantX:  -4096,
			wantY:  -1,
			wantZ:  1,
		},
	}
	d, a := initTestHMC5883LWithStubbedAdaptor()
	_ = d.Start()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			a.written = []byte{} // reset writes of former test and start
			// arrange reads
			returnRead := append(append(tc.inputX, tc.inputZ...), tc.inputY...)
			numCallsRead := 0
			a.i2cReadImpl = func(b []byte) (int, error) {
				numCallsRead++
				copy(b, returnRead)
				return len(b), nil
			}
			// act
			gotX, gotY, gotZ, err := d.readRawData()
			// assert
			require.NoError(t, err)
			assert.Equal(t, tc.wantX, gotX)
			assert.Equal(t, tc.wantY, gotY)
			assert.Equal(t, tc.wantZ, gotZ)
			assert.Equal(t, 1, numCallsRead)
			assert.Len(t, a.written, 1)
			assert.Equal(t, uint8(hmc5883lAxisX), a.written[0])
		})
	}
}

func TestHMC5883L_initialize(t *testing.T) {
	// sequence to prepare read in Start():
	// * prepare config register A content (samples averaged, data output rate, measurement mode)
	// * prepare config register B content (gain)
	// * prepare mode register (continuous/single/idle)
	// * write registers A, B, mode
	// arrange
	d, a := initTestHMC5883LWithStubbedAdaptor()
	a.written = []byte{} // reset writes of former test
	wantRegA := uint8(0x70)
	wantRegB := uint8(0xA0)
	wantRegM := uint8(0x00)
	// act, assert - initialize() must be called on Start()
	err := d.Start()
	// assert
	require.NoError(t, err)
	assert.Len(t, a.written, 6)
	assert.Equal(t, uint8(hmc5883lRegA), a.written[0])
	assert.Equal(t, wantRegA, a.written[1])
	assert.Equal(t, uint8(hmc5883lRegB), a.written[2])
	assert.Equal(t, wantRegB, a.written[3])
	assert.Equal(t, uint8(hmc5883lRegMode), a.written[4])
	assert.Equal(t, wantRegM, a.written[5])
}
