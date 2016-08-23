package i2c

import (
	"log"
	"time"

	"github.com/hybridgroup/gobot"
)

var (
	_             gobot.Driver = (*AdafruitMotorHatDriver)(nil)
	adafruitDebug              = false // Set this to true to see debug output
)

const motorHatAddress = 0x60

// TODO: consider moving to other file, etc
const (
	// Registers
	_Mode1       = 0x00
	_Mode2       = 0x01
	_SubAdr1     = 0x02
	_SubAdr2     = 0x03
	_SubAdr3     = 0x04
	_Prescale    = 0xFE
	_LedZeroOnL  = 0x06
	_LedZeroOnH  = 0x07
	_LedZeroOffL = 0x08
	_LedZeroOffH = 0x09
	_AllLedOnL   = 0xFA
	_AllLedOnH   = 0xFB
	_AllLedOffL  = 0xFC
	_AllLedOffH  = 0xFD

	// Bits
	_Restart = 0x80
	_Sleep   = 0x10
	_AllCall = 0x01
	_Invrt   = 0x10
	_Outdrv  = 0x04
)

// AdafruitDirection declares a type for specification of the motor direction
type AdafruitDirection int

// AdafruitStepStyle declares a type for specification of the stepper motor rotation
type AdafruitStepStyle int

const (
	AdafruitForward  AdafruitDirection = iota // 0
	AdafruitBackward                          // 1
	AdafruitRelease                           // 2
)
const (
	AdafruitSingle     AdafruitStepStyle = iota // 0
	AdafruitDouble                              // 1
	AdafruitInterleave                          // 2
	AdafruitMicrostep                           // 3
)

var stepperMicrosteps = 8
var stepperMicrostepCurve = []int{0, 50, 98, 142, 180, 212, 236, 250, 255}

type adaFruitDCMotor struct {
	pwmPin, in1Pin, in2Pin byte
}
type adaFruitStepperMotor struct {
	pwmPinA, pwmPinB                   byte
	ain1, ain2                         byte
	bin1, bin2                         byte
	secPerStep                         float64
	currentStep, stepCounter, revSteps int
}

// AdafruitMotorHatDriver is a driver for the DC+Stepper Motor HAT from Adafruit.
// The HAT is a Raspberry Pi add-on that can drive up to 4 DC or 2 Stepper motors
// with full PWM speed control.  It has a dedicated PWM driver chip onboard to
// control both motor direction and speed over I2C.
type AdafruitMotorHatDriver struct {
	name       string
	connection I2c
	gobot.Commander
	dcMotors      []adaFruitDCMotor
	stepperMotors []adaFruitStepperMotor
}

// Name identifies this driver object
func (a *AdafruitMotorHatDriver) Name() string { return a.name }

// Connection identifies the particular adapter object
func (a *AdafruitMotorHatDriver) Connection() gobot.Connection { return a.connection.(gobot.Connection) }

// NewAdafruitMotorHatDriver initializes the internal DCMotor and StepperMotor types.
// Again the Adafruit Motor Hat supports up to four DC motors and up to two stepper motors.
func NewAdafruitMotorHatDriver(a I2c, name string) *AdafruitMotorHatDriver {
	var dc []adaFruitDCMotor
	var st []adaFruitStepperMotor
	for i := 0; i < 4; i++ {
		switch {
		case i == 0:
			dc = append(dc, adaFruitDCMotor{pwmPin: 8, in1Pin: 10, in2Pin: 9})
			st = append(st, adaFruitStepperMotor{pwmPinA: 8, pwmPinB: 13,
				ain1: 10, ain2: 9, bin1: 11, bin2: 12, revSteps: 200, secPerStep: 0.1})
		case i == 1:
			dc = append(dc, adaFruitDCMotor{pwmPin: 13, in1Pin: 11, in2Pin: 12})
			st = append(st, adaFruitStepperMotor{pwmPinA: 2, pwmPinB: 7,
				ain1: 4, ain2: 3, bin1: 5, bin2: 6, revSteps: 200, secPerStep: 0.1})
		case i == 2:
			dc = append(dc, adaFruitDCMotor{pwmPin: 2, in1Pin: 4, in2Pin: 3})
		case i == 3:
			dc = append(dc, adaFruitDCMotor{pwmPin: 7, in1Pin: 5, in2Pin: 6})
		}
	}
	driver := &AdafruitMotorHatDriver{
		name:          name,
		connection:    a,
		Commander:     gobot.NewCommander(),
		dcMotors:      dc,
		stepperMotors: st,
	}
	// TODO: add API funcs?
	return driver
}

// Start initializes the Adafruit Motor HAT driver
func (a *AdafruitMotorHatDriver) Start() (errs []error) {
	if err := a.connection.I2cStart(motorHatAddress); err != nil {
		return []error{err}
	}
	return
}

// Halt returns true if devices is halted successfully
func (a *AdafruitMotorHatDriver) Halt() (errs []error) { return }

//
func (a *AdafruitMotorHatDriver) setPWM(pin byte, on, off int32) (err error) {
	reg := _LedZeroOnL + 4*pin
	val := byte(on & 0xff)
	if err = a.connection.I2cWrite(motorHatAddress, []byte{reg, val}); err != nil {
		return
	}
	reg = _LedZeroOnH + 4*pin
	val = byte(on >> 8)
	if err = a.connection.I2cWrite(motorHatAddress, []byte{reg, val}); err != nil {
		return
	}
	reg = _LedZeroOffL + 4*pin
	val = byte(off & 0xff)
	if err = a.connection.I2cWrite(motorHatAddress, []byte{reg, val}); err != nil {
		return
	}
	reg = _LedZeroOffH + 4*pin
	val = byte(off >> 8)
	if err = a.connection.I2cWrite(motorHatAddress, []byte{reg, val}); err != nil {
		return
	}
	return
}

func (a *AdafruitMotorHatDriver) setPin(pin byte, value int32) (err error) {
	if value == 0 {
		return a.setPWM(pin, 0, 4096)
	}
	if value == 1 {
		return a.setPWM(pin, 4096, 0)
	}
	return nil
}

// SetDCMotorSpeed will set the appropriate pins to run the specified DC motor
// for the given speed.
func (a *AdafruitMotorHatDriver) SetDCMotorSpeed(dcMotor int, speed int32) (err error) {
	if err = a.setPWM(a.dcMotors[dcMotor].pwmPin, 0, speed*16); err != nil {
		return
	}
	return
}

// RunDCMotor will set the appropriate pins to run the specified DC motor for
// the given direction
func (a *AdafruitMotorHatDriver) RunDCMotor(dcMotor int, dir AdafruitDirection) (err error) {

	switch {
	case dir == AdafruitForward:
		if err = a.setPin(a.dcMotors[dcMotor].in2Pin, 0); err != nil {
			return
		}
		if err = a.setPin(a.dcMotors[dcMotor].in1Pin, 1); err != nil {
			return
		}
	case dir == AdafruitBackward:
		if err = a.setPin(a.dcMotors[dcMotor].in1Pin, 0); err != nil {
			return
		}
		if err = a.setPin(a.dcMotors[dcMotor].in2Pin, 1); err != nil {
			return
		}
	case dir == AdafruitRelease:
		if err = a.setPin(a.dcMotors[dcMotor].in1Pin, 0); err != nil {
			return
		}
		if err = a.setPin(a.dcMotors[dcMotor].in2Pin, 0); err != nil {
			return
		}
	}
	return
}
func (a *AdafruitMotorHatDriver) oneStep(motor int, dir AdafruitDirection, style AdafruitStepStyle) (steps int, err error) {
	pwmA := 255
	pwmB := 255

	// Determine the stepping procedure
	switch {
	// TODO: refactor...
	case style == AdafruitSingle:
		if (a.stepperMotors[motor].currentStep / (stepperMicrosteps / 2) % 2) != 0 {
			// we're at an odd step
			if dir == AdafruitForward {
				a.stepperMotors[motor].currentStep += stepperMicrosteps / 2
			} else {
				a.stepperMotors[motor].currentStep -= stepperMicrosteps / 2
			}
		} else {
			// go to next even step
			if dir == AdafruitForward {
				a.stepperMotors[motor].currentStep += stepperMicrosteps
			} else {
				a.stepperMotors[motor].currentStep -= stepperMicrosteps
			}
		}
	case style == AdafruitDouble:
		if (a.stepperMotors[motor].currentStep / (stepperMicrosteps / 2) % 2) == 0 {
			// we're at an even step, weird
			if dir == AdafruitForward {
				a.stepperMotors[motor].currentStep += stepperMicrosteps / 2
			} else {
				a.stepperMotors[motor].currentStep -= stepperMicrosteps / 2
			}
		} else {
			// go to next odd step
			if dir == AdafruitForward {
				a.stepperMotors[motor].currentStep += stepperMicrosteps
			} else {
				a.stepperMotors[motor].currentStep -= stepperMicrosteps
			}
		}
	case style == AdafruitInterleave:
		if dir == AdafruitForward {
			a.stepperMotors[motor].currentStep += stepperMicrosteps / 2
		} else {
			a.stepperMotors[motor].currentStep -= stepperMicrosteps / 2
		}
	case style == AdafruitMicrostep:
		if dir == AdafruitForward {
			a.stepperMotors[motor].currentStep++
		} else {
			a.stepperMotors[motor].currentStep--
		}
		// go to next step and wrap around
		a.stepperMotors[motor].currentStep += stepperMicrosteps * 4
		a.stepperMotors[motor].currentStep %= stepperMicrosteps * 4

		pwmA = 0
		pwmB = 0
		currStep := a.stepperMotors[motor].currentStep
		if currStep >= 0 && currStep < stepperMicrosteps {
			pwmA = stepperMicrostepCurve[stepperMicrosteps-currStep]
			pwmB = stepperMicrostepCurve[currStep]
		} else if currStep >= stepperMicrosteps && currStep < stepperMicrosteps*2 {
			pwmA = stepperMicrostepCurve[currStep-stepperMicrosteps]
			pwmB = stepperMicrostepCurve[stepperMicrosteps*2-currStep]
		} else if currStep >= stepperMicrosteps*2 && currStep < stepperMicrosteps*3 {
			pwmA = stepperMicrostepCurve[stepperMicrosteps*3-currStep]
			pwmB = stepperMicrostepCurve[currStep-stepperMicrosteps*2]
		} else if currStep >= stepperMicrosteps*3 && currStep < stepperMicrosteps*4 {
			pwmA = stepperMicrostepCurve[currStep-stepperMicrosteps*3]
			pwmB = stepperMicrostepCurve[stepperMicrosteps*4-currStep]
		}
	} //switch

	//go to next 'step' and wrap around
	a.stepperMotors[motor].currentStep += stepperMicrosteps * 4
	a.stepperMotors[motor].currentStep %= stepperMicrosteps * 4

	//only really used for microstepping, otherwise always on!
	if err = a.setPWM(a.stepperMotors[motor].pwmPinA, 0, int32(pwmA*16)); err != nil {
		return
	}
	if err = a.setPWM(a.stepperMotors[motor].pwmPinB, 0, int32(pwmB*16)); err != nil {
		return
	}
	var coils = []int32{0, 0, 0, 0}
	currStep := a.stepperMotors[motor].currentStep
	if style == AdafruitMicrostep {
		switch {
		case currStep >= 0 && currStep < stepperMicrosteps:
			coils = []int32{1, 1, 0, 0}
		case currStep >= stepperMicrosteps && currStep < stepperMicrosteps*2:
			coils = []int32{0, 1, 1, 0}
		case currStep >= stepperMicrosteps*2 && currStep < stepperMicrosteps*3:
			coils = []int32{0, 0, 1, 1}
		case currStep >= stepperMicrosteps*3 && currStep < stepperMicrosteps*4:
			coils = []int32{1, 0, 0, 1}
		}
	} else {
		step2coils := make(map[int][]int32)
		step2coils[0] = []int32{1, 0, 0, 0}
		step2coils[1] = []int32{1, 1, 0, 0}
		step2coils[2] = []int32{0, 1, 0, 0}
		step2coils[3] = []int32{0, 1, 1, 0}
		step2coils[4] = []int32{0, 0, 1, 0}
		step2coils[5] = []int32{0, 0, 1, 1}
		step2coils[6] = []int32{0, 0, 0, 1}
		step2coils[7] = []int32{1, 0, 0, 1}
		coils = step2coils[(currStep / (stepperMicrosteps / 2))]
	}
	if adafruitDebug {
		log.Printf("[adafruit_driver] currStep: %d, index into step2coils: %d\n",
			currStep, (currStep / (stepperMicrosteps / 2)))
		log.Printf("[adafruit_driver] coils state = %v", coils)
	}
	if err = a.setPin(a.stepperMotors[motor].ain2, coils[0]); err != nil {
		return
	}
	if err = a.setPin(a.stepperMotors[motor].bin1, coils[1]); err != nil {
		return
	}
	if err = a.setPin(a.stepperMotors[motor].ain1, coils[2]); err != nil {
		return
	}
	if err = a.setPin(a.stepperMotors[motor].bin2, coils[3]); err != nil {
		return
	}
	return a.stepperMotors[motor].currentStep, nil
}

// SetStepperMotorSpeed sets the seconds-per-step for the given Stepper Motor.
func (a *AdafruitMotorHatDriver) SetStepperMotorSpeed(stepperMotor int, rpm int) {
	revSteps := a.stepperMotors[stepperMotor].revSteps
	a.stepperMotors[stepperMotor].secPerStep = 60.0 / float64(revSteps*rpm)
	a.stepperMotors[stepperMotor].stepCounter = 0
}

// Step will rotate the stepper motor the given number of steps, in the given direction and step style.
func (a *AdafruitMotorHatDriver) Step(motor, steps int, dir AdafruitDirection, style AdafruitStepStyle) (err error) {
	secPerStep := a.stepperMotors[motor].secPerStep
	latestStep := 0
	if style == AdafruitInterleave {
		secPerStep = secPerStep / 2.0
	}
	if style == AdafruitMicrostep {
		secPerStep /= float64(stepperMicrosteps)
		steps *= stepperMicrosteps
	}
	if adafruitDebug {
		log.Printf("[adafruit_driver] %f seconds per step", secPerStep)
	}
	for i := 0; i < steps; i++ {
		if latestStep, err = a.oneStep(motor, dir, style); err != nil {
			return
		}
		<-time.After(time.Duration(secPerStep) * time.Second)
	}
	// As documented in the Adafruit python driver:
	// This is an edge case, if we are in between full steps, keep going to end on a full step
	for latestStep != 0 && latestStep != stepperMicrosteps {
		if latestStep, err = a.oneStep(motor, dir, style); err != nil {
			return
		}
		<-time.After(time.Duration(secPerStep) * time.Second)
	}
	return
}
