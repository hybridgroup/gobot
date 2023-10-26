package gpio

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*MotorDriver)(nil)

func initTestMotorDriver() *MotorDriver {
	return NewMotorDriver(newGpioTestAdaptor(), "1")
}

func TestMotorDriver(t *testing.T) {
	d := NewMotorDriver(newGpioTestAdaptor(), "1")
	assert.NotNil(t, d.Connection())
}

func TestMotorDriverStart(t *testing.T) {
	d := initTestMotorDriver()
	assert.NoError(t, d.Start())
}

func TestMotorDriverHalt(t *testing.T) {
	d := initTestMotorDriver()
	assert.NoError(t, d.Halt())
}

func TestMotorDriverIsOn(t *testing.T) {
	d := initTestMotorDriver()
	d.CurrentMode = "digital"
	d.CurrentState = 1
	assert.True(t, d.IsOn())
	d.CurrentMode = "analog"
	d.CurrentSpeed = 100
	assert.True(t, d.IsOn())
}

func TestMotorDriverIsOff(t *testing.T) {
	d := initTestMotorDriver()
	_ = d.Off()
	assert.True(t, d.IsOff())
}

func TestMotorDriverOn(t *testing.T) {
	d := initTestMotorDriver()
	d.CurrentMode = "digital"
	_ = d.On()
	assert.Equal(t, uint8(1), d.CurrentState)
	d.CurrentMode = "analog"
	d.CurrentSpeed = 0
	_ = d.On()
	assert.Equal(t, uint8(255), d.CurrentSpeed)
}

func TestMotorDriverOff(t *testing.T) {
	d := initTestMotorDriver()
	d.CurrentMode = "digital"
	_ = d.Off()
	assert.Equal(t, uint8(0), d.CurrentState)
	d.CurrentMode = "analog"
	d.CurrentSpeed = 100
	_ = d.Off()
	assert.Equal(t, uint8(0), d.CurrentSpeed)
}

func TestMotorDriverToggle(t *testing.T) {
	d := initTestMotorDriver()
	_ = d.Off()
	_ = d.Toggle()
	assert.True(t, d.IsOn())
	_ = d.Toggle()
	assert.False(t, d.IsOn())
}

func TestMotorDriverMin(t *testing.T) {
	d := initTestMotorDriver()
	_ = d.Min()
}

func TestMotorDriverMax(t *testing.T) {
	d := initTestMotorDriver()
	_ = d.Max()
}

func TestMotorDriverSpeed(t *testing.T) {
	d := initTestMotorDriver()
	_ = d.Speed(100)
}

func TestMotorDriverForward(t *testing.T) {
	d := initTestMotorDriver()
	_ = d.Forward(100)
	assert.Equal(t, uint8(100), d.CurrentSpeed)
	assert.Equal(t, "forward", d.CurrentDirection)
}

func TestMotorDriverBackward(t *testing.T) {
	d := initTestMotorDriver()
	_ = d.Backward(100)
	assert.Equal(t, uint8(100), d.CurrentSpeed)
	assert.Equal(t, "backward", d.CurrentDirection)
}

func TestMotorDriverDirection(t *testing.T) {
	d := initTestMotorDriver()
	_ = d.Direction("none")
	d.DirectionPin = "2"
	_ = d.Direction("forward")
	_ = d.Direction("backward")
}

func TestMotorDriverDigital(t *testing.T) {
	d := initTestMotorDriver()
	d.SpeedPin = "" // Disable speed
	d.CurrentMode = "digital"
	d.ForwardPin = "2"
	d.BackwardPin = "3"

	_ = d.On()
	assert.Equal(t, uint8(1), d.CurrentState)
	_ = d.Off()
	assert.Equal(t, uint8(0), d.CurrentState)
}

func TestMotorDriverDefaultName(t *testing.T) {
	d := initTestMotorDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "Motor"))
}

func TestMotorDriverSetName(t *testing.T) {
	d := initTestMotorDriver()
	d.SetName("mybot")
	assert.Equal(t, "mybot", d.Name())
}
