package gpio

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var a *AnalogSensorDriver

func init() {
	a = NewAnalogSensorDriver(TestAdaptor{}, "bot", "1")
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
func TestRead(t *testing.T) {
	gobot.Expect(t, a.Read(), 99)
}
