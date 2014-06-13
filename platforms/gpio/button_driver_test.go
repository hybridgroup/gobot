package gpio

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var a *AnalogSensorDriver

func init() {
	b = NewButtonDriver(TestAdaptor{}, "bot", "1")
}

func TestStart(t *testing.T) {
	gobot.Expect(t, a.Start(), true)
}

func TestHalt(t *testing.T) {
	gobot.Expect(t, a.Halt(), true)
}

func TestInit(t *testing.T) {
	gobot.Expect(t, a.Init(), true)
}

func TestReadState(t *testing.T) {
	gobot.Expect(t, b.readState(), 1)
}

func TestActive(t *testing.T) {
	b.update(1)
	gobot.Expect(t, b.Active, true)

	b.update(0)
	gobot.Expect(b.Active, false)
}
