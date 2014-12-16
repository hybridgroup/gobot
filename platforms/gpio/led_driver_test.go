package gpio

import (
	"errors"
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestLedDriver(conn DigitalWriter) *LedDriver {
	testAdaptorDigitalWrite = func() (err error) {
		return nil
	}
	testAdaptorPwmWrite = func() (err error) {
		return nil
	}
	return NewLedDriver(conn, "bot", "1")
}

func TestLedDriver(t *testing.T) {
	var err interface{}

	d := initTestLedDriver(newGpioTestAdaptor("adaptor"))

	gobot.Assert(t, d.Name(), "bot")
	gobot.Assert(t, d.Pin(), "1")
	gobot.Assert(t, d.Connection().Name(), "adaptor")

	testAdaptorDigitalWrite = func() (err error) {
		return errors.New("write error")
	}
	testAdaptorPwmWrite = func() (err error) {
		return errors.New("pwm error")
	}

	err = d.Command("Toggle")(nil)
	gobot.Assert(t, err.(error), errors.New("write error"))

	err = d.Command("On")(nil)
	gobot.Assert(t, err.(error), errors.New("write error"))

	err = d.Command("Off")(nil)
	gobot.Assert(t, err.(error), errors.New("write error"))

	err = d.Command("Brightness")(map[string]interface{}{"level": 100.0})
	gobot.Assert(t, err.(error), errors.New("pwm error"))

}

func TestLedDriverStart(t *testing.T) {
	d := initTestLedDriver(newGpioTestAdaptor("adaptor"))
	gobot.Assert(t, len(d.Start()), 0)
}

func TestLedDriverHalt(t *testing.T) {
	d := initTestLedDriver(newGpioTestAdaptor("adaptor"))
	gobot.Assert(t, len(d.Halt()), 0)
}

func TestLedDriverToggle(t *testing.T) {
	d := initTestLedDriver(newGpioTestAdaptor("adaptor"))
	d.Off()
	d.Toggle()
	gobot.Assert(t, d.State(), true)
	d.Toggle()
	gobot.Assert(t, d.State(), false)
}

func TestLedDriverBrightness(t *testing.T) {
	d := initTestLedDriver(&gpioTestDigitalWriter{})
	gobot.Assert(t, d.Brightness(150), ErrPwmWriteUnsupported)

	d = initTestLedDriver(newGpioTestAdaptor("adaptor"))
	testAdaptorPwmWrite = func() (err error) {
		err = errors.New("pwm error")
		return
	}
	gobot.Assert(t, d.Brightness(150), errors.New("pwm error"))
}
