package gpio

import (
	"github.com/hybridgroup/gobot"
)

type AnalogSensorDriver struct {
	gobot.Driver
}

func NewAnalogSensorDriver(a AnalogReader, name string, pin string) *AnalogSensorDriver {
	d := &AnalogSensorDriver{
		Driver: *gobot.NewDriver(
			name,
			"AnalogSensorDriver",
			a.(gobot.AdaptorInterface),
			pin,
		),
	}

	d.Driver.AddCommand("Read", func(params map[string]interface{}) interface{} {
		return d.Read()
	})

	return d
}

func (a *AnalogSensorDriver) adaptor() AnalogReader {
	return a.Driver.Adaptor().(AnalogReader)
}

func (a *AnalogSensorDriver) Start() bool { return true }
func (a *AnalogSensorDriver) Init() bool  { return true }
func (a *AnalogSensorDriver) Halt() bool  { return true }

func (a *AnalogSensorDriver) Read() int {
	return a.adaptor().AnalogRead(a.Pin())
}
