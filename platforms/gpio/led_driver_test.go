package gpio

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestLedDriver() *LedDriver {
	return NewLedDriver(newGpioTestAdaptor("adaptor"), "myLed", "1")
}

func TestLedDriverStart(t *testing.T) {
	d := initTestLedDriver()
	gobot.Expect(t, d.Start(), true)
}

func TestLedDriverHalt(t *testing.T) {
	d := initTestLedDriver()
	gobot.Expect(t, d.Halt(), true)
}

func TestLedDriverInit(t *testing.T) {
	d := initTestLedDriver()
	gobot.Expect(t, d.Init(), true)
}

func TestLedDriverOn(t *testing.T) {
	d := initTestLedDriver()
	gobot.Expect(t, d.On(), true)
	gobot.Expect(t, d.IsOn(), true)
}

func TestLedDriverOff(t *testing.T) {
	d := initTestLedDriver()
	gobot.Expect(t, d.Off(), true)
	gobot.Expect(t, d.IsOff(), true)
}

func TestLedDriverToggle(t *testing.T) {
	d := initTestLedDriver()
	d.Off()
	d.Toggle()
	gobot.Expect(t, d.IsOn(), true)
	d.Toggle()
	gobot.Expect(t, d.IsOff(), true)
}

func TestLedDriverBrightness(t *testing.T) {
	d := initTestLedDriver()
	d.Brightness(150)
}
