package aio

import (
	"time"
)

// GroveRotaryDriver represents an analog rotary dial with a Grove connector
type GroveRotaryDriver struct {
	*AnalogSensorDriver
}

// NewGroveRotaryDriver returns a new GroveRotaryDriver with a polling interval of
// 10 Milliseconds given an AnalogReader and pin.
//
// Optionally accepts:
// 	time.Duration: Interval at which the AnalogSensor is polled for new information
//
// Adds the following API Commands:
// 	"Read" - See AnalogSensor.Read
func NewGroveRotaryDriver(a AnalogReader, pin string, v ...time.Duration) *GroveRotaryDriver {
	return &GroveRotaryDriver{
		AnalogSensorDriver: NewAnalogSensorDriver(a, pin, v...),
	}
}

// GroveLightSensorDriver represents an analog light sensor
// with a Grove connector
type GroveLightSensorDriver struct {
	*AnalogSensorDriver
}

// NewGroveLightSensorDriver returns a new GroveLightSensorDriver with a polling interval of
// 10 Milliseconds given an AnalogReader and pin.
//
// Optionally accepts:
// 	time.Duration: Interval at which the AnalogSensor is polled for new information
//
// Adds the following API Commands:
// 	"Read" - See AnalogSensor.Read
func NewGroveLightSensorDriver(a AnalogReader, pin string, v ...time.Duration) *GroveLightSensorDriver {
	return &GroveLightSensorDriver{
		AnalogSensorDriver: NewAnalogSensorDriver(a, pin, v...),
	}
}

// GrovePiezoVibrationSensorDriver represents an analog vibration sensor
// with a Grove connector
type GrovePiezoVibrationSensorDriver struct {
	*AnalogSensorDriver
}

// NewGrovePiezoVibrationSensorDriver returns a new GrovePiezoVibrationSensorDriver with a polling interval of
// 10 Milliseconds given an AnalogReader and pin.
//
// Optionally accepts:
// 	time.Duration: Interval at which the AnalogSensor is polled for new information
//
// Adds the following API Commands:
// 	"Read" - See AnalogSensor.Read
func NewGrovePiezoVibrationSensorDriver(a AnalogReader, pin string, v ...time.Duration) *GrovePiezoVibrationSensorDriver {
	sensor := &GrovePiezoVibrationSensorDriver{
		AnalogSensorDriver: NewAnalogSensorDriver(a, pin, v...),
	}

	sensor.AddEvent(Vibration)

	sensor.On(sensor.Event(Data), func(data interface{}) {
		if data.(int) > 1000 {
			sensor.Publish(sensor.Event(Vibration), data)
		}
	})

	return sensor
}

// GroveSoundSensorDriver represents a analog sound sensor
// with a Grove connector
type GroveSoundSensorDriver struct {
	*AnalogSensorDriver
}

// NewGroveSoundSensorDriver returns a new GroveSoundSensorDriver with a polling interval of
// 10 Milliseconds given an AnalogReader and pin.
//
// Optionally accepts:
// 	time.Duration: Interval at which the AnalogSensor is polled for new information
//
// Adds the following API Commands:
// 	"Read" - See AnalogSensor.Read
func NewGroveSoundSensorDriver(a AnalogReader, pin string, v ...time.Duration) *GroveSoundSensorDriver {
	return &GroveSoundSensorDriver{
		AnalogSensorDriver: NewAnalogSensorDriver(a, pin, v...),
	}
}
