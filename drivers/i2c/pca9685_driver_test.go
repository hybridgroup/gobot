package i2c

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*PCA9685Driver)(nil)

// and also the PwmWriter and ServoWriter interfaces
var (
	_ gpio.PwmWriter   = (*PCA9685Driver)(nil)
	_ gpio.ServoWriter = (*PCA9685Driver)(nil)
)

func initTestPCA9685WithStubbedAdaptor() (*PCA9685Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewPCA9685Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewPCA9685Driver(t *testing.T) {
	// arrange & act
	d := NewPCA9685Driver(newI2cTestAdaptor())
	// assert
	assert.IsType(t, &PCA9685Driver{}, d)
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "PCA9685"))
	assert.Equal(t, 0x40, d.defaultAddress)
}

func TestPCA9685Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	// arrange & act
	d := NewPCA9685Driver(newI2cTestAdaptor(), WithBus(2))
	// assert
	assert.Equal(t, 2, d.GetBusOrDefault(1))
}

func TestPCA9685Start(t *testing.T) {
	// arrange
	a := newI2cTestAdaptor()
	d := NewPCA9685Driver(a)
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x01})
		return 1, nil
	}
	// act & assert
	require.NoError(t, d.Start())
}

func TestPCA9685StartError(t *testing.T) {
	// arrange
	a := newI2cTestAdaptor()
	d := NewPCA9685Driver(a)
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	// act & assert
	require.ErrorContains(t, d.Start(), "write error")
}

func TestPCA9685Halt(t *testing.T) {
	// arrange
	d, a := initTestPCA9685WithStubbedAdaptor()
	a.written = []byte{} // reset writes of former test
	// act
	err := d.Halt()
	// assert
	require.NoError(t, err)
	require.NoError(t, err)
	assert.Len(t, a.written, 2)
	assert.Equal(t, []byte{0xFD, 0x10}, a.written)
}

func TestPCA9685HaltError(t *testing.T) {
	// arrange
	d, a := initTestPCA9685WithStubbedAdaptor()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	// act & assert
	require.ErrorContains(t, d.Halt(), "write error")
}

func TestPCA9685SetPWM(t *testing.T) {
	// sequence to set PWM for PCA9685:
	// * set LEDn ON-time register (n=0: 0x06, 0x07, n=1: 0x0A, 0x0B ... n=14: 0x3E, 0x3F, n=15: 0x42, 0x43)
	// * set LEDn OFF-time register (n=0: 0x08, 0x09, n=1: 0x0C, 0x0D ... n=14: 0x40, 0x41, n=15: 0x44, 0x45)
	tests := map[string]struct {
		pin                     int
		onCounts                uint16
		offCounts               uint16
		wantLedOnTimeOffTimeSet []uint8
	}{
		"example1_datasheet": {
			pin:                     0,
			onCounts:                409,
			offCounts:               1228,
			wantLedOnTimeOffTimeSet: []uint8{0x06, 0x99, 0x07, 0x01, 0x08, 0xCC, 0x09, 0x04},
		},
		"example2_datasheet": {
			pin:                     4,
			onCounts:                3685,
			offCounts:               3275,
			wantLedOnTimeOffTimeSet: []uint8{0x16, 0x65, 0x17, 0x0E, 0x18, 0xCB, 0x19, 0x0C},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestPCA9685WithStubbedAdaptor()
			a.written = []byte{} // reset writes of former test
			// act
			err := d.SetPWM(tc.pin, tc.onCounts, tc.offCounts)
			// assert
			require.NoError(t, err)
			assert.Len(t, a.written, 8)
			for writeIdx, wantVal := range tc.wantLedOnTimeOffTimeSet {
				assert.Equal(t, wantVal, a.written[writeIdx], "index %d differs", writeIdx)
			}
		})
	}
}

func TestPCA9685SetPWMError(t *testing.T) {
	// arrange
	d, a := initTestPCA9685WithStubbedAdaptor()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	// act & assert
	require.ErrorContains(t, d.SetPWM(0, 0, 256), "write error")
}

func TestPCA9685SetAllPWM(t *testing.T) {
	// sequence to set PWM for PCA9685:
	// * set LEDn ON-time register (n=0: 0x06, 0x07, n=1: 0x0A, 0x0B ... n=14: 0x3E, 0x3F, n=15: 0x42, 0x43)
	// * set LEDn OFF-time register (n=0: 0x08, 0x09, n=1: 0x0C, 0x0D ... n=14: 0x40, 0x41, n=15: 0x44, 0x45)
	tests := map[string]struct {
		pin                     byte
		onCounts                uint16
		offCounts               uint16
		wantLedOnTimeOffTimeSet []uint8
	}{
		"example1_datasheet": {
			onCounts:                409,
			offCounts:               1228,
			wantLedOnTimeOffTimeSet: []uint8{0xFA, 0x99, 0xFB, 0x01, 0xFC, 0xCC, 0xFD, 0x04},
		},
		"example2_datasheet": {
			onCounts:                3685,
			offCounts:               3275,
			wantLedOnTimeOffTimeSet: []uint8{0xFA, 0x65, 0xFB, 0x0E, 0xFC, 0xCB, 0xFD, 0x0C},
		},
		"own_example": {
			onCounts:                1234,
			offCounts:               4321,
			wantLedOnTimeOffTimeSet: []uint8{0xFA, 0xD2, 0xFB, 0x04, 0xFC, 0xE1, 0xFD, 0x10},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestPCA9685WithStubbedAdaptor()
			a.written = []byte{} // reset writes of former test
			// act
			err := d.SetAllPWM(tc.onCounts, tc.offCounts)
			// assert
			require.NoError(t, err)
			assert.Len(t, a.written, 8)
			for writeIdx, wantVal := range tc.wantLedOnTimeOffTimeSet {
				assert.Equal(t, wantVal, a.written[writeIdx], "index %d differs", writeIdx)
			}
		})
	}
}

func TestPCA9685SetAllPWMError(t *testing.T) {
	// arrange
	d, a := initTestPCA9685WithStubbedAdaptor()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	// act & assert
	require.ErrorContains(t, d.SetAllPWM(0, 256), "write error")
}

func TestPCA9685SetPWMFreq(t *testing.T) {
	// sequence to set PWM frequency prescaler for PCA9685 (can only be set in sleep mode):
	// * read MODE1 register (0x00)
	// * prepare MODE1 register with sleep mode set (bit 4 - 0x10, no stopping of PWM channels done before)
	// * write MODE1 register
	// * write the prescaler value to PRE_SCALE register (0xFE)
	// * prepare MIODE1 register with sleep mode bit reset
	// * write MODE1 register
	// * wait > 500us
	// * prepare the MODE1 register with set of reset bit
	// * write MODE1 register
	const readMode1Val = 0x0F // to check for only sleep mode bit (0x10) or reset bit (0x80) will be set
	var (
		wantMode1SleepSequence   = []uint8{0x00, 0x1F}
		wantMode1NoSleepSequence = []uint8{0x00, readMode1Val}
		wantMode1ResetSequence   = []uint8{0x00, 0x8F}
	)
	tests := map[string]struct {
		freq                  float32
		wantPrescalerSequence []uint8
	}{
		"example_datasheet": {
			freq:                  200,
			wantPrescalerSequence: []uint8{0xFE, 0x1E},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestPCA9685WithStubbedAdaptor()
			// arrange read for MODE1 register
			numCallsRead := 0
			a.i2cReadImpl = func(b []byte) (int, error) {
				numCallsRead++
				b[0] = readMode1Val
				return len(b), nil
			}
			a.written = []byte{} // reset writes of former test
			// act
			err := d.SetPWMFreq(tc.freq)
			// assert
			require.NoError(t, err)
			assert.Len(t, a.written, 9)
			var writeIdx int
			// for read old mode:
			assert.Equal(t, wantMode1SleepSequence[0], a.written[writeIdx], "index %d differs", writeIdx)
			writeIdx++
			for idx, wantVal := range wantMode1SleepSequence {
				assert.Equal(t, wantVal, a.written[writeIdx], "index %d (%d) differs", writeIdx, idx)
				writeIdx++
			}
			for idx, wantVal := range tc.wantPrescalerSequence {
				assert.Equal(t, wantVal, a.written[writeIdx], "index %d (%d) differs", writeIdx, idx)
				writeIdx++
			}
			for idx, wantVal := range wantMode1NoSleepSequence {
				assert.Equal(t, wantVal, a.written[writeIdx], "index %d (%d) differs", writeIdx, idx)
				writeIdx++
			}
			for idx, wantVal := range wantMode1ResetSequence {
				assert.Equal(t, wantVal, a.written[writeIdx], "index %d (%d) differs", writeIdx, idx)
				writeIdx++
			}
			assert.Equal(t, 1, numCallsRead)
		})
	}
}

func TestPCA9685SetPWMFreqReadError(t *testing.T) {
	// arrange
	d, a := initTestPCA9685WithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}
	// act & assert
	require.ErrorContains(t, d.SetPWMFreq(60), "read error")
}

func TestPCA9685SetPWMFreqWriteError(t *testing.T) {
	// arrange
	d, a := initTestPCA9685WithStubbedAdaptor()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	// act & assert
	require.ErrorContains(t, d.SetPWMFreq(60), "write error")
}

func TestPCA9685Commands(t *testing.T) {
	// arrange
	d, _ := initTestPCA9685WithStubbedAdaptor()
	// act & assert
	assert.Nil(t, d.Command("PwmWrite")(map[string]interface{}{"pin": "1", "val": "1"}))
	assert.Nil(t, d.Command("ServoWrite")(map[string]interface{}{"pin": "1", "val": "1"}))
	assert.Nil(t, d.Command("SetPWM")(map[string]interface{}{"channel": "1", "on": "0", "off": "1024"}))
	assert.Nil(t, d.Command("SetPWMFreq")(map[string]interface{}{"freq": "60"}))
}

func TestPCA9685_initialize(t *testing.T) {
	// sequence to reset the PCA9685 in initialize():
	// * set all LED ON-time and OFF-time registers (0xFA..0xFD for 16 channels, each for low and high byte)
	// * set MODE2 register (0x01) to defaults: not inverted, outputs change on stop, OE reaction to 0, totem-pole
	// * set MODE1 register (0x00) to defaults, except sleep: no restart, internal clock, no AI, no sleep (not default),
	//   no response to sub address 1, 2 or 3, activate response to all-call
	// * wait > 500us, read back the MODE1 register
	// * prepare the MODE1 register with set of reset bit
	// * write MODE1 register
	// arrange
	d, a := initTestPCA9685WithStubbedAdaptor()
	a.written = []byte{} // reset writes of former test
	wantAllLedOnTimeOffTimeSequence := []uint8{0xFA, 0x00, 0xFB, 0x00, 0xFC, 0x00, 0xFD, 0x00}
	wantMode2RegSetDefaultsSequence := []uint8{0x01, 0x04}
	wantMode1RegSetDefaultsNoSleepSequence := []uint8{0x00, 0x01}
	wantMode1RegResetSequence := []uint8{0x00, 0x81}
	// arrange read for MODE1 register
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		b[0] = wantMode1RegSetDefaultsNoSleepSequence[1]
		return len(b), nil
	}
	// act, assert - initialize() must be called on Start()
	err := d.Start()
	// assert
	require.NoError(t, err)
	assert.Len(t, a.written, 15)
	var writeIdx int
	for idx, wantVal := range wantAllLedOnTimeOffTimeSequence {
		assert.Equal(t, wantVal, a.written[writeIdx], "index %d (%d) differs", writeIdx, idx)
		writeIdx++
	}
	for idx, wantVal := range wantMode2RegSetDefaultsSequence {
		assert.Equal(t, wantVal, a.written[writeIdx], "index %d (%d) differs", writeIdx, idx)
		writeIdx++
	}
	for idx, wantVal := range wantMode1RegSetDefaultsNoSleepSequence {
		assert.Equal(t, wantVal, a.written[writeIdx], "index %d (%d) differs", writeIdx, idx)
		writeIdx++
	}
	// for read old mode:
	assert.Equal(t, wantMode1RegSetDefaultsNoSleepSequence[0], a.written[writeIdx], "index %d differs", writeIdx)
	writeIdx++
	for idx, wantVal := range wantMode1RegResetSequence {
		assert.Equal(t, wantVal, a.written[writeIdx], "index %d (%d) differs", writeIdx, idx)
		writeIdx++
	}
	assert.Equal(t, 1, numCallsRead)
}
