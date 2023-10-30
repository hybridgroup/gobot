package i2c

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

// this ensures that the implementation implements the gobot.Driver interface
var _ gobot.Driver = (*Adafruit2348Driver)(nil)

func initTestAdafruit2348WithStubbedAdaptor() (*Adafruit2348Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewAdafruit2348Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewAdafruit2348Driver(t *testing.T) {
	// arrange & act
	d := NewAdafruit2348Driver(newI2cTestAdaptor())
	// assert
	assert.IsType(t, &Adafruit2348Driver{}, d)
	assert.True(t, strings.HasPrefix(d.Name(), "Adafruit2348MotorHat"))
	assert.Equal(t, 0x40, d.defaultAddress)                               // the default address of PCA9685 driver
	assert.Equal(t, 0x60, d.Config.GetAddressOrDefault(d.defaultAddress)) // the really used address
}

func TestAdafruit2348Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	// arrange & act
	d := NewAdafruit2348Driver(newI2cTestAdaptor(), WithBus(2), WithAddress(0x45))
	// assert
	assert.Equal(t, 2, d.GetBusOrDefault(1))
	assert.Equal(t, 0x45, d.GetAddressOrDefault(2))
}

func TestAdafruit2348SetDCMotorSpeed(t *testing.T) {
	// arrange
	d, a := initTestAdafruit2348WithStubbedAdaptor()
	a.written = []byte{} // reset writes of former test
	const (
		dcMotor = 1
		speed   = 255
	)
	// act
	err := d.SetDCMotorSpeed(dcMotor, speed)
	// assert
	assert.NoError(t, err)
	assert.Equal(t, 8, len(a.written)) // detailed test, see "TestPCA9685SetPWM"
}

func TestAdafruit2348SetDCMotorSpeedError(t *testing.T) {
	// arrange
	d, a := initTestAdafruit2348WithStubbedAdaptor()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	// act & assert
	assert.ErrorContains(t, d.SetDCMotorSpeed(1, 255), "write error")
}

func TestAdafruit2348RunDCMotor(t *testing.T) {
	// arrange
	d, _ := initTestAdafruit2348WithStubbedAdaptor()
	const dcMotor = 1
	// act & assert
	assert.NoError(t, d.RunDCMotor(dcMotor, Adafruit2348Forward))
	assert.NoError(t, d.RunDCMotor(dcMotor, Adafruit2348Backward))
	assert.NoError(t, d.RunDCMotor(dcMotor, Adafruit2348Release))
}

func TestAdafruit2348RunDCMotorError(t *testing.T) {
	// arrange
	d, a := initTestAdafruit2348WithStubbedAdaptor()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	const dcMotor = 1
	// act & assert
	assert.ErrorContains(t, d.RunDCMotor(dcMotor, Adafruit2348Forward), "write error")
	assert.ErrorContains(t, d.RunDCMotor(dcMotor, Adafruit2348Backward), "write error")
	assert.ErrorContains(t, d.RunDCMotor(dcMotor, Adafruit2348Release), "write error")
}

func TestAdafruit2348SetStepperMotorSpeed(t *testing.T) {
	// arrange
	d, _ := initTestAdafruit2348WithStubbedAdaptor()
	const (
		stepperMotor = 1
		rpm          = 30
	)
	// act & assert
	assert.NoError(t, d.SetStepperMotorSpeed(stepperMotor, rpm))
	assert.Equal(t, 0.01, d.stepperMotors[stepperMotor].secPerStep) // 60/(revSteps*rpm), revSteps=200
}

func TestAdafruit2348StepperSingleStep(t *testing.T) {
	// arrange
	d, _ := initTestAdafruit2348WithStubbedAdaptor()
	const (
		stepperMotor = 0
		steps        = 50
		back         = 1
		single       = 0
	)
	// act
	err := d.Step(stepperMotor, steps, back, single)
	// assert
	assert.NoError(t, err)
}

func TestAdafruit2348StepperDoubleStep(t *testing.T) {
	// arrange
	d, _ := initTestAdafruit2348WithStubbedAdaptor()
	const (
		stepperMotor = 0
		steps        = 50
		back         = 1
		double       = 1
	)
	// act
	err := d.Step(stepperMotor, steps, back, double)
	// assert
	assert.NoError(t, err)
}

func TestAdafruit2348StepperInterleaveStep(t *testing.T) {
	// arrange
	d, _ := initTestAdafruit2348WithStubbedAdaptor()
	const (
		stepperMotor = 0
		steps        = 50
		back         = 1
		interleave   = 2
	)
	// act
	err := d.Step(stepperMotor, steps, back, interleave)
	// assert
	assert.NoError(t, err)
}

func TestAdafruit2348StepperMicroStep(t *testing.T) {
	// arrange
	d, _ := initTestAdafruit2348WithStubbedAdaptor()
	const (
		stepperMotor = 0
		steps        = 50
		back         = 1
		micro        = 3
	)
	// act
	err := d.Step(stepperMotor, steps, back, micro)
	// assert
	assert.NoError(t, err)
}
