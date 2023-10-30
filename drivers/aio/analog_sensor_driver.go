package aio

import (
	"sync"
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
	mutex    *sync.Mutex // to prevent data race between cyclic and single shot write/read to values and scaler
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
		mutex:      &sync.Mutex{},
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
			rawValue, value, err := a.analogRead()
			if err != nil {
				a.Publish(a.Event(Error), err)
			} else {
				if rawValue != oldRawValue && rawValue != -1 {
					a.Publish(a.Event(Data), rawValue)
					oldRawValue = rawValue
				}
				if value != oldValue && value != -1 {
					a.Publish(a.Event(Value), value)
					oldValue = value
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

// Read returns the current reading from the sensor, scaled by the current scaler
func (a *AnalogSensorDriver) Read() (float64, error) {
	_, value, err := a.analogRead()
	return value, err
}

// ReadRaw returns the current reading from the sensor without scaling
func (a *AnalogSensorDriver) ReadRaw() (int, error) {
	rawValue, _, err := a.analogRead()
	return rawValue, err
}

// SetScaler substitute the default 1:1 return value function by a new scaling function
func (a *AnalogSensorDriver) SetScaler(scaler func(int) float64) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.scale = scaler
}

// Value returns the last read value from the sensor
func (a *AnalogSensorDriver) Value() float64 {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.value
}

// RawValue returns the last read raw value from the sensor
func (a *AnalogSensorDriver) RawValue() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.rawValue
}

// analogRead performs an reading from the sensor and sets the internal attributes and returns the raw and scaled value
func (a *AnalogSensorDriver) analogRead() (int, float64, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	rawValue, err := a.connection.AnalogRead(a.Pin())
	if err != nil {
		return 0, 0, err
	}

	a.rawValue = rawValue
	a.value = a.scale(a.rawValue)
	return a.rawValue, a.value, nil
}

// AnalogSensorLinearScaler creates a linear scaler function from the given values.
func AnalogSensorLinearScaler(fromMin, fromMax int, toMin, toMax float64) func(int) float64 {
	m := (toMax - toMin) / float64(fromMax-fromMin)
	n := toMin - m*float64(fromMin)
	return func(input int) float64 {
		if input <= fromMin {
			return toMin
		}
		if input >= fromMax {
			return toMax
		}
		return float64(input)*m + n
	}
}
