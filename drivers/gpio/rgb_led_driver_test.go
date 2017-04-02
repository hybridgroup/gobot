package gpio

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*RgbLedDriver)(nil)

func initTestRgbLedDriver() *RgbLedDriver {
	a := newGpioTestAdaptor()
	a.testAdaptorDigitalWrite = func() (err error) {
		return nil
	}
	a.testAdaptorPwmWrite = func() (err error) {
		return nil
	}
	return NewRgbLedDriver(a, "1", "2", "3")
}

func TestRgbLedDriver(t *testing.T) {
	var err interface{}

	a := newGpioTestAdaptor()
	d := NewRgbLedDriver(a, "1", "2", "3")

	gobottest.Assert(t, d.Pin(), "r=1, g=2, b=3")
	gobottest.Assert(t, d.RedPin(), "1")
	gobottest.Assert(t, d.GreenPin(), "2")
	gobottest.Assert(t, d.BluePin(), "3")
	gobottest.Refute(t, d.Connection(), nil)

	a.testAdaptorDigitalWrite = func() (err error) {
		return errors.New("write error")
	}
	a.testAdaptorPwmWrite = func() (err error) {
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
	d := initTestRgbLedDriver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestRgbLedDriverHalt(t *testing.T) {
	d := initTestRgbLedDriver()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestRgbLedDriverToggle(t *testing.T) {
	d := initTestRgbLedDriver()
	d.Off()
	d.Toggle()
	gobottest.Assert(t, d.State(), true)
	d.Toggle()
	gobottest.Assert(t, d.State(), false)
}

func TestRgbLedDriverSetLevel(t *testing.T) {
	a := newGpioTestAdaptor()
	d := NewRgbLedDriver(a, "1", "2", "3")
	gobottest.Assert(t, d.SetLevel("1", 150), nil)

	d = NewRgbLedDriver(a, "1", "2", "3")
	a.testAdaptorPwmWrite = func() (err error) {
		err = errors.New("pwm error")
		return
	}
	gobottest.Assert(t, d.SetLevel("1", 150), errors.New("pwm error"))
}

func TestRgbLedDriverDefaultName(t *testing.T) {
	a := newGpioTestAdaptor()
	d := NewRgbLedDriver(a, "1", "2", "3")
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "RGB"), true)
}

func TestRgbLedDriverSetName(t *testing.T) {
	a := newGpioTestAdaptor()
	d := NewRgbLedDriver(a, "1", "2", "3")
	d.SetName("mybot")
	gobottest.Assert(t, d.Name(), "mybot")
}
