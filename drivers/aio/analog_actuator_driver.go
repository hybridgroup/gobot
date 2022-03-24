package aio

import (
	"strconv"

	"gobot.io/x/gobot"
)

// AnalogActuatorDriver represents an analog actuator
type AnalogActuatorDriver struct {
	name       string
	pin        string
	connection AnalogWriter
	gobot.Eventer
	gobot.Commander
	scale        func(input float64) (value int)
	lastValue    float64
	lastRawValue int
}

// NewAnalogActuatorDriver returns a new AnalogActuatorDriver given by an AnalogWriter and pin.
// The driver supports customizable scaling from given float64 value to written int.
// The default scaling is 1:1. An adjustable linear scaler is provided by the driver.
//
// Adds the following API Commands:
// 	"Write"    - See AnalogActuator.Write
// 	"RawWrite" - See AnalogActuator.RawWrite
func NewAnalogActuatorDriver(a AnalogWriter, pin string) *AnalogActuatorDriver {
	d := &AnalogActuatorDriver{
		name:       gobot.DefaultName("AnalogActuator"),
		connection: a,
		pin:        pin,
		Commander:  gobot.NewCommander(),
		scale:      func(input float64) (value int) { return int(input) },
	}

	d.AddCommand("Write", func(params map[string]interface{}) interface{} {
		val, err := strconv.ParseFloat(params["val"].(string), 64)
		if err != nil {
			return err
		}
		return d.Write(val)
	})

	d.AddCommand("RawWrite", func(params map[string]interface{}) interface{} {
		val, _ := strconv.Atoi(params["val"].(string))
		return d.RawWrite(val)
	})

	return d
}

// Start starts driver
func (a *AnalogActuatorDriver) Start() (err error) { return }

// Halt is for halt
func (a *AnalogActuatorDriver) Halt() (err error) { return }

// Name returns the drivers name
func (a *AnalogActuatorDriver) Name() string { return a.name }

// SetName sets the drivers name
func (a *AnalogActuatorDriver) SetName(n string) { a.name = n }

// Pin returns the drivers pin
func (a *AnalogActuatorDriver) Pin() string { return a.pin }

// Connection returns the drivers Connection
func (a *AnalogActuatorDriver) Connection() gobot.Connection { return a.connection.(gobot.Connection) }

// RawWrite write the given raw value to the actuator
func (a *AnalogActuatorDriver) RawWrite(val int) (err error) {
	a.lastRawValue = val
	return a.connection.AnalogWrite(a.Pin(), val)
}

// SetScaler substitute the default 1:1 return value function by a new scaling function
func (a *AnalogActuatorDriver) SetScaler(scaler func(float64) int) {
	a.scale = scaler
}

// Write writes the given value to the actuator
func (a *AnalogActuatorDriver) Write(val float64) (err error) {
	a.lastValue = val
	rawValue := a.scale(val)
	return a.RawWrite(rawValue)
}

// RawValue returns the last written raw value
func (a *AnalogActuatorDriver) RawValue() (val int) {
	return a.lastRawValue
}

// Value returns the last written value
func (a *AnalogActuatorDriver) Value() (val float64) {
	return a.lastValue
}

func AnalogActuatorLinearScaler(fromMin, fromMax float64, toMin, toMax int) func(input float64) (value int) {
	m := float64(toMax-toMin) / (fromMax - fromMin)
	n := float64(toMin) - m*fromMin
	return func(input float64) (value int) {
		if input <= fromMin {
			return toMin
		}
		if input >= fromMax {
			return toMax
		}
		return int(input*m + n)
	}
}
