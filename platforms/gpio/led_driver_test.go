package gpio

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var l *LedDriver

func init() {
	l = NewLedDriver(TestAdaptor{}, "myLed", "1")
}

func TestLedStart(t *testing.T) {
	gobot.Expect(t, l.Start(), true)
}

func TestLedHalt(t *testing.T) {
	gobot.Expect(t, l.Halt(), true)
}

func TestLedInit(t *testing.T) {
	gobot.Expect(t, l.Init(), true)
}

func TestLedOn(t *testing.T) {
	gobot.Expect(t, l.On(), true)
	gobot.Expect(t, l.IsOn(), true)
}

func TestLedOff(t *testing.T) {
	gobot.Expect(t, l.Off(), true)
	gobot.Expect(t, l.IsOff(), true)
}

func TestLedToggle(t *testing.T) {
	l.Off()
	l.Toggle()
	gobot.Expect(t, l.IsOn(), true)
	l.Toggle()
	gobot.Expect(t, l.IsOff(), true)
}

func TestLedBrightness(t *testing.T) {
	l.Brightness(150)
}
