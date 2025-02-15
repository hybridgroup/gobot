package gpio

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"

	"gobot.io/x/gobot/v2"
)

const (
	stepperDriverDebug = false

	// StepperDriverForward is to set the stepper to run in forward direction (e.g. turn clock wise)
	StepperDriverForward = "forward"
	// StepperDriverBackward is to set the stepper to run in backward direction (e.g. turn counter clock wise)
	StepperDriverBackward = "backward"
)

type phase [][4]byte

// StepperModes to decide on Phase and Stepping
var StepperModes = struct {
	SinglePhaseStepping phase
	DualPhaseStepping   phase
	HalfStepping        phase
}{
	// 1 cycle = 4 steps with lesser torque
	SinglePhaseStepping: phase{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	},
	// 1 cycle = 4 steps with higher torque and current
	DualPhaseStepping: phase{
		{1, 0, 0, 1},
		{1, 1, 0, 0},
		{0, 1, 1, 0},
		{0, 0, 1, 1},
	},
	// 1 cycle = 8 steps with lesser torque than full stepping
	HalfStepping: phase{
		{1, 0, 0, 1},
		{1, 0, 0, 0},
		{1, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 1, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 1},
		{0, 0, 0, 1},
	},
}

// StepperDriver is a common driver for stepper motors. It supports 3 different stepping modes.
type StepperDriver struct {
	*driver

	pins        [4]string
	phase       phase
	stepsPerRev float32

	stepperDebug   bool
	speedRpm       uint
	direction      string
	skipStepErrors bool
	haltIfRunning  bool // stop automatically if run is called
	disabled       bool
	valueMutex     *sync.Mutex // to ensure that read and write of values do not interfere

	stepFunc          func() error
	sleepFunc         func() error
	stepNum           int
	stopAsynchRunFunc func(bool) error
}

// NewStepperDriver returns a new StepperDriver given a DigitalWriter
// Pins - To which the stepper is connected
// Phase - Defined by StepperModes {SinglePhaseStepping, DualPhaseStepping, HalfStepping}
// Steps - No of steps per revolution of Stepper motor
//
// Supported options:
//
//	"WithName"
func NewStepperDriver(
	a DigitalWriter,
	pins [4]string,
	phase phase,
	stepsPerRev uint,
	opts ...interface{},
) *StepperDriver {
	if stepsPerRev <= 0 {
		panic("steps per revolution needs to be greater than zero")
	}
	//nolint:forcetypeassert // no error return value, so there is no better way
	d := &StepperDriver{
		driver:         newDriver(a.(gobot.Connection), "Stepper", opts...),
		pins:           pins,
		phase:          phase,
		stepsPerRev:    float32(stepsPerRev),
		stepperDebug:   stepperDriverDebug,
		skipStepErrors: false,
		haltIfRunning:  true,
		direction:      StepperDriverForward,
		stepNum:        0,
		speedRpm:       1,
		valueMutex:     &sync.Mutex{},
	}
	d.speedRpm = d.MaxSpeed()
	d.stepFunc = d.phasedStepping
	d.sleepFunc = d.sleepOuputs
	d.beforeHalt = d.shutdown

	//nolint:forcetypeassert // ok here
	d.AddCommand("MoveDeg", func(params map[string]interface{}) interface{} {
		degs, _ := strconv.Atoi(params["degs"].(string))
		return d.MoveDeg(degs)
	})
	//nolint:forcetypeassert // ok here
	d.AddCommand("Move", func(params map[string]interface{}) interface{} {
		steps, _ := strconv.Atoi(params["steps"].(string))
		return d.Move(steps)
	})
	d.AddCommand("Step", func(_ map[string]interface{}) interface{} {
		return d.Move(1)
	})
	d.AddCommand("Run", func(_ map[string]interface{}) interface{} {
		return d.Run()
	})
	d.AddCommand("Sleep", func(_ map[string]interface{}) interface{} {
		return d.Sleep()
	})
	d.AddCommand("Stop", func(_ map[string]interface{}) interface{} {
		return d.Stop()
	})
	d.AddCommand("Halt", func(_ map[string]interface{}) interface{} {
		return d.Halt()
	})

	return d
}

// Move moves the motor for given number of steps.
func (d *StepperDriver) Move(stepsToMove int) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.stepAsynch(float64(stepsToMove)); err != nil {
		// something went wrong with preparation
		return err
	}

	err := d.stopAsynchRunFunc(false) // wait to finish with err or nil
	d.stopAsynchRunFunc = nil

	return err
}

// MoveDeg moves the motor given number of degrees at current speed. Negative values cause to move backward.
func (d *StepperDriver) MoveDeg(degs int) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	stepsToMove := float64(degs) * float64(d.stepsPerRev) / 360

	if err := d.stepAsynch(stepsToMove); err != nil {
		// something went wrong with preparation
		return err
	}

	err := d.stopAsynchRunFunc(false) // wait to finish with err or nil
	d.stopAsynchRunFunc = nil

	return err
}

// Run runs the stepper continuously. Stop needs to be done with call Stop().
func (d *StepperDriver) Run() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.stepAsynch(float64(math.MaxInt) + 1)
}

// IsMoving returns a bool stating whether motor is currently in motion
func (d *StepperDriver) IsMoving() bool {
	return d.stopAsynchRunFunc != nil
}

// Stop running the stepper
func (d *StepperDriver) Stop() error {
	if d.stopAsynchRunFunc == nil {
		return fmt.Errorf("'%s' is not yet started", d.driverCfg.name)
	}

	err := d.stopAsynchRunFunc(true)
	d.stopAsynchRunFunc = nil

	return err
}

// Sleep release all pins to the same output level, so no current is consumed anymore.
func (d *StepperDriver) Sleep() error {
	return d.sleepFunc()
}

// SetDirection sets the direction in which motor should be moving, default is forward.
// Changing the direction affects the next step, also for asynchronous running.
func (d *StepperDriver) SetDirection(direction string) error {
	direction = strings.ToLower(direction)
	if direction != StepperDriverForward && direction != StepperDriverBackward {
		return fmt.Errorf("Invalid direction '%s'. Value should be '%s' or '%s'",
			direction, StepperDriverForward, StepperDriverBackward)
	}

	// ensure that write of variable can not interfere with read in step()
	d.valueMutex.Lock()
	defer d.valueMutex.Unlock()
	d.direction = direction

	return nil
}

// MaxSpeed gives the max RPM of motor
// max. speed is limited by:
// * motor friction, inertia and inductance, load inertia
// * full step rate is normally below 1000 per second (1kHz), typically not more than ~400 per second
// * mostly not more than 1000-2000rpm (20-40 revolutions per second) are possible
// * higher values can be achieved only by ramp-up the velocity
// * duration of GPIO write (PI1 can reach up to 70kHz, typically 20kHz, so this is most likely not the limiting factor)
// * the hardware driver, to force the high current transitions for the max. speed
// * there are CNC steppers with 1000..20.000 steps per revolution, which works with faster step rates (e.g. 200kHz)
func (d *StepperDriver) MaxSpeed() uint {
	const maxStepsPerSecond = 700 // a typical value for a normal, lightly loaded motor
	return uint(float32(60*maxStepsPerSecond) / d.stepsPerRev)
}

// SetSpeed sets the rpm for the next move or run. A valid value is between 1 and MaxSpeed().
// The run needs to be stopped and called again after set this value.
func (d *StepperDriver) SetSpeed(rpm uint) error {
	var err error
	if rpm <= 0 {
		rpm = 0
		err = fmt.Errorf("RPM (%d) cannot be a zero or negative value", rpm)
	}

	maxRpm := d.MaxSpeed()
	if rpm > maxRpm {
		rpm = maxRpm
		err = fmt.Errorf("RPM (%d) cannot be greater then maximal value %d", rpm, maxRpm)
	}

	d.valueMutex.Lock()
	defer d.valueMutex.Unlock()
	d.speedRpm = rpm

	return err
}

// CurrentStep gives the current step of motor
func (d *StepperDriver) CurrentStep() int {
	// ensure that read can not interfere with write in step()
	d.valueMutex.Lock()
	defer d.valueMutex.Unlock()

	return d.stepNum
}

// SetHaltIfRunning with the given value. Normally a call of Run() returns an error if already running. If set this
// to true, the next call of Run() cause a automatic stop before.
func (d *StepperDriver) SetHaltIfRunning(val bool) {
	d.haltIfRunning = val
}

// shutdown the driver
func (d *StepperDriver) shutdown() error {
	// stops the continuous motion of the stepper, if running
	return d.stopIfRunning()
}

func (d *StepperDriver) stepAsynch(stepsToMove float64) error {
	if d.disabled {
		return fmt.Errorf("'%s' is disabled and can not be running or moving", d.driverCfg.name)
	}

	// if running, return error or stop automatically
	if d.stopAsynchRunFunc != nil {
		if !d.haltIfRunning {
			return fmt.Errorf("'%s' already running or moving", d.driverCfg.name)
		}
		d.debug("stop former run forcefully")
		if err := d.stopAsynchRunFunc(true); err != nil {
			d.stopAsynchRunFunc = nil
			return err
		}
	}

	// prepare stepping behavior
	stepsLeft := uint64(math.Abs(stepsToMove))
	if stepsLeft == 0 {
		return fmt.Errorf("no steps to do for '%s'", d.driverCfg.name)
	}

	// t [min] = steps [st] / (steps_per_revolution [st/u] * speed [u/min]) or
	// t [min] = steps [st] * delay_per_step [min/st], use safety factor 2 and a small offset of 100 ms
	// prepare this timeout outside of stop function to prevent data race with stepsLeft
	//nolint:gosec // TODO: fix later
	stopTimeout := time.Duration(2*stepsLeft)*d.getDelayPerStep() + 100*time.Millisecond
	endlessMovement := false

	if stepsLeft > math.MaxInt {
		stopTimeout = 100 * time.Millisecond
		endlessMovement = true
	} else {
		d.direction = "forward"
		if stepsToMove < 0 {
			d.direction = "backward"
		}
	}

	// prepare new asynchronous stepping
	onceDoneChan := make(chan struct{})
	runStopChan := make(chan struct{})
	runErrChan := make(chan error)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	d.stopAsynchRunFunc = func(forceStop bool) error {
		defer func() {
			d.debug("RUN: cleanup stop channel")
			if runStopChan != nil {
				close(runStopChan)
			}
			runStopChan = nil
			d.debug("STOP: cleanup err channel")
			if runErrChan != nil {
				close(runErrChan)
			}
			runErrChan = nil
			d.debug("STOP: cleanup done")
		}()

		d.debug("STOP: wait for once done")
		<-onceDoneChan // wait for the first step was called

		// send stop for endless movement or a forceful stop happen
		if endlessMovement || forceStop {
			d.debug("STOP: send stop channel")
			runStopChan <- struct{}{}
		}

		if !endlessMovement && forceStop {
			// do not wait if an normal movement was stopped forcefully
			log.Printf("'%s' was forcefully stopped\n", d.driverCfg.name)
			return nil
		}

		// wait for go routine is finished and cleanup
		d.debug(fmt.Sprintf("STOP: wait %s for err channel", stopTimeout))
		select {
		case err := <-runErrChan:
			return err
		case <-time.After(stopTimeout):
			return fmt.Errorf("'%s' was not finished in %s", d.driverCfg.name, stopTimeout)
		}
	}

	d.debug(fmt.Sprintf("going to start go routine - endless=%t, steps=%d", endlessMovement, stepsLeft))
	go func(name string) {
		var err error
		var onceDone bool
		defer func() {
			// some cases here:
			// * stop by stop channel: error should be send as nil
			// * count of steps reached: error should be send as nil
			// * write error occurred
			//    * for Run(): caller needs to send stop channel and read the error
			//    * for Move(): caller waits for the error, but don't send stop channel
			//
			d.debug(fmt.Sprintf("RUN: write '%v' to err channel", err))
			runErrChan <- err
		}()
		for stepsLeft > 0 {
			select {
			case <-sigChan:
				d.debug("RUN: OS signal received")
				err = fmt.Errorf("OS signal received")
				return
			case <-runStopChan:
				d.debug("RUN: stop channel received")
				return
			default:
				if err == nil {
					err = d.stepFunc()
					if err != nil {
						if d.skipStepErrors {
							fmt.Printf("step skipped for '%s': %v\n", name, err)
							err = nil
						} else {
							d.debug("RUN: write error occurred")
						}
					}
					if !onceDone {
						close(onceDoneChan) // to inform that we are ready for stop now
						onceDone = true
						d.debug("RUN: once done")
					}
					if !endlessMovement {
						if err != nil {
							return
						}
						stepsLeft--
					}
				}
			}
		}
	}(d.driverCfg.name)

	return nil
}

// getDelayPerStep gives the delay per step
// formula: delay_per_step [min] = 1/(steps_per_revolution * speed [rpm])
func (d *StepperDriver) getDelayPerStep() time.Duration {
	// considering a max. speed of 1000 rpm and max. 1000 steps per revolution, a microsecond resolution is needed
	// if the motor or application needs bigger values, switch to nanosecond is needed
	return time.Duration(60*1000*1000/(d.stepsPerRev*float32(d.speedRpm))) * time.Microsecond
}

// phasedStepping moves the motor one step with the configured speed and direction. The speed can be adjusted
// by SetSpeed() and the direction can be changed by SetDirection() asynchronously.
func (d *StepperDriver) phasedStepping() error {
	// ensure that read and write of variables (direction, stepNum) can not interfere
	d.valueMutex.Lock()
	defer d.valueMutex.Unlock()

	oldStepNum := d.stepNum

	if d.direction == StepperDriverForward {
		d.stepNum++
	} else {
		d.stepNum--
	}

	if d.stepNum >= int(d.stepsPerRev) {
		d.stepNum = 0
	} else if d.stepNum < 0 {
		d.stepNum = int(d.stepsPerRev) - 1
	}

	r := int(math.Abs(float64(d.stepNum))) % len(d.phase)

	for i, v := range d.phase[r] {
		if err := d.digitalWrite(d.pins[i], v); err != nil {
			d.stepNum = oldStepNum
			return err
		}
	}

	delay := d.getDelayPerStep()
	time.Sleep(delay)

	return nil
}

func (d *StepperDriver) sleepOuputs() error {
	for _, pin := range d.pins {
		if err := d.digitalWrite(pin, 0); err != nil {
			return err
		}
	}
	return nil
}

// stopIfRunning stop the stepper if moving or running
func (d *StepperDriver) stopIfRunning() error {
	// stops the continuous motion of the stepper, if running
	if d.stopAsynchRunFunc == nil {
		return nil
	}

	err := d.stopAsynchRunFunc(true)
	d.stopAsynchRunFunc = nil

	return err
}

func (d *StepperDriver) debug(text string) {
	if d.stepperDebug {
		fmt.Println(text)
	}
}
