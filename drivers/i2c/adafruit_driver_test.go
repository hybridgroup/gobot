package i2c

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*AdafruitMotorHatDriver)(nil)

// --------- HELPERS
func initTestAdafruitMotorHatDriver() (driver *AdafruitMotorHatDriver) {
	driver, _ = initTestAdafruitMotorHatDriverWithStubbedAdaptor()
	return
}

func initTestAdafruitMotorHatDriverWithStubbedAdaptor() (*AdafruitMotorHatDriver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewAdafruitMotorHatDriver(adaptor), adaptor
}

// --------- TESTS
func TestNewAdafruitMotorHatDriver(t *testing.T) {
	var adafruit interface{} = NewAdafruitMotorHatDriver(newI2cTestAdaptor())
	_, ok := adafruit.(*AdafruitMotorHatDriver)
	if !ok {
		t.Errorf("AdafruitMotorHatDriver() should have returned a *AdafruitMotorHatDriver")
	}

	a := NewAdafruitMotorHatDriver(newI2cTestAdaptor())
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "AdafruitMotorHat"), true)
}

// Methods
func TestAdafruitMotorHatDriverStart(t *testing.T) {
	ada, _ := initTestAdafruitMotorHatDriverWithStubbedAdaptor()
	gobottest.Refute(t, ada.Connection(), nil)
	gobottest.Assert(t, ada.Start(), nil)
}

func TestAdafruitMotorHatDriverStartWriteError(t *testing.T) {
	d, adaptor := initTestAdafruitMotorHatDriverWithStubbedAdaptor()
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, d.Start(), errors.New("write error"))
}

func TestAdafruitMotorHatDriverStartReadError(t *testing.T) {
	d, adaptor := initTestAdafruitMotorHatDriverWithStubbedAdaptor()
	adaptor.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}
	gobottest.Assert(t, d.Start(), errors.New("read error"))
}

func TestAdafruitMotorHatDriverStartConnectError(t *testing.T) {
	d, adaptor := initTestAdafruitMotorHatDriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestAdafruitMotorHatDriverHalt(t *testing.T) {
	ada, _ := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	gobottest.Assert(t, ada.Halt(), nil)
}

func TestSetHatAddresses(t *testing.T) {
	ada, _ := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	motorHatAddr := 0x61
	servoHatAddr := 0x41
	gobottest.Assert(t, ada.SetMotorHatAddress(motorHatAddr), nil)
	gobottest.Assert(t, ada.SetServoHatAddress(servoHatAddr), nil)
}

func TestAdafruitMotorHatDriverSetServoMotorFreq(t *testing.T) {
	ada, _ := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	gobottest.Assert(t, ada.Start(), nil)

	freq := 60.0
	err := ada.SetServoMotorFreq(freq)
	gobottest.Assert(t, err, nil)
}

func TestAdafruitMotorHatDriverSetServoMotorFreqError(t *testing.T) {
	ada, a := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	gobottest.Assert(t, ada.Start(), nil)
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	freq := 60.0
	gobottest.Assert(t, ada.SetServoMotorFreq(freq), errors.New("write error"))
}

func TestAdafruitMotorHatDriverSetServoMotorPulse(t *testing.T) {
	ada, _ := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	gobottest.Assert(t, ada.Start(), nil)

	var channel byte = 7
	var on int32 = 1234
	var off int32 = 4321
	err := ada.SetServoMotorPulse(channel, on, off)
	gobottest.Assert(t, err, nil)
}

func TestAdafruitMotorHatDriverSetServoMotorPulseError(t *testing.T) {
	ada, a := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	gobottest.Assert(t, ada.Start(), nil)
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	var channel byte = 7
	var on int32 = 1234
	var off int32 = 4321
	gobottest.Assert(t, ada.SetServoMotorPulse(channel, on, off), errors.New("write error"))
}

func TestAdafruitMotorHatDriverSetDCMotorSpeed(t *testing.T) {
	ada, _ := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	gobottest.Assert(t, ada.Start(), nil)

	dcMotor := 1
	var speed int32 = 255
	err := ada.SetDCMotorSpeed(dcMotor, speed)
	gobottest.Assert(t, err, nil)
}

func TestAdafruitMotorHatDriverSetDCMotorSpeedError(t *testing.T) {
	ada, a := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	gobottest.Assert(t, ada.Start(), nil)
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	gobottest.Assert(t, ada.SetDCMotorSpeed(1, 255), errors.New("write error"))
}

func TestAdafruitMotorHatDriverRunDCMotor(t *testing.T) {
	ada, _ := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	gobottest.Assert(t, ada.Start(), nil)

	dcMotor := 1
	gobottest.Assert(t, ada.RunDCMotor(dcMotor, AdafruitForward), nil)
	gobottest.Assert(t, ada.RunDCMotor(dcMotor, AdafruitBackward), nil)
	gobottest.Assert(t, ada.RunDCMotor(dcMotor, AdafruitRelease), nil)
}

func TestAdafruitMotorHatDriverRunDCMotorError(t *testing.T) {
	ada, a := initTestAdafruitMotorHatDriverWithStubbedAdaptor()
	gobottest.Assert(t, ada.Start(), nil)
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	dcMotor := 1
	gobottest.Assert(t, ada.RunDCMotor(dcMotor, AdafruitForward), errors.New("write error"))
	gobottest.Assert(t, ada.RunDCMotor(dcMotor, AdafruitBackward), errors.New("write error"))
	gobottest.Assert(t, ada.RunDCMotor(dcMotor, AdafruitRelease), errors.New("write error"))
}

func TestAdafruitMotorHatDriverSetStepperMotorSpeed(t *testing.T) {
	ada, _ := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	gobottest.Assert(t, ada.Start(), nil)

	stepperMotor := 1
	rpm := 30
	gobottest.Assert(t, ada.SetStepperMotorSpeed(stepperMotor, rpm), nil)
}

func TestAdafruitMotorHatDriverStepperMicroStep(t *testing.T) {
	ada, _ := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	gobottest.Assert(t, ada.Start(), nil)

	// NOTE: not using the direction and style constants to prevent importing
	// the i2c package
	stepperMotor := 0
	steps := 50
	err := ada.Step(stepperMotor, steps, 1, 3)
	gobottest.Assert(t, err, nil)
}

func TestAdafruitMotorHatDriverStepperSingleStep(t *testing.T) {
	ada, _ := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	gobottest.Assert(t, ada.Start(), nil)

	// NOTE: not using the direction and style constants to prevent importing
	// the i2c package
	stepperMotor := 0
	steps := 50
	err := ada.Step(stepperMotor, steps, 1, 0)
	gobottest.Assert(t, err, nil)
}

func TestAdafruitMotorHatDriverStepperDoubleStep(t *testing.T) {
	ada, _ := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	gobottest.Assert(t, ada.Start(), nil)

	// NOTE: not using the direction and style constants to prevent importing
	// the i2c package
	stepperMotor := 0
	steps := 50
	err := ada.Step(stepperMotor, steps, 1, 1)
	gobottest.Assert(t, err, nil)
}

func TestAdafruitMotorHatDriverStepperInterleaveStep(t *testing.T) {
	ada, _ := initTestAdafruitMotorHatDriverWithStubbedAdaptor()

	gobottest.Assert(t, ada.Start(), nil)

	// NOTE: not using the direction and style constants to prevent importing
	// the i2c package
	stepperMotor := 0
	steps := 50
	err := ada.Step(stepperMotor, steps, 1, 2)
	gobottest.Assert(t, err, nil)
}

func TestAdafruitMotorHatDriverSetName(t *testing.T) {
	d := initTestAdafruitMotorHatDriver()
	d.SetName("TESTME")
	gobottest.Assert(t, d.Name(), "TESTME")
}

func TestAdafruitMotorHatDriverOptions(t *testing.T) {
	d := NewAdafruitMotorHatDriver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}
