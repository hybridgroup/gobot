package gpio

import (
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot/v2"
)

// EasyDriver object
type EasyDriver struct {
	gobot.Commander

	name       string
	connection DigitalWriter
	stepPin    string
	dirPin     string
	enPin      string
	sleepPin   string

	angle    float32
	rpm      uint
	dir      int8
	moving   bool
	stepNum  int
	enabled  bool
	sleeping bool
}

// NewEasyDriver returns a new EasyDriver from SparkFun (https://www.sparkfun.com/products/12779)
// TODO: Support selecting phase input instead of hard-wiring MS1 and MS2 to board truth table
// This should also work for the BigEasyDriver (untested)
// A - DigitalWriter
// stepPin - Pin corresponding to step input on EasyDriver
// dirPin - Pin corresponding to dir input on EasyDriver.  Optional
// enPin - Pin corresponding to enabled input on EasyDriver.  Optional
// sleepPin - Pin corresponding to sleep input on EasyDriver.  Optional
// angle - Step angle of motor
func NewEasyDriver(a DigitalWriter, angle float32, stepPin string, dirPin string, enPin string, sleepPin string) *EasyDriver {
	d := &EasyDriver{
		Commander:  gobot.NewCommander(),
		name:       gobot.DefaultName("EasyDriver"),
		connection: a,
		stepPin:    stepPin,
		dirPin:     dirPin,
		enPin:      enPin,
		sleepPin:   sleepPin,

		angle:    angle,
		rpm:      1,
		dir:      1,
		enabled:  true,
		sleeping: false,
	}

	// panic if step pin isn't set
	if stepPin == "" {
		panic("Step pin is not set")
	}

	// 1/4 of max speed.  Not too fast, not too slow
	d.rpm = d.GetMaxSpeed() / 4

	d.AddCommand("Move", func(params map[string]interface{}) interface{} {
		degs, _ := strconv.Atoi(params["degs"].(string))
		return d.Move(degs)
	})
	d.AddCommand("Step", func(params map[string]interface{}) interface{} {
		return d.Step()
	})
	d.AddCommand("Run", func(params map[string]interface{}) interface{} {
		return d.Run()
	})
	d.AddCommand("Stop", func(params map[string]interface{}) interface{} {
		return d.Stop()
	})

	return d
}

// Name of EasyDriver
func (d *EasyDriver) Name() string { return d.name }

// SetName sets name for EasyDriver
func (d *EasyDriver) SetName(n string) { d.name = n }

// Connection returns EasyDriver's connection
func (d *EasyDriver) Connection() gobot.Connection { return d.connection.(gobot.Connection) }

// Start implements the Driver interface
func (d *EasyDriver) Start() error { return nil }

// Halt implements the Driver interface; stops running the stepper
func (d *EasyDriver) Halt() error {
	return d.Stop()
}

// Move the motor given number of degrees at current speed.
func (d *EasyDriver) Move(degs int) error {
	if d.moving {
		// don't do anything if already moving
		return nil
	}

	d.moving = true

	steps := int(float32(degs) / d.angle)
	for i := 0; i < steps; i++ {
		if !d.moving {
			// don't continue to step if driver is stopped
			break
		}

		if err := d.Step(); err != nil {
			return err
		}
	}

	d.moving = false

	return nil
}

// Step the stepper 1 step
func (d *EasyDriver) Step() error {
	stepsPerRev := d.GetMaxSpeed()

	// a valid steps occurs for a low to high transition
	if err := d.connection.DigitalWrite(d.stepPin, 0); err != nil {
		return err
	}
	// 1 minute / steps per revolution / revolutions per minute
	// let's keep it as Microseconds so we only have to do integer math
	time.Sleep(time.Duration(60*1000*1000/stepsPerRev/d.rpm) * time.Microsecond)
	if err := d.connection.DigitalWrite(d.stepPin, 1); err != nil {
		return err
	}

	// increment or decrement the number of steps by 1
	d.stepNum += int(d.dir)

	return nil
}

// Run the stepper continuously
func (d *EasyDriver) Run() error {
	if d.moving {
		// don't do anything if already moving
		return nil
	}

	d.moving = true

	go func() {
		for d.moving {
			if err := d.Step(); err != nil {
				panic(err)
			}
		}
	}()

	return nil
}

// Stop running the stepper
func (d *EasyDriver) Stop() error {
	d.moving = false
	return nil
}

// SetDirection sets the direction to be moving.  Valid directions are "cw" or "ccw"
func (d *EasyDriver) SetDirection(dir string) error {
	// can't change direct if dirPin isn't set
	if d.dirPin == "" {
		return errors.New("dirPin is not set")
	}

	if dir == "ccw" {
		d.dir = -1
		// high is ccw
		return d.connection.DigitalWrite(d.dirPin, 1)
	}

	// default to cw, even if user specified wrong value
	d.dir = 1
	// low is cw
	return d.connection.DigitalWrite(d.dirPin, 0)
}

// SetSpeed sets the speed of the motor in RPMs.  1 is the lowest and GetMaxSpeed is the highest
func (d *EasyDriver) SetSpeed(rpm uint) error {
	if rpm < 1 {
		d.rpm = 1
	} else if rpm > d.GetMaxSpeed() {
		d.rpm = d.GetMaxSpeed()
	} else {
		d.rpm = rpm
	}

	return nil
}

// GetMaxSpeed returns the max speed of the stepper
func (d *EasyDriver) GetMaxSpeed() uint {
	return uint(360 / d.angle)
}

// GetCurrentStep returns current step number
func (d *EasyDriver) GetCurrentStep() int {
	return d.stepNum
}

// IsMoving returns a bool stating whether motor is currently in motion
func (d *EasyDriver) IsMoving() bool {
	return d.moving
}

// Enable enables all motor output
func (d *EasyDriver) Enable() error {
	// can't enable if enPin isn't set.  This is fine normally since it will be enabled by default
	if d.enPin == "" {
		return errors.New("enPin is not set.  Board is enabled by default")
	}

	// enPin is active low
	if err := d.connection.DigitalWrite(d.enPin, 0); err != nil {
		return err
	}

	d.enabled = true
	return nil
}

// Disable disables all motor output
func (d *EasyDriver) Disable() error {
	// can't disable if enPin isn't set
	if d.enPin == "" {
		return errors.New("enPin is not set")
	}

	// let's stop the motor first, but do not return on error
	err := d.Stop()

	// enPin is active low
	if e := d.connection.DigitalWrite(d.enPin, 1); e != nil {
		err = multierror.Append(err, e)
	} else {
		d.enabled = false
	}

	return err
}

// IsEnabled returns a bool stating whether motor is enabled
func (d *EasyDriver) IsEnabled() bool {
	return d.enabled
}

// Sleep puts the driver to sleep and disables all motor output.  Low power mode.
func (d *EasyDriver) Sleep() error {
	// can't sleep if sleepPin isn't set
	if d.sleepPin == "" {
		return errors.New("sleepPin is not set")
	}

	// let's stop the motor first
	err := d.Stop()

	// sleepPin is active low
	if e := d.connection.DigitalWrite(d.sleepPin, 0); e != nil {
		err = multierror.Append(err, e)
	} else {
		d.sleeping = true
	}

	return err
}

// Wake wakes up the driver
func (d *EasyDriver) Wake() error {
	// can't wake if sleepPin isn't set
	if d.sleepPin == "" {
		return errors.New("sleepPin is not set")
	}

	// sleepPin is active low
	if err := d.connection.DigitalWrite(d.sleepPin, 1); err != nil {
		return err
	}

	d.sleeping = false

	// we need to wait 1ms after sleeping before doing a step to charge the step pump (according to data sheet)
	// this will ensure that happens
	time.Sleep(1 * time.Millisecond)

	return nil
}

// IsSleeping returns a bool stating whether motor is enabled
func (d *EasyDriver) IsSleeping() bool {
	return d.sleeping
}
