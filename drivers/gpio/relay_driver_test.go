package gpio

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*RelayDriver)(nil)

func initTestRelayDriver(conn DigitalWriter) *RelayDriver {
	testAdaptorDigitalWrite = func() (err error) {
		return nil
	}
	testAdaptorPwmWrite = func() (err error) {
		return nil
	}
	return NewRelayDriver(conn, "1")
}

func TestRelayDriverName(t *testing.T) {
	g := initTestRelayDriver(newGpioTestAdaptor())
	gobottest.Refute(t, g.Connection(), nil)
	gobottest.Assert(t, strings.HasPrefix(g.Name(), "Relay"), true)
}

func TestRelayDriverStart(t *testing.T) {
	d := initTestRelayDriver(newGpioTestAdaptor())
	gobottest.Assert(t, d.Start(), nil)
}

func TestRelayDriverHalt(t *testing.T) {
	d := initTestRelayDriver(newGpioTestAdaptor())
	gobottest.Assert(t, d.Halt(), nil)
}

func TestRelayDriverToggle(t *testing.T) {
	d := initTestRelayDriver(newGpioTestAdaptor())
	d.Off()
	d.Toggle()
	gobottest.Assert(t, d.State(), true)
	d.Toggle()
	gobottest.Assert(t, d.State(), false)
}
