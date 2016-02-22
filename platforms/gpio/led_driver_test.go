package gpio

import (
	"errors"
	"testing"

	"github.com/hybridgroup/gobot/gobottest"
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

	gobottest.Assert(t, d.Name(), "bot")
	gobottest.Assert(t, d.Pin(), "1")
	gobottest.Assert(t, d.Connection().Name(), "adaptor")

	testAdaptorDigitalWrite = func() (err error) {
		return errors.New("write error")
	}
	testAdaptorPwmWrite = func() (err error) {
		return errors.New("pwm error")
	}

	err = d.Command("Toggle")(nil)
	gobottest.Assert(t, err.(error), errors.New("write error"))

	err = d.Command("On")(nil)
	gobottest.Assert(t, err.(error), errors.New("write error"))

	err = d.Command("Off")(nil)
	gobottest.Assert(t, err.(error), errors.New("write error"))

	err = d.Command("Brightness")(map[string]interface{}{"level": 100.0})
	gobottest.Assert(t, err.(error), errors.New("pwm error"))

}

func TestLedDriverStart(t *testing.T) {
	d := initTestLedDriver(newGpioTestAdaptor("adaptor"))
	gobottest.Assert(t, len(d.Start()), 0)
}

func TestLedDriverHalt(t *testing.T) {
	d := initTestLedDriver(newGpioTestAdaptor("adaptor"))
	gobottest.Assert(t, len(d.Halt()), 0)
}

func TestLedDriverToggle(t *testing.T) {
	d := initTestLedDriver(newGpioTestAdaptor("adaptor"))
	d.Off()
	d.Toggle()
	gobottest.Assert(t, d.State(), true)
	d.Toggle()
	gobottest.Assert(t, d.State(), false)
}

func TestLedDriverBrightness(t *testing.T) {
	d := initTestLedDriver(&gpioTestDigitalWriter{})
	gobottest.Assert(t, d.Brightness(150), ErrPwmWriteUnsupported)

	d = initTestLedDriver(newGpioTestAdaptor("adaptor"))
	testAdaptorPwmWrite = func() (err error) {
		err = errors.New("pwm error")
		return
	}
	gobottest.Assert(t, d.Brightness(150), errors.New("pwm error"))
}
