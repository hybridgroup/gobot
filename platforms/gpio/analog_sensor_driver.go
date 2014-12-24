package gpio

import (
	"time"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*AnalogSensorDriver)(nil)

// Represents an Analog Sensor
type AnalogSensorDriver struct {
	name       string
	pin        string
	interval   time.Duration
	connection AnalogReader
	gobot.Eventer
	gobot.Commander
}

// NewAnalogSensorDriver returns a new AnalogSensorDriver given an AnalogReader, name and pin.
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

// Starts the AnalogSensorDriver and reads the Analog Sensor at the given Driver.Interval().
// Returns true on successful start of the driver.
// Emits the Events:
//	"data" int - Event is emitted on change and represents the current reading from the sensor.
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
			<-time.After(a.interval)
		}
	}()
	return
}

// Halt returns true on a successful halt of the driver
func (a *AnalogSensorDriver) Halt() (errs []error)         { return }
func (a *AnalogSensorDriver) Name() string                 { return a.name }
func (a *AnalogSensorDriver) Pin() string                  { return a.pin }
func (a *AnalogSensorDriver) Connection() gobot.Connection { return a.connection.(gobot.Connection) }

// Read returns the current reading from the Analog Sensor
func (a *AnalogSensorDriver) Read() (val int, err error) {
	return a.connection.AnalogRead(a.Pin())
}
