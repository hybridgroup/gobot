package gpio

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestAnalogSensorDriver() *AnalogSensorDriver {
	return NewAnalogSensorDriver(newGpioTestAdaptor("adaptor"), "bot", "1")
}

func TestAnalogSensorDriverStart(t *testing.T) {
	d := initTestAnalogSensorDriver()
	gobot.Assert(t, d.Start(), nil)
}

func TestAnalogSensorDriverHalt(t *testing.T) {
	d := initTestAnalogSensorDriver()
	gobot.Assert(t, d.Halt(), nil)
}

func TestAnalogSensorDriverRead(t *testing.T) {
	d := initTestAnalogSensorDriver()
	gobot.Assert(t, d.Read(), 99)
}
