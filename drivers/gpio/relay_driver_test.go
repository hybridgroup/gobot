package gpio

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*RelayDriver)(nil)

// Helper to return low/high value for testing
func (l *RelayDriver) High() bool { return l.high }

func initTestRelayDriver() (*RelayDriver, *gpioTestAdaptor) {
	a := newGpioTestAdaptor()
	a.testAdaptorDigitalWrite = func(string, byte) (err error) {
		return nil
	}
	a.testAdaptorPwmWrite = func(string, byte) (err error) {
		return nil
	}
	return NewRelayDriver(a, "1"), a
}

func TestRelayDriverDefaultName(t *testing.T) {
	g, _ := initTestRelayDriver()
	gobottest.Refute(t, g.Connection(), nil)
	gobottest.Assert(t, strings.HasPrefix(g.Name(), "Relay"), true)
}

func TestRelayDriverSetName(t *testing.T) {
	g, _ := initTestRelayDriver()
	g.SetName("mybot")
	gobottest.Assert(t, g.Name(), "mybot")
}

func TestRelayDriverStart(t *testing.T) {
	d, _ := initTestRelayDriver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestRelayDriverHalt(t *testing.T) {
	d, _ := initTestRelayDriver()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestRelayDriverToggle(t *testing.T) {
	d, a := initTestRelayDriver()
	var lastVal byte
	a.TestAdaptorDigitalWrite(func(pin string, val byte) error {
		lastVal = val
		return nil
	})

	d.Off()
	gobottest.Assert(t, d.State(), false)
	gobottest.Assert(t, lastVal, byte(0))
	d.Toggle()
	gobottest.Assert(t, d.State(), true)
	gobottest.Assert(t, lastVal, byte(1))
	d.Toggle()
	gobottest.Assert(t, d.State(), false)
	gobottest.Assert(t, lastVal, byte(0))
}

func TestRelayDriverToggleInverted(t *testing.T) {
	d, a := initTestRelayDriver()
	var lastVal byte
	a.TestAdaptorDigitalWrite(func(pin string, val byte) error {
		lastVal = val
		return nil
	})

	d.Inverted = true
	d.Off()
	gobottest.Assert(t, d.State(), false)
	gobottest.Assert(t, lastVal, byte(1))
	d.Toggle()
	gobottest.Assert(t, d.State(), true)
	gobottest.Assert(t, lastVal, byte(0))
	d.Toggle()
	gobottest.Assert(t, d.State(), false)
	gobottest.Assert(t, lastVal, byte(1))
}

func TestRelayDriverCommands(t *testing.T) {
	d, a := initTestRelayDriver()
	var lastVal byte
	a.TestAdaptorDigitalWrite(func(pin string, val byte) error {
		lastVal = val
		return nil
	})

	gobottest.Assert(t, d.Command("Off")(nil), nil)
	gobottest.Assert(t, d.State(), false)
	gobottest.Assert(t, lastVal, byte(0))

	gobottest.Assert(t, d.Command("On")(nil), nil)
	gobottest.Assert(t, d.State(), true)
	gobottest.Assert(t, lastVal, byte(1))

	gobottest.Assert(t, d.Command("Toggle")(nil), nil)
	gobottest.Assert(t, d.State(), false)
	gobottest.Assert(t, lastVal, byte(0))
}

func TestRelayDriverCommandsInverted(t *testing.T) {
	d, a := initTestRelayDriver()
	var lastVal byte
	a.TestAdaptorDigitalWrite(func(pin string, val byte) error {
		lastVal = val
		return nil
	})
	d.Inverted = true

	gobottest.Assert(t, d.Command("Off")(nil), nil)
	gobottest.Assert(t, d.High(), true)
	gobottest.Assert(t, d.State(), false)
	gobottest.Assert(t, lastVal, byte(1))

	gobottest.Assert(t, d.Command("On")(nil), nil)
	gobottest.Assert(t, d.High(), false)
	gobottest.Assert(t, d.State(), true)
	gobottest.Assert(t, lastVal, byte(0))

	gobottest.Assert(t, d.Command("Toggle")(nil), nil)
	gobottest.Assert(t, d.High(), true)
	gobottest.Assert(t, d.State(), false)
	gobottest.Assert(t, lastVal, byte(1))
}
