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

func TestBuzzerDriverDefaultName(t *testing.T) {
	g := initTestBuzzerDriver(newGpioTestAdaptor())
	gobottest.Assert(t, strings.HasPrefix(g.Name(), "Buzzer"), true)
}

func TestBuzzerDriverSetName(t *testing.T) {
	g := initTestBuzzerDriver(newGpioTestAdaptor())
	g.SetName("mybot")
	gobottest.Assert(t, g.Name(), "mybot")
}

func TestBuzzerDriverStart(t *testing.T) {
	d := initTestBuzzerDriver(newGpioTestAdaptor())
	gobottest.Assert(t, d.Start(), nil)
}

func TestBuzzerDriverHalt(t *testing.T) {
	d := initTestBuzzerDriver(newGpioTestAdaptor())
	gobottest.Assert(t, d.Halt(), nil)
}

func TestBuzzerDriverToggle(t *testing.T) {
	d := initTestBuzzerDriver(newGpioTestAdaptor())
	d.Off()
	d.Toggle()
	gobottest.Assert(t, d.State(), true)
	d.Toggle()
	gobottest.Assert(t, d.State(), false)
}

func TestBuzzerDriverTone(t *testing.T) {
	d := initTestBuzzerDriver(newGpioTestAdaptor())
	gobottest.Assert(t, d.Tone(100, 0.01), nil)
}
