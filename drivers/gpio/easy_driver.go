package gpio

import (
	"fmt"
	"strings"
	"time"

	"gobot.io/x/gobot/v2"
)

const easyDriverDebug = false

// EasyDriver is an driver for stepper hardware board from SparkFun (https://www.sparkfun.com/products/12779)
// This should also work for the BigEasyDriver (untested). It is basically a wrapper for the common StepperDriver{}
// with the specific additions for the board, e.g. direction, enable and sleep outputs.
type EasyDriver struct {
	*StepperDriver

	stepPin      string
	dirPin       string
	enPin        string
	sleepPin     string
	anglePerStep float32

	sleeping bool
}

// NewEasyDriver returns a new driver
// TODO: Support selecting phase input instead of hard-wiring MS1 and MS2 to board truth table
// A - DigitalWriter
// anglePerStep - Step angle of motor
// stepPin - Pin corresponding to step input on EasyDriver
// dirPin - Pin corresponding to dir input on EasyDriver.  Optional
// enPin - Pin corresponding to enabled input on EasyDriver.  Optional
// sleepPin - Pin corresponding to sleep input on EasyDriver.  Optional
func NewEasyDriver(
	a DigitalWriter,
	anglePerStep float32,
	stepPin string,
	dirPin string,
	enPin string,
	sleepPin string,
) *EasyDriver {
	if anglePerStep <= 0 {
		panic("angle per step needs to be greater than zero")
	}
	// panic if step pin isn't set
	if stepPin == "" {
		panic("Step pin is not set")
	}

	stepper := NewStepperDriver(a, [4]string{}, nil, 1)
	stepper.name = gobot.DefaultName("EasyDriver")
	stepper.stepperDebug = easyDriverDebug
	stepper.haltIfRunning = false
	stepper.stepsPerRev = 360.0 / anglePerStep
	d := &EasyDriver{
		StepperDriver: stepper,
		stepPin:       stepPin,
		dirPin:        dirPin,
		enPin:         enPin,
		sleepPin:      sleepPin,
		anglePerStep:  anglePerStep,

		sleeping: false,
	}
	d.stepFunc = d.onePinStepping
	d.sleepFunc = d.sleepWithSleepPin
	d.beforeHalt = d.shutdown

	// 1/4 of max speed. Not too fast, not too slow
	d.speedRpm = d.MaxSpeed() / 4

	return d
}

// SetDirection sets the direction to be moving.
func (d *EasyDriver) SetDirection(direction string) error {
	direction = strings.ToLower(direction)
	if direction != StepperDriverForward && direction != StepperDriverBackward {
		return fmt.Errorf("Invalid direction '%s'. Value should be '%s' or '%s'",
			direction, StepperDriverForward, StepperDriverBackward)
	}

	if d.dirPin == "" {
		return fmt.Errorf("dirPin is not set for '%s'", d.name)
	}

	writeVal := byte(0) // low is forward
	if direction == StepperDriverBackward {
		writeVal = 1 // high is backward
	}

	if err := d.connection.(DigitalWriter).DigitalWrite(d.dirPin, writeVal); err != nil {
		return err
	}

	// ensure that write of variable can not interfere with read in step()
	d.valueMutex.Lock()
	defer d.valueMutex.Unlock()
	d.direction = direction

	return nil
}

// Enable enables all motor output
func (d *EasyDriver) Enable() error {
	if d.enPin == "" {
		d.disabled = false
		return fmt.Errorf("enPin is not set - board '%s' is enabled by default", d.name)
	}

	// enPin is active low
	if err := d.connection.(DigitalWriter).DigitalWrite(d.enPin, 0); err != nil {
		return err
	}

	d.disabled = false
	return nil
}

// Disable disables all motor output
func (d *EasyDriver) Disable() error {
	if d.enPin == "" {
		return fmt.Errorf("enPin is not set for '%s'", d.name)
	}

	_ = d.stopIfRunning() // drop step errors

	// enPin is active low
	if err := d.connection.(DigitalWriter).DigitalWrite(d.enPin, 1); err != nil {
		return err
	}
	d.disabled = true

	return nil
}

// IsEnabled returns a bool stating whether motor is enabled
func (d *EasyDriver) IsEnabled() bool {
	return !d.disabled
}

// Wake wakes up the driver
func (d *EasyDriver) Wake() error {
	if d.sleepPin == "" {
		return fmt.Errorf("sleepPin is not set for '%s'", d.name)
	}

	// sleepPin is active low
	if err := d.connection.(DigitalWriter).DigitalWrite(d.sleepPin, 1); err != nil {
		return err
	}

	d.sleeping = false

	// we need to wait 1ms after sleeping before doing a step to charge the step pump (according to data sheet)
	time.Sleep(1 * time.Millisecond)

	return nil
}

// IsSleeping returns a bool stating whether motor is sleeping
func (d *EasyDriver) IsSleeping() bool {
	return d.sleeping
}

func (d *EasyDriver) onePinStepping() error {
	// ensure that read and write of variables (direction, stepNum) can not interfere
	d.valueMutex.Lock()
	defer d.valueMutex.Unlock()

	// a valid steps occurs for a low to high transition
	if err := d.connection.(DigitalWriter).DigitalWrite(d.stepPin, 0); err != nil {
		return err
	}

	time.Sleep(d.getDelayPerStep())
	if err := d.connection.(DigitalWriter).DigitalWrite(d.stepPin, 1); err != nil {
		return err
	}

	if d.direction == StepperDriverForward {
		d.stepNum++
	} else {
		d.stepNum--
	}

	return nil
}

// sleepWithSleepPin puts the driver to sleep and disables all motor output.  Low power mode.
func (d *EasyDriver) sleepWithSleepPin() error {
	if d.sleepPin == "" {
		return fmt.Errorf("sleepPin is not set for '%s'", d.name)
	}

	_ = d.stopIfRunning() // drop step errors

	// sleepPin is active low
	if err := d.connection.(DigitalWriter).DigitalWrite(d.sleepPin, 0); err != nil {
		return err
	}
	d.sleeping = true

	return nil
}
