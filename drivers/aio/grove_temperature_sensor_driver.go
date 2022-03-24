package aio

import (
	"time"

	"gobot.io/x/gobot"
)

var _ gobot.Driver = (*GroveTemperatureSensorDriver)(nil)

// GroveTemperatureSensorDriver represents a temperature sensor
// The temperature is reported in degree Celsius
type GroveTemperatureSensorDriver struct {
	*TemperatureSensorDriver
}

// NewGroveTemperatureSensorDriver returns a new GroveTemperatureSensorDriver with a polling interval of
// 10 Milliseconds given an AnalogReader and pin.
//
// Optionally accepts:
// 	time.Duration: Interval at which the sensor is polled for new information (given 0 switch the polling off)
//
// Adds the following API Commands:
// 	"Read"      - See AnalogDriverSensor.Read
// 	"ReadValue" - See AnalogDriverSensor.ReadValue
func NewGroveTemperatureSensorDriver(a AnalogReader, pin string, v ...time.Duration) *GroveTemperatureSensorDriver {
	t := NewTemperatureSensorDriver(a, pin, v...)
	ntc := TemperatureSensorNtcConf{TC0: 25, R0: 10000.0, B: 3975} //Ohm, R25=10k
	t.SetNtcScaler(1023, 10000, false, ntc)                        //Ohm, reference value: 1023, series R: 10k

	d := &GroveTemperatureSensorDriver{
		TemperatureSensorDriver: t,
	}
	d.SetName(gobot.DefaultName("GroveTemperatureSensor"))

	return d
}

// Temperature returns the last read temperature from the sensor.
func (t *TemperatureSensorDriver) Temperature() (val float64) {
	return t.Value()
}
