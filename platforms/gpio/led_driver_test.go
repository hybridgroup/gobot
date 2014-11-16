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
	gobot.Assert(t, d.Start(), nil)
}

func TestLedDriverHalt(t *testing.T) {
	d := initTestLedDriver()
	gobot.Assert(t, d.Halt(), nil)
}

func TestLedDriverOn(t *testing.T) {
	d := initTestLedDriver()
	d.On()
	gobot.Assert(t, d.State(), true)
}

func TestLedDriverOff(t *testing.T) {
	d := initTestLedDriver()
	d.Off()
	gobot.Assert(t, d.State(), false)
}

func TestLedDriverToggle(t *testing.T) {
	d := initTestLedDriver()
	d.Off()
	d.Toggle()
	gobot.Assert(t, d.State(), true)
	d.Toggle()
	gobot.Assert(t, d.State(), false)
}

func TestLedDriverBrightness(t *testing.T) {
	d := initTestLedDriver()
	d.Brightness(150)
}
