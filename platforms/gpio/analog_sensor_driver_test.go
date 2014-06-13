package gpio

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var a *AnalogSensorDriver

func init() {
	a = NewAnalogSensorDriver(TestAdaptor{}, "bot", "1")
}

func TestAnalogSensorStart(t *testing.T) {
	gobot.Expect(t, a.Start(), true)
}

func TestAnalogSensorHalt(t *testing.T) {
	gobot.Expect(t, a.Halt(), true)
}

func TestAnalogSensorInit(t *testing.T) {
	gobot.Expect(t, a.Init(), true)
}

func TestAnalogSensorRead(t *testing.T) {
	gobot.Expect(t, a.Read(), 99)
}
