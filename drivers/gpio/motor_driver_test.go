package gpio

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
)

var _ gobot.Driver = (*MotorDriver)(nil)

func initTestMotorDriver() *MotorDriver {
	return NewMotorDriver(newGpioTestAdaptor(), "1")
}

func TestMotorDriver(t *testing.T) {
	d := NewMotorDriver(newGpioTestAdaptor(), "1")
	assert.NotNil(t, d.Connection())
}

func TestNewMotorDriver(t *testing.T) {
	// arrange
	a := newGpioTestAdaptor()
	// act
	d := NewMotorDriver(a, "10")
	// assert
	assert.IsType(t, &MotorDriver{}, d)
	// assert: gpio.driver attributes
	require.NotNil(t, d.driver)
	assert.True(t, strings.HasPrefix(d.driverCfg.name, "Motor"))
	assert.Equal(t, "10", d.driverCfg.pin)
	assert.Equal(t, a, d.connection)
	assert.NotNil(t, d.afterStart)
	assert.NotNil(t, d.beforeHalt)
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
	// assert: driver specific attributes
	assert.Equal(t, "", d.motorCfg.directionPin)
	assert.Equal(t, "", d.motorCfg.forwardPin)
	assert.Equal(t, "", d.motorCfg.backwardPin)
	assert.False(t, d.motorCfg.modeIsAnalog)
	assert.Equal(t, uint8(0), d.currentState)
	assert.Equal(t, uint8(0), d.currentSpeed)
	assert.Equal(t, "forward", d.currentDirection)
}

func TestNewMotorDriver_options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option, least one
	// option of this driver and one of another driver (which should lead to panic). Further tests for options can also
	// be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const (
		myName = "move analog"
	)
	panicFunc := func() {
		NewMotorDriver(newGpioTestAdaptor(), "1", WithName("crazy"), aio.WithActuatorScaler(func(float64) int { return 0 }))
	}
	// act
	d := NewMotorDriver(newGpioTestAdaptor(), "1", WithName(myName), WithMotorAnalog())
	// assert
	assert.True(t, d.motorCfg.modeIsAnalog)
	assert.Equal(t, myName, d.Name())
	assert.PanicsWithValue(t, "'scaler option for analog actuators' can not be applied on 'crazy'", panicFunc)
}

func TestMotor_WithMotorDirectionPin(t *testing.T) {
	// arrange
	const myPin = "8"
	cfg := motorConfiguration{directionPin: "old_pin"}
	// act
	WithMotorDirectionPin(myPin).apply(&cfg)
	// assert
	assert.Equal(t, myPin, cfg.directionPin)
}

func TestMotor_WithMotorForwardPin(t *testing.T) {
	// arrange
	const myPin = "3"
	cfg := motorConfiguration{directionPin: "old_pin"}
	// act
	WithMotorForwardPin(myPin).apply(&cfg)
	// assert
	assert.Equal(t, myPin, cfg.forwardPin)
}

func TestMotor_WithMotorBackwardPin(t *testing.T) {
	// arrange
	const myPin = "6"
	cfg := motorConfiguration{directionPin: "old_pin"}
	// act
	WithMotorBackwardPin(myPin).apply(&cfg)
	// assert
	assert.Equal(t, myPin, cfg.backwardPin)
}

func TestMotorIsOn(t *testing.T) {
	d := initTestMotorDriver()
	d.motorCfg.modeIsAnalog = false
	d.currentState = 1
	assert.True(t, d.IsDigital())
	assert.True(t, d.IsOn())
	d.motorCfg.modeIsAnalog = true
	d.currentSpeed = 100
	assert.False(t, d.IsDigital())
	assert.True(t, d.IsOn())
}

func TestMotorIsOff(t *testing.T) {
	d := initTestMotorDriver()
	require.NoError(t, d.Off())
	assert.True(t, d.IsOff())
}

func TestMotorOn(t *testing.T) {
	d := initTestMotorDriver()
	d.motorCfg.modeIsAnalog = false
	assert.True(t, d.IsDigital())
	require.NoError(t, d.On())
	assert.Equal(t, uint8(1), d.currentState)
	d.motorCfg.modeIsAnalog = true
	d.currentSpeed = 0
	assert.False(t, d.IsDigital())
	require.NoError(t, d.On())
	assert.Equal(t, uint8(255), d.currentSpeed)
}

func TestMotorOff(t *testing.T) {
	d := initTestMotorDriver()
	d.motorCfg.modeIsAnalog = false
	assert.True(t, d.IsDigital())
	require.NoError(t, d.Off())
	assert.Equal(t, uint8(0), d.currentState)
	d.motorCfg.modeIsAnalog = true
	d.currentSpeed = 100
	assert.False(t, d.IsDigital())
	require.NoError(t, d.Off())
	assert.Equal(t, uint8(0), d.currentSpeed)
}

func TestMotorToggle(t *testing.T) {
	d := initTestMotorDriver()
	require.NoError(t, d.Off())
	require.NoError(t, d.Toggle())
	assert.True(t, d.IsOn())
	require.NoError(t, d.Toggle())
	assert.False(t, d.IsOn())
}

func TestMotorRunMin(t *testing.T) {
	d := initTestMotorDriver()
	require.NoError(t, d.RunMin())
}

func TestMotorRunMax(t *testing.T) {
	d := initTestMotorDriver()
	require.NoError(t, d.RunMax())
}

func TestMotorSetSpeed(t *testing.T) {
	d := initTestMotorDriver()
	require.NoError(t, d.SetSpeed(100))
}

func TestMotorForward(t *testing.T) {
	d := initTestMotorDriver()
	require.NoError(t, d.Forward(100))
	assert.Equal(t, uint8(100), d.currentSpeed)
	assert.Equal(t, "forward", d.currentDirection)
}

func TestMotorBackward(t *testing.T) {
	d := initTestMotorDriver()
	require.NoError(t, d.Backward(100))
	assert.Equal(t, uint8(100), d.currentSpeed)
	assert.Equal(t, "backward", d.currentDirection)
}

func TestMotorSetDirection(t *testing.T) {
	d := initTestMotorDriver()
	require.NoError(t, d.SetDirection("none"))
	d.motorCfg.directionPin = "2"
	require.NoError(t, d.SetDirection("forward"))
	require.NoError(t, d.SetDirection("backward"))
}

func TestMotorDigital(t *testing.T) {
	d := initTestMotorDriver()
	d.driverCfg.pin = "" // Disable speed
	d.motorCfg.modeIsAnalog = false
	d.motorCfg.forwardPin = "2"
	d.motorCfg.backwardPin = "3"

	require.NoError(t, d.On())
	assert.Equal(t, uint8(1), d.currentState)
	require.NoError(t, d.Off())
	assert.Equal(t, uint8(0), d.currentState)
}
