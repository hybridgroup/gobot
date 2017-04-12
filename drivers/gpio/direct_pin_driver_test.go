package gpio

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*DirectPinDriver)(nil)

func initTestDirectPinDriver() *DirectPinDriver {
	a := newGpioTestAdaptor()
	a.testAdaptorDigitalRead = func() (val int, err error) {
		val = 1
		return
	}
	a.testAdaptorDigitalWrite = func() (err error) {
		return errors.New("write error")
	}
	a.testAdaptorPwmWrite = func() (err error) {
		return errors.New("write error")
	}
	a.testAdaptorServoWrite = func() (err error) {
		return errors.New("write error")
	}
	return NewDirectPinDriver(a, "1")
}

func TestDirectPinDriver(t *testing.T) {
	var ret map[string]interface{}
	var err interface{}

	d := initTestDirectPinDriver()
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
	d := initTestDirectPinDriver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestDirectPinDriverHalt(t *testing.T) {
	d := initTestDirectPinDriver()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestDirectPinDriverOff(t *testing.T) {
	d := initTestDirectPinDriver()
	gobottest.Refute(t, d.Off(), nil)

	a := newGpioTestAdaptor()
	d = NewDirectPinDriver(a, "1")
	gobottest.Assert(t, d.Off(), nil)
}

func TestDirectPinDriverOffNotSupported(t *testing.T) {
	a := &gpioTestBareAdaptor{}
	d := NewDirectPinDriver(a, "1")
	gobottest.Assert(t, d.Off(), errors.New("DigitalWrite is not supported by this platform"))
}

func TestDirectPinDriverOn(t *testing.T) {
	a := newGpioTestAdaptor()
	d := NewDirectPinDriver(a, "1")
	gobottest.Assert(t, d.On(), nil)
}

func TestDirectPinDriverOnError(t *testing.T) {
	d := initTestDirectPinDriver()
	gobottest.Refute(t, d.On(), nil)
}

func TestDirectPinDriverOnNotSupported(t *testing.T) {
	a := &gpioTestBareAdaptor{}
	d := NewDirectPinDriver(a, "1")
	gobottest.Assert(t, d.On(), errors.New("DigitalWrite is not supported by this platform"))
}

func TestDirectPinDriverDigitalWrite(t *testing.T) {
	adaptor := newGpioTestAdaptor()
	d := NewDirectPinDriver(adaptor, "1")
	gobottest.Assert(t, d.DigitalWrite(1), nil)
}

func TestDirectPinDriverDigitalWriteNotSupported(t *testing.T) {
	a := &gpioTestBareAdaptor{}
	d := NewDirectPinDriver(a, "1")
	gobottest.Assert(t, d.DigitalWrite(1), errors.New("DigitalWrite is not supported by this platform"))
}

func TestDirectPinDriverDigitalWriteError(t *testing.T) {
	d := initTestDirectPinDriver()
	gobottest.Refute(t, d.DigitalWrite(1), nil)
}

func TestDirectPinDriverDigitalRead(t *testing.T) {
	d := initTestDirectPinDriver()
	ret, err := d.DigitalRead()
	gobottest.Assert(t, ret, 1)
	gobottest.Assert(t, err, nil)
}

func TestDirectPinDriverDigitalReadNotSupported(t *testing.T) {
	a := &gpioTestBareAdaptor{}
	d := NewDirectPinDriver(a, "1")
	_, e := d.DigitalRead()
	gobottest.Assert(t, e, errors.New("DigitalRead is not supported by this platform"))
}

func TestDirectPinDriverPwmWrite(t *testing.T) {
	a := newGpioTestAdaptor()
	d := NewDirectPinDriver(a, "1")
	gobottest.Assert(t, d.PwmWrite(1), nil)
}

func TestDirectPinDriverPwmWriteNotSupported(t *testing.T) {
	a := &gpioTestBareAdaptor{}
	d := NewDirectPinDriver(a, "1")
	gobottest.Assert(t, d.PwmWrite(1), errors.New("PwmWrite is not supported by this platform"))
}

func TestDirectPinDriverPwmWriteError(t *testing.T) {
	d := initTestDirectPinDriver()
	gobottest.Refute(t, d.PwmWrite(1), nil)
}

func TestDirectPinDriverServoWrite(t *testing.T) {
	a := newGpioTestAdaptor()
	d := NewDirectPinDriver(a, "1")
	gobottest.Assert(t, d.ServoWrite(1), nil)
}

func TestDirectPinDriverServoWriteNotSupported(t *testing.T) {
	a := &gpioTestBareAdaptor{}
	d := NewDirectPinDriver(a, "1")
	gobottest.Assert(t, d.ServoWrite(1), errors.New("ServoWrite is not supported by this platform"))
}

func TestDirectPinDriverServoWriteError(t *testing.T) {
	d := initTestDirectPinDriver()
	gobottest.Refute(t, d.ServoWrite(1), nil)
}

func TestDirectPinDriverDefaultName(t *testing.T) {
	d := initTestDirectPinDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Direct"), true)
}

func TestDirectPinDriverSetName(t *testing.T) {
	d := initTestDirectPinDriver()
	d.SetName("mybot")
	gobottest.Assert(t, d.Name(), "mybot")
}
