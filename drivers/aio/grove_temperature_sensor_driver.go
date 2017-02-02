package aio

import (
	"math"
	"time"

	"gobot.io/x/gobot"
)

var _ gobot.Driver = (*GroveTemperatureSensorDriver)(nil)

// GroveTemperatureSensorDriver represents a Temperature Sensor
// The temperature is reported in degree Celsius
type GroveTemperatureSensorDriver struct {
	name        string
	pin         string
	halt        chan bool
	temperature float64
	interval    time.Duration
	connection  AnalogReader
	gobot.Eventer
}

// NewGroveTemperatureSensorDriver returns a new GroveTemperatureSensorDriver with a polling interval of
// 10 Milliseconds given an AnalogReader and pin.
//
// Optionally accepts:
// 	time.Duration: Interval at which the TemperatureSensor is polled for new information
//
// Adds the following API Commands:
// 	"Read" - See AnalogSensor.Read
func NewGroveTemperatureSensorDriver(a AnalogReader, pin string, v ...time.Duration) *GroveTemperatureSensorDriver {
	d := &GroveTemperatureSensorDriver{
		name:       gobot.DefaultName("GroveTemperatureSensor"),
		connection: a,
		pin:        pin,
		Eventer:    gobot.NewEventer(),
		interval:   10 * time.Millisecond,
		halt:       make(chan bool),
	}

	if len(v) > 0 {
		d.interval = v[0]
	}

	d.AddEvent(Data)
	d.AddEvent(Error)

	return d
}

// Start starts the GroveTemperatureSensorDriver and reads the Sensor at the given interval.
// Emits the Events:
//	Data int - Event is emitted on change and represents the current temperature in celsius from the sensor.
//	Error error - Event is emitted on error reading from the sensor.
func (a *GroveTemperatureSensorDriver) Start() (err error) {
	thermistor := 3975.0
	a.temperature = 0

	go func() {
		for {
			rawValue, err := a.Read()

			resistance := float64(1023.0-rawValue) * 10000 / float64(rawValue)
			newValue := 1/(math.Log(resistance/10000.0)/thermistor+1/298.15) - 273.15

			if err != nil {
				a.Publish(Error, err)
			} else if newValue != a.temperature && newValue != -1 {
				a.temperature = newValue
				a.Publish(Data, a.temperature)
			}
			select {
			case <-time.After(a.interval):
			case <-a.halt:
				return
			}
		}
	}()
	return
}

// Halt stops polling the analog sensor for new information
func (a *GroveTemperatureSensorDriver) Halt() (err error) {
	a.halt <- true
	return
}

// Name returns the GroveTemperatureSensorDrivers name
func (a *GroveTemperatureSensorDriver) Name() string { return a.name }

// SetName sets the GroveTemperatureSensorDrivers name
func (a *GroveTemperatureSensorDriver) SetName(n string) { a.name = n }

// Pin returns the GroveTemperatureSensorDrivers pin
func (a *GroveTemperatureSensorDriver) Pin() string { return a.pin }

// Connection returns the GroveTemperatureSensorDrivers Connection
func (a *GroveTemperatureSensorDriver) Connection() gobot.Connection {
	return a.connection.(gobot.Connection)
}

// Read returns the current Temperature from the Sensor
func (a *GroveTemperatureSensorDriver) Temperature() (val float64) {
	return a.temperature
}

// Read returns the raw reading from the Sensor
func (a *GroveTemperatureSensorDriver) Read() (val int, err error) {
	return a.connection.AnalogRead(a.Pin())
}
