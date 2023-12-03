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

var _ gobot.Driver = (*ServoDriver)(nil)

func initTestServoDriver() *ServoDriver {
	return NewServoDriver(newGpioTestAdaptor(), "1")
}

func TestNewServoDriver(t *testing.T) {
	// arrange
	a := newGpioTestAdaptor()
	// act
	d := NewServoDriver(a, "10")
	// assert
	assert.IsType(t, &ServoDriver{}, d)
	// assert: gpio.driver attributes
	require.NotNil(t, d.driver)
	assert.True(t, strings.HasPrefix(d.driverCfg.name, "Servo"))
	assert.Equal(t, "10", d.driverCfg.pin)
	assert.Equal(t, a, d.connection)
	assert.NotNil(t, d.afterStart)
	assert.NotNil(t, d.beforeHalt)
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
	// assert: driver specific attributes
	assert.Equal(t, uint8(0), d.currentAngle)
}

func TestNewServoDriver_options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option, least one
	// option of this driver and one of another driver (which should lead to panic). Further tests for options can also
	// be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const (
		myName = "left wheel"
	)
	panicFunc := func() {
		NewServoDriver(newGpioTestAdaptor(), "1", WithName("crazy"),
			aio.WithActuatorScaler(func(float64) int { return 0 }))
	}
	// act
	d := NewServoDriver(newGpioTestAdaptor(), "1", WithName(myName))
	// assert
	assert.Equal(t, myName, d.Name())
	assert.PanicsWithValue(t, "'scaler option for analog actuators' can not be applied on 'crazy'", panicFunc)
}

func TestServo_Commands(t *testing.T) {
	var err interface{}

	a := newGpioTestAdaptor()
	d := NewServoDriver(a, "1")

	a.servoWriteFunc = func(string, byte) error {
		return errors.New("pwm error")
	}

	err = d.Command("ToMin")(nil)
	require.EqualError(t, err.(error), "pwm error")

	err = d.Command("ToCenter")(nil)
	require.EqualError(t, err.(error), "pwm error")

	err = d.Command("ToMax")(nil)
	require.EqualError(t, err.(error), "pwm error")

	err = d.Command("Move")(map[string]interface{}{"angle": 100.0})
	require.EqualError(t, err.(error), "pwm error")
}

func TestServoMove(t *testing.T) {
	d := initTestServoDriver()
	_ = d.Move(100)
	assert.Equal(t, uint8(100), d.currentAngle)
	err := d.Move(200)
	require.EqualError(t, err, "servo angle (200) must be between 0-180")
}

func TestServoMin(t *testing.T) {
	d := initTestServoDriver()
	_ = d.ToMin()
	assert.Equal(t, uint8(0), d.currentAngle)
	assert.Equal(t, d.currentAngle, d.Angle())
}

func TestServoMax(t *testing.T) {
	d := initTestServoDriver()
	_ = d.ToMax()
	assert.Equal(t, uint8(180), d.currentAngle)
}

func TestServoCenter(t *testing.T) {
	d := initTestServoDriver()
	_ = d.ToCenter()
	assert.Equal(t, uint8(90), d.currentAngle)
}
