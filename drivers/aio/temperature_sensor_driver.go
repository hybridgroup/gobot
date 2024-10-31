package aio

import (
	"fmt"
	"math"
	"time"

	"gobot.io/x/gobot/v2"
)

const kelvinOffset = 273.15

// TemperatureSensorNtcConf contains all attributes to calculate key parameters of a NTC thermistor.
type TemperatureSensorNtcConf struct {
	TC0 int     // °C
	R0  float64 // same unit as resistance of NTC (Ohm is recommended)
	B   float64 // 2000..5000K
	TC1 int     // used if B is not given, °C
	R1  float64 // used if B is not given, same unit as R0 needed
	t0  float64
	r   float64
}

// TemperatureSensorDriver represents an Analog Sensor
type TemperatureSensorDriver struct {
	*AnalogSensorDriver
}

// NewTemperatureSensorDriver is a driver for analog temperature sensors, given an AnalogReader and pin.
// Linear scaling and NTC scaling is supported.
//
// Supported options: see [aio.NewAnalogSensorDriver]
// Adds the following API Commands: see [aio.NewAnalogSensorDriver]
func NewTemperatureSensorDriver(a AnalogReader, pin string, opts ...interface{}) *TemperatureSensorDriver {
	d := &TemperatureSensorDriver{AnalogSensorDriver: NewAnalogSensorDriver(a, pin)}
	d.driverCfg.name = gobot.DefaultName("TemperatureSensor")

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

	return d
}

// SetNtcScaler sets a function for typical NTC scaling the read value.
// The read value is related to the voltage over the thermistor in an series connection to a resistor.
// If the thermistor is connected to ground, the reverse flag must be set to true.
// This means the voltage decreases when temperature gets higher.
// Currently no negative values for voltage are supported.
// If the scaler is not changed after initialization, prefer to use [aio.WithSensorScaler] instead.
func (t *TemperatureSensorDriver) SetNtcScaler(vRef uint, rOhm uint, reverse bool, ntc TemperatureSensorNtcConf) {
	t.SetScaler(TemperatureSensorNtcScaler(vRef, rOhm, reverse, ntc))
}

// SetLinearScaler sets a function for linear scaling the read value.
// This can be applied for some silicon based PTC sensors or e.g. PT100,
// and in a small temperature range also for NTC.
// If the scaler is not changed after initialization, prefer to use [aio.WithSensorScaler] instead.
func (t *TemperatureSensorDriver) SetLinearScaler(fromMin, fromMax int, toMin, toMax float64) {
	t.SetScaler(AnalogSensorLinearScaler(fromMin, fromMax, toMin, toMax))
}

// TemperatureSensorNtcScaler creates a function for typical NTC scaling the read value.
// The read value is related to the voltage over the thermistor in an series connection to a resistor.
// If the thermistor is connected to ground, the reverse flag must be set to true.
// This means the voltage decreases when temperature gets higher.
// Currently no negative values for voltage are supported.
func TemperatureSensorNtcScaler(
	vRef uint,
	rOhm uint,
	reverse bool,
	ntc TemperatureSensorNtcConf,
) func(input int) (value float64) {
	ntc.initialize()
	return (func(input int) float64 {
		if input < 0 {
			input = 0
		}
		rTherm := temperaturSensorGetResistance(uint(input), vRef, rOhm, reverse) //nolint:gosec // checked before
		temp := ntc.getTemp(rTherm)
		return temp
	})
}

// getResistance calculates the value of the series thermistor by given value
// and reference value (e.g. the voltage value over the complete series circuit)
// The unit of the returned thermistor value equals the given series resistor unit.
func temperaturSensorGetResistance(value uint, vRef uint, rSeries uint, reverse bool) float64 {
	if value > vRef {
		value = vRef
	}
	valDiff := vRef - value
	if reverse {
		//        rSeries  thermistor
		// vRef o--|==|--o--|=/=|----| GND
		//               |-> value <-|
		if value == 0 {
			// prevent jump to -273.15
			value = 1
		}
		return float64(rSeries*value) / float64(valDiff)
	}

	//      thermistor  rSeries
	// vRef o--|=/=|--o--|==|-----| GND
	//                |-> value <-|
	if valDiff == 0 {
		// prevent jump to -273.15
		valDiff = 1
	}
	return float64(rSeries*valDiff) / float64(value)
}

// getTemp calculates the temperature from the given resistance of the NTC resistor
func (n *TemperatureSensorNtcConf) getTemp(rntc float64) float64 {
	// 1/T = 1/T0 + 1/B * ln(R/R0)
	//
	// B/T = B/T0 + ln(R/R0) = k, B/T0 = r
	// T = B/k, Tc = T - 273

	k := n.r + math.Log(rntc/n.R0)
	return n.B/k - kelvinOffset
}

// initialize is used to calculate some constants for the NTC algorithm.
// If B is unknown (given as 0), the function needs a second pair to calculate
// B from the both pairs (R1, TC1), (R0, TC0)
func (n *TemperatureSensorNtcConf) initialize() {
	n.t0 = float64(n.TC0) + kelvinOffset
	if n.B <= 0 {
		// B=[ln(R0)-ln(R1)]/(1/T0-1/T1)
		T1 := float64(n.TC1) + kelvinOffset
		n.B = (1/n.t0 - 1/T1)
		n.B = (math.Log(n.R0) - math.Log(n.R1)) / n.B // 2000K...5000K
	}
	n.r = n.B / n.t0
}
