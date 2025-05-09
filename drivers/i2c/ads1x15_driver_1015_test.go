package i2c

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initTestADS1015DriverWithStubbedAdaptor() (*ADS1x15Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewADS1015Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewADS1015Driver(t *testing.T) {
	var di interface{} = NewADS1015Driver(newI2cTestAdaptor())
	d, ok := di.(*ADS1x15Driver)
	if !ok {
		require.Fail(t, "NewADS1015Driver() should have returned a *ADS1x15Driver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "ADS1015"))
	for i := 0; i <= 3; i++ {
		assert.Equal(t, 1, d.channelCfgs[i].gain)
		assert.Equal(t, 1600, d.channelCfgs[i].dataRate)
	}
}

func TestADS1015Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewADS1015Driver(newI2cTestAdaptor(), WithBus(2), WithADS1x15Gain(2), WithADS1x15DataRate(920))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
	for i := 0; i <= 3; i++ {
		assert.Equal(t, 2, d.channelCfgs[i].gain)
		assert.Equal(t, 920, d.channelCfgs[i].dataRate)
	}
}

func TestADS1015WithADS1x15BestGainForVoltage(t *testing.T) {
	d, _ := initTestADS1015DriverWithStubbedAdaptor()
	WithADS1x15BestGainForVoltage(1.01)(d)
	for i := 0; i <= 3; i++ {
		assert.Equal(t, 3, d.channelCfgs[i].gain)
	}
}

func TestADS1015WithADS1x15ChannelBestGainForVoltage(t *testing.T) {
	d, _ := initTestADS1015DriverWithStubbedAdaptor()
	WithADS1x15ChannelBestGainForVoltage(0, 1.0)(d)
	WithADS1x15ChannelBestGainForVoltage(1, 2.5)(d)
	WithADS1x15ChannelBestGainForVoltage(2, 3.3)(d)
	WithADS1x15ChannelBestGainForVoltage(3, 5.0)(d)
	assert.Equal(t, 3, d.channelCfgs[0].gain)
	assert.Equal(t, 1, d.channelCfgs[1].gain)
	assert.Equal(t, 1, d.channelCfgs[2].gain)
	assert.Equal(t, 0, d.channelCfgs[3].gain)
}

func TestADS1015AnalogRead(t *testing.T) {
	d, a := initTestADS1015DriverWithStubbedAdaptor()
	WithADS1x15WaitSingleCycle()(d)

	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x7F, 0xFF})
		return 2, nil
	}

	val, err := d.AnalogRead("0")
	assert.Equal(t, 32767, val)
	require.NoError(t, err)

	val, err = d.AnalogRead("1")
	assert.Equal(t, 32767, val)
	require.NoError(t, err)

	val, err = d.AnalogRead("2")
	assert.Equal(t, 32767, val)
	require.NoError(t, err)

	val, err = d.AnalogRead("3")
	assert.Equal(t, 32767, val)
	require.NoError(t, err)

	val, err = d.AnalogRead("0-1")
	assert.Equal(t, 32767, val)
	require.NoError(t, err)

	val, err = d.AnalogRead("0-3")
	assert.Equal(t, 32767, val)
	require.NoError(t, err)

	val, err = d.AnalogRead("1-3")
	assert.Equal(t, 32767, val)
	require.NoError(t, err)

	val, err = d.AnalogRead("2-3")
	assert.Equal(t, 32767, val)
	require.NoError(t, err)

	_, err = d.AnalogRead("3-2")
	require.Error(t, err)
}

func TestADS1x15AnalogReadError(t *testing.T) {
	d, a := initTestADS1015DriverWithStubbedAdaptor()

	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}

	_, err := d.AnalogRead("0")
	require.ErrorContains(t, err, "read error")
}

func TestADS1x15AnalogReadInvalidPin(t *testing.T) {
	d, _ := initTestADS1015DriverWithStubbedAdaptor()

	_, err := d.AnalogRead("99")
	require.ErrorContains(t, err, "Invalid channel (99), must be between 0 and 3")
}

func TestADS1x15AnalogReadWriteError(t *testing.T) {
	d, a := initTestADS1015DriverWithStubbedAdaptor()

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	_, err := d.AnalogRead("0")
	require.ErrorContains(t, err, "write error")

	_, err = d.AnalogRead("0-1")
	require.ErrorContains(t, err, "write error")

	_, err = d.AnalogRead("2-3")
	require.ErrorContains(t, err, "write error")
}

func TestADS1x15ReadInvalidChannel(t *testing.T) {
	d, _ := initTestADS1015DriverWithStubbedAdaptor()

	_, err := d.Read(9, 1, 1600)
	require.ErrorContains(t, err, "Invalid channel (9), must be between 0 and 3")
}

func TestADS1x15ReadInvalidGain(t *testing.T) {
	d, _ := initTestADS1015DriverWithStubbedAdaptor()

	_, err := d.Read(0, 8, 1600)
	require.ErrorContains(t, err, "Gain (8) must be one of: [0 1 2 3 4 5 6 7]")
}

func TestADS1x15ReadInvalidDataRate(t *testing.T) {
	d, _ := initTestADS1015DriverWithStubbedAdaptor()

	_, err := d.Read(0, 1, 321)
	require.ErrorContains(t, err, "Invalid data rate (321). Accepted values: [128 250 490 920 1600 2400 3300]")
}

func TestADS1x15ReadDifferenceInvalidChannel(t *testing.T) {
	d, _ := initTestADS1015DriverWithStubbedAdaptor()

	_, err := d.ReadDifference(9, 1, 1600)
	require.ErrorContains(t, err, "Invalid channel (9), must be between 0 and 3")
}

func TestADS1015_rawRead(t *testing.T) {
	// sequence to read:
	// * prepare config register content (mode, input, gain, data rate, comparator)
	// * write config register (16 bit, MSByte first)
	// * read config register (16 bit, MSByte first) and wait for bit 15 is set
	// * read conversion register (16 bit, MSByte first) for the value
	// * apply two's complement converter, relates to one digit resolution (1/2^15), voltage multiplier
	tests := map[string]struct {
		input      []uint8
		gain       int
		dataRate   int
		want       int
		wantConfig []uint8
	}{
		"+FS": {
			input:      []uint8{0x7F, 0xFF},
			gain:       0,
			dataRate:   128,
			want:       (1<<15 - 1),
			wantConfig: []uint8{0x91, 0x03},
		},
		"+1": {
			input:      []uint8{0x00, 0x01},
			gain:       0,
			dataRate:   250,
			want:       1,
			wantConfig: []uint8{0x91, 0x23},
		},
		"+-0": {
			input:      []uint8{0x00, 0x00},
			gain:       0,
			dataRate:   490,
			want:       0,
			wantConfig: []uint8{0x91, 0x43},
		},
		"-1": {
			input:      []uint8{0xFF, 0xFF},
			gain:       0,
			dataRate:   920,
			want:       -1,
			wantConfig: []uint8{0x91, 0x63},
		},
		"-FS": {
			input:      []uint8{0x80, 0x00},
			gain:       0,
			dataRate:   1600,
			want:       -(1 << 15),
			wantConfig: []uint8{0x91, 0x83},
		},
		"+FS gain 1": {
			input:      []uint8{0x7F, 0xFF},
			gain:       1,
			dataRate:   2400,
			want:       (1<<15 - 1),
			wantConfig: []uint8{0x93, 0xA3},
		},
		"+FS gain 3": {
			input:      []uint8{0x7F, 0xFF},
			gain:       3,
			dataRate:   3300,
			want:       (1<<15 - 1),
			wantConfig: []uint8{0x97, 0xC3},
		},
		"+FS gain 5": {
			input:      []uint8{0x7F, 0xFF},
			gain:       5,
			dataRate:   2400,
			want:       (1<<15 - 1),
			wantConfig: []uint8{0x9B, 0xA3},
		},
		"+FS gain 7": {
			input:      []uint8{0x7F, 0xFF},
			gain:       7,
			dataRate:   1600,
			want:       (1<<15 - 1),
			wantConfig: []uint8{0x9F, 0x83},
		},
	}
	d, a := initTestADS1015DriverWithStubbedAdaptor()
	// arrange
	channel := 0
	channelOffset := 1
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			a.written = []byte{} // reset writes of Start() and former test
			// arrange reads
			conversion := []uint8{0x00, 0x00}   // a conversion is in progress
			noConversion := []uint8{0x80, 0x00} // no conversion in progress
			returnRead := [3][]uint8{conversion, noConversion, tc.input}
			numCallsRead := 0
			a.i2cReadImpl = func(b []byte) (int, error) {
				numCallsRead++
				retRead := returnRead[numCallsRead-1]
				copy(b, retRead)
				return len(b), nil
			}
			// act
			got, err := d.rawRead(channel, channelOffset, tc.gain, tc.dataRate)
			// assert
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, 3, numCallsRead)
			assert.Len(t, a.written, 6)
			assert.Equal(t, uint8(ads1x15PointerConfig), a.written[0])
			assert.Equal(t, tc.wantConfig[0], a.written[1])            // MSByte: OS, MUX, PGA, MODE
			assert.Equal(t, tc.wantConfig[1], a.written[2])            // LSByte: DR, COMP_*
			assert.Equal(t, uint8(ads1x15PointerConfig), a.written[3]) // first check for no conversion
			assert.Equal(t, uint8(ads1x15PointerConfig), a.written[4]) // second check for no conversion
			assert.Equal(t, uint8(ads1x15PointerConversion), a.written[5])
		})
	}
}
