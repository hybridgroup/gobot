package gpio

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestButtonDriver() *ButtonDriver {
	return NewButtonDriver(newGpioTestAdaptor("adaptor"), "bot", "1")
}

func TestButtonDriverStart(t *testing.T) {
	d := initTestButtonDriver()
	gobot.Assert(t, d.Start(), true)
}

func TestButtonDriverHalt(t *testing.T) {
	d := initTestButtonDriver()
	gobot.Assert(t, d.Halt(), true)
}

func TestButtonDriverInit(t *testing.T) {
	d := initTestButtonDriver()
	gobot.Assert(t, d.Init(), true)
}

func TestButtonDriverReadState(t *testing.T) {
	d := initTestButtonDriver()
	gobot.Assert(t, d.readState(), 1)
}

func TestButtonDriverActive(t *testing.T) {
	d := initTestButtonDriver()
	d.update(1)
	gobot.Assert(t, d.Active, true)

	d.update(0)
	gobot.Assert(t, d.Active, false)
}
