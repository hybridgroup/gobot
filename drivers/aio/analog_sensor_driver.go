package aio

import (
	"time"

	"gobot.io/x/gobot/v2"
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
	rawValue int
	value    float64
	scale    func(input int) (value float64)
}

// NewAnalogSensorDriver returns a new AnalogSensorDriver with a polling interval of
// 10 Milliseconds given an AnalogReader and pin.
// The driver supports customizable scaling from read int value to returned float64.
// The default scaling is 1:1. An adjustable linear scaler is provided by the driver.
//
// Optionally accepts:
//
//	time.Duration: Interval at which the AnalogSensor is polled for new information
//
// Adds the following API Commands:
//
//	"Read"    - See AnalogDriverSensor.Read
//	"ReadRaw" - See AnalogDriverSensor.ReadRaw
func NewAnalogSensorDriver(a AnalogReader, pin string, v ...time.Duration) *AnalogSensorDriver {
	d := &AnalogSensorDriver{
		name:       gobot.DefaultName("AnalogSensor"),
		connection: a,
		pin:        pin,
		Eventer:    gobot.NewEventer(),
		Commander:  gobot.NewCommander(),
		interval:   10 * time.Millisecond,
		halt:       make(chan bool),
		scale:      func(input int) (value float64) { return float64(input) },
	}

	if len(v) > 0 {
		d.interval = v[0]
	}

	d.AddEvent(Data)
	d.AddEvent(Value)
	d.AddEvent(Error)

	d.AddCommand("Read", func(params map[string]interface{}) interface{} {
		val, err := d.Read()
		return map[string]interface{}{"val": val, "err": err}
	})

	d.AddCommand("ReadRaw", func(params map[string]interface{}) interface{} {
		val, err := d.ReadRaw()
		return map[string]interface{}{"val": val, "err": err}
	})

	return d
}

// Start starts the AnalogSensorDriver and reads the sensor at the given interval.
// Emits the Events:
//
//	Data int - Event is emitted on change and represents the current raw reading from the sensor.
//	Value float64 - Event is emitted on change and represents the current reading from the sensor.
//	Error error - Event is emitted on error reading from the sensor.
func (a *AnalogSensorDriver) Start() (err error) {
	if a.interval == 0 {
		// cyclic reading deactivated
		return
	}
	oldRawValue := 0
	oldValue := 0.0
	go func() {
		timer := time.NewTimer(a.interval)
		timer.Stop()
		for {
			_, err := a.Read()
			if err != nil {
				a.Publish(a.Event(Error), err)
			} else {
				if a.rawValue != oldRawValue && a.rawValue != -1 {
					a.Publish(a.Event(Data), a.rawValue)
					oldRawValue = a.rawValue
				}
				if a.value != oldValue && a.value != -1 {
					a.Publish(a.Event(Value), a.value)
					oldValue = a.value
				}
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
	if a.interval == 0 {
		// cyclic reading deactivated
		return
	}
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

// Read returns the current reading from the sensor
func (a *AnalogSensorDriver) Read() (val float64, err error) {
	if a.rawValue, err = a.ReadRaw(); err != nil {
		return
	}
	a.value = a.scale(a.rawValue)
	return a.value, nil
}

// ReadRaw returns the current reading from the sensor without scaling
func (a *AnalogSensorDriver) ReadRaw() (val int, err error) {
	return a.connection.AnalogRead(a.Pin())
}

// SetScaler substitute the default 1:1 return value function by a new scaling function
func (a *AnalogSensorDriver) SetScaler(scaler func(int) float64) {
	a.scale = scaler
}

// Value returns the last read value from the sensor
func (a *AnalogSensorDriver) Value() float64 {
	return a.value
}

// RawValue returns the last read raw value from the sensor
func (a *AnalogSensorDriver) RawValue() int {
	return a.rawValue
}

// AnalogSensorLinearScaler creates a linear scaler function from the given values.
func AnalogSensorLinearScaler(fromMin, fromMax int, toMin, toMax float64) func(input int) (value float64) {
	m := (toMax - toMin) / float64(fromMax-fromMin)
	n := toMin - m*float64(fromMin)
	return func(input int) (value float64) {
		if input <= fromMin {
			return toMin
		}
		if input >= fromMax {
			return toMax
		}
		return float64(input)*m + n
	}
}
