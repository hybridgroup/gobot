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
	gobot.Assert(t, len(d.Start()), 0)
}

func TestDirectPinDriverHalt(t *testing.T) {
	d := initTestDirectPinDriver()
	gobot.Assert(t, len(d.Halt()), 0)
}

func TestDirectPinDriverDigitalRead(t *testing.T) {
	d := initTestDirectPinDriver()
	val, _ := d.DigitalRead()
	gobot.Assert(t, val, 1)
}

func TestDirectPinDriverDigitalWrite(t *testing.T) {
	d := initTestDirectPinDriver()
	d.DigitalWrite(1)
}

func TestDirectPinDriverAnalogRead(t *testing.T) {
	d := initTestDirectPinDriver()
	val, _ := d.AnalogRead()
	gobot.Assert(t, val, 99)
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
