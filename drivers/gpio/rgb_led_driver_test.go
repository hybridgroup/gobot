package gpio

import (
	"errors"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gobottest"
)

var _ gobot.Driver = (*RgbLedDriver)(nil)

func initTestRgbLedDriver(conn DigitalWriter) *RgbLedDriver {
	testAdaptorDigitalWrite = func() (err error) {
		return nil
	}
	testAdaptorPwmWrite = func() (err error) {
		return nil
	}
	return NewRgbLedDriver(conn, "bot", "1", "2", "3")
}

func TestRgbLedDriver(t *testing.T) {
	var err interface{}

	d := initTestRgbLedDriver(newGpioTestAdaptor("adaptor"))

	gobottest.Assert(t, d.Name(), "bot")
	gobottest.Assert(t, d.Pin(), "r=1, g=2, b=3")
	gobottest.Assert(t, d.RedPin(), "1")
	gobottest.Assert(t, d.GreenPin(), "2")
	gobottest.Assert(t, d.BluePin(), "3")
	gobottest.Assert(t, d.Connection().Name(), "adaptor")

	testAdaptorDigitalWrite = func() (err error) {
		return errors.New("write error")
	}
	testAdaptorPwmWrite = func() (err error) {
		return errors.New("pwm error")
	}

	err = d.Command("Toggle")(nil)
	gobottest.Assert(t, err.(error), errors.New("pwm error"))

	err = d.Command("On")(nil)
	gobottest.Assert(t, err.(error), errors.New("pwm error"))

	err = d.Command("Off")(nil)
	gobottest.Assert(t, err.(error), errors.New("pwm error"))

	err = d.Command("SetRGB")(map[string]interface{}{"r": 0xff, "g": 0xff, "b": 0xff})
	gobottest.Assert(t, err.(error), errors.New("pwm error"))

}

func TestRgbLedDriverStart(t *testing.T) {
	d := initTestRgbLedDriver(newGpioTestAdaptor("adaptor"))
	gobottest.Assert(t, len(d.Start()), 0)
}

func TestRgbLedDriverHalt(t *testing.T) {
	d := initTestRgbLedDriver(newGpioTestAdaptor("adaptor"))
	gobottest.Assert(t, len(d.Halt()), 0)
}

func TestRgbLedDriverToggle(t *testing.T) {
	d := initTestRgbLedDriver(newGpioTestAdaptor("adaptor"))
	d.Off()
	d.Toggle()
	gobottest.Assert(t, d.State(), true)
	d.Toggle()
	gobottest.Assert(t, d.State(), false)
}

func TestRgbLedDriverSetLevel(t *testing.T) {
	d := initTestRgbLedDriver(&gpioTestDigitalWriter{})
	gobottest.Assert(t, d.SetLevel("1", 150), ErrPwmWriteUnsupported)

	d = initTestRgbLedDriver(newGpioTestAdaptor("adaptor"))
	testAdaptorPwmWrite = func() (err error) {
		err = errors.New("pwm error")
		return
	}
	gobottest.Assert(t, d.SetLevel("1", 150), errors.New("pwm error"))
}
