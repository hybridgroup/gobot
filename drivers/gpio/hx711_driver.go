package gpio

import (
	"fmt"
	"time"

	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/hx711"

	gobot "gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/system"
)

const (
	defaultGainAndChannelCfg = hx711.A128
	defaultReadTimeout       = 2 * time.Second
	defaultReadTriesMax      = 100
	defaultTickSleep         = 0 * time.Microsecond // 0..50 us, depending on CPU performance and scheduling
)

// hx711OptionApplier needs to be implemented by each configurable option type
type hx711OptionApplier interface {
	apply(cfg *hx711Configuration)
}

// hx711Configuration contains all changeable attributes of the driver.
type hx711Configuration struct {
	gainAndChannelCfg hx711.GainAndChannelCfg
	readTimeout       time.Duration
	readTriesMax      uint8 // how often a check for ready state is done until the read timeout is reached
	tickSleep         time.Duration
}

// hx711GainAndChannelCfgOption is the type for applying another gain and channel configuration.
type hx711GainAndChannelCfgOption hx711.GainAndChannelCfg

// hx711ReadTimeoutOption is the type for applying another duration for the read timeout.
type hx711ReadTimeoutOption time.Duration

// hx711ReadTriesMaxOption is the type for applying another value for maximum read tries.
type hx711ReadTriesMaxOption uint8

// hx711TickSleepOption is the type for applying another duration for the tick H-level-time.
type hx711TickSleepOption time.Duration

type sensorer interface {
	Configure(cfg *hx711.DeviceConfig) error
	Zero(secondReading bool) error
	Calibrate(setValue int32, secondReading bool) error
	Update(m drivers.Measurement) error
	Values() (int64, int64)
	OffsetAndCalibrationFactor(secondReading bool) (int32, float32)
	SetOffsetAndCalibrationFactor(offset int32, calibrationFactor float32, secondReading bool) error
}

// HX711 is the gobot driver for the HX711 chip.
type HX711Driver struct {
	*driver

	hx711Cfg  *hx711Configuration
	newSensor func(clockPin, dataPin gobot.DigitalPinner) sensorer
	sensor    sensorer
}

// NewHX711 creates a new driver for HX711 24 bit, 2 channel, configurable ADC with serial output to measure small
// differential voltages. The device is handy for load cells but can be used to read all kind of Wheatstone bridges.
// Therefore the usage of the phrases "mass", "weight" or "load" are prevented in this driver - "value" is used instead.
// The implementation just wraps the according TinyGo-driver.
//
// Datasheet: https://cdn.sparkfun.com/datasheets/Sensors/ForceFlex/hx711_english.pdf
//
// Supported options:
//
//	"WithName"
//	"WithHX711GainAndChannelCfg"
//	"WithHX711ReadTimeout"
//	"WithHX711ReadTriesMax"
//	"WithHX711TickSleep"
func NewHX711Driver(a gobot.Connection, clockPinID, dataPinID string, opts ...interface{}) *HX711Driver {
	d := HX711Driver{
		driver: newDriver(a, "HX711"),
		hx711Cfg: &hx711Configuration{
			gainAndChannelCfg: defaultGainAndChannelCfg,
			readTimeout:       defaultReadTimeout,
			readTriesMax:      defaultReadTriesMax,
			tickSleep:         defaultTickSleep,
		},
	}

	for _, opt := range opts {
		switch o := opt.(type) {
		case optionApplier:
			o.apply(d.driverCfg)
		case hx711OptionApplier:
			o.apply(d.hx711Cfg)
		default:
			panic(fmt.Sprintf("'%s' can not be applied on '%s'", opt, d.driverCfg.name))
		}
	}

	// preparation in this way opens the possibility to change the sensor for unit tests before Connect()
	d.newSensor = func(clockPin, dataPin gobot.DigitalPinner) sensorer {
		clockPinSetState := func(v bool) error {
			if v {
				return clockPin.Write(1)
			}

			return clockPin.Write(0)
		}

		dataPinState := func() (bool, error) {
			val, err := dataPin.Read()

			return val > 0, err
		}

		return hx711.New(clockPinSetState, dataPinState, d.hx711Cfg.gainAndChannelCfg)
	}

	// note: Unexport() of all pins will be done on adaptor.Finalize()
	d.afterStart = func() error {
		clockPin, err := a.(gobot.DigitalPinnerProvider).DigitalPin(clockPinID) //nolint:forcetypeassert // ok here
		if err != nil {
			return fmt.Errorf("error on get clock pin: %v", err)
		}
		if err := clockPin.ApplyOptions(system.WithPinDirectionOutput(0)); err != nil {
			return fmt.Errorf("error on apply output for clock pin: %v", err)
		}

		// pins are inputs by default
		dataPin, err := a.(gobot.DigitalPinnerProvider).DigitalPin(dataPinID) //nolint:forcetypeassert // ok here
		if err != nil {
			return fmt.Errorf("error on get data pin: %v", err)
		}

		if err := dataPin.ApplyOptions(system.WithPinPullUp()); err != nil {
			return fmt.Errorf("error on apply pull up option for data pin: %v", err)
		}

		d.sensor = d.newSensor(clockPin, dataPin)

		cfg := hx711.DefaultConfig
		cfg.ReadTimeout = d.hx711Cfg.readTimeout
		cfg.ReadTriesMax = d.hx711Cfg.readTriesMax
		cfg.TickSleep = d.hx711Cfg.tickSleep

		// this leads to an active reset of the sensor hardware, so pins must be ready at this point
		return d.sensor.Configure(&cfg)
	}

	return &d
}

// WithHX711GainAndChannelCfg use the given value for the gain and channel configuration instead the default hx711.A128.
func WithHX711GainAndChannelCfg(gc hx711.GainAndChannelCfg) hx711OptionApplier {
	return hx711GainAndChannelCfgOption(gc)
}

// WithHX711ReadTimeout use the given duration for read timeout instead the default of 2 seconds.
func WithHX711ReadTimeout(timeout time.Duration) hx711OptionApplier {
	return hx711ReadTimeoutOption(timeout)
}

// WithHX711ReadTriesMax use the given value for maximum read tries instead the default of 100. A value of zero will be
// automatically adjusted to the minimum of 1 try by the underlying sensor driver.
func WithHX711ReadTriesMax(tries uint8) hx711OptionApplier {
	return hx711ReadTriesMaxOption(tries)
}

// WithHX711TickSleep use the given duration for the H-level of the clock tick. The default is set to 0, which performs
// best for most CPUs. The value must be smaller than 60us, otherwise a hardware reset on each measure will occur. This
// will be automatically adjusted by the underlying sensor driver to this maximum, but the CPU will not rely on that
// value in any case, depending on its performance and scheduling. E.g. on a tinkerboard a value of 1us leads to a
// real duration of H-level between 10 and 70us, which causes wrong measures with high occurrence.
func WithHX711TickSleep(tickHLevelDuration time.Duration) hx711OptionApplier {
	return hx711TickSleepOption(tickHLevelDuration)
}

// Zero sets the offset for the reading. If the given flag is true, this is done for the second reading.
func (d *HX711Driver) Zero(secondReading bool) error {
	// locking concurrent calls is done at sensor side
	return d.sensor.Zero(secondReading)
}

// Calibrate calculates, after a measurement of the set value is done, a factor for linear scaling the values of the
// subsequent measurements. The unit of the given set value define the unit of the measurement result later. Before
// using this function, the offset value should be obtained by calling Zero() function with no load.
// If the given flag is true, this is done for the second reading.
func (d *HX711Driver) Calibrate(setValue int32, secondReading bool) error {
	// locking concurrent calls is done at sensor side
	return d.sensor.Calibrate(setValue, secondReading)
}

// OffsetAndCalibrationFactor returns linear correction values, used for reading.
// If the given flag is true, this values are related to the second reading.
func (d *HX711Driver) OffsetAndCalibrationFactor(secondReading bool) (int32, float32) {
	// locking concurrent calls is done at sensor side
	return d.sensor.OffsetAndCalibrationFactor(secondReading)
}

// SetOffsetAndCalibrationFactor sets linear correction values, used for reading.
// If the given flag is true, this values are related to the second reading.
func (d *HX711Driver) SetOffsetAndCalibrationFactor(offset int32, calibrationFactor float32, secondReading bool) error {
	// locking concurrent calls is done at sensor side
	return d.sensor.SetOffsetAndCalibrationFactor(offset, calibrationFactor, secondReading)
}

// Measure retrieves the values of the sensor and returns the measure. If only channel A is configured, channel B just
// returns always zero. Because the timing is somewhat difficult on CPUs regarding the low time of the clock pin (>60us
// leads to a reset of the hardware), it is suggested to implement a validation for the values and drop or repeat the
// measure in case of a value is unexpected.
func (d *HX711Driver) Measure() (int64, int64, error) {
	// locking concurrent calls is done at sensor side
	if err := d.sensor.Update(0); err != nil {
		return 0, 0, err
	}

	v1, v2 := d.sensor.Values()

	return v1, v2, nil
}

func (o hx711GainAndChannelCfgOption) String() string {
	return "hx711 gain and channel configuration option"
}

func (o hx711GainAndChannelCfgOption) apply(cfg *hx711Configuration) {
	cfg.gainAndChannelCfg = hx711.GainAndChannelCfg(o)
}

func (o hx711ReadTimeoutOption) String() string {
	return "hx711 read timeout option"
}

func (o hx711ReadTimeoutOption) apply(cfg *hx711Configuration) {
	cfg.readTimeout = time.Duration(o)
}

func (o hx711ReadTriesMaxOption) String() string {
	return "hx711 maximum read tries option"
}

func (o hx711ReadTriesMaxOption) apply(cfg *hx711Configuration) {
	cfg.readTriesMax = uint8(o)
}

func (o hx711TickSleepOption) String() string {
	return "hx711 tick sleep option"
}

func (o hx711TickSleepOption) apply(cfg *hx711Configuration) {
	cfg.tickSleep = time.Duration(o)
}
