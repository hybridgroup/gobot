package gpio

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*BuzzerDriver)(nil)

func initTestBuzzerDriver(conn DigitalWriter) *BuzzerDriver {
	testAdaptorDigitalWrite = func() (err error) {
		return nil
	}
	testAdaptorPwmWrite = func() (err error) {
		return nil
	}
	return NewBuzzerDriver(conn, "1")
}

func TestBuzzerDriverName(t *testing.T) {
	g := initTestBuzzerDriver(newGpioTestAdaptor())
	gobottest.Refute(t, g.Connection(), nil)
	gobottest.Assert(t, strings.HasPrefix(g.Name(), "Buzzer"), true)
}

func TestBuzzerDriverStart(t *testing.T) {
	d := initTestBuzzerDriver(newGpioTestAdaptor())
	gobottest.Assert(t, d.Start(), nil)
}

func TestBuzzerDriverHalt(t *testing.T) {
	d := initTestBuzzerDriver(newGpioTestAdaptor())
	gobottest.Assert(t, d.Halt(), nil)
}
