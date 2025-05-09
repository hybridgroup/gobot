package aio

import (
	"fmt"
	"strconv"
)

// actuatorOptionApplier needs to be implemented by each configurable option type
type actuatorOptionApplier interface {
	apply(cfg *actuatorConfiguration)
}

// actuatorConfiguration contains all changeable attributes of the driver.
type actuatorConfiguration struct {
	scale func(input float64) (value int)
}

// actuatorScaleOption is the type for applying another scaler to the configuration
type actuatorScaleOption struct {
	scaler func(input float64) (value int)
}

// AnalogActuatorDriver represents an analog actuator
type AnalogActuatorDriver struct {
	*driver
	pin          string
	actuatorCfg  *actuatorConfiguration
	lastValue    float64
	lastRawValue int
}

// NewAnalogActuatorDriver returns a new driver for analog actuator, given by an AnalogWriter and pin.
// The driver supports customizable scaling from given float64 value to written int.
// The default scaling is 1:1. An adjustable linear scaler is provided by the driver.
//
// Supported options:
//
//	"WithName"
//	"WithActuatorScaler"
//
// Adds the following API Commands:
//
//	"Write"    - See AnalogActuator.Write
//	"WriteRaw" - See AnalogActuator.WriteRaw
func NewAnalogActuatorDriver(a AnalogWriter, pin string, opts ...interface{}) *AnalogActuatorDriver {
	d := &AnalogActuatorDriver{
		driver:      newDriver(a, "AnalogActuator"),
		pin:         pin,
		actuatorCfg: &actuatorConfiguration{scale: func(input float64) int { return int(input) }},
	}

	for _, opt := range opts {
		switch o := opt.(type) {
		case optionApplier:
			o.apply(d.driverCfg)
		case actuatorOptionApplier:
			o.apply(d.actuatorCfg)
		default:
			panic(fmt.Sprintf("'%s' can not be applied on '%s'", opt, d.driverCfg.name))
		}
	}

	//nolint:forcetypeassert // ok here
	d.AddCommand("Write", func(params map[string]interface{}) interface{} {
		val, err := strconv.ParseFloat(params["val"].(string), 64)
		if err != nil {
			return err
		}
		return d.Write(val)
	})
	//nolint:forcetypeassert // ok here
	d.AddCommand("WriteRaw", func(params map[string]interface{}) interface{} {
		val, _ := strconv.Atoi(params["val"].(string))
		return d.WriteRaw(val)
	})

	return d
}

// WithActuatorScaler substitute the default 1:1 return value function by a new scaling function
func WithActuatorScaler(scaler func(input float64) (value int)) actuatorOptionApplier {
	return actuatorScaleOption{scaler: scaler}
}

// SetScaler substitute the default 1:1 return value function by a new scaling function
// If the scaler is not changed after initialization, prefer to use [aio.WithActuatorScaler] instead.
func (a *AnalogActuatorDriver) SetScaler(scaler func(float64) int) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	WithActuatorScaler(scaler).apply(a.actuatorCfg)
}

// Pin returns the drivers pin
func (a *AnalogActuatorDriver) Pin() string { return a.pin }

// Write writes the given value to the actuator
func (a *AnalogActuatorDriver) Write(val float64) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	rawValue := a.actuatorCfg.scale(val)
	if err := a.WriteRaw(rawValue); err != nil {
		return err
	}
	a.lastValue = val
	return nil
}

// RawWrite write the given raw value to the actuator
// Deprecated: Please use [aio.WriteRaw] instead.
func (a *AnalogActuatorDriver) RawWrite(val int) error {
	return a.WriteRaw(val)
}

// WriteRaw write the given raw value to the actuator
func (a *AnalogActuatorDriver) WriteRaw(val int) error {
	writer, ok := a.connection.(AnalogWriter)
	if !ok {
		return fmt.Errorf("AnalogWrite is not supported by the platform '%s'", a.Connection().Name())
	}
	if err := writer.AnalogWrite(a.Pin(), val); err != nil {
		return err
	}
	a.lastRawValue = val
	return nil
}

// Value returns the last written value
func (a *AnalogActuatorDriver) Value() float64 {
	return a.lastValue
}

// RawValue returns the last written raw value
func (a *AnalogActuatorDriver) RawValue() int {
	return a.lastRawValue
}

func (o actuatorScaleOption) String() string {
	return "scaler option for analog actuators"
}

func (o actuatorScaleOption) apply(cfg *actuatorConfiguration) {
	cfg.scale = o.scaler
}

// AnalogActuatorLinearScaler creates a linear scaler function from the given values.
func AnalogActuatorLinearScaler(fromMin, fromMax float64, toMin, toMax int) func(input float64) (value int) {
	m := float64(toMax-toMin) / (fromMax - fromMin)
	n := float64(toMin) - m*fromMin
	return func(input float64) int {
		if input <= fromMin {
			return toMin
		}
		if input >= fromMax {
			return toMax
		}
		return int(input*m + n)
	}
}
