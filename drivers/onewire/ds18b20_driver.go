package onewire

import (
	"fmt"
	"math"
	"time"
)

const (
	ds18b20DefaultResolution     = 12
	ds18b20DefaultConversionTime = 750

	temperatureCommand = "temperature"
	extPowerCommand    = "ext_power"
	resolutionCommand  = "resolution"
	convTimeCommand    = "conv_time"
)

// ds18b20OptionApplier needs to be implemented by each configurable option type
type ds18b20OptionApplier interface {
	apply(cfg *ds18b20Configuration)
}

// ds18b20Configuration contains all changeable attributes of the driver.
type ds18b20Configuration struct {
	scaleUnit      func(int) float32
	resolution     uint8
	conversionTime uint16
}

// ds18b20UnitscalerOption is the type for applying another unit scaler to the configuration
type ds18b20UnitscalerOption struct {
	unitscaler func(int) float32
}

type ds18b20ResolutionOption uint8

type ds18b20ConversionTimeOption uint16

// DS18B20Driver is a driver for the DS18B20 1-wire temperature sensor.
type DS18B20Driver struct {
	*driver
	ds18b20Cfg *ds18b20Configuration
}

// NewDS18B20Driver creates a new Gobot Driver for DS18B20 one wire temperature sensor.
//
// Params:
//
//	a *Adaptor - the Adaptor to use with this Driver.
//	serial number int - the serial number of the device, without the family code
//
// Optional params:
//
// onewire.WithFahrenheit()
// onewire.WithResolution(byte)
// onewire.WithConversionTime(uint16)
func NewDS18B20Driver(a connector, serialNumber uint64, opts ...interface{}) *DS18B20Driver {
	d := &DS18B20Driver{
		driver: newDriver(a, "DS18B20", 0x28, serialNumber),
		ds18b20Cfg: &ds18b20Configuration{
			scaleUnit:      func(input int) float32 { return float32(input) / 1000 }, // 1000:1 in °C
			resolution:     ds18b20DefaultResolution,
			conversionTime: ds18b20DefaultConversionTime,
		},
	}
	d.afterStart = d.initialize
	d.beforeHalt = d.shutdown
	for _, opt := range opts {
		switch o := opt.(type) {
		case optionApplier:
			o.apply(d.driverCfg)
		case ds18b20OptionApplier:
			o.apply(d.ds18b20Cfg)
		default:
			panic(fmt.Sprintf("'%s' can not be applied on '%s'", opt, d.driverCfg.name))
		}
	}
	return d
}

// WithFahrenheit substitute the default °C scaler by a scaler for °F
func WithFahrenheit() ds18b20OptionApplier {
	// (1°C × 9/5) + 32 = 33,8°F
	unitscaler := func(input int) float32 { return float32(input)/1000*9.0/5.0 + 32.0 }
	return ds18b20UnitscalerOption{unitscaler: unitscaler}
}

// WithResolution substitute the default 12 bit resolution by the given one (9, 10, 11). The device will adjust
// the conversion time automatically. Each smaller resolution will decrease the conversion time by a factor of 2.
// Note: some devices are fixed in 12 bit mode only and do not support this feature (I/O error or just ignore it).
// WithConversionTime() is most likely supported.
func WithResolution(resolution uint8) ds18b20OptionApplier {
	return ds18b20ResolutionOption(resolution)
}

// WithConversionTime substitute the default 750 ms by the given one (93, 187, 375, 750).
// Note: Devices will not adjust the resolution automatically. Some devices accept conversion time values different
// from common specification. E.g. 10...1000, which leads to real conversion time of conversionTime+50ms. This needs
// to be tested for your device and measured for your needs, e.g. by DebugConversionTime(0, 500, 5, true).
func WithConversionTime(conversionTime uint16) ds18b20OptionApplier {
	return ds18b20ConversionTimeOption(conversionTime)
}

// Temperature returns the current temperature, in celsius degrees, if the default unit scaler is used.
func (d *DS18B20Driver) Temperature() (float32, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	val, err := d.connection.ReadInteger(temperatureCommand)
	if err != nil {
		return 0, err
	}

	return d.ds18b20Cfg.scaleUnit(val), nil
}

// Resolution returns the current resolution in bits (9, 10, 11, 12)
func (d *DS18B20Driver) Resolution() (uint8, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	val, err := d.connection.ReadInteger(resolutionCommand)
	if err != nil {
		return 0, err
	}

	if val < 9 || val > 12 {
		return 0, fmt.Errorf("the read value '%d' is out of range (9, 10, 11, 12)", val)
	}

	return uint8(val), nil
}

// IsExternalPowered returns whether the device is external or parasitic powered
func (d *DS18B20Driver) IsExternalPowered() (bool, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	val, err := d.connection.ReadInteger(extPowerCommand)
	if err != nil {
		return false, err
	}

	return val > 0, nil
}

// ConversionTime returns the conversion time in ms
func (d *DS18B20Driver) ConversionTime() (uint16, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	val, err := d.connection.ReadInteger(convTimeCommand)
	if err != nil {
		return 0, err
	}

	if val < 0 || val > math.MaxUint16 {
		return 0, fmt.Errorf("the read value '%d' is out of range (uint16)", val)
	}

	return uint16(val), nil
}

// DebugConversionTime try to set the conversion time and compare with real time to read temperature.
func (d *DS18B20Driver) DebugConversionTime(start, end uint16, stepwide uint16, skipInvalid bool) {
	r, _ := d.Resolution()
	fmt.Printf("\n---- Conversion time check for '%s'@%dbit %d..%d +%d ----\n",
		d.connection.ID(), r, start, end, stepwide)
	fmt.Println("|r1(err)\t|w(err)\t\t|r2(err)\t|T(err)\t\t|real\t\t|diff\t\t|")
	fmt.Println("--------------------------------------------------------------------------------")
	for ct := start; ct < end; ct += stepwide {
		r1, e1 := d.ConversionTime()
		ew := d.connection.WriteInteger(convTimeCommand, int(ct))
		r2, e2 := d.ConversionTime()
		time.Sleep(100 * time.Millisecond) // relax the system
		start := time.Now()
		temp, err := d.Temperature()
		dur := time.Since(start)
		valid := ct == r2
		if valid || !skipInvalid {
			diff := dur - time.Duration(r2)*time.Millisecond
			fmt.Printf("|%d(%t)\t|%d(%t)\t|%d(%t)\t|%v(%t)\t|%s\t|%s\t|\n",
				r1, e1 != nil, ct, ew != nil, r2, e2 != nil, temp, err != nil, dur, diff)
		}
	}
}

func (d *DS18B20Driver) initialize() error {
	if d.ds18b20Cfg.resolution != ds18b20DefaultResolution {
		if err := d.connection.WriteInteger(resolutionCommand, int(d.ds18b20Cfg.resolution)); err != nil {
			return err
		}
	}

	if d.ds18b20Cfg.conversionTime != ds18b20DefaultConversionTime {
		return d.connection.WriteInteger(convTimeCommand, int(d.ds18b20Cfg.conversionTime))
	}

	return nil
}

func (d *DS18B20Driver) shutdown() error {
	if d.ds18b20Cfg.resolution != ds18b20DefaultResolution {
		if err := d.connection.WriteInteger(resolutionCommand, ds18b20DefaultResolution); err != nil {
			return err
		}
	}

	if d.ds18b20Cfg.conversionTime != ds18b20DefaultConversionTime {
		return d.connection.WriteInteger(convTimeCommand, int(ds18b20DefaultConversionTime))
	}

	return nil
}

func (o ds18b20UnitscalerOption) apply(cfg *ds18b20Configuration) {
	cfg.scaleUnit = o.unitscaler
}

func (o ds18b20ResolutionOption) apply(cfg *ds18b20Configuration) {
	cfg.resolution = uint8(o)
}

func (o ds18b20ConversionTimeOption) apply(cfg *ds18b20Configuration) {
	cfg.conversionTime = uint16(o)
}
