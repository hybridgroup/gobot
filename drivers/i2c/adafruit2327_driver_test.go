package i2c

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

// this ensures that the implementation implements the gobot.Driver interface
var _ gobot.Driver = (*Adafruit2327Driver)(nil)

func initTestAdafruit2327WithStubbedAdaptor() (*Adafruit2327Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewAdafruit2327Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewAdafruit2327Driver(t *testing.T) {
	// arrange & act
	d := NewAdafruit2327Driver(newI2cTestAdaptor())
	// assert
	assert.IsType(t, &Adafruit2327Driver{}, d)
	assert.True(t, strings.HasPrefix(d.Name(), "Adafruit2327ServoHat"))
	assert.Equal(t, 0x40, d.defaultAddress)
}

func TestAdafruit2327Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	// arrange & act
	d := NewAdafruit2327Driver(newI2cTestAdaptor(), WithBus(2))
	// assert
	assert.Equal(t, 2, d.GetBusOrDefault(1))
}

func TestAdafruit2327SetServoMotorFreq(t *testing.T) {
	// arrange
	d, a := initTestAdafruit2327WithStubbedAdaptor()
	a.written = []byte{} // reset writes of former test
	const freq = 60.0
	// act
	err := d.SetServoMotorFreq(freq)
	// assert
	require.NoError(t, err)
	assert.Len(t, a.written, 9) // detailed test, see "TestPCA9685SetPWMFreq"
}

func TestAdafruit2327SetServoMotorFreqError(t *testing.T) {
	// arrange
	d, a := initTestAdafruit2327WithStubbedAdaptor()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	const freq = 60.0
	// act & assert
	require.ErrorContains(t, d.SetServoMotorFreq(freq), "write error")
}

func TestAdafruit2327SetServoMotorPulse(t *testing.T) {
	// arrange
	d, a := initTestAdafruit2327WithStubbedAdaptor()
	a.written = []byte{} // reset writes of former test
	const (
		channel byte  = 7
		on      int32 = 1234
		off     int32 = 4321
	)
	// act
	err := d.SetServoMotorPulse(channel, on, off)
	// assert
	require.NoError(t, err)
	assert.Len(t, a.written, 8) // detailed test, see "TestPCA9685SetPWM"
}

func TestAdafruit2327SetServoMotorPulseError(t *testing.T) {
	// arrange
	d, a := initTestAdafruit2327WithStubbedAdaptor()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	const (
		channel byte  = 7
		on      int32 = 1234
		off     int32 = 4321
	)
	// act & assert
	require.ErrorContains(t, d.SetServoMotorPulse(channel, on, off), "write error")
}
