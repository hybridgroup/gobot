package gpio

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestDirectPinDriver() *DirectPinDriver {
	return NewDirectPinDriver(TestAdaptor{}, "bot", "1")
}

func TestDirectPinDriverStart(t *testing.T) {
	d := initTestDirectPinDriver()
	gobot.Expect(t, d.Start(), true)
}

func TestDirectPinDriverHalt(t *testing.T) {
	d := initTestDirectPinDriver()
	gobot.Expect(t, d.Halt(), true)
}

func TestDirectPinDriverInit(t *testing.T) {
	d := initTestDirectPinDriver()
	gobot.Expect(t, d.Init(), true)
}

func TestDirectPinDriverDigitalRead(t *testing.T) {
	d := initTestDirectPinDriver()
	gobot.Expect(t, d.DigitalRead(), 1)
}

func TestDirectPinDriverDigitalWrite(t *testing.T) {
	d := initTestDirectPinDriver()
	d.DigitalWrite(1)
}

func TestDirectPinDriverAnalogRead(t *testing.T) {
	d := initTestDirectPinDriver()
	gobot.Expect(t, d.AnalogRead(), 99)
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
