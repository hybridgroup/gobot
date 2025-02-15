package aio

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
)

// sensorOptionApplier needs to be implemented by each configurable option type
type sensorOptionApplier interface {
	apply(cfg *sensorConfiguration)
}

// sensorConfiguration contains all changeable attributes of the driver.
type sensorConfiguration struct {
	readInterval time.Duration
	scale        func(input int) (value float64)
}

// sensorReadIntervalOption is the type for applying another read interval to the configuration
type sensorReadIntervalOption time.Duration

// sensorScaleOption is the type for applying another scaler to the configuration
type sensorScaleOption struct {
	scaler func(input int) (value float64)
}

// AnalogSensorDriver represents an analog sensor
type AnalogSensorDriver struct {
	*driver
	sensorCfg *sensorConfiguration
	pin       string
	halt      chan struct{}
	gobot.Eventer
	lastRawValue int
	lastValue    float64
	analogRead   func() (int, float64, error)
}

// NewAnalogSensorDriver returns a new driver for analog sensors, given an AnalogReader and pin.
// The driver supports cyclic reading and customizable scaling from read int value to returned float64.
// The default scaling is 1:1. An adjustable linear scaler is provided by the driver.
//
// Supported options:
//
//	"WithName"
//	"WithSensorCyclicRead"
//	"WithSensorScaler"
//
// Adds the following API Commands:
//
//	"Read"    - See AnalogDriverSensor.Read
//	"ReadRaw" - See AnalogDriverSensor.ReadRaw
func NewAnalogSensorDriver(a AnalogReader, pin string, opts ...interface{}) *AnalogSensorDriver {
	d := &AnalogSensorDriver{
		driver:    newDriver(a, "AnalogSensor"),
		sensorCfg: &sensorConfiguration{scale: func(input int) float64 { return float64(input) }},
		pin:       pin,
		Eventer:   gobot.NewEventer(), // needed early due to grove vibration sensor driver
	}
	d.afterStart = d.initialize
	d.beforeHalt = d.shutdown
	d.analogRead = d.analogSensorRead

	for _, opt := range opts {
		switch o := opt.(type) {
		case optionApplier:
			o.apply(d.driverCfg)
		case sensorOptionApplier:
			o.apply(d.sensorCfg)
		case time.Duration:
			// TODO this is only for backward compatibility and will be removed after version 2.x
			d.sensorCfg.readInterval = o
		default:
			panic(fmt.Sprintf("'%s' can not be applied on '%s'", opt, d.driverCfg.name))
		}
	}

	d.AddCommand("Read", func(_ map[string]interface{}) interface{} {
		val, err := d.Read()
		return map[string]interface{}{"val": val, "err": err}
	})

	d.AddCommand("ReadRaw", func(_ map[string]interface{}) interface{} {
		val, err := d.ReadRaw()
		return map[string]interface{}{"val": val, "err": err}
	})

	return d
}

// WithSensorCyclicRead add a asynchronous cyclic reading functionality to the sensor with the given read interval.
func WithSensorCyclicRead(interval time.Duration) sensorOptionApplier {
	return sensorReadIntervalOption(interval)
}

// WithSensorScaler substitute the default 1:1 return value function by a new scaling function
func WithSensorScaler(scaler func(input int) (value float64)) sensorOptionApplier {
	return sensorScaleOption{scaler: scaler}
}

// SetScaler substitute the default 1:1 return value function by a new scaling function
// If the scaler is not changed after initialization, prefer to use [aio.WithSensorScaler] instead.
func (a *AnalogSensorDriver) SetScaler(scaler func(int) float64) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	WithSensorScaler(scaler).apply(a.sensorCfg)
}

// Pin returns the AnalogSensorDrivers pin
func (a *AnalogSensorDriver) Pin() string { return a.pin }

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

// Value returns the last read value from the sensor
func (a *AnalogSensorDriver) Value() float64 {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.lastValue
}

// RawValue returns the last read raw value from the sensor
func (a *AnalogSensorDriver) RawValue() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.lastRawValue
}

// initialize the AnalogSensorDriver and if the cyclic reading is active, reads the sensor at the given interval.
// Emits the Events:
//
//	Data int - Event is emitted on change and represents the current raw reading from the sensor.
//	Value float64 - Event is emitted on change and represents the current reading from the sensor.
//	Error error - Event is emitted on error reading from the sensor.
func (a *AnalogSensorDriver) initialize() error {
	if a.sensorCfg.readInterval == 0 {
		// cyclic reading deactivated
		return nil
	}

	a.AddEvent(Data)
	a.AddEvent(Value)
	a.AddEvent(Error)

	// A small buffer is needed to prevent mutex-channel-deadlock between Halt() and analogRead().
	// This can happen, if the shutdown is in progress (mutex passed) and the go routine is calling
	// the analogRead() in between, before the halt can be evaluated by the select statement.
	// In this case the mutex of analogRead() blocks the reading of the halt channel and, without a small buffer,
	// the writing to halt is blocked because there is no immediate read from channel.
	// Please note, that this is special behavior caused by the first read is done immediately before the select
	// statement.
	a.halt = make(chan struct{}, 1)

	oldRawValue := 0
	oldValue := 0.0
	go func() {
		timer := time.NewTimer(a.sensorCfg.readInterval)
		timer.Stop()

		for {
			// please note, that this ensures the first read is done immediately, but has drawbacks, see notes above
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

			timer.Reset(a.sensorCfg.readInterval) // ensure that after each read is a wait, independent of duration of read
			select {
			case <-timer.C:
			case <-a.halt:
				timer.Stop()
				return
			}
		}
	}()
	return nil
}

// shutdown stops polling the analog sensor for new information
func (a *AnalogSensorDriver) shutdown() error {
	if a.sensorCfg.readInterval == 0 || a.halt == nil {
		// cyclic reading deactivated
		return nil
	}
	close(a.halt) // broadcast halt, also to the test
	return nil
}

// analogSensorRead performs an reading from the sensor, sets the internal attributes and returns
// the raw and scaled value
func (a *AnalogSensorDriver) analogSensorRead() (int, float64, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	reader, ok := a.connection.(AnalogReader)
	if !ok {
		return 0, 0, fmt.Errorf("AnalogRead is not supported by the platform '%s'", a.Connection().Name())
	}

	rawValue, err := reader.AnalogRead(a.Pin())
	if err != nil {
		return 0, 0, err
	}

	a.lastRawValue = rawValue
	a.lastValue = a.sensorCfg.scale(a.lastRawValue)
	return a.lastRawValue, a.lastValue, nil
}

func (o sensorReadIntervalOption) String() string {
	return "read interval option for analog sensors"
}

func (o sensorScaleOption) String() string {
	return "scaler option for analog sensors"
}

func (o sensorReadIntervalOption) apply(cfg *sensorConfiguration) {
	cfg.readInterval = time.Duration(o)
}

func (o sensorScaleOption) apply(cfg *sensorConfiguration) {
	cfg.scale = o.scaler
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
