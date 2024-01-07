package aio

import "gobot.io/x/gobot/v2"

// GroveRotaryDriver represents an analog rotary dial with a Grove connector
type GroveRotaryDriver struct {
	*AnalogSensorDriver
}

// NewGroveRotaryDriver returns a new driver for grove rotary dial, given an AnalogReader and pin.
//
// Supported options: see [aio.NewAnalogSensorDriver]
// Adds the following API Commands: see [aio.NewAnalogSensorDriver]
func NewGroveRotaryDriver(a AnalogReader, pin string, opts ...interface{}) *GroveRotaryDriver {
	d := GroveRotaryDriver{
		AnalogSensorDriver: NewAnalogSensorDriver(a, pin, opts...),
	}
	d.driverCfg.name = gobot.DefaultName("GroveRotary")

	return &d
}

// GroveLightSensorDriver represents an analog light sensor
// with a Grove connector
type GroveLightSensorDriver struct {
	*AnalogSensorDriver
}

// NewGroveLightSensorDriver returns a new driver for grove light sensor, given an AnalogReader and pin.
//
// Supported options: see [aio.NewAnalogSensorDriver]
// Adds the following API Commands: see [aio.NewAnalogSensorDriver]
func NewGroveLightSensorDriver(a AnalogReader, pin string, opts ...interface{}) *GroveLightSensorDriver {
	d := GroveLightSensorDriver{
		AnalogSensorDriver: NewAnalogSensorDriver(a, pin, opts...),
	}
	d.driverCfg.name = gobot.DefaultName("GroveLightSensor")

	return &d
}

// GrovePiezoVibrationSensorDriver represents an analog vibration sensor with a Grove connector
type GrovePiezoVibrationSensorDriver struct {
	*AnalogSensorDriver
}

// NewGrovePiezoVibrationSensorDriver returns a new driver for grove piezo vibration sensor, given an AnalogReader
// and pin.
//
// Supported options: see [aio.NewAnalogSensorDriver]
// Adds the following API Commands: see [aio.NewAnalogSensorDriver]
func NewGrovePiezoVibrationSensorDriver(
	a AnalogReader,
	pin string,
	opts ...interface{},
) *GrovePiezoVibrationSensorDriver {
	d := &GrovePiezoVibrationSensorDriver{
		AnalogSensorDriver: NewAnalogSensorDriver(a, pin, opts...),
	}
	d.driverCfg.name = gobot.DefaultName("GrovePiezoVibrationSensor")

	d.AddEvent(Vibration)

	if err := d.On(d.Event(Data), func(data interface{}) {
		if data.(int) > 1000 { //nolint:forcetypeassert // no error return value, so there is no better way
			d.Publish(d.Event(Vibration), data)
		}
	}); err != nil {
		panic(err)
	}

	return d
}

// GroveSoundSensorDriver represents a analog sound sensor with a Grove connector
type GroveSoundSensorDriver struct {
	*AnalogSensorDriver
}

// NewGroveSoundSensorDriver returns a new driver for grove sound sensor, given an AnalogReader and pin.
//
// Supported options: see [aio.NewAnalogSensorDriver]
// Adds the following API Commands: see [aio.NewAnalogSensorDriver]
func NewGroveSoundSensorDriver(a AnalogReader, pin string, opts ...interface{}) *GroveSoundSensorDriver {
	d := GroveSoundSensorDriver{
		AnalogSensorDriver: NewAnalogSensorDriver(a, pin, opts...),
	}
	d.driverCfg.name = gobot.DefaultName("GroveSoundSensor")

	return &d
}
