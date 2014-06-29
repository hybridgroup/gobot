package gpio

import (
	"github.com/hybridgroup/gobot"
)

type AnalogSensorDriver struct {
	gobot.Driver
}

func NewAnalogSensorDriver(a AnalogReader, name string, pin string) *AnalogSensorDriver {
	d := &AnalogSensorDriver{
		Driver: gobot.Driver{
			Name: name,
			Pin:  pin,
			Events: map[string]*gobot.Event{
				"data": gobot.NewEvent(),
			},
			Commands: make(map[string]func(map[string]interface{}) interface{}),
			Adaptor:  a.(gobot.AdaptorInterface),
		},
	}
	d.Driver.AddCommand("Read", func(params map[string]interface{}) interface{} {
		return d.Read()
	})
	return d
}

func (a *AnalogSensorDriver) adaptor() AnalogReader {
	return a.Driver.Adaptor.(AnalogReader)
}

func (a *AnalogSensorDriver) Start() bool {
	value := 0
	gobot.Every(a.Interval, func() {
		newValue := a.Read()
		if newValue != value && newValue != -1 {
			value = newValue
			gobot.Publish(a.Events["data"], value)
		}
	})
	return true
}
func (a *AnalogSensorDriver) Init() bool { return true }
func (a *AnalogSensorDriver) Halt() bool { return true }

func (a *AnalogSensorDriver) Read() int {
	return a.adaptor().AnalogRead(a.Pin)
}
