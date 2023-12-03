package gpio

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
)

var _ gobot.Driver = (*BuzzerDriver)(nil)

func initTestBuzzerDriver(conn DigitalWriter) *BuzzerDriver {
	return NewBuzzerDriver(conn, "1")
}

func TestNewBuzzerDriver(t *testing.T) {
	// arrange
	a := newGpioTestAdaptor()
	// act
	d := NewBuzzerDriver(a, "10")
	// assert
	assert.IsType(t, &BuzzerDriver{}, d)
	// assert: gpio.driver attributes
	require.NotNil(t, d.driver)
	assert.True(t, strings.HasPrefix(d.driverCfg.name, "Buzzer"))
	assert.Equal(t, "10", d.driverCfg.pin)
	assert.Equal(t, a, d.connection)
	assert.NotNil(t, d.afterStart)
	assert.NotNil(t, d.beforeHalt)
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
	// assert: driver specific attributes
	assert.False(t, d.high)
	assert.InDelta(t, 96, d.bpm, 0.0)
}

func TestNewBuzzerDriver_options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option, least one
	// option of this driver and one of another driver (which should lead to panic). Further tests for options can also
	// be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const (
		myName = "song player"
	)
	panicFunc := func() {
		NewBuzzerDriver(newGpioTestAdaptor(), "1", WithName("crazy"),
			aio.WithActuatorScaler(func(float64) int { return 0 }))
	}
	// act
	d := NewBuzzerDriver(newGpioTestAdaptor(), "1", WithName(myName))
	// assert
	assert.Equal(t, myName, d.Name())
	assert.PanicsWithValue(t, "'scaler option for analog actuators' can not be applied on 'crazy'", panicFunc)
}

func TestBuzzerToggle(t *testing.T) {
	d := initTestBuzzerDriver(newGpioTestAdaptor())
	require.NoError(t, d.Off())
	require.NoError(t, d.Toggle())
	assert.True(t, d.State())
	require.NoError(t, d.Toggle())
	assert.False(t, d.State())
}

func TestBuzzerTone(t *testing.T) {
	d := initTestBuzzerDriver(newGpioTestAdaptor())
	require.NoError(t, d.Tone(100, 0.01))
}

func TestBuzzerOnError(t *testing.T) {
	a := newGpioTestAdaptor()
	d := initTestBuzzerDriver(a)
	a.digitalWriteFunc = func(string, byte) error {
		return errors.New("write error")
	}

	require.EqualError(t, d.On(), "write error")
}

func TestBuzzerOffError(t *testing.T) {
	a := newGpioTestAdaptor()
	d := initTestBuzzerDriver(a)
	a.digitalWriteFunc = func(string, byte) error {
		return errors.New("write error")
	}

	require.EqualError(t, d.Off(), "write error")
}

func TestBuzzerToneError(t *testing.T) {
	a := newGpioTestAdaptor()
	d := initTestBuzzerDriver(a)
	a.digitalWriteFunc = func(string, byte) error {
		return errors.New("write error")
	}

	require.EqualError(t, d.Tone(100, 0.01), "write error")
}
