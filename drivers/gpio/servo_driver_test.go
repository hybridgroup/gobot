package gpio

import (
	"errors"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gobottest"
)

var _ gobot.Driver = (*ServoDriver)(nil)

func initTestServoDriver() *ServoDriver {
	return NewServoDriver(newGpioTestAdaptor("adaptor"), "bot", "1")
}

func TestServoDriver(t *testing.T) {
	var err interface{}

	d := initTestServoDriver()

	gobottest.Assert(t, d.Name(), "bot")
	gobottest.Assert(t, d.Pin(), "1")
	gobottest.Assert(t, d.Connection().Name(), "adaptor")

	testAdaptorServoWrite = func() (err error) {
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
	gobottest.Assert(t, len(d.Start()), 0)
}

func TestServoDriverHalt(t *testing.T) {
	d := initTestServoDriver()
	gobottest.Assert(t, len(d.Halt()), 0)
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
