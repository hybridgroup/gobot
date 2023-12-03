//nolint:forcetypeassert // ok here
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

var _ gobot.Driver = (*LedDriver)(nil)

func initTestLedDriver() *LedDriver {
	a := newGpioTestAdaptor()
	a.digitalWriteFunc = func(string, byte) error {
		return nil
	}
	a.pwmWriteFunc = func(string, byte) error {
		return nil
	}
	return NewLedDriver(a, "1")
}

func TestNewLedDriver(t *testing.T) {
	// arrange
	a := newGpioTestAdaptor()
	// act
	d := NewLedDriver(a, "10")
	// assert
	assert.IsType(t, &LedDriver{}, d)
	// assert: gpio.driver attributes
	require.NotNil(t, d.driver)
	assert.True(t, strings.HasPrefix(d.driverCfg.name, "LED"))
	assert.Equal(t, "10", d.driverCfg.pin)
	assert.Equal(t, a, d.connection)
	assert.NotNil(t, d.afterStart)
	assert.NotNil(t, d.beforeHalt)
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
	// assert: driver specific attributes
	assert.False(t, d.high)
}

func TestNewLedDriver_options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option, least one
	// option of this driver and one of another driver (which should lead to panic). Further tests for options can also
	// be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const (
		myName = "back light"
	)
	panicFunc := func() {
		NewLedDriver(newGpioTestAdaptor(), "1", WithName("crazy"), aio.WithActuatorScaler(func(float64) int { return 0 }))
	}
	// act
	d := NewLedDriver(newGpioTestAdaptor(), "1", WithName(myName))
	// assert
	assert.Equal(t, myName, d.Name())
	assert.PanicsWithValue(t, "'scaler option for analog actuators' can not be applied on 'crazy'", panicFunc)
}

func TestLed_Commands(t *testing.T) {
	var err interface{}
	a := newGpioTestAdaptor()
	d := NewLedDriver(a, "1")

	a.digitalWriteFunc = func(string, byte) error {
		return errors.New("write error")
	}
	a.pwmWriteFunc = func(string, byte) error {
		return errors.New("pwm error")
	}

	err = d.Command("Toggle")(nil)
	require.EqualError(t, err.(error), "write error")

	err = d.Command("On")(nil)
	require.EqualError(t, err.(error), "write error")

	err = d.Command("Off")(nil)
	require.EqualError(t, err.(error), "write error")

	err = d.Command("Brightness")(map[string]interface{}{"level": 100.0})
	require.EqualError(t, err.(error), "pwm error")
}

func TestLedToggle(t *testing.T) {
	d := initTestLedDriver()
	require.NoError(t, d.Off())
	require.NoError(t, d.Toggle())
	assert.True(t, d.State())
	require.NoError(t, d.Toggle())
	assert.False(t, d.State())
}

func TestLedBrightness(t *testing.T) {
	a := newGpioTestAdaptor()
	d := NewLedDriver(a, "1")
	a.pwmWriteFunc = func(string, byte) error {
		return errors.New("pwm error")
	}
	require.EqualError(t, d.Brightness(150), "pwm error")
}
