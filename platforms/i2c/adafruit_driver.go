package i2c

import "github.com/hybridgroup/gobot"

var _ gobot.Driver = (*AdafruitMotorHatDriver)(nil)

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

const (
	AdafruitForward  AdafruitDirection = iota // 0
	AdafruitBackward                          // 1
	AdafruitRelease                           // 2
)

type adaFruitDCMotor struct {
	pwmPin, in1Pin, in2Pin byte
}

// AdafruitMotorHatDriver is a driver for the DC+Stepper Motor HAT from Adafruit.
// The HAT is a Raspberry Pi add-on that can drive up to 4 DC or 2 Stepper motors
// with full PWM speed control.  It has a dedicated PWM driver chip onboard to
// control both motor direction and speed over I2C.
type AdafruitMotorHatDriver struct {
	name       string
	connection I2c
	gobot.Commander
	dcMotors []adaFruitDCMotor
}

// Name identifies this driver object
func (a *AdafruitMotorHatDriver) Name() string { return a.name }

// Connection identifies the particular adapter object
func (a *AdafruitMotorHatDriver) Connection() gobot.Connection { return a.connection.(gobot.Connection) }

// NewAdafruitMotorHatDriver initializes the internal DCMotor and StepperMotor types
func NewAdafruitMotorHatDriver(a I2c, name string) *AdafruitMotorHatDriver {
	var dc []adaFruitDCMotor
	for i := 0; i < 4; i++ {
		switch {
		case i == 0:
			dc = append(dc, adaFruitDCMotor{pwmPin: 8, in1Pin: 10, in2Pin: 9})
		case i == 1:
			dc = append(dc, adaFruitDCMotor{pwmPin: 13, in1Pin: 11, in2Pin: 12})
		case i == 2:
			dc = append(dc, adaFruitDCMotor{pwmPin: 2, in1Pin: 4, in2Pin: 3})
		case i == 3:
			dc = append(dc, adaFruitDCMotor{pwmPin: 7, in1Pin: 5, in2Pin: 6})
		}
	}
	driver := &AdafruitMotorHatDriver{
		name:       name,
		connection: a,
		Commander:  gobot.NewCommander(),
		dcMotors:   dc,
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
