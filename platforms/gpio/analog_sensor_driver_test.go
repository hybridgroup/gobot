package gpio

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestAnalogSensorDriver() *AnalogSensorDriver {
	return NewAnalogSensorDriver(newGpioTestAdaptor("adaptor"), "bot", "1")
}

func TestAnalogSensorDriverStart(t *testing.T) {
	d := initTestAnalogSensorDriver()
	gobot.Assert(t, d.Start(), true)
}

func TestAnalogSensorDriverHalt(t *testing.T) {
	d := initTestAnalogSensorDriver()
	gobot.Assert(t, d.Halt(), true)
}

func TestAnalogSensorDriverInit(t *testing.T) {
	d := initTestAnalogSensorDriver()
	gobot.Assert(t, d.Init(), true)
}

func TestAnalogSensorDriverRead(t *testing.T) {
	d := initTestAnalogSensorDriver()
	gobot.Assert(t, d.Read(), 99)
}
