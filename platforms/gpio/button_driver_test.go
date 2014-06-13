package gpio

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestButtonDriver() *ButtonDriver {
	return NewButtonDriver(TestAdaptor{}, "bot", "1")
}

func TestButtonDriverStart(t *testing.T) {
	d := initTestButtonDriver()
	gobot.Expect(t, d.Start(), true)
}

func TestButtonDriverHalt(t *testing.T) {
	d := initTestButtonDriver()
	gobot.Expect(t, d.Halt(), true)
}

func TestButtonDriverInit(t *testing.T) {
	d := initTestButtonDriver()
	gobot.Expect(t, d.Init(), true)
}

func TestButtonDriverReadState(t *testing.T) {
	d := initTestButtonDriver()
	gobot.Expect(t, d.readState(), 1)
}

func TestButtonDriverActive(t *testing.T) {
	d := initTestButtonDriver()
	d.update(1)
	gobot.Expect(t, d.Active, true)

	d.update(0)
	gobot.Expect(t, d.Active, false)
}
