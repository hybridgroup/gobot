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

var _ gobot.Driver = (*RgbLedDriver)(nil)

func initTestRgbLedDriver() *RgbLedDriver {
	a := newGpioTestAdaptor()
	a.digitalWriteFunc = func(string, byte) error {
		return nil
	}
	a.pwmWriteFunc = func(string, byte) error {
		return nil
	}
	return NewRgbLedDriver(a, "1", "2", "3")
}

func TestNewRgbLedDriver(t *testing.T) {
	// arrange
	a := newGpioTestAdaptor()
	// act
	d := NewRgbLedDriver(a, "10", "20", "30")
	// assert
	assert.IsType(t, &RgbLedDriver{}, d)
	// assert: gpio.driver attributes
	require.NotNil(t, d.driver)
	assert.True(t, strings.HasPrefix(d.driverCfg.name, "RGBLED"))
	assert.Equal(t, a, d.connection)
	assert.NotNil(t, d.afterStart)
	assert.NotNil(t, d.beforeHalt)
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
	// assert: driver specific attributes
	assert.Equal(t, "10", d.RedPin())
	assert.Equal(t, "20", d.GreenPin())
	assert.Equal(t, "30", d.BluePin())
	assert.Equal(t, "r=10, g=20, b=30", d.Pin())
	assert.Equal(t, uint8(0), d.redColor)
	assert.Equal(t, uint8(0), d.greenColor)
	assert.Equal(t, uint8(0), d.blueColor)
	assert.False(t, d.high)
}

func TestNewRgbLedDriver_options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option, least one
	// option of this driver and one of another driver (which should lead to panic). Further tests for options can also
	// be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const (
		myName = "colored light"
	)
	panicFunc := func() {
		NewRgbLedDriver(newGpioTestAdaptor(), "1", "2", "3", WithName("crazy"),
			aio.WithActuatorScaler(func(float64) int { return 0 }))
	}
	// act
	d := NewRgbLedDriver(newGpioTestAdaptor(), "1", "2", "3", WithName(myName))
	// assert
	assert.Equal(t, myName, d.Name())
	assert.PanicsWithValue(t, "'scaler option for analog actuators' can not be applied on 'crazy'", panicFunc)
}

func TestRgbLed_Commands(t *testing.T) {
	var err interface{}

	a := newGpioTestAdaptor()
	d := NewRgbLedDriver(a, "1", "2", "3")

	a.digitalWriteFunc = func(string, byte) error {
		return errors.New("write error")
	}
	a.pwmWriteFunc = func(string, byte) error {
		return errors.New("pwm error")
	}

	err = d.Command("Toggle")(nil)
	require.EqualError(t, err.(error), "pwm error")

	err = d.Command("On")(nil)
	require.EqualError(t, err.(error), "pwm error")

	err = d.Command("Off")(nil)
	require.EqualError(t, err.(error), "pwm error")

	err = d.Command("SetRGB")(map[string]interface{}{"r": 0xff, "g": 0xff, "b": 0xff})
	require.EqualError(t, err.(error), "pwm error")
}

func TestRgbLedDriverToggle(t *testing.T) {
	d := initTestRgbLedDriver()
	_ = d.Off()
	_ = d.Toggle()
	assert.True(t, d.State())
	_ = d.Toggle()
	assert.False(t, d.State())
}

func TestRgbLedSetLevel(t *testing.T) {
	a := newGpioTestAdaptor()
	d := NewRgbLedDriver(a, "1", "2", "3")
	require.NoError(t, d.SetLevel("1", 150))

	d = NewRgbLedDriver(a, "1", "2", "3")
	a.pwmWriteFunc = func(string, byte) error {
		return errors.New("pwm error")
	}
	require.EqualError(t, d.SetLevel("1", 150), "pwm error")
}
