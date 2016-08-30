package i2c

import (
	"errors"
	"testing"

	"github.com/jfinken/gobot/gobottest"
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
