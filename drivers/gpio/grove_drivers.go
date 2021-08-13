package gpio

import (
	"time"
)

// GroveRelayDriver represents a Relay with a Grove connector
type GroveRelayDriver struct {
	*RelayDriver
}

// NewGroveRelayDriver return a new GroveRelayDriver given a DigitalWriter and pin.
//
// Adds the following API Commands:
//	"Toggle" - See RelayDriver.Toggle
//	"On" - See RelayDriver.On
//	"Off" - See RelayDriver.Off
func NewGroveRelayDriver(a DigitalWriter, pin string) *GroveRelayDriver {
	return &GroveRelayDriver{
		RelayDriver: NewRelayDriver(a, pin),
	}
}

// GroveLedDriver represents an LED with a Grove connector
type GroveLedDriver struct {
	*LedDriver
}

// NewGroveLedDriver return a new GroveLedDriver given a DigitalWriter and pin.
//
// Adds the following API Commands:
//	"Brightness" - See LedDriver.Brightness
//	"Toggle" - See LedDriver.Toggle
//	"On" - See LedDriver.On
//	"Off" - See LedDriver.Off
func NewGroveLedDriver(a DigitalWriter, pin string) *GroveLedDriver {
	return &GroveLedDriver{
		LedDriver: NewLedDriver(a, pin),
	}
}

// GroveBuzzerDriver represents a buzzer
// with a Grove connector
type GroveBuzzerDriver struct {
	*BuzzerDriver
}

// NewGroveBuzzerDriver return a new GroveBuzzerDriver given a DigitalWriter and pin.
func NewGroveBuzzerDriver(a DigitalWriter, pin string) *GroveBuzzerDriver {
	return &GroveBuzzerDriver{
		BuzzerDriver: NewBuzzerDriver(a, pin),
	}
}

// GroveButtonDriver represents a button sensor
// with a Grove connector
type GroveButtonDriver struct {
	*ButtonDriver
}

// NewGroveButtonDriver returns a new GroveButtonDriver with a polling interval of
// 10 Milliseconds given a DigitalReader and pin.
//
// Optionally accepts:
//  time.Duration: Interval at which the ButtonDriver is polled for new information
func NewGroveButtonDriver(a DigitalReader, pin string, v ...time.Duration) *GroveButtonDriver {
	return &GroveButtonDriver{
		ButtonDriver: NewButtonDriver(a, pin, v...),
	}
}

// GroveTouchDriver represents a touch button sensor
// with a Grove connector
type GroveTouchDriver struct {
	*ButtonDriver
}

// NewGroveTouchDriver returns a new GroveTouchDriver with a polling interval of
// 10 Milliseconds given a DigitalReader and pin.
//
// Optionally accepts:
//  time.Duration: Interval at which the ButtonDriver is polled for new information
func NewGroveTouchDriver(a DigitalReader, pin string, v ...time.Duration) *GroveTouchDriver {
	return &GroveTouchDriver{
		ButtonDriver: NewButtonDriver(a, pin, v...),
	}
}

// GroveMagneticSwitchDriver represent a magnetic
// switch sensor with a Grove connector
type GroveMagneticSwitchDriver struct {
	*ButtonDriver
}

// NewGroveMagneticSwitchDriver returns a new GroveMagneticSwitchDriver with a polling interval of
// 10 Milliseconds given a DigitalReader, name and pin.
//
// Optionally accepts:
//  time.Duration: Interval at which the ButtonDriver is polled for new information
func NewGroveMagneticSwitchDriver(a DigitalReader, pin string, v ...time.Duration) *GroveMagneticSwitchDriver {
	return &GroveMagneticSwitchDriver{
		ButtonDriver: NewButtonDriver(a, pin, v...),
	}
}

// GroveHumidityTemperatureSensorDriver represents a humidity and
// temperature sensor with a Grove connector.
type GroveHumidityTemperatureSensorDriver struct {
	*GroveDHT11SensorDriver
}

// NewGroveHumidityTemperatureSensorDriver returns a new GroveHumidityTemperatureSensorDriver with
// a polling interval of 1000ms given a DHTReader and pin. Currently underyling sensor should
// support DHT11 and DHT22 sensors. DHT11 interval should not be less than 1000ms.
//
// Optionally accepts:
//  time.Duration: Interval at which the sensor is polled for new information
func NewGroveHumidityTemperatureSensorDriver(dht DHTReader, pin string, v ...time.Duration) *GroveHumidityTemperatureSensorDriver {
	var opts []GroveDHT11SensorOption
	if len(v) > 0 {
		opts = append(opts, WithGroveDHT11SensorInterval(v[0]))
	}

	return &GroveHumidityTemperatureSensorDriver{
		GroveDHT11SensorDriver: NewGroveDHT11SensorDriver(dht, pin, opts...),
	}
}

// DHTReader interface represents an Adaptor which has capabilities to
// read digital humidity and temperature for the GrovePi DHT sensors.
type DHTReader interface {
	ReadDHT(pin string) (temperature, humidity float32, err error)
}
