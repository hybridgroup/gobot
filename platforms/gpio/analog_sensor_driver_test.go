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
	gobot.Assert(t, len(d.Start()), 0)
}

func TestAnalogSensorDriverHalt(t *testing.T) {
	d := initTestAnalogSensorDriver()
	gobot.Assert(t, len(d.Halt()), 0)
}

func TestAnalogSensorDriverRead(t *testing.T) {
	d := initTestAnalogSensorDriver()
	val, _ := d.Read()
	gobot.Assert(t, val, 99)
}
