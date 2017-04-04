package gpio

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*BuzzerDriver)(nil)

func initTestBuzzerDriver(conn DigitalWriter) *BuzzerDriver {
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

func TestBuzzerDriverOnError(t *testing.T) {
	a := newGpioTestAdaptor()
	d := initTestBuzzerDriver(a)
	a.TestAdaptorDigitalWrite(func() (err error) {
		return errors.New("write error")
	})

	gobottest.Assert(t, d.On(), errors.New("write error"))
}

func TestBuzzerDriverOffError(t *testing.T) {
	a := newGpioTestAdaptor()
	d := initTestBuzzerDriver(a)
	a.TestAdaptorDigitalWrite(func() (err error) {
		return errors.New("write error")
	})

	gobottest.Assert(t, d.Off(), errors.New("write error"))
}

func TestBuzzerDriverToneError(t *testing.T) {
	a := newGpioTestAdaptor()
	d := initTestBuzzerDriver(a)
	a.TestAdaptorDigitalWrite(func() (err error) {
		return errors.New("write error")
	})

	gobottest.Assert(t, d.Tone(100, 0.01), errors.New("write error"))
}
