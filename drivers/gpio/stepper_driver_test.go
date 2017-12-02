package gpio

import (
	"errors"
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot/gobottest"
)

const (
	stepsInRev = 32
)

func initStepperMotorDriver() *StepperDriver {
	return NewStepperDriver(newGpioTestAdaptor(), [4]string{"7", "11", "13", "15"}, StepperModes.DualPhaseStepping, stepsInRev)
}

func TestStepperDriverRun(t *testing.T) {
	d := initStepperMotorDriver()
	d.Run()
	gobottest.Assert(t, d.IsMoving(), true)
}

func TestStepperDriverHalt(t *testing.T) {
	d := initStepperMotorDriver()
	d.Run()
	time.Sleep(200 * time.Millisecond)
	d.Halt()
	gobottest.Assert(t, d.IsMoving(), false)
}

func TestStepperDriverDefaultName(t *testing.T) {
	d := initStepperMotorDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Stepper"), true)
}

func TestStepperDriverSetName(t *testing.T) {
	name := "SomeStepperSriver"
	d := initStepperMotorDriver()
	d.SetName(name)
	gobottest.Assert(t, d.Name(), name)
}

func TestStepperDriverSetDirection(t *testing.T) {
	dir := "backward"
	d := initStepperMotorDriver()
	d.SetDirection(dir)
	gobottest.Assert(t, d.direction, dir)
}

func TestStepperDriverDefaultDirection(t *testing.T) {
	d := initStepperMotorDriver()
	gobottest.Assert(t, d.direction, "forward")
}

func TestStepperDriverInvalidDirection(t *testing.T) {
	d := initStepperMotorDriver()
	err := d.SetDirection("reverse")
	gobottest.Assert(t, err.(error), errors.New("Invalid direction. Value should be forward or backward"))
}

func TestStepperDriverMoveForward(t *testing.T) {
	d := initStepperMotorDriver()
	d.Move(1)
	gobottest.Assert(t, d.GetCurrentStep(), 1)

	d.Move(10)
	gobottest.Assert(t, d.GetCurrentStep(), 11)
}

func TestStepperDriverMoveBackward(t *testing.T) {
	d := initStepperMotorDriver()
	d.Move(-1)
	gobottest.Assert(t, d.GetCurrentStep(), stepsInRev-1)

	d.Move(-10)
	gobottest.Assert(t, d.GetCurrentStep(), stepsInRev-11)
}

func TestStepperDriverMoveFullRotation(t *testing.T) {
	d := initStepperMotorDriver()
	d.Move(stepsInRev)
	gobottest.Assert(t, d.GetCurrentStep(), 0)
}

func TestStepperDriverMotorSetSpeedMoreThanMax(t *testing.T) {
	d := initStepperMotorDriver()
	m := d.GetMaxSpeed()

	d.SetSpeed(m + 1)
	gobottest.Assert(t, m, d.speed)
}

func TestStepperDriverMotorSetSpeedLessOrEqualMax(t *testing.T) {
	d := initStepperMotorDriver()
	m := d.GetMaxSpeed()

	d.SetSpeed(m - 1)
	gobottest.Assert(t, m-1, d.speed)

	d.SetSpeed(m)
	gobottest.Assert(t, m, d.speed)
}
