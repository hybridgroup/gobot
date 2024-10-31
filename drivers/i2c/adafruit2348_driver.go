package i2c

import (
	"errors"
	"log"
	"sync"
	"time"

	"gobot.io/x/gobot/v2"
)

const adafruit2348Debug = false // Set this to true to see debug output

const (
	adafruit2348MotorHatDefaultAddress = 0x60
	adafruit2348StepperMicrosteps      = 8
)

// Adafruit2348Direction declares a type for specification of the motor direction
type Adafruit2348Direction int

// Adafruit2348StepStyle declares a type for specification of the stepper motor rotation
type Adafruit2348StepStyle int

// constants for running the motor in different direction or release
const (
	Adafruit2348Forward  Adafruit2348Direction = iota // 0
	Adafruit2348Backward                              // 1
	Adafruit2348Release                               // 2
)

// constants for running the stepper motor in different modes
const (
	Adafruit2348Single     Adafruit2348StepStyle = iota // 0
	Adafruit2348Double                                  // 1
	Adafruit2348Interleave                              // 2
	Adafruit2348Microstep                               // 3
)

type adafruit2348DCMotor struct {
	pwmPin, in1Pin, in2Pin byte
}

type adafruit2348StepperMotor struct {
	pwmPinA, pwmPinB      byte
	ain1, ain2            byte
	bin1, bin2            byte
	secPerStep            float64
	currentStep, revSteps int
}

// Adafruit2348Driver is a driver for Adafruit DC and Stepper Motor HAT - a Raspberry Pi add-on, based on PCA9685.
// The HAT can drive up to 4 DC or 2 stepper motors with full PWM speed and direction control over I2C.
// This driver wraps the PCA9685Driver.
// Stacking 32 of them is possible (addresses 0x60..0x80), for controlling up to 64 stepper motors or 128 DC motors.
// datasheet:
// https://cdn-learn.adafruit.com/downloads/pdf/adafruit2348-dc-and-stepper-motor-hat-for-raspberry-pi.pdf
type Adafruit2348Driver struct {
	*PCA9685Driver
	dcMotors              []adafruit2348DCMotor
	stepperMotors         []adafruit2348StepperMotor
	stepperMicrostepCurve []int
	step2coils            map[int][]int32
	stepperSpeedMutex     *sync.Mutex
}

// NewAdafruit2348Driver initializes a new driver for DC and stepper motors.
// Params:
//
//	c Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	    bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
func NewAdafruit2348Driver(c Connector, options ...func(Config)) *Adafruit2348Driver {
	var dc []adafruit2348DCMotor
	var st []adafruit2348StepperMotor
	for i := 0; i < 4; i++ {
		switch {
		case i == 0:
			dc = append(dc, adafruit2348DCMotor{pwmPin: 8, in1Pin: 10, in2Pin: 9})
			st = append(st, adafruit2348StepperMotor{
				pwmPinA: 8, pwmPinB: 13,
				ain1: 10, ain2: 9, bin1: 11, bin2: 12, revSteps: 200, secPerStep: 0.1,
			})
		case i == 1:
			dc = append(dc, adafruit2348DCMotor{pwmPin: 13, in1Pin: 11, in2Pin: 12})
			st = append(st, adafruit2348StepperMotor{
				pwmPinA: 2, pwmPinB: 7,
				ain1: 4, ain2: 3, bin1: 5, bin2: 6, revSteps: 200, secPerStep: 0.1,
			})
		case i == 2:
			dc = append(dc, adafruit2348DCMotor{pwmPin: 2, in1Pin: 4, in2Pin: 3})
		case i == 3:
			dc = append(dc, adafruit2348DCMotor{pwmPin: 7, in1Pin: 5, in2Pin: 6})
		}
	}

	// external given address must be able to override the default one
	o := append([]func(Config){WithAddress(adafruit2348MotorHatDefaultAddress)}, options...)
	pca := NewPCA9685Driver(c, o...)
	pca.SetName(gobot.DefaultName("Adafruit2348MotorHat"))
	d := &Adafruit2348Driver{
		PCA9685Driver:         pca,
		dcMotors:              dc,
		stepperMotors:         st,
		stepperMicrostepCurve: []int{0, 50, 98, 142, 180, 212, 236, 250, 255},
		step2coils: map[int][]int32{
			0: {1, 0, 0, 0},
			1: {1, 1, 0, 0},
			2: {0, 1, 0, 0},
			3: {0, 1, 1, 0},
			4: {0, 0, 1, 0},
			5: {0, 0, 1, 1},
			6: {0, 0, 0, 1},
			7: {1, 0, 0, 1},
		},
		stepperSpeedMutex: &sync.Mutex{},
	}

	// TODO: add API funcs
	return d
}

// SetDCMotorSpeed will set the appropriate pins to run the specified DC motor for the given speed.
func (a *Adafruit2348Driver) SetDCMotorSpeed(dcMotor int, speed int32) error {
	return a.SetPWM(int(a.dcMotors[dcMotor].pwmPin), 0, uint16(speed*16)) //nolint:gosec // TODO: fix later
}

// RunDCMotor will set the appropriate pins to run the specified DC motor for the given direction.
func (a *Adafruit2348Driver) RunDCMotor(dcMotor int, dir Adafruit2348Direction) error {
	switch {
	case dir == Adafruit2348Forward:
		if err := a.setPin(a.dcMotors[dcMotor].in2Pin, 0); err != nil {
			return err
		}
		if err := a.setPin(a.dcMotors[dcMotor].in1Pin, 1); err != nil {
			return err
		}
	case dir == Adafruit2348Backward:
		if err := a.setPin(a.dcMotors[dcMotor].in1Pin, 0); err != nil {
			return err
		}
		if err := a.setPin(a.dcMotors[dcMotor].in2Pin, 1); err != nil {
			return err
		}
	case dir == Adafruit2348Release:
		if err := a.setPin(a.dcMotors[dcMotor].in1Pin, 0); err != nil {
			return err
		}
		if err := a.setPin(a.dcMotors[dcMotor].in2Pin, 0); err != nil {
			return err
		}
	}
	return nil
}

// SetStepperMotorSpeed sets the seconds-per-step for the given stepper motor. It is applied in the next cycle.
func (a *Adafruit2348Driver) SetStepperMotorSpeed(stepperMotor int, rpm int) error {
	a.stepperSpeedMutex.Lock()
	defer a.stepperSpeedMutex.Unlock()

	revSteps := a.stepperMotors[stepperMotor].revSteps
	a.stepperMotors[stepperMotor].secPerStep = 60.0 / float64(revSteps*rpm)
	return nil
}

// Step will rotate the stepper motor the given number of steps, in the given direction and step style.
func (a *Adafruit2348Driver) Step(motor, steps int, dir Adafruit2348Direction, style Adafruit2348StepStyle) error {
	a.stepperSpeedMutex.Lock()
	defer a.stepperSpeedMutex.Unlock()

	secPerStep := a.stepperMotors[motor].secPerStep
	var latestStep int
	var err error
	if style == Adafruit2348Interleave {
		secPerStep = secPerStep / 2.0
	}
	if style == Adafruit2348Microstep {
		secPerStep /= float64(adafruit2348StepperMicrosteps)
		steps *= adafruit2348StepperMicrosteps
	}
	if adafruit2348Debug {
		log.Printf("[adafruit2348_driver] %f seconds per step", secPerStep)
	}
	for i := 0; i < steps; i++ {
		if latestStep, err = a.oneStep(motor, dir, style); err != nil {
			return err
		}
		time.Sleep(time.Duration(secPerStep) * time.Second)
	}
	// As documented in the Adafruit python driver:
	// This is an edge case, if we are in between full steps, keep going to end on a full step
	if style == Adafruit2348Microstep {
		for latestStep != 0 && latestStep != adafruit2348StepperMicrosteps {
			if latestStep, err = a.oneStep(motor, dir, style); err != nil {
				return err
			}
			time.Sleep(time.Duration(secPerStep) * time.Second)
		}
	}
	return nil
}

func (a *Adafruit2348Driver) oneStep(motor int, dir Adafruit2348Direction, style Adafruit2348StepStyle) (int, error) {
	pwmA := 255
	pwmB := 255

	// Determine the stepping procedure
	switch {
	case style == Adafruit2348Single:
		if (a.stepperMotors[motor].currentStep / (adafruit2348StepperMicrosteps / 2) % 2) != 0 {
			// we're at an odd step
			if dir == Adafruit2348Forward {
				a.stepperMotors[motor].currentStep += adafruit2348StepperMicrosteps / 2
			} else {
				a.stepperMotors[motor].currentStep -= adafruit2348StepperMicrosteps / 2
			}
		} else {
			// go to next even step
			if dir == Adafruit2348Forward {
				a.stepperMotors[motor].currentStep += adafruit2348StepperMicrosteps
			} else {
				a.stepperMotors[motor].currentStep -= adafruit2348StepperMicrosteps
			}
		}
	case style == Adafruit2348Double:
		if (a.stepperMotors[motor].currentStep / (adafruit2348StepperMicrosteps / 2) % 2) == 0 {
			// we're at an even step, weird
			if dir == Adafruit2348Forward {
				a.stepperMotors[motor].currentStep += adafruit2348StepperMicrosteps / 2
			} else {
				a.stepperMotors[motor].currentStep -= adafruit2348StepperMicrosteps / 2
			}
		} else {
			// go to next odd step
			if dir == Adafruit2348Forward {
				a.stepperMotors[motor].currentStep += adafruit2348StepperMicrosteps
			} else {
				a.stepperMotors[motor].currentStep -= adafruit2348StepperMicrosteps
			}
		}
	case style == Adafruit2348Interleave:
		if dir == Adafruit2348Forward {
			a.stepperMotors[motor].currentStep += adafruit2348StepperMicrosteps / 2
		} else {
			a.stepperMotors[motor].currentStep -= adafruit2348StepperMicrosteps / 2
		}
	case style == Adafruit2348Microstep:
		if dir == Adafruit2348Forward {
			a.stepperMotors[motor].currentStep++
		} else {
			a.stepperMotors[motor].currentStep--
		}
		// go to next step and wrap around
		a.stepperMotors[motor].currentStep += adafruit2348StepperMicrosteps * 4
		a.stepperMotors[motor].currentStep %= adafruit2348StepperMicrosteps * 4

		pwmA = 0
		pwmB = 0
		currStep := a.stepperMotors[motor].currentStep
		switch {
		case currStep >= 0 && currStep < adafruit2348StepperMicrosteps:
			pwmA = a.stepperMicrostepCurve[adafruit2348StepperMicrosteps-currStep]
			pwmB = a.stepperMicrostepCurve[currStep]
		case currStep >= adafruit2348StepperMicrosteps && currStep < adafruit2348StepperMicrosteps*2:
			pwmA = a.stepperMicrostepCurve[currStep-adafruit2348StepperMicrosteps]
			pwmB = a.stepperMicrostepCurve[adafruit2348StepperMicrosteps*2-currStep]
		case currStep >= adafruit2348StepperMicrosteps*2 && currStep < adafruit2348StepperMicrosteps*3:
			pwmA = a.stepperMicrostepCurve[adafruit2348StepperMicrosteps*3-currStep]
			pwmB = a.stepperMicrostepCurve[currStep-adafruit2348StepperMicrosteps*2]
		case currStep >= adafruit2348StepperMicrosteps*3 && currStep < adafruit2348StepperMicrosteps*4:
			pwmA = a.stepperMicrostepCurve[currStep-adafruit2348StepperMicrosteps*3]
			pwmB = a.stepperMicrostepCurve[adafruit2348StepperMicrosteps*4-currStep]
		}
	} // switch

	// go to next 'step' and wrap around
	a.stepperMotors[motor].currentStep += adafruit2348StepperMicrosteps * 4
	a.stepperMotors[motor].currentStep %= adafruit2348StepperMicrosteps * 4

	// only really used for microstepping, otherwise always on!
	//nolint:gosec // TODO: fix later
	if err := a.SetPWM(int(a.stepperMotors[motor].pwmPinA), 0, uint16(pwmA*16)); err != nil {
		return 0, err
	}
	//nolint:gosec // TODO: fix later
	if err := a.SetPWM(int(a.stepperMotors[motor].pwmPinB), 0, uint16(pwmB*16)); err != nil {
		return 0, err
	}
	var coils []int32
	currStep := a.stepperMotors[motor].currentStep
	if style == Adafruit2348Microstep {
		switch {
		case currStep >= 0 && currStep < adafruit2348StepperMicrosteps:
			coils = []int32{1, 1, 0, 0}
		case currStep >= adafruit2348StepperMicrosteps && currStep < adafruit2348StepperMicrosteps*2:
			coils = []int32{0, 1, 1, 0}
		case currStep >= adafruit2348StepperMicrosteps*2 && currStep < adafruit2348StepperMicrosteps*3:
			coils = []int32{0, 0, 1, 1}
		case currStep >= adafruit2348StepperMicrosteps*3 && currStep < adafruit2348StepperMicrosteps*4:
			coils = []int32{1, 0, 0, 1}
		}
	} else {
		// step-2-coils is initialized in init()
		coils = a.step2coils[(currStep / (adafruit2348StepperMicrosteps / 2))]
	}
	if adafruit2348Debug {
		log.Printf("[adafruit2348_driver] currStep: %d, index into step2coils: %d\n",
			currStep, (currStep / (adafruit2348StepperMicrosteps / 2)))
		log.Printf("[adafruit2348_driver] coils state = %v", coils)
	}
	if err := a.setPin(a.stepperMotors[motor].ain2, coils[0]); err != nil {
		return 0, err
	}
	if err := a.setPin(a.stepperMotors[motor].bin1, coils[1]); err != nil {
		return 0, err
	}
	if err := a.setPin(a.stepperMotors[motor].ain1, coils[2]); err != nil {
		return 0, err
	}
	if err := a.setPin(a.stepperMotors[motor].bin2, coils[3]); err != nil {
		return 0, err
	}
	return a.stepperMotors[motor].currentStep, nil
}

func (a *Adafruit2348Driver) setPin(pin byte, value int32) error {
	if value == 0 {
		return a.SetPWM(int(pin), 0, 4096)
	}
	if value == 1 {
		return a.SetPWM(int(pin), 4096, 0)
	}
	return errors.New("Invalid pin")
}
