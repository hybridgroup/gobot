package gpio

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	stepsInRev = 32
)

func initStepperMotorDriver() *StepperDriver {
	return NewStepperDriver(newGpioTestAdaptor(), [4]string{"7", "11", "13", "15"}, StepperModes.DualPhaseStepping, stepsInRev)
}

func TestStepperDriverRun(t *testing.T) {
	d := initStepperMotorDriver()
	_ = d.Run()
	assert.True(t, d.IsMoving())
}

func TestStepperDriverHalt(t *testing.T) {
	d := initStepperMotorDriver()
	_ = d.Run()
	time.Sleep(200 * time.Millisecond)
	_ = d.Halt()
	assert.False(t, d.IsMoving())
}

func TestStepperDriverDefaultName(t *testing.T) {
	d := initStepperMotorDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "Stepper"))
}

func TestStepperDriverSetName(t *testing.T) {
	name := "SomeStepperSriver"
	d := initStepperMotorDriver()
	d.SetName(name)
	assert.Equal(t, name, d.Name())
}

func TestStepperDriverSetDirection(t *testing.T) {
	dir := "backward"
	d := initStepperMotorDriver()
	_ = d.SetDirection(dir)
	assert.Equal(t, dir, d.direction)
}

func TestStepperDriverDefaultDirection(t *testing.T) {
	d := initStepperMotorDriver()
	assert.Equal(t, "forward", d.direction)
}

func TestStepperDriverInvalidDirection(t *testing.T) {
	d := initStepperMotorDriver()
	err := d.SetDirection("reverse")
	assert.ErrorContains(t, err, "Invalid direction. Value should be forward or backward")
}

func TestStepperDriverMoveForward(t *testing.T) {
	d := initStepperMotorDriver()
	_ = d.Move(1)
	assert.Equal(t, 1, d.GetCurrentStep())

	_ = d.Move(10)
	assert.Equal(t, 11, d.GetCurrentStep())
}

func TestStepperDriverMoveBackward(t *testing.T) {
	d := initStepperMotorDriver()
	_ = d.Move(-1)
	assert.Equal(t, stepsInRev-1, d.GetCurrentStep())

	_ = d.Move(-10)
	assert.Equal(t, stepsInRev-11, d.GetCurrentStep())
}

func TestStepperDriverMoveFullRotation(t *testing.T) {
	d := initStepperMotorDriver()
	_ = d.Move(stepsInRev)
	assert.Equal(t, 0, d.GetCurrentStep())
}

func TestStepperDriverMotorSetSpeedMoreThanMax(t *testing.T) {
	d := initStepperMotorDriver()
	m := d.GetMaxSpeed()

	_ = d.SetSpeed(m + 1)
	assert.Equal(t, d.speed, m)
}

func TestStepperDriverMotorSetSpeedLessOrEqualMax(t *testing.T) {
	d := initStepperMotorDriver()
	m := d.GetMaxSpeed()

	_ = d.SetSpeed(m - 1)
	assert.Equal(t, d.speed, m-1)

	_ = d.SetSpeed(m)
	assert.Equal(t, d.speed, m)
}
