package gpio

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*RelayDriver)(nil)

func initTestRelayDriver() *RelayDriver {
	a := newGpioTestAdaptor()
	a.testAdaptorDigitalWrite = func() (err error) {
		return nil
	}
	a.testAdaptorPwmWrite = func() (err error) {
		return nil
	}
	return NewRelayDriver(a, "1")
}

func TestRelayDriverDefaultName(t *testing.T) {
	g := initTestRelayDriver()
	gobottest.Refute(t, g.Connection(), nil)
	gobottest.Assert(t, strings.HasPrefix(g.Name(), "Relay"), true)
}

func TestRelayDriverSetName(t *testing.T) {
	g := initTestRelayDriver()
	g.SetName("mybot")
	gobottest.Assert(t, g.Name(), "mybot")
}

func TestRelayDriverStart(t *testing.T) {
	d := initTestRelayDriver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestRelayDriverHalt(t *testing.T) {
	d := initTestRelayDriver()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestRelayDriverToggle(t *testing.T) {
	d := initTestRelayDriver()
	d.Off()
	d.Toggle()
	gobottest.Assert(t, d.State(), true)
	d.Toggle()
	gobottest.Assert(t, d.State(), false)
}

func TestRelayDriverCommands(t *testing.T) {
	d := initTestRelayDriver()
	gobottest.Assert(t, d.Command("Off")(nil), nil)
	gobottest.Assert(t, d.State(), false)

	gobottest.Assert(t, d.Command("On")(nil), nil)
	gobottest.Assert(t, d.State(), true)

	gobottest.Assert(t, d.Command("Toggle")(nil), nil)
	gobottest.Assert(t, d.State(), false)
}
