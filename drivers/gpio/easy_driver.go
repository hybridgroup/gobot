package gpio

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot/v2"
)

// EasyDriver object
type EasyDriver struct {
	*Driver

	stepPin  string
	dirPin   string
	enPin    string
	sleepPin string

	angle            float32
	rpm              uint
	dir              int8
	stepNum          int
	enabled          bool
	sleeping         bool
	runStopChan      chan struct{}
	runStopWaitGroup *sync.WaitGroup
}

// NewEasyDriver returns a new driver for EasyDriver from SparkFun (https://www.sparkfun.com/products/12779)
// TODO: Support selecting phase input instead of hard-wiring MS1 and MS2 to board truth table
// This should also work for the BigEasyDriver (untested)
// A - DigitalWriter
// stepPin - Pin corresponding to step input on EasyDriver
// dirPin - Pin corresponding to dir input on EasyDriver.  Optional
// enPin - Pin corresponding to enabled input on EasyDriver.  Optional
// sleepPin - Pin corresponding to sleep input on EasyDriver.  Optional
// angle - Step angle of motor
func NewEasyDriver(
	a DigitalWriter,
	angle float32,
	stepPin string,
	dirPin string,
	enPin string,
	sleepPin string,
) *EasyDriver {
	if angle <= 0 {
		panic("angle needs to be greater than zero")
	}
	d := &EasyDriver{
		Driver:   NewDriver(a.(gobot.Connection), "EasyDriver"),
		stepPin:  stepPin,
		dirPin:   dirPin,
		enPin:    enPin,
		sleepPin: sleepPin,
		angle:    angle,
		rpm:      1,
		dir:      1,
		enabled:  true,
		sleeping: false,
	}
	d.beforeHalt = func() error {
		if err := d.Stop(); err != nil {
			fmt.Printf("no need to stop motion: %v\n", err)
		}

		return nil
	}

	// panic if step pin isn't set
	if stepPin == "" {
		panic("Step pin is not set")
	}

	// 1/4 of max speed.  Not too fast, not too slow
	d.rpm = d.MaxSpeed() / 4

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

// Move the motor given number of degrees at current speed. The move can be stopped asynchronously.
func (d *EasyDriver) Move(degs int) error {
	// ensure that move and run can not interfere
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if !d.enabled {
		return fmt.Errorf("motor '%s' is disabled and can not be running", d.name)
	}

	if d.runStopChan != nil {
		return fmt.Errorf("motor '%s' already running or moving", d.name)
	}

	d.runStopChan = make(chan struct{})
	d.runStopWaitGroup = &sync.WaitGroup{}
	d.runStopWaitGroup.Add(1)

	defer func() {
		close(d.runStopChan)
		d.runStopChan = nil
		d.runStopWaitGroup.Done()
	}()

	steps := int(float32(degs) / d.angle)
	if steps <= 0 {
		fmt.Printf("steps are smaller than zero, no move for '%s'\n", d.name)
	}

	for i := 0; i < steps; i++ {
		select {
		case <-d.runStopChan:
			// don't continue to step if driver is stopped
			log.Println("stop happen")
			return nil
		default:
			if err := d.step(); err != nil {
				return err
			}
		}
	}

	return nil
}

// Run the stepper continuously.
func (d *EasyDriver) Run() error {
	// ensure that run, can not interfere with step or move
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if !d.enabled {
		return fmt.Errorf("motor '%s' is disabled and can not be moving", d.name)
	}

	if d.runStopChan != nil {
		return fmt.Errorf("motor '%s' already running or moving", d.name)
	}

	d.runStopChan = make(chan struct{})
	d.runStopWaitGroup = &sync.WaitGroup{}
	d.runStopWaitGroup.Add(1)

	go func(name string) {
		defer d.runStopWaitGroup.Done()
		for {
			select {
			case <-d.runStopChan:
				d.runStopChan = nil
				return
			default:
				if err := d.step(); err != nil {
					fmt.Printf("motor step skipped for '%s': %v\n", name, err)
				}
			}
		}
	}(d.name)

	return nil
}

// IsMoving returns a bool stating whether motor is currently in motion
func (d *EasyDriver) IsMoving() bool {
	return d.runStopChan != nil
}

// Stop running the stepper
func (d *EasyDriver) Stop() error {
	if !d.IsMoving() {
		return fmt.Errorf("motor '%s' is not yet started", d.name)
	}

	d.runStopChan <- struct{}{}
	d.runStopWaitGroup.Wait()

	return nil
}

// Step the stepper 1 step
func (d *EasyDriver) Step() error {
	// ensure that move and step can not interfere
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.IsMoving() {
		return fmt.Errorf("motor '%s' already running or moving", d.name)
	}

	return d.step()
}

// SetDirection sets the direction to be moving.  Valid directions are "cw" or "ccw"
func (d *EasyDriver) SetDirection(dir string) error {
	// can't change direct if dirPin isn't set
	if d.dirPin == "" {
		return fmt.Errorf("dirPin is not set for '%s'", d.name)
	}

	if dir == "ccw" {
		d.dir = -1
		// high is ccw
		return d.connection.(DigitalWriter).DigitalWrite(d.dirPin, 1)
	}

	// default to cw, even if user specified wrong value
	d.dir = 1
	// low is cw
	return d.connection.(DigitalWriter).DigitalWrite(d.dirPin, 0)
}

// SetSpeed sets the speed of the motor in RPMs. 1 is the lowest and GetMaxSpeed is the highest
func (d *EasyDriver) SetSpeed(rpm uint) error {
	if rpm < 1 {
		d.rpm = 1
	} else if rpm > d.MaxSpeed() {
		d.rpm = d.MaxSpeed()
	} else {
		d.rpm = rpm
	}

	return nil
}

// MaxSpeed returns the max speed of the stepper
func (d *EasyDriver) MaxSpeed() uint {
	return uint(360 / d.angle)
}

// CurrentStep returns current step number
func (d *EasyDriver) CurrentStep() int {
	return d.stepNum
}

// Enable enables all motor output
func (d *EasyDriver) Enable() error {
	// can't enable if enPin isn't set.  This is fine normally since it will be enabled by default
	if d.enPin == "" {
		d.enabled = true
		return fmt.Errorf("enPin is not set - board '%s' is enabled by default", d.name)
	}

	// enPin is active low
	if err := d.connection.(DigitalWriter).DigitalWrite(d.enPin, 0); err != nil {
		return err
	}

	d.enabled = true
	return nil
}

// Disable disables all motor output
func (d *EasyDriver) Disable() error {
	// can't disable if enPin isn't set
	if d.enPin == "" {
		return fmt.Errorf("enPin is not set for '%s'", d.name)
	}

	// stop the motor if running
	err := d.tryStop()

	// enPin is active low
	if e := d.connection.(DigitalWriter).DigitalWrite(d.enPin, 1); e != nil {
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
		return fmt.Errorf("sleepPin is not set for '%s'", d.name)
	}

	// stop the motor if running
	err := d.tryStop()

	// sleepPin is active low
	if e := d.connection.(DigitalWriter).DigitalWrite(d.sleepPin, 0); e != nil {
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

func (d *EasyDriver) step() error {
	stepsPerRev := d.MaxSpeed()

	// a valid steps occurs for a low to high transition
	if err := d.connection.(DigitalWriter).DigitalWrite(d.stepPin, 0); err != nil {
		return err
	}
	// 1 minute / steps per revolution / revolutions per minute
	// let's keep it as Microseconds so we only have to do integer math
	time.Sleep(time.Duration(60*1000*1000/stepsPerRev/d.rpm) * time.Microsecond)
	if err := d.connection.(DigitalWriter).DigitalWrite(d.stepPin, 1); err != nil {
		return err
	}

	// increment or decrement the number of steps by 1
	d.stepNum += int(d.dir)

	return nil
}

// tryStop stop the stepper if moving or running
func (d *EasyDriver) tryStop() error {
	if !d.IsMoving() {
		return nil
	}

	return d.Stop()
}
