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
var _ gobot.Driver = (*CCS811Driver)(nil)

func initTestCCS811WithStubbedAdaptor() (*CCS811Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	return NewCCS811Driver(a), a
}

func TestNewCCS811Driver(t *testing.T) {
	var di interface{} = NewCCS811Driver(newI2cTestAdaptor())
	d, ok := di.(*CCS811Driver)
	if !ok {
		require.Fail(t, "NewCCS811Driver() should have returned a *CCS811Driver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "CCS811"))
	assert.Equal(t, 0x5A, d.defaultAddress)
	assert.NotNil(t, d.measMode)
	assert.Equal(t, uint32(100000), d.ntcResistanceValue)
}

func TestCCS811Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewCCS811Driver(newI2cTestAdaptor(), WithBus(2), WithAddress(0xFF), WithCCS811NTCResistance(0xFF))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
	assert.Equal(t, 0xFF, d.GetAddressOrDefault(0x5a))
	assert.Equal(t, uint32(0xFF), d.ntcResistanceValue)
}

func TestCCS811WithCCS811MeasMode(t *testing.T) {
	d := NewCCS811Driver(newI2cTestAdaptor(), WithCCS811MeasMode(CCS811DriveMode10Sec))
	assert.Equal(t, CCS811DriveMode10Sec, d.measMode.driveMode)
}

func TestCCS811GetGasData(t *testing.T) {
	tests := map[string]struct {
		readReturn func([]byte) (int, error)
		eco2       uint16
		tvoc       uint16
		err        error
	}{
		"ideal values taken from the bus": {
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{1, 156, 0, 86})
				return 4, nil
			},
			eco2: 412,
			tvoc: 86,
			err:  nil,
		},
		"max values possible taken from the bus": {
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{255, 255, 255, 255})
				return 4, nil
			},
			eco2: 65535,
			tvoc: 65535,
			err:  nil,
		},
		"error when the i2c operation fails": {
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{255, 255, 255, 255})
				return 4, errors.New("Error")
			},
			eco2: 0,
			tvoc: 0,
			err:  errors.New("Error"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestCCS811WithStubbedAdaptor()
			// Create stub function as it is needed by read submethod in driver code
			a.i2cWriteImpl = func([]byte) (int, error) { return 0, nil }
			_ = d.Start()
			a.i2cReadImpl = tc.readReturn
			// act
			eco2, tvoc, err := d.GetGasData()
			// assert
			assert.Equal(t, tc.eco2, eco2)
			assert.Equal(t, tc.tvoc, tvoc)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestCCS811GetTemperature(t *testing.T) {
	tests := map[string]struct {
		readReturn func([]byte) (int, error)
		temp       float32
		err        error
	}{
		"ideal values taken from the bus": {
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{10, 197, 0, 248})
				return 4, nil
			},
			temp: 27.811005,
			err:  nil,
		},
		"without bus values overflowing": {
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{129, 197, 10, 248})
				return 4, nil
			},
			temp: 29.48822,
			err:  nil,
		},
		"negative temperature": {
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{255, 255, 255, 255})
				return 4, nil
			},
			temp: -25.334152,
			err:  nil,
		},
		"error if the i2c bus errors": {
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{129, 197, 0, 248})
				return 4, errors.New("Error")
			},
			temp: 0,
			err:  errors.New("Error"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestCCS811WithStubbedAdaptor()
			// Create stub function as it is needed by read submethod in driver code
			a.i2cWriteImpl = func([]byte) (int, error) { return 0, nil }
			_ = d.Start()
			a.i2cReadImpl = tc.readReturn
			// act
			temp, err := d.GetTemperature()
			// assert
			assert.InDelta(t, tc.temp, temp, 0.0)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestCCS811HasData(t *testing.T) {
	tests := map[string]struct {
		readReturn func([]byte) (int, error)
		result     bool
		err        error
	}{
		"true for HasError=0 and DataRdy=1": {
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{0x08})
				return 1, nil
			},
			result: true,
			err:    nil,
		},
		"false for HasError=1 and DataRdy=1": {
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{0x09})
				return 1, nil
			},
			result: false,
			err:    nil,
		},
		"false for HasError=1 and DataRdy=0": {
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{0x01})
				return 1, nil
			},
			result: false,
			err:    nil,
		},
		"false for HasError=0 and DataRdy=0": {
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{0x00})
				return 1, nil
			},
			result: false,
			err:    nil,
		},
		"error when the i2c read operation fails": {
			readReturn: func(b []byte) (int, error) {
				copy(b, []byte{0x00})
				return 1, errors.New("Error")
			},
			result: false,
			err:    errors.New("Error"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestCCS811WithStubbedAdaptor()
			// Create stub function as it is needed by read submethod in driver code
			a.i2cWriteImpl = func([]byte) (int, error) { return 0, nil }
			_ = d.Start()
			a.i2cReadImpl = tc.readReturn
			// act
			result, err := d.HasData()
			// assert
			assert.Equal(t, tc.result, result)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestCCS811_initialize(t *testing.T) {
	// sequence for initialization the device on Start()
	// * write hardware ID register (0x20)
	// * read the ID
	// * prepare software reset register content: a sequence of four bytes must
	//   be written to this register in a single I²C sequence: 0x11, 0xE5, 0x72, 0x8A
	// * write software reset register content (0xFF)
	// * write application start register (0xF4)
	// * prepare measurement mode register content
	//  * INT_THRESH = 0 (normal mode)
	//  * INT_DATARDY = 0 (disable interrupt mode)
	//  * DRIVE_MODE = 0x01 (constant power, value every 1 sec)
	// * write measure mode register content (0x01)
	//
	// arrange
	d, a := initTestCCS811WithStubbedAdaptor()
	a.written = []byte{} // reset writes of former test
	const (
		wantChipIDReg    = uint8(0x20)
		wantChipIDRegVal = uint8(0x20)
		wantResetReg     = uint8(0xFF)
		wantAppStartReg  = uint8(0xF4)
		wantMeasReg      = uint8(0x01)
		wantMeasRegVal   = uint8(0x10)
	)
	wantResetRegSequence := []byte{0x11, 0xE5, 0x72, 0x8A}
	// arrange reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		// chip ID
		b[0] = 0x81
		return len(b), nil
	}
	// arrange, act - initialize() must be called on Start()
	err := d.Start()
	// assert
	require.NoError(t, err)
	assert.Equal(t, 1, numCallsRead)
	assert.Len(t, a.written, 9)
	assert.Equal(t, wantChipIDReg, a.written[0])
	assert.Equal(t, wantResetReg, a.written[1])
	assert.Equal(t, wantResetRegSequence, a.written[2:6])
	assert.Equal(t, wantAppStartReg, a.written[6])
	assert.Equal(t, wantMeasReg, a.written[7])
	assert.Equal(t, wantMeasRegVal, a.written[8])
}
