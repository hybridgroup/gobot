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

// TODO: package-wide, thus make more descriptive names
type Direction int

const (
	Forward  Direction = iota // 0
	Backward                  // 1
	Release                   // 2
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

// NewAdafruitDriver adds the following API commands: TODO
func NewAdafruitDriver(a I2c, name string) *AdafruitMotorHatDriver {
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

// Start initializes the Adafruit Motor HAT
func (a *AdafruitMotorHatDriver) Start() (errs []error) {
	if err := a.connection.I2cStart(motorHatAddress); err != nil {
		return []error{err}
	}
	/*
		  def __init__(self, address=0x40, debug=False):
			self.i2c = Adafruit_I2C(address)
			self.i2c.debug = debug
			self.address = address
			self.debug = debug
			if (self.debug):
			print "Reseting PCA9685 MODE1 (without SLEEP) and MODE2"
			self.setAllPWM(0, 0)
			self.i2c.write8(self.__MODE2, self.__OUTDRV)
			self.i2c.write8(self.__MODE1, self.__ALLCALL)
			time.sleep(0.005)                                       # wait for oscillator

			mode1 = self.i2c.readU8(self.__MODE1)
			mode1 = mode1 & ~self.__SLEEP                 # wake up (reset sleep)
			self.i2c.write8(self.__MODE1, mode1)
			time.sleep(0.005)                             # wait for oscillator
	*/
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

/*

Test Runner:
	For motor 3:
	        pwm = 2
            in2 = 3
            in1 = 4

	1. SetSpeed(setPWM(self.PWMpin, 0, speed*16))
	2. Run(Adafruit_MotorHAT.FORWARD):
    	self.MC.setPin(self.IN2pin, 0)
    	self.MC.setPin(self.IN1pin, 1)

	Run():
	  Where setPin(pin, value) is:
		if (value == 0):
        	self._pwm.setPWM(pin, 0, 4096)
        if (value == 1):
            self._pwm.setPWM(pin, 4096, 0)

*/
func (a *AdafruitMotorHatDriver) SetDCMotorSpeed(dcMotor int, speed int32) (err error) {
	if err = a.setPWM(a.dcMotors[dcMotor].pwmPin, 0, speed*16); err != nil {
		return
	}
	return
}

// RunDCMotor will set the appropriate pins to run the specific DC motor for
// the given direction
func (a *AdafruitMotorHatDriver) RunDCMotor(dcMotor int, dir Direction) (err error) {

	switch {
	case dir == Forward:
		if err = a.setPin(a.dcMotors[dcMotor].in2Pin, 0); err != nil {
			return
		}
		if err = a.setPin(a.dcMotors[dcMotor].in1Pin, 1); err != nil {
			return
		}
	case dir == Backward:
		if err = a.setPin(a.dcMotors[dcMotor].in1Pin, 0); err != nil {
			return
		}
		if err = a.setPin(a.dcMotors[dcMotor].in2Pin, 1); err != nil {
			return
		}
	case dir == Release:
		if err = a.setPin(a.dcMotors[dcMotor].in1Pin, 0); err != nil {
			return
		}
		if err = a.setPin(a.dcMotors[dcMotor].in2Pin, 0); err != nil {
			return
		}
	}
	return
}
