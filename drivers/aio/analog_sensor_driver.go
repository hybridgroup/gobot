package aio

import (
	"time"

	"gobot.io/x/gobot"
)

// AnalogSensorDriver represents an Analog Sensor
type AnalogSensorDriver struct {
	name       string
	pin        string
	halt       chan bool
	interval   time.Duration
	connection AnalogReader
	gobot.Eventer
	gobot.Commander
}

// NewAnalogSensorDriver returns a new AnalogSensorDriver with a polling interval of
// 10 Milliseconds given an AnalogReader and pin.
//
// Optionally accepts:
// 	time.Duration: Interval at which the AnalogSensor is polled for new information
//
// Adds the following API Commands:
// 	"Read" - See AnalogSensor.Read
func NewAnalogSensorDriver(a AnalogReader, pin string, v ...time.Duration) *AnalogSensorDriver {
	d := &AnalogSensorDriver{
		name:       gobot.DefaultName("AnalogSensor"),
		connection: a,
		pin:        pin,
		Eventer:    gobot.NewEventer(),
		Commander:  gobot.NewCommander(),
		interval:   10 * time.Millisecond,
		halt:       make(chan bool),
	}

	if len(v) > 0 {
		d.interval = v[0]
	}

	d.AddEvent(Data)
	d.AddEvent(Error)

	d.AddCommand("Read", func(params map[string]interface{}) interface{} {
		val, err := d.Read()
		return map[string]interface{}{"val": val, "err": err}
	})

	return d
}

// Start starts the AnalogSensorDriver and reads the Analog Sensor at the given interval.
// Emits the Events:
//	Data int - Event is emitted on change and represents the current reading from the sensor.
//	Error error - Event is emitted on error reading from the sensor.
func (a *AnalogSensorDriver) Start() (err error) {
	var value int = 0
	go func() {
		timer := time.NewTimer(a.interval)
		timer.Stop()
		for {
			newValue, err := a.Read()
			if err != nil {
				a.Publish(a.Event(Error), err)
			} else if newValue != value && newValue != -1 {
				value = newValue
				a.Publish(a.Event(Data), value)
			}

			timer.Reset(a.interval)
			select {
			case <-timer.C:
			case <-a.halt:
				timer.Stop()
				return
			}
		}
	}()
	return
}

// Halt stops polling the analog sensor for new information
func (a *AnalogSensorDriver) Halt() (err error) {
	a.halt <- true
	return
}

// Name returns the AnalogSensorDrivers name
func (a *AnalogSensorDriver) Name() string { return a.name }

// SetName sets the AnalogSensorDrivers name
func (a *AnalogSensorDriver) SetName(n string) { a.name = n }

// Pin returns the AnalogSensorDrivers pin
func (a *AnalogSensorDriver) Pin() string { return a.pin }

// Connection returns the AnalogSensorDrivers Connection
func (a *AnalogSensorDriver) Connection() gobot.Connection { return a.connection.(gobot.Connection) }

// Read returns the current reading from the Analog Sensor
func (a *AnalogSensorDriver) Read() (val int, err error) {
	return a.connection.AnalogRead(a.Pin())
}
