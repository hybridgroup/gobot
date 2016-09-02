package i2c

import (
	"errors"
	"testing"

	"github.com/hybridgroup/gobot/gobottest"
)

// --------- HELPERS
func initTestAdafruitMotorHatDriver() (driver *AdafruitMotorHatDriver) {
	driver, _ = initTestAdafruitMotorHatDriverWithStubbedAdaptor()
	return
}

func initTestAdafruitMotorHatDriverWithStubbedAdaptor() (*AdafruitMotorHatDriver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor("adaptor")
	return NewAdafruitMotorHatDriver(adaptor, "bot"), adaptor
}

// --------- TESTS
func TestNewAdafruitMotorHatDriver(t *testing.T) {
	var adafruit interface{} = NewAdafruitMotorHatDriver(newI2cTestAdaptor("adaptor"), "bot")
	_, ok := adafruit.(*AdafruitMotorHatDriver)
	if !ok {
		t.Errorf("AdafruitMotorHatDriver() should have returned a *AdafruitMotorHatDriver")
	}

	a := NewAdafruitMotorHatDriver(newI2cTestAdaptor("adaptor"), "bot")
	gobottest.Assert(t, a.Name(), "bot")
	gobottest.Assert(t, a.Connection().Name(), "adaptor")
}

// Methods
func TestAdafruitMotorHatDriverStart(t *testing.T) {
	ada, adaptor := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	gobottest.Assert(t, len(ada.Start()), 0)

	adaptor.i2cStartImpl = func() error {
		return errors.New("start error")
	}
	err := ada.Start()
	gobottest.Assert(t, err[0], errors.New("start error"))

}

func TestAdafruitMotorHatDriverHalt(t *testing.T) {
	ada, _ := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	gobottest.Assert(t, len(ada.Halt()), 0)
}
func TestSetHatAddresses(t *testing.T) {
	ada, _ := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	motorHatAddr := 0x61
	servoHatAddr := 0x41
	gobottest.Assert(t, len(ada.SetMotorHatAddress(motorHatAddr)), 0)
	gobottest.Assert(t, len(ada.SetServoHatAddress(servoHatAddr)), 0)
}

func TestAdafruitMotorHatDriverSetServoMotorFreq(t *testing.T) {
	ada, _ := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	freq := 60.0
	err := ada.SetServoMotorFreq(freq)
	gobottest.Assert(t, err, nil)
}

func TestAdafruitMotorHatDriverSetServoMotorPulse(t *testing.T) {
	ada, _ := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	var channel byte = 7
	var on int32 = 1234
	var off int32 = 4321
	err := ada.SetServoMotorPulse(channel, on, off)
	gobottest.Assert(t, err, nil)
}

func TestAdafruitMotorHatDriverSetDCMotorSpeed(t *testing.T) {
	ada, _ := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	dcMotor := 1
	var speed int32 = 255
	err := ada.SetDCMotorSpeed(dcMotor, speed)
	gobottest.Assert(t, err, nil)
}

func TestAdafruitMotorHatDriverRunDCMotor(t *testing.T) {
	ada, _ := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	dcMotor := 1
	// NOTE: not using the direction constant to prevent importing
	// the i2c package
	err := ada.RunDCMotor(dcMotor, 1)
	gobottest.Assert(t, err, nil)
}

func TestAdafruitMotorHatDriverSetStepperMotorSpeed(t *testing.T) {
	ada, _ := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	stepperMotor := 1
	rpm := 30
	gobottest.Assert(t, len(ada.SetStepperMotorSpeed(stepperMotor, rpm)), 0)
}

func TestAdafruitMotorHatDriverStepperStep(t *testing.T) {
	ada, _ := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	// NOTE: not using the direction and style constants to prevent importing
	// the i2c package
	stepperMotor := 0
	steps := 50
	err := ada.Step(stepperMotor, steps, 1, 3)
	gobottest.Assert(t, err, nil)
}
