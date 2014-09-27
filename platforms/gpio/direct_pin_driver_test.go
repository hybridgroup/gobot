package gpio

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestDirectPinDriver() *DirectPinDriver {
	return NewDirectPinDriver(newGpioTestAdaptor("adaptor"), "bot", "1")
}

func TestDirectPinDriverStart(t *testing.T) {
	d := initTestDirectPinDriver()
	gobot.Assert(t, d.Start(), true)
}

func TestDirectPinDriverHalt(t *testing.T) {
	d := initTestDirectPinDriver()
	gobot.Assert(t, d.Halt(), true)
}

func TestDirectPinDriverDigitalRead(t *testing.T) {
	d := initTestDirectPinDriver()
	gobot.Assert(t, d.DigitalRead(), 1)
}

func TestDirectPinDriverDigitalWrite(t *testing.T) {
	d := initTestDirectPinDriver()
	d.DigitalWrite(1)
}

func TestDirectPinDriverAnalogRead(t *testing.T) {
	d := initTestDirectPinDriver()
	gobot.Assert(t, d.AnalogRead(), 99)
}

func TestDirectPinDriverAnalogWrite(t *testing.T) {
	d := initTestDirectPinDriver()
	d.AnalogWrite(100)
}

func TestDirectPinDriverPwmWrite(t *testing.T) {
	d := initTestDirectPinDriver()
	d.PwmWrite(100)
}

func TestDirectPinDriverServoWrite(t *testing.T) {
	d := initTestDirectPinDriver()
	d.ServoWrite(100)
}
