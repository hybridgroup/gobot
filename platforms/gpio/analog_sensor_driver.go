package gpio

import (
	"github.com/hybridgroup/gobot"
)

var _ gobot.DriverInterface = (*AnalogSensorDriver)(nil)

// Represents an Analog Sensor
type AnalogSensorDriver struct {
	gobot.Driver
}

// NewAnalogSensorDriver returns a new AnalogSensorDriver given an AnalogReader, name and pin.
//
// Adds the following API Commands:
// 	"Read" - See AnalogSensor.Read
func NewAnalogSensorDriver(a AnalogReader, name string, pin string) *AnalogSensorDriver {
	d := &AnalogSensorDriver{
		Driver: *gobot.NewDriver(
			name,
			"AnalogSensorDriver",
			a.(gobot.AdaptorInterface),
			pin,
		),
	}

	d.AddEvent("data")
	d.AddCommand("Read", func(params map[string]interface{}) interface{} {
		return d.Read()
	})

	return d
}

func (a *AnalogSensorDriver) adaptor() AnalogReader {
	return a.Adaptor().(AnalogReader)
}

// Starts the AnalogSensorDriver and reads the Analog Sensor at the given Driver.Interval().
// Returns true on successful start of the driver.
// Emits the Events:
//	"data" int - Event is emitted on change and represents the current reading from the sensor.
func (a *AnalogSensorDriver) Start() error {
	value := 0
	gobot.Every(a.Interval(), func() {
		newValue := a.Read()
		if newValue != value && newValue != -1 {
			value = newValue
			gobot.Publish(a.Event("data"), value)
		}
	})
	return nil
}

// Halt returns true on a successful halt of the driver
func (a *AnalogSensorDriver) Halt() error { return nil }

// Read returns the current reading from the Analog Sensor
func (a *AnalogSensorDriver) Read() int {
	return a.adaptor().AnalogRead(a.Pin())
}
