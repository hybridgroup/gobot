package gpio

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*DirectPinDriver)(nil)

func initTestDirectPinDriver(conn gobot.Connection) *DirectPinDriver {
	testAdaptorDigitalRead = func() (val int, err error) {
		val = 1
		return
	}
	testAdaptorDigitalWrite = func() (err error) {
		return errors.New("write error")
	}
	testAdaptorPwmWrite = func() (err error) {
		return errors.New("write error")
	}
	testAdaptorServoWrite = func() (err error) {
		return errors.New("write error")
	}
	return NewDirectPinDriver(conn, "1")
}

func TestDirectPinDriver(t *testing.T) {
	var ret map[string]interface{}
	var err interface{}

	d := initTestDirectPinDriver(newGpioTestAdaptor())
	gobottest.Assert(t, d.Pin(), "1")
	gobottest.Refute(t, d.Connection(), nil)

	ret = d.Command("DigitalRead")(nil).(map[string]interface{})

	gobottest.Assert(t, ret["val"].(int), 1)
	gobottest.Assert(t, ret["err"], nil)

	err = d.Command("DigitalWrite")(map[string]interface{}{"level": "1"})
	gobottest.Assert(t, err.(error), errors.New("write error"))

	err = d.Command("PwmWrite")(map[string]interface{}{"level": "1"})
	gobottest.Assert(t, err.(error), errors.New("write error"))

	err = d.Command("ServoWrite")(map[string]interface{}{"level": "1"})
	gobottest.Assert(t, err.(error), errors.New("write error"))
}

func TestDirectPinDriverStart(t *testing.T) {
	d := initTestDirectPinDriver(newGpioTestAdaptor())
	gobottest.Assert(t, d.Start(), nil)
}

func TestDirectPinDriverHalt(t *testing.T) {
	d := initTestDirectPinDriver(newGpioTestAdaptor())
	gobottest.Assert(t, d.Halt(), nil)
}

func TestDirectPinDriverOff(t *testing.T) {
	d := initTestDirectPinDriver(newGpioTestAdaptor())
	gobottest.Refute(t, d.DigitalWrite(0), nil)

	d = initTestDirectPinDriver(&gpioTestBareAdaptor{})
	gobottest.Assert(t, d.DigitalWrite(0), ErrDigitalWriteUnsupported)
}

func TestDirectPinDriverOn(t *testing.T) {
	d := initTestDirectPinDriver(newGpioTestAdaptor())
	gobottest.Refute(t, d.DigitalWrite(1), nil)

	d = initTestDirectPinDriver(&gpioTestBareAdaptor{})
	gobottest.Assert(t, d.DigitalWrite(1), ErrDigitalWriteUnsupported)
}

func TestDirectPinDriverDigitalWrite(t *testing.T) {
	d := initTestDirectPinDriver(newGpioTestAdaptor())
	gobottest.Refute(t, d.DigitalWrite(1), nil)

	d = initTestDirectPinDriver(&gpioTestBareAdaptor{})
	gobottest.Assert(t, d.DigitalWrite(1), ErrDigitalWriteUnsupported)
}

func TestDirectPinDriverDigitalRead(t *testing.T) {
	d := initTestDirectPinDriver(newGpioTestAdaptor())
	ret, err := d.DigitalRead()
	gobottest.Assert(t, ret, 1)

	d = initTestDirectPinDriver(&gpioTestBareAdaptor{})
	ret, err = d.DigitalRead()
	gobottest.Assert(t, err, ErrDigitalReadUnsupported)
}

func TestDirectPinDriverPwmWrite(t *testing.T) {
	d := initTestDirectPinDriver(newGpioTestAdaptor())
	gobottest.Refute(t, d.PwmWrite(1), nil)

	d = initTestDirectPinDriver(&gpioTestBareAdaptor{})
	gobottest.Assert(t, d.PwmWrite(1), ErrPwmWriteUnsupported)
}

func TestDirectPinDriverServoWrite(t *testing.T) {
	d := initTestDirectPinDriver(newGpioTestAdaptor())
	gobottest.Refute(t, d.ServoWrite(1), nil)

	d = initTestDirectPinDriver(&gpioTestBareAdaptor{})
	gobottest.Assert(t, d.ServoWrite(1), ErrServoWriteUnsupported)
}

func TestDirectPinDriverDefaultName(t *testing.T) {
	d := initTestDirectPinDriver(newGpioTestAdaptor())
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Direct"), true)
}

func TestDirectPinDriverSetName(t *testing.T) {
	d := initTestDirectPinDriver(newGpioTestAdaptor())
	d.SetName("mybot")
	gobottest.Assert(t, d.Name(), "mybot")
}
