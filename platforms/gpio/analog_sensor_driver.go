package gpio

import (
	"github.com/hybridgroup/gobot"
)

type AnalogSensorDriver struct {
	gobot.Driver
	Adaptor AnalogReader
}

func NewAnalogSensor(a AnalogReader) *AnalogSensorDriver {
	return &AnalogSensorDriver{
		Driver: gobot.Driver{
			Commands: []string{
				"ReadC",
			},
		},
		Adaptor: a,
	}
}

func (a *AnalogSensorDriver) Start() bool { return true }
func (a *AnalogSensorDriver) Init() bool  { return true }
func (a *AnalogSensorDriver) Halt() bool  { return true }

func (a *AnalogSensorDriver) Read() int {
	return a.Adaptor.AnalogRead(a.Pin)
}
