package gobotGPIO

import (
	"github.com/hybridgroup/gobot"
)

type AnalogSensorInterface interface {
	AnalogRead(string) int
}

type AnalogSensor struct {
	gobot.Driver
	Adaptor AnalogSensorInterface
}

func NewAnalogSensor(a AnalogSensorInterface) *AnalogSensor {
	b := new(AnalogSensor)
	b.Adaptor = a
	b.Events = make(map[string]chan interface{})
	b.Commands = []string{
		"ReadC",
	}
	return b
}

func (a *AnalogSensor) Start() bool { return true }
func (a *AnalogSensor) Init() bool  { return true }
func (a *AnalogSensor) Halt() bool  { return true }

func (a *AnalogSensor) Read() int {
	return a.Adaptor.AnalogRead(a.Pin)
}
