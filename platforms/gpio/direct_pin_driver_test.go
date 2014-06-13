package gpio

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var d *DirectPinDriver

func init() {
	d = NewDirectPinDriver(TestAdaptor{}, "bot", "1")
}

func TestStart(t *testing.T) {
	gobot.Expect(t, a.Start(), true)
}

func TestHalt(t *testing.T) {
	gobot.Expect(t, a.Halt(), true)
}

func TestInit(t *testing.T) {
	gobot.Expect(t, a.Init(), true)
}

func TestDigitalRead(t *testing.T) {
	gobot.Expect(t, d.DigitalRead(), 1)
}

func TestDigitalWrite(t *testing.T) {
	d.DigitalWrite(1)
}

func TestAnalogRead(t *testing.T) {
	gobot.Expect(t, d.AnalogRead(), 99)
}

func TestAnalogWrite(t *testing.T) {
	d.AnalogWrite(100)
}

func TestPwmWrite(t *testing.T) {
	d.PwmWrite(100)
}

func TestServoWrite(t *testing.T) {
	d.ServoWrite(100)
}
