package gpio

import (
	"errors"
	"math"
	"strings"
	"sync"
	"time"

	"gobot.io/x/gobot"
)

type phase [][4]byte

// StepperModes to decide on Phase and Stepping
var StepperModes = struct {
	SinglePhaseStepping [][4]byte
	DualPhaseStepping   [][4]byte
	HalfStepping        [][4]byte
}{
	//1 cycle = 4 steps with lesser torque
	SinglePhaseStepping: [][4]byte{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	},
	//1 cycle = 4 steps with higher torque and current
	DualPhaseStepping: [][4]byte{
		{1, 0, 0, 1},
		{1, 1, 0, 0},
		{0, 1, 1, 0},
		{0, 0, 1, 1},
	},
	//1 cycle = 8 steps with lesser torque than full stepping
	HalfStepping: [][4]byte{
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

// StepperDriver object
type StepperDriver struct {
	name        string
	pins        [4]string
	connection  DigitalWriter
	phase       phase
	stepsPerRev uint
	moving      bool
	direction   string
	stepNum     int
	speed       uint
	mutex       *sync.Mutex
}

// NewStepperDriver returns a new StepperDriver given a
// DigitalWriter
// Pins - To which the stepper is connected
// Phase - Defined by StepperModes {SinglePhaseStepping, DualPhaseStepping, HalfStepping}
// Steps - No of steps per revolution of Stepper motor
func NewStepperDriver(a DigitalWriter, pins [4]string, phase phase, stepsPerRev uint) *StepperDriver {
	s := &StepperDriver{
		name:        gobot.DefaultName("Stepper"),
		connection:  a,
		pins:        pins,
		phase:       phase,
		stepsPerRev: stepsPerRev,
		moving:      false,
		direction:   "forward",
		stepNum:     0,
		speed:       1,
		mutex:       &sync.Mutex{},
	}

	s.speed = s.GetMaxSpeed()
	return s
}

// Name of StepperDriver
func (s *StepperDriver) Name() string { return s.name }

// SetName sets name for StepperDriver
func (s *StepperDriver) SetName(n string) { s.name = n }

// Connection returns StepperDriver's connection
func (s *StepperDriver) Connection() gobot.Connection { return s.connection.(gobot.Connection) }

// Start implements the Driver interface and keeps running the stepper till halt is called
func (s *StepperDriver) Start() (err error) { return }

// Run continuously runs the stepper
func (s *StepperDriver) Run() (err error) {
	//halt if already moving
	if s.moving == true {
		s.Halt()
	}

	s.mutex.Lock()
	s.moving = true
	s.mutex.Unlock()

	go func() {
		for {
			if s.moving == false {
				break
			}
			s.step()
		}
	}()

	return
}

// Halt implements the Driver interface and halts the motion of the Stepper
func (s *StepperDriver) Halt() (err error) {
	s.mutex.Lock()
	s.moving = false
	s.mutex.Unlock()
	return nil
}

// SetDirection sets the direction in which motor should be moving, Default is forward
func (s *StepperDriver) SetDirection(direction string) error {
	direction = strings.ToLower(direction)
	if direction != "forward" && direction != "backward" {
		return errors.New("Invalid direction. Value should be forward or backward")
	}

	s.mutex.Lock()
	s.direction = direction
	s.mutex.Unlock()
	return nil
}

// IsMoving returns a bool stating whether motor is currently in motion
func (s *StepperDriver) IsMoving() bool {
	return s.moving
}

// Step moves motor one step in giving direction
func (s *StepperDriver) step() error {
	if s.direction == "forward" {
		s.stepNum++
	} else {
		s.stepNum--
	}

	if s.stepNum >= int(s.stepsPerRev) {
		s.stepNum = 0
	} else if s.stepNum < 0 {
		s.stepNum = int(s.stepsPerRev) - 1
	}

	r := int(math.Abs(float64(s.stepNum))) % len(s.phase)

	for i, v := range s.phase[r] {
		if err := s.connection.DigitalWrite(s.pins[i], v); err != nil {
			return err
		}
	}

	return nil
}

// Move moves the motor for given number of steps
func (s *StepperDriver) Move(stepsToMove int) error {
	if stepsToMove == 0 {
		return s.Halt()
	}

	if s.moving == true {
		//stop previous motion
		s.Halt()
	}

	s.mutex.Lock()
	s.moving = true
	s.direction = "forward"

	if stepsToMove < 0 {
		s.direction = "backward"
	}
	s.mutex.Unlock()

	stepsLeft := int64(math.Abs(float64(stepsToMove)))
	//Do not remove *1000 and change duration to time.Millisecond. It has been done for a reason
	delay := time.Duration(60000*1000/(s.stepsPerRev*s.speed)) * time.Microsecond

	for stepsLeft > 0 {
		if err := s.step(); err != nil {
			return err
		}
		stepsLeft--
		time.Sleep(delay)
	}

	s.moving = false
	return nil
}

// GetCurrentStep gives the current step of motor
func (s *StepperDriver) GetCurrentStep() int {
	return s.stepNum
}

// GetMaxSpeed gives the max RPM of motor
func (s *StepperDriver) GetMaxSpeed() uint {
	//considering time for 1 rev as no of steps per rev * 1.5 (min time req between each step)
	return uint(60000 / (float64(s.stepsPerRev) * 1.5))
}

// SetSpeed sets the rpm
func (s *StepperDriver) SetSpeed(rpm uint) error {
	if rpm <= 0 {
		return errors.New("RPM cannot be a zero or negative value")
	}

	m := s.GetMaxSpeed()
	if rpm > m {
		rpm = m
	}

	s.speed = rpm
	return nil
}
