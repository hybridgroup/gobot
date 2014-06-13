package gpio

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var l *LedDriver

func init() {
	l = NewLedDriver(TestAdaptor{}, "myLed", "1")
}

func TestStart(t *testing.T) {
	gobot.Expect(t, l.Start(), true)
}

func TestHalt(t *testing.T) {
	gobot.Expect(t, l.Halt(), true)
}

func TestInit(t *testing.T) {
	gobot.Expect(t, l.Init(), true)
}

func TestOn(t *testing.T) {
	gobot.Expect(t, l.On(), true)
	gobot.Expect(t, l.IsOn(), true)
}

func TestOff(t *testing.T) {
	gobot.Expect(t, l.Off(), true)
	gobot.Expect(t, l.IsOff(), true)
}

func TestToggle(t *testing.T) {
	l.Off()
	l.Toggle()
	gobot.Expect(t, l.IsOn(), true)
	l.Toggle()
	gobot.Expect(l.IsOff(), true)
}

func TestBrightness(t *testing.T) {
	l.Brightness(150)
}
