package gpio

import (
	"errors"
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestDirectPinDriver(conn gobot.Connection) *DirectPinDriver {
	testAdaptorDigitalRead = func() (val int, err error) {
		val = 1
		return
	}
	testAdaptorDigitalWrite = func() (err error) {
		return errors.New("write error")
	}
	testAdaptorAnalogRead = func() (val int, err error) {
		val = 80
		return
	}
	testAdaptorPwmWrite = func() (err error) {
		return errors.New("write error")
	}
	testAdaptorServoWrite = func() (err error) {
		return errors.New("write error")
	}
	return NewDirectPinDriver(conn, "bot", "1")
}

func TestDirectPinDriver(t *testing.T) {
	var ret map[string]interface{}
	var err interface{}

	d := initTestDirectPinDriver(newGpioTestAdaptor("adaptor"))
	gobot.Assert(t, d.Name(), "bot")
	gobot.Assert(t, d.Pin(), "1")
	gobot.Assert(t, d.Connection().Name(), "adaptor")

	ret = d.Command("DigitalRead")(nil).(map[string]interface{})

	gobot.Assert(t, ret["val"].(int), 1)
	gobot.Assert(t, ret["err"], nil)

	err = d.Command("DigitalWrite")(map[string]interface{}{"level": "1"})
	gobot.Assert(t, err.(error), errors.New("write error"))

	ret = d.Command("AnalogRead")(nil).(map[string]interface{})

	gobot.Assert(t, ret["val"].(int), 80)
	gobot.Assert(t, ret["err"], nil)

	err = d.Command("PwmWrite")(map[string]interface{}{"level": "1"})
	gobot.Assert(t, err.(error), errors.New("write error"))

	err = d.Command("ServoWrite")(map[string]interface{}{"level": "1"})
	gobot.Assert(t, err.(error), errors.New("write error"))
}

func TestDirectPinDriverStart(t *testing.T) {
	d := initTestDirectPinDriver(newGpioTestAdaptor("adaptor"))
	gobot.Assert(t, len(d.Start()), 0)
}

func TestDirectPinDriverHalt(t *testing.T) {
	d := initTestDirectPinDriver(newGpioTestAdaptor("adaptor"))
	gobot.Assert(t, len(d.Halt()), 0)
}

func TestDirectPinDriverDigitalWrite(t *testing.T) {
	d := initTestDirectPinDriver(newGpioTestAdaptor("adaptor"))
	gobot.Refute(t, d.DigitalWrite(1), nil)

	d = initTestDirectPinDriver(&gpioTestBareAdaptor{})
	gobot.Assert(t, d.DigitalWrite(1), ErrDigitalWriteUnsupported)
}

func TestDirectPinDriverDigitalRead(t *testing.T) {
	d := initTestDirectPinDriver(newGpioTestAdaptor("adaptor"))
	ret, err := d.DigitalRead()
	gobot.Assert(t, ret, 1)

	d = initTestDirectPinDriver(&gpioTestBareAdaptor{})
	ret, err = d.DigitalRead()
	gobot.Assert(t, err, ErrDigitalReadUnsupported)
}

func TestDirectPinDriverAnalogRead(t *testing.T) {
	d := initTestDirectPinDriver(newGpioTestAdaptor("adaptor"))
	ret, err := d.AnalogRead()
	gobot.Assert(t, ret, 80)

	d = initTestDirectPinDriver(&gpioTestBareAdaptor{})
	ret, err = d.AnalogRead()
	gobot.Assert(t, err, ErrAnalogReadUnsupported)
}

func TestDirectPinDriverPwmWrite(t *testing.T) {
	d := initTestDirectPinDriver(newGpioTestAdaptor("adaptor"))
	gobot.Refute(t, d.PwmWrite(1), nil)

	d = initTestDirectPinDriver(&gpioTestBareAdaptor{})
	gobot.Assert(t, d.PwmWrite(1), ErrPwmWriteUnsupported)
}
func TestDirectPinDriverDigitalWrie(t *testing.T) {
	d := initTestDirectPinDriver(newGpioTestAdaptor("adaptor"))
	gobot.Refute(t, d.ServoWrite(1), nil)

	d = initTestDirectPinDriver(&gpioTestBareAdaptor{})
	gobot.Assert(t, d.ServoWrite(1), ErrServoWriteUnsupported)
}
