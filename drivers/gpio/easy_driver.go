package gpio

import (
	"errors"
	"strconv"
	"time"

	"gobot.io/x/gobot"
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
func (d *EasyDriver) Start() (err error) { return }

// Halt implements the Driver interface; stops running the stepper
func (d *EasyDriver) Halt() (err error) {
	d.Stop()
	return
}

// Move the motor given number of degrees at current speed.
func (d *EasyDriver) Move(degs int) (err error) {
	if d.moving {
		// don't do anything if already moving
		return
	}

	d.moving = true

	steps := int(float32(degs) / d.angle)
	for i := 0; i < steps; i++ {
		if !d.moving {
			// don't continue to step if driver is stopped
			break
		}

		d.Step()
	}

	d.moving = false

	return
}

// Step the stepper 1 step
func (d *EasyDriver) Step() (err error) {
	stepsPerRev := d.GetMaxSpeed()

	// a valid steps occurs for a low to high transition
	d.connection.DigitalWrite(d.stepPin, 0)
	// 1 minute / steps per revolution / revolutions per minute
	// let's keep it as Microseconds so we only have to do integer math
	time.Sleep(time.Duration(60*1000*1000/stepsPerRev/d.rpm) * time.Microsecond)
	d.connection.DigitalWrite(d.stepPin, 1)

	// increment or decrement the number of steps by 1
	d.stepNum += int(d.dir)

	return
}

// Run the stepper continuously
func (d *EasyDriver) Run() (err error) {
	if d.moving {
		// don't do anything if already moving
		return
	}

	d.moving = true

	go func() {
		for d.moving {
			d.Step()
		}
	}()

	return
}

// Stop running the stepper
func (d *EasyDriver) Stop() (err error) {
	d.moving = false
	return
}

// SetDirection sets the direction to be moving.  Valid directions are "cw" or "ccw"
func (d *EasyDriver) SetDirection(dir string) (err error) {
	// can't change direct if dirPin isn't set
	if d.dirPin == "" {
		return errors.New("dirPin is not set")
	}

	if dir == "ccw" {
		d.dir = -1
		d.connection.DigitalWrite(d.dirPin, 1) // high is ccw
	} else { // default to cw, even if user specified wrong value
		d.dir = 1
		d.connection.DigitalWrite(d.dirPin, 0) // low is cw
	}

	return
}

// SetSpeed sets the speed of the motor in RPMs.  1 is the lowest and GetMaxSpeed is the highest
func (d *EasyDriver) SetSpeed(rpm uint) (err error) {
	if rpm < 1 {
		d.rpm = 1
	} else if rpm > d.GetMaxSpeed() {
		d.rpm = d.GetMaxSpeed()
	} else {
		d.rpm = rpm
	}

	return
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
func (d *EasyDriver) Enable() (err error) {
	// can't enable if enPin isn't set.  This is fine normally since it will be enabled by default
	if d.enPin == "" {
		return errors.New("enPin is not set.  Board is enabled by default")
	}

	d.enabled = true
	d.connection.DigitalWrite(d.enPin, 0) // enPin is active low

	return
}

// Disable disables all motor output
func (d *EasyDriver) Disable() (err error) {
	// can't disable if enPin isn't set
	if d.enPin == "" {
		return errors.New("enPin is not set")
	}

	// let's stop the motor first
	d.Stop()

	d.enabled = false
	d.connection.DigitalWrite(d.enPin, 1) // enPin is active low

	return
}

// IsEnabled returns a bool stating whether motor is enabled
func (d *EasyDriver) IsEnabled() bool {
	return d.enabled
}

// Sleep puts the driver to sleep and disables all motor output.  Low power mode.
func (d *EasyDriver) Sleep() (err error) {
	// can't sleep if sleepPin isn't set
	if d.sleepPin == "" {
		return errors.New("sleepPin is not set")
	}

	// let's stop the motor first
	d.Stop()

	d.sleeping = true
	d.connection.DigitalWrite(d.sleepPin, 0) // sleepPin is active low

	return
}

// Wake wakes up the driver
func (d *EasyDriver) Wake() (err error) {
	// can't wake if sleepPin isn't set
	if d.sleepPin == "" {
		return errors.New("sleepPin is not set")
	}

	d.sleeping = false
	d.connection.DigitalWrite(d.sleepPin, 1) // sleepPin is active low

	// we need to wait 1ms after sleeping before doing a step to charge the step pump (according to data sheet)
	// this will ensure that happens
	time.Sleep(1 * time.Millisecond)

	return
}

// IsSleeping returns a bool stating whether motor is enabled
func (d *EasyDriver) IsSleeping() bool {
	return d.sleeping
}
