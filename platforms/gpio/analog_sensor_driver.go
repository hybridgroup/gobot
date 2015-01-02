package gpio

import (
	"time"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*AnalogSensorDriver)(nil)

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
// 10 Milliseconds given an AnalogReader, name and pin.
//
// Optinally accepts:
// 	time.Duration: Interval at which the AnalogSensor is polled for new information
//
// Adds the following API Commands:
// 	"Read" - See AnalogSensor.Read
func NewAnalogSensorDriver(a AnalogReader, name string, pin string, v ...time.Duration) *AnalogSensorDriver {
	d := &AnalogSensorDriver{
		name:       name,
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
func (a *AnalogSensorDriver) Start() (errs []error) {
	value := 0
	go func() {
		for {
			newValue, err := a.Read()
			if err != nil {
				gobot.Publish(a.Event(Error), err)
			} else if newValue != value && newValue != -1 {
				value = newValue
				gobot.Publish(a.Event(Data), value)
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
func (a *AnalogSensorDriver) Halt() (errs []error) {
	a.halt <- true
	return
}

// Name returns the AnalogSensorDrivers name
func (a *AnalogSensorDriver) Name() string { return a.name }

// Pin returns the AnalogSensorDrivers pin
func (a *AnalogSensorDriver) Pin() string { return a.pin }

// Connection returns the AnalogSensorDrivers Connection
func (a *AnalogSensorDriver) Connection() gobot.Connection { return a.connection.(gobot.Connection) }

// Read returns the current reading from the Analog Sensor
func (a *AnalogSensorDriver) Read() (val int, err error) {
	return a.connection.AnalogRead(a.Pin())
}
