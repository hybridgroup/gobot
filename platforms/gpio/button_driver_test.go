package gpio

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestButtonDriver() *ButtonDriver {
	return NewButtonDriver(newGpioTestAdaptor("adaptor"), "bot", "1")
}

func TestButtonDriverStart(t *testing.T) {
	d := initTestButtonDriver()
	gobot.Assert(t, len(d.Start()), 0)
}

func TestButtonDriverHalt(t *testing.T) {
	d := initTestButtonDriver()
	gobot.Assert(t, len(d.Halt()), 0)
}

func TestButtonDriverActive(t *testing.T) {
	d := initTestButtonDriver()
	d.update(1)
	gobot.Assert(t, d.Active, true)

	d.update(0)
	gobot.Assert(t, d.Active, false)
}
