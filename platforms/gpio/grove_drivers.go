package gpio

import (
	"time"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*GroveTouchDriver)(nil)
var _ gobot.Driver = (*GroveSoundSensorDriver)(nil)
var _ gobot.Driver = (*GroveButtonDriver)(nil)
var _ gobot.Driver = (*GroveBuzzerDriver)(nil)
var _ gobot.Driver = (*GroveLightSensorDriver)(nil)
var _ gobot.Driver = (*GrovePiezoVibrationSensorDriver)(nil)
var _ gobot.Driver = (*GroveLedDriver)(nil)
var _ gobot.Driver = (*GroveRotaryDriver)(nil)
var _ gobot.Driver = (*GroveRelayDriver)(nil)

// GroveRelayDriver represents a Relay with a Grove connector
type GroveRelayDriver struct {
	*RelayDriver
}

// NewGroveRelayDriver return a new GroveRelayDriver given a DigitalWriter, name and pin.
//
// Adds the following API Commands:
//	"Toggle" - See RelayDriver.Toggle
//	"On" - See RelayDriver.On
//	"Off" - See RelayDriver.Off
func NewGroveRelayDriver(a DigitalWriter, name string, pin string) *GroveRelayDriver {
	return &GroveRelayDriver{
		RelayDriver: NewRelayDriver(a, name, pin),
	}
}

// GroveRotaryDriver represents an analog rotary dial with a Grove connector
type GroveRotaryDriver struct {
	*AnalogSensorDriver
}

// NewGroveRotaryDriver returns a new GroveRotaryDriver with a polling interval of
// 10 Milliseconds given an AnalogReader, name and pin.
//
// Optinally accepts:
// 	time.Duration: Interval at which the AnalogSensor is polled for new information
//
// Adds the following API Commands:
// 	"Read" - See AnalogSensor.Read
func NewGroveRotaryDriver(a AnalogReader, name string, pin string, v ...time.Duration) *GroveRotaryDriver {
	return &GroveRotaryDriver{
		AnalogSensorDriver: NewAnalogSensorDriver(a, name, pin, v...),
	}
}

// GroveLedDriver represents an LED with a Grove connector
type GroveLedDriver struct {
	*LedDriver
}

// NewGroveLedDriver return a new GroveLedDriver given a DigitalWriter, name and pin.
//
// Adds the following API Commands:
//	"Brightness" - See LedDriver.Brightness
//	"Toggle" - See LedDriver.Toggle
//	"On" - See LedDriver.On
//	"Off" - See LedDriver.Off
func NewGroveLedDriver(a DigitalWriter, name string, pin string) *GroveLedDriver {
	return &GroveLedDriver{
		LedDriver: NewLedDriver(a, name, pin),
	}
}

// GroveLightSensorDriver represents an analog light sensor
// with a Grove connector
type GroveLightSensorDriver struct {
	*AnalogSensorDriver
}

// NewGroveLightSensorDriver returns a new GroveLightSensorDriver with a polling interval of
// 10 Milliseconds given an AnalogReader, name and pin.
//
// Optinally accepts:
// 	time.Duration: Interval at which the AnalogSensor is polled for new information
//
// Adds the following API Commands:
// 	"Read" - See AnalogSensor.Read
func NewGroveLightSensorDriver(a AnalogReader, name string, pin string, v ...time.Duration) *GroveLightSensorDriver {
	return &GroveLightSensorDriver{
		AnalogSensorDriver: NewAnalogSensorDriver(a, name, pin, v...),
	}
}

// GrovePiezoVibrationSensorDriver represents an analog vibration sensor
// with a Grove connector
type GrovePiezoVibrationSensorDriver struct {
	*AnalogSensorDriver
}

// NewGrovePiezoVibrationSensorDriver returns a new GrovePiezoVibrationSensorDriver with a polling interval of
// 10 Milliseconds given an AnalogReader, name and pin.
//
// Optinally accepts:
// 	time.Duration: Interval at which the AnalogSensor is polled for new information
//
// Adds the following API Commands:
// 	"Read" - See AnalogSensor.Read
func NewGrovePiezoVibrationSensorDriver(a AnalogReader, name string, pin string, v ...time.Duration) *GrovePiezoVibrationSensorDriver {
	sensor := &GrovePiezoVibrationSensorDriver{
		AnalogSensorDriver: NewAnalogSensorDriver(a, name, pin, v...),
	}

	sensor.AddEvent(Vibration)

	gobot.On(sensor.Event(Data), func(data interface{}) {
		if data.(int) > 1000 {
			gobot.Publish(sensor.Event(Vibration), data)
		}
	})

	return sensor
}

// GroveBuzzerDriver represents a buzzer
// with a Grove connector
type GroveBuzzerDriver struct {
	*BuzzerDriver
}

// NewGroveBuzzerDriver return a new GroveBuzzerDriver given a DigitalWriter, name and pin.
func NewGroveBuzzerDriver(a DigitalWriter, name string, pin string) *GroveBuzzerDriver {
	return &GroveBuzzerDriver{
		BuzzerDriver: NewBuzzerDriver(a, name, pin),
	}
}

// GroveButtonDriver represents a button sensor
// with a Grove connector
type GroveButtonDriver struct {
	*ButtonDriver
}

// NewGroveButtonDriver returns a new GroveButtonDriver with a polling interval of
// 10 Milliseconds given a DigitalReader, name and pin.
//
// Optinally accepts:
//  time.Duration: Interval at which the ButtonDriver is polled for new information
func NewGroveButtonDriver(a DigitalReader, name string, pin string, v ...time.Duration) *GroveButtonDriver {
	return &GroveButtonDriver{
		ButtonDriver: NewButtonDriver(a, name, pin, v...),
	}
}

// GroveSoundSensorDriver represents a analog sound sensor
// with a Grove connector
type GroveSoundSensorDriver struct {
	*AnalogSensorDriver
}

// NewGroveSoundSensorDriver returns a new GroveSoundSensorDriver with a polling interval of
// 10 Milliseconds given an AnalogReader, name and pin.
//
// Optinally accepts:
// 	time.Duration: Interval at which the AnalogSensor is polled for new information
//
// Adds the following API Commands:
// 	"Read" - See AnalogSensor.Read
func NewGroveSoundSensorDriver(a AnalogReader, name string, pin string, v ...time.Duration) *GroveSoundSensorDriver {
	return &GroveSoundSensorDriver{
		AnalogSensorDriver: NewAnalogSensorDriver(a, name, pin, v...),
	}
}

// GroveTouchDriver represents a touch button sensor
// with a Grove connector
type GroveTouchDriver struct {
	*ButtonDriver
}

// NewGroveTouchDriver returns a new GroveTouchDriver with a polling interval of
// 10 Milliseconds given a DigitalReader, name and pin.
//
// Optinally accepts:
//  time.Duration: Interval at which the ButtonDriver is polled for new information
func NewGroveTouchDriver(a DigitalReader, name string, pin string, v ...time.Duration) *GroveTouchDriver {
	return &GroveTouchDriver{
		ButtonDriver: NewButtonDriver(a, name, pin, v...),
	}
}
