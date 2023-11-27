package aio

import (
	"fmt"

	"gobot.io/x/gobot/v2"
)

// thermalZoneOptionApplier needs to be implemented by each configurable option type
type thermalZoneOptionApplier interface {
	apply(cfg *thermalZoneConfiguration)
}

// thermalZoneConfiguration contains all changeable attributes of the driver.
type thermalZoneConfiguration struct {
	scaleUnit func(float64) float64
}

// thermalZoneUnitscalerOption is the type for applying another unit scaler to the configuration
type thermalZoneUnitscalerOption struct {
	unitscaler func(float64) float64
}

// ThermalZoneDriver represents an driver for reading the system thermal zone temperature
type ThermalZoneDriver struct {
	*AnalogSensorDriver
	thermalZoneCfg *thermalZoneConfiguration
}

// NewThermalZoneDriver is a driver for reading the system thermal zone temperature, given an AnalogReader and zone id.
//
// Supported options: see also [aio.NewAnalogSensorDriver]
//
//	"WithFahrenheit()"
//
// Adds the following API Commands: see [aio.NewAnalogSensorDriver]
func NewThermalZoneDriver(a AnalogReader, zoneID string, opts ...interface{}) *ThermalZoneDriver {
	degreeScaler := func(input int) float64 { return float64(input) / 1000 }
	d := ThermalZoneDriver{
		AnalogSensorDriver: NewAnalogSensorDriver(a, zoneID, WithSensorScaler(degreeScaler)),
		thermalZoneCfg: &thermalZoneConfiguration{
			scaleUnit: func(input float64) float64 { return input }, // 1:1 in °C
		},
	}
	d.driverCfg.name = gobot.DefaultName("ThermalZone")
	d.analogRead = d.thermalZoneRead

	for _, opt := range opts {
		switch o := opt.(type) {
		case optionApplier:
			o.apply(d.driverCfg)
		case sensorOptionApplier:
			o.apply(d.sensorCfg)
		case thermalZoneOptionApplier:
			o.apply(d.thermalZoneCfg)
		default:
			panic(fmt.Sprintf("'%s' can not be applied on '%s'", opt, d.driverCfg.name))
		}
	}

	return &d
}

// WithFahrenheit substitute the default 1:1 °C scaler by a scaler for °F
func WithFahrenheit() thermalZoneOptionApplier {
	// (1°C × 9/5) + 32 = 33,8°F
	unitscaler := func(input float64) float64 { return input*9.0/5.0 + 32.0 }
	return thermalZoneUnitscalerOption{unitscaler: unitscaler}
}

// thermalZoneRead overrides and extends the analogSensorRead() function to add the unit scaler
func (d *ThermalZoneDriver) thermalZoneRead() (int, float64, error) {
	if _, _, err := d.analogSensorRead(); err != nil {
		return 0, 0, err
	}
	// apply unit scaler on value
	d.lastValue = d.thermalZoneCfg.scaleUnit(d.lastValue)
	return d.lastRawValue, d.lastValue, nil
}

func (o thermalZoneUnitscalerOption) apply(cfg *thermalZoneConfiguration) {
	cfg.scaleUnit = o.unitscaler
}
