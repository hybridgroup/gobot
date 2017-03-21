package gpio

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*LedDriver)(nil)

func initTestLedDriver(conn DigitalWriter) *LedDriver {
	testAdaptorDigitalWrite = func() (err error) {
		return nil
	}
	testAdaptorPwmWrite = func() (err error) {
		return nil
	}
	return NewLedDriver(conn, "1")
}

func TestLedDriver(t *testing.T) {
	var err interface{}

	d := initTestLedDriver(newGpioTestAdaptor())

	gobottest.Assert(t, d.Pin(), "1")
	gobottest.Refute(t, d.Connection(), nil)

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
	d := initTestLedDriver(newGpioTestAdaptor())
	gobottest.Assert(t, d.Start(), nil)
}

func TestLedDriverHalt(t *testing.T) {
	d := initTestLedDriver(newGpioTestAdaptor())
	gobottest.Assert(t, d.Halt(), nil)
}

func TestLedDriverToggle(t *testing.T) {
	d := initTestLedDriver(newGpioTestAdaptor())
	d.Off()
	d.Toggle()
	gobottest.Assert(t, d.State(), true)
	d.Toggle()
	gobottest.Assert(t, d.State(), false)
}

func TestLedDriverBrightness(t *testing.T) {
	d := initTestLedDriver(&gpioTestDigitalWriter{})
	gobottest.Assert(t, d.Brightness(150), ErrPwmWriteUnsupported)

	d = initTestLedDriver(newGpioTestAdaptor())
	testAdaptorPwmWrite = func() (err error) {
		err = errors.New("pwm error")
		return
	}
	gobottest.Assert(t, d.Brightness(150), errors.New("pwm error"))
}

func TestLEDDriverDefaultName(t *testing.T) {
	d := initTestLedDriver(&gpioTestDigitalWriter{})
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "LED"), true)
}

func TestLEDDriverSetName(t *testing.T) {
	d := initTestLedDriver(&gpioTestDigitalWriter{})
	d.SetName("mybot")
	gobottest.Assert(t, d.Name(), "mybot")
}
