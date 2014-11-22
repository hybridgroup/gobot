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
	connection gobot.Connection
	gobot.Eventer
	gobot.Commander
}

// NewAnalogSensorDriver returns a new AnalogSensorDriver given an AnalogReader, name and pin.
//
// Adds the following API Commands:
// 	"Read" - See AnalogSensor.Read
func NewAnalogSensorDriver(a AnalogReader, name string, pin string) *AnalogSensorDriver {
	d := &AnalogSensorDriver{
		name:       name,
		connection: a.(gobot.Connection),
		pin:        pin,
		Eventer:    gobot.NewEventer(),
		Commander:  gobot.NewCommander(),
		interval:   10 * time.Millisecond,
	}

	d.AddEvent("data")
	d.AddEvent("error")
	d.AddCommand("Read", func(params map[string]interface{}) interface{} {
		val, err := d.Read()
		return map[string]interface{}{"val": val, "err": err}
	})

	return d
}

func (a *AnalogSensorDriver) conn() AnalogReader {
	return a.Connection().(AnalogReader)
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
				gobot.Publish(a.Event("error"), err)
			} else if newValue != value && newValue != -1 {
				value = newValue
				gobot.Publish(a.Event("data"), value)
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
func (a *AnalogSensorDriver) Connection() gobot.Connection { return a.connection }

// Read returns the current reading from the Analog Sensor
func (a *AnalogSensorDriver) Read() (val int, err error) {
	return a.conn().AnalogRead(a.Pin())
}
