package gpio

import (
	"fmt"
	"strings"
	"time"

	"gobot.io/x/gobot/v2"
)

const easyDriverDebug = false

// easyOptionApplier needs to be implemented by each configurable option type
type easyOptionApplier interface {
	apply(cfg *easyConfiguration)
}

// easyConfiguration contains all changeable attributes of the driver.
type easyConfiguration struct {
	dirPin   string
	enPin    string
	sleepPin string
}

// easyDirPinOption is the type for applying a pin for change direction
type easyDirPinOption string

// easyEnPinOption is the type for applying a pin for device disabling/enabling
type easyEnPinOption string

// easySleepPinOption is the type for applying a pin for setting device to sleep/wake
type easySleepPinOption string

// EasyDriver is an driver for stepper hardware board from SparkFun (https://www.sparkfun.com/products/12779)
// This should also work for the BigEasyDriver (untested). It is basically a wrapper for the common StepperDriver{}
// with the specific additions for the board, e.g. direction, enable and sleep outputs.
type EasyDriver struct {
	*StepperDriver
	easyCfg      *easyConfiguration
	stepPin      string
	anglePerStep float32
	sleeping     bool
}

// NewEasyDriver returns a new driver
// TODO: Support selecting phase input instead of hard-wiring MS1 and MS2 to board truth table
// A - DigitalWriter
// anglePerStep - Step angle of motor
// stepPin - Pin corresponding to step input on EasyDriver
//
// Supported options:
//
//	"WithName"
//	"WithEasyDirectionPin"
//	"WithEasyEnablePin"
//	"WithEasySleepPin"
func NewEasyDriver(a DigitalWriter, anglePerStep float32, stepPin string, opts ...interface{}) *EasyDriver {
	if anglePerStep <= 0 {
		panic("angle per step needs to be greater than zero")
	}

	if stepPin == "" {
		panic("step pin is mandatory for easy driver")
	}

	stepper := NewStepperDriver(a, [4]string{}, nil, 1)
	stepper.driverCfg.name = gobot.DefaultName("EasyDriver")
	stepper.stepperDebug = easyDriverDebug
	stepper.haltIfRunning = false
	stepper.stepsPerRev = 360.0 / anglePerStep
	d := &EasyDriver{
		StepperDriver: stepper,
		easyCfg:       &easyConfiguration{},
		stepPin:       stepPin,
		anglePerStep:  anglePerStep,
	}
	d.stepFunc = d.onePinStepping
	d.sleepFunc = d.sleepWithSleepPin
	d.beforeHalt = d.shutdown

	// 1/4 of max speed. Not too fast, not too slow
	d.speedRpm = d.MaxSpeed() / 4

	for _, opt := range opts {
		switch o := opt.(type) {
		case optionApplier:
			o.apply(d.driverCfg)
		case easyOptionApplier:
			o.apply(d.easyCfg)
		default:
			oNames := []string{"WithEasyDirectionPin", "WithEasyEnablePin", "WithEasySleepPin"}
			msg := fmt.Sprintf("'%s' can not be applied on '%s', consider to use one of the options instead: %s",
				opt, d.driverCfg.name, strings.Join(oNames, ", "))
			panic(msg)
		}
	}

	return d
}

// WithEasyDirectionPin configure a pin for change the moving direction.
func WithEasyDirectionPin(pin string) easyOptionApplier {
	return easyDirPinOption(pin)
}

// WithEasyEnablePin configure a pin for disabling/enabling the driver.
func WithEasyEnablePin(pin string) easyOptionApplier {
	return easyEnPinOption(pin)
}

// WithEasySleepPin configure a pin for sleep/wake the driver.
func WithEasySleepPin(pin string) easyOptionApplier {
	return easySleepPinOption(pin)
}

// SetDirection sets the direction to be moving.
func (d *EasyDriver) SetDirection(direction string) error {
	if d.easyCfg.dirPin == "" {
		return fmt.Errorf("dirPin is not set for '%s'", d.driverCfg.name)
	}

	direction = strings.ToLower(direction)
	if direction != StepperDriverForward && direction != StepperDriverBackward {
		return fmt.Errorf("Invalid direction '%s'. Value should be '%s' or '%s'",
			direction, StepperDriverForward, StepperDriverBackward)
	}

	writeVal := byte(0) // low is forward
	if direction == StepperDriverBackward {
		writeVal = 1 // high is backward
	}

	if err := d.digitalWrite(d.easyCfg.dirPin, writeVal); err != nil {
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
	if d.easyCfg.enPin == "" {
		d.disabled = false
		return fmt.Errorf("enPin is not set - board '%s' is enabled by default", d.driverCfg.name)
	}

	// enPin is active low
	if err := d.digitalWrite(d.easyCfg.enPin, 0); err != nil {
		return err
	}

	d.disabled = false
	return nil
}

// Disable disables all motor output
func (d *EasyDriver) Disable() error {
	if d.easyCfg.enPin == "" {
		return fmt.Errorf("enPin is not set for '%s'", d.driverCfg.name)
	}

	_ = d.stopIfRunning() // drop step errors

	// enPin is active low
	if err := d.digitalWrite(d.easyCfg.enPin, 1); err != nil {
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
	if d.easyCfg.sleepPin == "" {
		return fmt.Errorf("sleepPin is not set for '%s'", d.driverCfg.name)
	}

	// sleepPin is active low
	if err := d.digitalWrite(d.easyCfg.sleepPin, 1); err != nil {
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
	if err := d.digitalWrite(d.stepPin, 0); err != nil {
		return err
	}

	time.Sleep(d.getDelayPerStep())
	if err := d.digitalWrite(d.stepPin, 1); err != nil {
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
	if d.easyCfg.sleepPin == "" {
		return fmt.Errorf("sleepPin is not set for '%s'", d.driverCfg.name)
	}

	_ = d.stopIfRunning() // drop step errors

	// sleepPin is active low
	if err := d.digitalWrite(d.easyCfg.sleepPin, 0); err != nil {
		return err
	}
	d.sleeping = true

	return nil
}

func (o easyDirPinOption) String() string {
	return "direction pin option easy driver"
}

func (o easyEnPinOption) String() string {
	return "enable pin option easy driver"
}

func (o easySleepPinOption) String() string {
	return "sleep pin option easy driver"
}

func (o easyDirPinOption) apply(cfg *easyConfiguration) {
	cfg.dirPin = string(o)
}

func (o easyEnPinOption) apply(cfg *easyConfiguration) {
	cfg.enPin = string(o)
}

func (o easySleepPinOption) apply(cfg *easyConfiguration) {
	cfg.sleepPin = string(o)
}
