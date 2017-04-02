package gpio

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*ServoDriver)(nil)

func initTestServoDriver() *ServoDriver {
	return NewServoDriver(newGpioTestAdaptor(), "1")
}

func TestServoDriver(t *testing.T) {
	var err interface{}

	a := newGpioTestAdaptor()
	d := NewServoDriver(a, "1")

	gobottest.Assert(t, d.Pin(), "1")
	gobottest.Refute(t, d.Connection(), nil)

	a.testAdaptorServoWrite = func() (err error) {
		return errors.New("pwm error")
	}

	err = d.Command("Min")(nil)
	gobottest.Assert(t, err.(error), errors.New("pwm error"))

	err = d.Command("Center")(nil)
	gobottest.Assert(t, err.(error), errors.New("pwm error"))

	err = d.Command("Max")(nil)
	gobottest.Assert(t, err.(error), errors.New("pwm error"))

	err = d.Command("Move")(map[string]interface{}{"angle": 100.0})
	gobottest.Assert(t, err.(error), errors.New("pwm error"))
}

func TestServoDriverStart(t *testing.T) {
	d := initTestServoDriver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestServoDriverHalt(t *testing.T) {
	d := initTestServoDriver()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestServoDriverMove(t *testing.T) {
	d := initTestServoDriver()
	d.Move(100)
	gobottest.Assert(t, d.CurrentAngle, uint8(100))
	err := d.Move(200)
	gobottest.Assert(t, err, ErrServoOutOfRange)
}

func TestServoDriverMin(t *testing.T) {
	d := initTestServoDriver()
	d.Min()
	gobottest.Assert(t, d.CurrentAngle, uint8(0))
}

func TestServoDriverMax(t *testing.T) {
	d := initTestServoDriver()
	d.Max()
	gobottest.Assert(t, d.CurrentAngle, uint8(180))
}

func TestServoDriverCenter(t *testing.T) {
	d := initTestServoDriver()
	d.Center()
	gobottest.Assert(t, d.CurrentAngle, uint8(90))
}

func TestServoDriverDefaultName(t *testing.T) {
	d := initTestServoDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Servo"), true)
}

func TestServoDriverSetName(t *testing.T) {
	d := initTestServoDriver()
	d.SetName("mybot")
	gobottest.Assert(t, d.Name(), "mybot")
}
