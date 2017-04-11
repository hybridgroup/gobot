package i2c

import (
	"errors"
	"log"
	"math"
	"time"

	"gobot.io/x/gobot"
)

// AdafruitDirection declares a type for specification of the motor direction
type AdafruitDirection int

// AdafruitStepStyle declares a type for specification of the stepper motor rotation
type AdafruitStepStyle int

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
	name               string
	connector          Connector
	motorHatConnection Connection
	servoHatConnection Connection
	Config
	gobot.Commander
	dcMotors      []adaFruitDCMotor
	stepperMotors []adaFruitStepperMotor
}

var adafruitDebug = false // Set this to true to see debug output

var (
	// Each Adafruit HAT must have a unique I2C address. The default address for
	// the DC and Stepper Motor HAT is 0x60. The addresses of the Motor HATs can
	// range from 0x60 to 0x80 for a total of 32 unique addresses.
	// The base address for the Adafruit PWM-Servo HAT is 0x40.  Please consult
	// the Adafruit documentation for soldering and addressing stacked HATs.
	motorHatAddress       = 0x60
	servoHatAddress       = 0x40
	stepperMicrosteps     = 8
	stepperMicrostepCurve []int
	step2coils            = make(map[int][]int32)
)

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

// NewAdafruitMotorHatDriver initializes the internal DCMotor and StepperMotor types.
// Again the Adafruit Motor Hat supports up to four DC motors and up to two stepper motors.
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewAdafruitMotorHatDriver(conn Connector, options ...func(Config)) *AdafruitMotorHatDriver {
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
		name:          gobot.DefaultName("AdafruitMotorHat"),
		connector:     conn,
		Config:        NewConfig(),
		Commander:     gobot.NewCommander(),
		dcMotors:      dc,
		stepperMotors: st,
	}

	for _, option := range options {
		option(driver)
	}

	// TODO: add API funcs
	return driver
}

// SetMotorHatAddress sets the I2C address for the DC and Stepper Motor HAT.
// This addressing flexibility empowers "stacking" the HATs.
func (a *AdafruitMotorHatDriver) SetMotorHatAddress(addr int) (err error) {
	motorHatAddress = addr
	return
}

// SetServoHatAddress sets the I2C address for the PWM-Servo Motor HAT.
// This addressing flexibility empowers "stacking" the HATs.
func (a *AdafruitMotorHatDriver) SetServoHatAddress(addr int) (err error) {
	servoHatAddress = addr
	return
}

// Name identifies this driver object
func (a *AdafruitMotorHatDriver) Name() string { return a.name }

// SetName sets nae for driver
func (a *AdafruitMotorHatDriver) SetName(n string) { a.name = n }

// Connection identifies the particular adapter object
func (a *AdafruitMotorHatDriver) Connection() gobot.Connection { return a.connector.(gobot.Connection) }

func (a *AdafruitMotorHatDriver) startDriver(connection Connection) (err error) {
	if err = a.setAllPWM(connection, 0, 0); err != nil {
		return
	}
	reg := byte(_Mode2)
	val := byte(_Outdrv)
	if _, err = connection.Write([]byte{reg, val}); err != nil {
		return
	}
	reg = byte(_Mode1)
	val = byte(_AllCall)
	if _, err = connection.Write([]byte{reg, val}); err != nil {
		return
	}
	time.Sleep(5 * time.Millisecond)

	// Read a byte from the I2C device.  Note: no ability to read from a specified reg?
	mode1 := []byte{0}
	_, rerr := connection.Read(mode1)
	if rerr != nil {
		return rerr
	}
	if len(mode1) > 0 {
		reg = byte(_Mode1)
		val = mode1[0] & _Sleep
		if _, err = connection.Write([]byte{reg, val}); err != nil {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}

	return
}

// Start initializes both I2C-addressable Adafruit Motor HAT drivers
func (a *AdafruitMotorHatDriver) Start() (err error) {
	bus := a.GetBusOrDefault(a.connector.GetDefaultBus())

	if a.servoHatConnection, err = a.connector.GetConnection(servoHatAddress, bus); err != nil {
		return
	}

	if err = a.startDriver(a.servoHatConnection); err != nil {
		return
	}

	if a.motorHatConnection, err = a.connector.GetConnection(motorHatAddress, bus); err != nil {
		return
	}

	if err = a.startDriver(a.motorHatConnection); err != nil {
		return
	}

	return
}

// Halt returns true if devices is halted successfully
func (a *AdafruitMotorHatDriver) Halt() (err error) { return }

// setPWM sets the start (on) and end (off) of the high-segment of the PWM pulse
// on the specific channel (pin).
func (a *AdafruitMotorHatDriver) setPWM(conn Connection, pin byte, on, off int32) (err error) {
	// register and values to be written to that register
	regVals := make(map[int][]byte)
	regVals[0] = []byte{byte(_LedZeroOnL + 4*pin), byte(on & 0xff)}
	regVals[1] = []byte{byte(_LedZeroOnH + 4*pin), byte(on >> 8)}
	regVals[2] = []byte{byte(_LedZeroOffL + 4*pin), byte(off & 0xff)}
	regVals[3] = []byte{byte(_LedZeroOffH + 4*pin), byte(off >> 8)}
	for i := 0; i < len(regVals); i++ {
		if _, err = conn.Write(regVals[i]); err != nil {
			return
		}
	}
	return
}

// SetServoMotorFreq sets the frequency for the currently addressed PWM Servo HAT.
func (a *AdafruitMotorHatDriver) SetServoMotorFreq(freq float64) (err error) {
	if err = a.setPWMFreq(a.servoHatConnection, freq); err != nil {
		return
	}
	return
}

// SetServoMotorPulse is a convenience function to specify the 'tick' value,
// between 0-4095, when the signal will turn on, and when it will turn off.
func (a *AdafruitMotorHatDriver) SetServoMotorPulse(channel byte, on, off int32) (err error) {
	if err = a.setPWM(a.servoHatConnection, channel, on, off); err != nil {
		return
	}
	return
}

// setPWMFreq adjusts the PWM frequency which determines how many full
// pulses per second are generated by the integrated circuit.  The frequency
// determines how "long" each pulse is in duration from start to finish,
// taking into account the high and low segments of the pulse.
func (a *AdafruitMotorHatDriver) setPWMFreq(conn Connection, freq float64) (err error) {
	// 25MHz
	preScaleVal := 25000000.0
	// 12-bit
	preScaleVal /= 4096.0
	preScaleVal /= freq
	preScaleVal -= 1.0
	preScale := math.Floor(preScaleVal + 0.5)
	if adafruitDebug {
		log.Printf("Setting PWM frequency to:	%.2f Hz", freq)
		log.Printf("Estimated pre-scale: 		%.2f", preScaleVal)
		log.Printf("Final pre-scale: 			%.2f", preScale)
	}
	// default (and only) reads register 0
	oldMode := []byte{0}
	_, err = conn.Read(oldMode)
	if err != nil {
		return
	}
	// sleep?
	if len(oldMode) > 0 {
		newMode := (oldMode[0] & 0x7F) | 0x10
		reg := byte(_Mode1)
		if _, err = conn.Write([]byte{reg, newMode}); err != nil {
			return
		}
		reg = byte(_Prescale)
		val := byte(math.Floor(preScale))
		if _, err = conn.Write([]byte{reg, val}); err != nil {
			return
		}
		reg = byte(_Mode1)
		if _, err = conn.Write([]byte{reg, oldMode[0]}); err != nil {
			return
		}
		time.Sleep(5 * time.Millisecond)
		if _, err = conn.Write([]byte{reg, (oldMode[0] | 0x80)}); err != nil {
			return
		}
	}
	return
}

// setAllPWM sets all PWM channels for the given address
func (a *AdafruitMotorHatDriver) setAllPWM(conn Connection, on, off int32) (err error) {
	// register and values to be written to that register
	regVals := make(map[int][]byte)
	regVals[0] = []byte{byte(_AllLedOnL), byte(on & 0xff)}
	regVals[1] = []byte{byte(_AllLedOnH), byte(on >> 8)}
	regVals[2] = []byte{byte(_AllLedOffL), byte(off & 0xFF)}
	regVals[3] = []byte{byte(_AllLedOffH), byte(off >> 8)}
	for i := 0; i < len(regVals); i++ {
		if _, err = conn.Write(regVals[i]); err != nil {
			return
		}
	}
	return
}

func (a *AdafruitMotorHatDriver) setPin(conn Connection, pin byte, value int32) (err error) {
	if value == 0 {
		return a.setPWM(conn, pin, 0, 4096)
	}
	if value == 1 {
		return a.setPWM(conn, pin, 4096, 0)
	}
	return errors.New("Invalid pin")
}

// SetDCMotorSpeed will set the appropriate pins to run the specified DC motor
// for the given speed.
func (a *AdafruitMotorHatDriver) SetDCMotorSpeed(dcMotor int, speed int32) (err error) {
	if err = a.setPWM(a.motorHatConnection, a.dcMotors[dcMotor].pwmPin, 0, speed*16); err != nil {
		return
	}
	return
}

// RunDCMotor will set the appropriate pins to run the specified DC motor for
// the given direction
func (a *AdafruitMotorHatDriver) RunDCMotor(dcMotor int, dir AdafruitDirection) (err error) {

	switch {
	case dir == AdafruitForward:
		if err = a.setPin(a.motorHatConnection, a.dcMotors[dcMotor].in2Pin, 0); err != nil {
			return
		}
		if err = a.setPin(a.motorHatConnection, a.dcMotors[dcMotor].in1Pin, 1); err != nil {
			return
		}
	case dir == AdafruitBackward:
		if err = a.setPin(a.motorHatConnection, a.dcMotors[dcMotor].in1Pin, 0); err != nil {
			return
		}
		if err = a.setPin(a.motorHatConnection, a.dcMotors[dcMotor].in2Pin, 1); err != nil {
			return
		}
	case dir == AdafruitRelease:
		if err = a.setPin(a.motorHatConnection, a.dcMotors[dcMotor].in1Pin, 0); err != nil {
			return
		}
		if err = a.setPin(a.motorHatConnection, a.dcMotors[dcMotor].in2Pin, 0); err != nil {
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
	if err = a.setPWM(a.motorHatConnection, a.stepperMotors[motor].pwmPinA, 0, int32(pwmA*16)); err != nil {
		return
	}
	if err = a.setPWM(a.motorHatConnection, a.stepperMotors[motor].pwmPinB, 0, int32(pwmB*16)); err != nil {
		return
	}
	var coils []int32
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
		// step-2-coils is initialized in init()
		coils = step2coils[(currStep / (stepperMicrosteps / 2))]
	}
	if adafruitDebug {
		log.Printf("[adafruit_driver] currStep: %d, index into step2coils: %d\n",
			currStep, (currStep / (stepperMicrosteps / 2)))
		log.Printf("[adafruit_driver] coils state = %v", coils)
	}
	if err = a.setPin(a.motorHatConnection, a.stepperMotors[motor].ain2, coils[0]); err != nil {
		return
	}
	if err = a.setPin(a.motorHatConnection, a.stepperMotors[motor].bin1, coils[1]); err != nil {
		return
	}
	if err = a.setPin(a.motorHatConnection, a.stepperMotors[motor].ain1, coils[2]); err != nil {
		return
	}
	if err = a.setPin(a.motorHatConnection, a.stepperMotors[motor].bin2, coils[3]); err != nil {
		return
	}
	return a.stepperMotors[motor].currentStep, nil
}

// SetStepperMotorSpeed sets the seconds-per-step for the given Stepper Motor.
func (a *AdafruitMotorHatDriver) SetStepperMotorSpeed(stepperMotor int, rpm int) (err error) {
	revSteps := a.stepperMotors[stepperMotor].revSteps
	a.stepperMotors[stepperMotor].secPerStep = 60.0 / float64(revSteps*rpm)
	a.stepperMotors[stepperMotor].stepCounter = 0
	return
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
		time.Sleep(time.Duration(secPerStep) * time.Second)
	}
	// As documented in the Adafruit python driver:
	// This is an edge case, if we are in between full steps, keep going to end on a full step
	if style == AdafruitMicrostep {
		for latestStep != 0 && latestStep != stepperMicrosteps {
			if latestStep, err = a.oneStep(motor, dir, style); err != nil {
				return
			}
			time.Sleep(time.Duration(secPerStep) * time.Second)
		}
	}
	return
}

func init() {
	stepperMicrostepCurve = []int{0, 50, 98, 142, 180, 212, 236, 250, 255}
	step2coils[0] = []int32{1, 0, 0, 0}
	step2coils[1] = []int32{1, 1, 0, 0}
	step2coils[2] = []int32{0, 1, 0, 0}
	step2coils[3] = []int32{0, 1, 1, 0}
	step2coils[4] = []int32{0, 0, 1, 0}
	step2coils[5] = []int32{0, 0, 1, 1}
	step2coils[6] = []int32{0, 0, 0, 1}
	step2coils[7] = []int32{1, 0, 0, 1}
}
