package aio

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
)

// GroveTemperatureSensorDriver represents a temperature sensor
// The temperature is reported in degree Celsius
type GroveTemperatureSensorDriver struct {
	*TemperatureSensorDriver
}

// NewGroveTemperatureSensorDriver returns a new driver for grove temperature sensor, given an AnalogReader and pin.
//
// Supported options: see [aio.NewAnalogSensorDriver]
// Adds the following API Commands: see [aio.NewAnalogSensorDriver]
func NewGroveTemperatureSensorDriver(a AnalogReader, pin string, opts ...interface{}) *GroveTemperatureSensorDriver {
	t := NewTemperatureSensorDriver(a, pin, opts...)
	ntc := TemperatureSensorNtcConf{TC0: 25, R0: 10000.0, B: 3975} // Ohm, R25=10k
	t.SetNtcScaler(1023, 10000, false, ntc)                        // Ohm, reference value: 1023, series R: 10k

	d := &GroveTemperatureSensorDriver{
		TemperatureSensorDriver: t,
	}
	d.driverCfg.name = gobot.DefaultName("GroveTemperatureSensor")

	for _, opt := range opts {
		switch o := opt.(type) {
		case optionApplier:
			o.apply(d.driverCfg)
		case sensorOptionApplier:
			o.apply(d.sensorCfg)
		case time.Duration:
			// TODO this is only for backward compatibility and will be removed after version 2.x
			d.sensorCfg.readInterval = o
		default:
			panic(fmt.Sprintf("'%s' can not be applied on '%s'", opt, d.driverCfg.name))
		}
	}

	return d
}

// Temperature returns the last read temperature from the sensor.
func (t *TemperatureSensorDriver) Temperature() float64 {
	return t.Value()
}
