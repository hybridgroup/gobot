package gpio

import (
	"errors"
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestServoDriver() *ServoDriver {
	return NewServoDriver(newGpioTestAdaptor("adaptor"), "bot", "1")
}

func TestServoDriver(t *testing.T) {
	var err interface{}

	d := initTestServoDriver()

	gobot.Assert(t, d.Name(), "bot")
	gobot.Assert(t, d.Pin(), "1")
	gobot.Assert(t, d.Connection().Name(), "adaptor")

	testAdaptorServoWrite = func() (err error) {
		return errors.New("pwm error")
	}

	err = d.Command("Min")(nil)
	gobot.Assert(t, err.(error), errors.New("pwm error"))

	err = d.Command("Center")(nil)
	gobot.Assert(t, err.(error), errors.New("pwm error"))

	err = d.Command("Max")(nil)
	gobot.Assert(t, err.(error), errors.New("pwm error"))

	err = d.Command("Move")(map[string]interface{}{"angle": 100.0})
	gobot.Assert(t, err.(error), errors.New("pwm error"))

}

func TestServoDriverStart(t *testing.T) {
	d := initTestServoDriver()
	gobot.Assert(t, len(d.Start()), 0)
}

func TestServoDriverHalt(t *testing.T) {
	d := initTestServoDriver()
	gobot.Assert(t, len(d.Halt()), 0)
}

func TestServoDriverMove(t *testing.T) {
	d := initTestServoDriver()
	d.Move(100)
	gobot.Assert(t, d.CurrentAngle, uint8(100))
	err := d.Move(200)
	gobot.Assert(t, err, ErrServoOutOfRange)
}

func TestServoDriverMin(t *testing.T) {
	d := initTestServoDriver()
	d.Min()
	gobot.Assert(t, d.CurrentAngle, uint8(0))
}

func TestServoDriverMax(t *testing.T) {
	d := initTestServoDriver()
	d.Max()
	gobot.Assert(t, d.CurrentAngle, uint8(180))
}

func TestServoDriverCenter(t *testing.T) {
	d := initTestServoDriver()
	d.Center()
	gobot.Assert(t, d.CurrentAngle, uint8(90))
}
