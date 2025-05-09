package i2c

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"time"
)

const ads1x15DefaultAddress = 0x48

const (
	ads1x15Debug          = false
	ads1x15WaitMaxCount   = 200
	ads1x15FullScaleValue = 0x7FFF // same as 32768, 1<<15 or 64
)

const (
	// Address pointer registers
	ads1x15PointerConversion    = 0x00
	ads1x15PointerConfig        = 0x01
	ads1x15PointerLowThreshold  = 0x02
	ads1x15PointerHighThreshold = 0x03

	// Values for config register
	ads1x15ConfigCompQueDisable = 0x0003
	ads1x15ConfigCompLatching   = 0x0004
	ads1x15ConfigCompActiveHigh = 0x0008
	ads1x15ConfigCompWindow     = 0x0010

	ads1x15ConfigModeContinuous = 0x0000
	ads1x15ConfigModeSingle     = 0x0100 // single shot mode
	ads1x15ConfigOsSingle       = 0x8000 // write: set to start a single-conversion, read: wait for finished
	ads1x15ConfigMuxOffset      = 12
	ads1x15ConfigPgaOffset      = 9
)

type ads1x15ChanCfg struct {
	gain     int
	dataRate int
}

// ADS1x15Driver is the Gobot driver for the ADS1015/ADS1115 ADC
// datasheet:
// https://www.ti.com/lit/gpn/ads1115
//
// reference implementations:
// * https://github.com/adafruit/Adafruit_Python_ADS1x15
// * https://github.com/Wh1teRabbitHU/ADS1115-Driver
type ADS1x15Driver struct {
	*Driver
	dataRates        map[int]uint16
	channelCfgs      map[int]*ads1x15ChanCfg
	waitOnlyOneCycle bool
}

var ads1x15FullScaleRange = map[int]float64{
	0: 6.144,
	1: 4.096,
	2: 2.048,
	3: 1.024,
	4: 0.512,
	5: 0.256,
	6: 0.256,
	7: 0.256,
}

// NewADS1015Driver creates a new driver for the ADS1015 (12-bit ADC)
func NewADS1015Driver(a Connector, options ...func(Config)) *ADS1x15Driver {
	dataRates := map[int]uint16{
		128:  0x0000,
		250:  0x0020,
		490:  0x0040,
		920:  0x0060,
		1600: 0x0080,
		2400: 0x00A0,
		3300: 0x00C0,
	}
	defaultDataRate := 1600

	return newADS1x15Driver(a, "ADS1015", dataRates, defaultDataRate, options...)
}

// NewADS1115Driver creates a new driver for the ADS1115 (16-bit ADC)
func NewADS1115Driver(a Connector, options ...func(Config)) *ADS1x15Driver {
	dataRates := map[int]uint16{
		8:   0x0000,
		16:  0x0020,
		32:  0x0040,
		64:  0x0060,
		128: 0x0080,
		250: 0x00A0,
		475: 0x00C0,
		860: 0x00E0,
	}
	defaultDataRate := 128

	return newADS1x15Driver(a, "ADS1115", dataRates, defaultDataRate, options...)
}

func newADS1x15Driver(c Connector, name string, drs map[int]uint16, ddr int, options ...func(Config)) *ADS1x15Driver {
	ccs := map[int]*ads1x15ChanCfg{0: {1, ddr}, 1: {1, ddr}, 2: {1, ddr}, 3: {1, ddr}}
	d := &ADS1x15Driver{
		Driver:      NewDriver(c, name, ads1x15DefaultAddress),
		dataRates:   drs,
		channelCfgs: ccs,
	}

	for _, option := range options {
		option(d)
	}

	//nolint:forcetypeassert // ok here
	d.AddCommand("ReadDifferenceWithDefaults", func(params map[string]interface{}) interface{} {
		channel := params["channel"].(int)
		val, err := d.ReadDifferenceWithDefaults(channel)
		return map[string]interface{}{"val": val, "err": err}
	})

	//nolint:forcetypeassert // ok here
	d.AddCommand("ReadDifference", func(params map[string]interface{}) interface{} {
		channel := params["channel"].(int)
		gain := params["gain"].(int)
		dataRate := params["dataRate"].(int)
		val, err := d.ReadDifference(channel, gain, dataRate)
		return map[string]interface{}{"val": val, "err": err}
	})

	//nolint:forcetypeassert // ok here
	d.AddCommand("ReadWithDefaults", func(params map[string]interface{}) interface{} {
		channel := params["channel"].(int)
		val, err := d.ReadWithDefaults(channel)
		return map[string]interface{}{"val": val, "err": err}
	})

	//nolint:forcetypeassert // ok here
	d.AddCommand("Read", func(params map[string]interface{}) interface{} {
		channel := params["channel"].(int)
		gain := params["gain"].(int)
		dataRate := params["dataRate"].(int)
		val, err := d.Read(channel, gain, dataRate)
		return map[string]interface{}{"val": val, "err": err}
	})

	//nolint:forcetypeassert // ok here
	d.AddCommand("AnalogRead", func(params map[string]interface{}) interface{} {
		pin := params["pin"].(string)
		val, err := d.AnalogRead(pin)
		return map[string]interface{}{"val": val, "err": err}
	})

	return d
}

// WithADS1x15BestGainForVoltage option sets the ADS1x15Driver best gain for all channels.
func WithADS1x15BestGainForVoltage(voltage float64) func(Config) {
	return func(c Config) {
		d, ok := c.(*ADS1x15Driver)
		if ok {
			// validate the given value
			bestGain, err := ads1x15BestGainForVoltage(voltage)
			if err != nil {
				panic(err)
			}
			WithADS1x15Gain(bestGain)(d)
		} else if ads1x15Debug {
			log.Printf("Trying to set best gain for voltage for non-ADS1x15Driver %v", c)
		}
	}
}

// WithADS1x15ChannelBestGainForVoltage option sets the ADS1x15Driver best gain for one channel.
func WithADS1x15ChannelBestGainForVoltage(channel int, voltage float64) func(Config) {
	return func(c Config) {
		d, ok := c.(*ADS1x15Driver)
		if ok {
			// validate the given value
			bestGain, err := ads1x15BestGainForVoltage(voltage)
			if err != nil {
				panic(err)
			}
			WithADS1x15ChannelGain(channel, bestGain)(d)
		} else if ads1x15Debug {
			log.Printf("Trying to set channel best gain for voltage for non-ADS1x15Driver %v", c)
		}
	}
}

// WithADS1x15Gain option sets the ADS1x15Driver gain for all channels.
// Valid gain settings are any of the PGA values (0..7).
func WithADS1x15Gain(val int) func(Config) {
	return func(c Config) {
		d, ok := c.(*ADS1x15Driver)
		if ok {
			// validate the given value
			if _, err := ads1x15GetFullScaleRange(val); err != nil {
				panic(err)
			}
			d.setChannelGains(val)
		} else if ads1x15Debug {
			log.Printf("Trying to set gain for non-ADS1x15Driver %v", c)
		}
	}
}

// WithADS1x15ChannelGain option sets the ADS1x15Driver gain for one channel.
// Valid gain settings are any of the PGA values (0..7).
func WithADS1x15ChannelGain(channel int, val int) func(Config) {
	return func(c Config) {
		d, ok := c.(*ADS1x15Driver)
		if ok {
			// validate the given value
			if _, err := ads1x15GetFullScaleRange(val); err != nil {
				panic(err)
			}
			if err := d.checkChannel(channel); err != nil {
				panic(err)
			}
			d.channelCfgs[channel].gain = val
		} else if ads1x15Debug {
			log.Printf("Trying to set channel gain for non-ADS1x15Driver %v", c)
		}
	}
}

// WithADS1x15DataRate option sets the ADS1x15Driver data rate for all channels.
// Valid gain settings are any of the DR values in SPS.
func WithADS1x15DataRate(val int) func(Config) {
	return func(c Config) {
		d, ok := c.(*ADS1x15Driver)
		if ok {
			// validate the given value
			if _, err := ads1x15GetDataRateBits(d.dataRates, val); err != nil {
				panic(err)
			}
			d.setChannelDataRates(val)
		} else if ads1x15Debug {
			log.Printf("Trying to set data rate for non-ADS1x15Driver %v", c)
		}
	}
}

// WithADS1x15ChannelDataRate option sets the ADS1x15Driver data rate for one channel.
// Valid gain settings are any of the DR values in SPS.
func WithADS1x15ChannelDataRate(channel int, val int) func(Config) {
	return func(c Config) {
		d, ok := c.(*ADS1x15Driver)
		if ok {
			// validate the given values
			if _, err := ads1x15GetDataRateBits(d.dataRates, val); err != nil {
				panic(err)
			}
			if err := d.checkChannel(channel); err != nil {
				panic(err)
			}
			d.channelCfgs[channel].dataRate = val
		} else if ads1x15Debug {
			log.Printf("Trying to set channel data rate for non-ADS1x15Driver %v", c)
		}
	}
}

// WithADS1x15WaitSingleCycle option sets the ADS1x15Driver to wait only a single cycle for conversion. According to the
// specification, chapter "Output Data Rate and Conversion Time", the device normally finishes the conversion within one
// cycle (after wake up). The cycle time depends on configured data rate and will be calculated. For unknown reasons
// some devices do not work with this setting. So the default behavior for single shot mode is to wait for a conversion
// is finished by reading the configuration register bit 15. Activating this option will switch off this behavior and
// will possibly create faster response. But, if multiple inputs are used and some inputs calculates the same result,
// most likely the device is not working with this option.
func WithADS1x15WaitSingleCycle() func(Config) {
	return func(c Config) {
		d, ok := c.(*ADS1x15Driver)
		if ok {
			d.waitOnlyOneCycle = true
		} else if ads1x15Debug {
			log.Printf("Trying to set wait single cycle for non-ADS1x15Driver %v", c)
		}
	}
}

// ReadDifferenceWithDefaults reads the difference in V between 2 inputs. It uses the default gain and data rate
// diff can be:
// * 0: Channel 0 - channel 1
// * 1: Channel 0 - channel 3
// * 2: Channel 1 - channel 3
// * 3: Channel 2 - channel 3
func (d *ADS1x15Driver) ReadDifferenceWithDefaults(diff int) (float64, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.checkChannel(diff); err != nil {
		return 0, err
	}
	return d.readVoltage(diff, 0, d.channelCfgs[diff].gain, d.channelCfgs[diff].dataRate)
}

// ReadDifference reads the difference in V between 2 inputs.
// diff can be:
// * 0: Channel 0 - channel 1
// * 1: Channel 0 - channel 3
// * 2: Channel 1 - channel 3
// * 3: Channel 2 - channel 3
func (d *ADS1x15Driver) ReadDifference(diff int, gain int, dataRate int) (float64, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.checkChannel(diff); err != nil {
		return 0, err
	}
	return d.readVoltage(diff, 0, gain, dataRate)
}

// ReadWithDefaults reads the voltage at the specified channel (between 0 and 3).
// Default values are used for the gain and data rate. The result is in V.
func (d *ADS1x15Driver) ReadWithDefaults(channel int) (float64, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.checkChannel(channel); err != nil {
		return 0, err
	}
	return d.readVoltage(channel, 0x04, d.channelCfgs[channel].gain, d.channelCfgs[channel].dataRate)
}

// Read reads the voltage at the specified channel (between 0 and 3). The result is in V.
func (d *ADS1x15Driver) Read(channel int, gain int, dataRate int) (float64, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.checkChannel(channel); err != nil {
		return 0, err
	}
	return d.readVoltage(channel, 0x04, gain, dataRate)
}

// AnalogRead returns value from analog reading of specified pin using the default values.
func (d *ADS1x15Driver) AnalogRead(pin string) (int, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	var channel int
	var channelOffset int

	// Check for the ADC is used in difference mode
	switch pin {
	case "0-1":
		channel = 0
	case "0-3":
		channel = 1
	case "1-3":
		channel = 2
	case "2-3":
		channel = 3
	default:
		// read the voltage at a specific pin, compared to the ground
		var err error
		channel, err = strconv.Atoi(pin)
		if err != nil {
			return 0, err
		}
		channelOffset = 0x04
	}

	if err := d.checkChannel(channel); err != nil {
		return 0, err
	}

	return d.rawRead(channel, channelOffset, d.channelCfgs[channel].gain, d.channelCfgs[channel].dataRate)
}

func (d *ADS1x15Driver) readVoltage(channel int, channelOffset int, gain int, dataRate int) (float64, error) {
	fsr, err := ads1x15GetFullScaleRange(gain)
	if err != nil {
		return 0, err
	}

	rawValue, err := d.rawRead(channel, channelOffset, gain, dataRate)

	// Calculate return value in V
	value := float64(rawValue) / float64(1<<15) * fsr

	return value, err
}

func (d *ADS1x15Driver) rawRead(channel int, channelOffset int, gain int, dataRate int) (int, error) {
	// Validate the passed in data rate (differs between ADS1015 and ADS1115).
	dataRateBits, err := ads1x15GetDataRateBits(d.dataRates, dataRate)
	if err != nil {
		return 0, err
	}

	var config uint16
	// Go out of power-down mode for conversion.
	config = ads1x15ConfigOsSingle

	// Specify mux value.
	mux := channel + channelOffset
	config |= uint16((mux & 0x07) << ads1x15ConfigMuxOffset) //nolint:gosec // TODO: fix later

	// Set the programmable gain amplifier bits.
	config |= uint16(gain) << ads1x15ConfigPgaOffset //nolint:gosec // TODO: fix later

	// Set the mode (continuous or single shot).
	config |= ads1x15ConfigModeSingle

	// Set the data rate.
	config |= dataRateBits

	// Disable comparator mode.
	config |= ads1x15ConfigCompQueDisable

	// Send the config value to start the ADC conversion.
	if err := d.writeWordBigEndian(ads1x15PointerConfig, config); err != nil {
		return 0, err
	}

	// Wait for the ADC sample to finish based on the sample rate plus a
	// small offset to be sure (0.1 millisecond).
	delay := time.Duration(1000000/dataRate+100) * time.Microsecond
	if err := d.waitForConversionFinished(delay); err != nil {
		return 0, err
	}

	// Retrieve the result.
	udata, err := d.readWordBigEndian(ads1x15PointerConversion)
	if err != nil {
		return 0, err
	}

	// Handle negative values as two's complement
	return int(twosComplement16Bit(udata)), nil
}

func (d *ADS1x15Driver) checkChannel(channel int) error {
	if channel < 0 || channel > 3 {
		return fmt.Errorf("Invalid channel (%d), must be between 0 and 3", channel)
	}
	return nil
}

func (d *ADS1x15Driver) waitForConversionFinished(delay time.Duration) error {
	start := time.Now()

	for i := 0; i < ads1x15WaitMaxCount; i++ {
		if i == ads1x15WaitMaxCount-1 {
			// most likely the last try will also not finish, so we stop with an error
			return fmt.Errorf("The conversion is not finished within %s", time.Since(start))
		}

		data, err := d.readWordBigEndian(ads1x15PointerConfig)
		if err != nil {
			return err
		}
		if ads1x15Debug {
			log.Printf("ADS1x15Driver: config register state: 0x%X\n", data)
		}
		// the highest bit 15: 0-device perform a conversion, 1-no conversion in progress
		if data&ads1x15ConfigOsSingle > 0 {
			break
		}
		time.Sleep(delay)
		if d.waitOnlyOneCycle {
			break
		}
	}

	if ads1x15Debug {
		elapsed := time.Since(start)
		log.Printf("conversion takes %s", elapsed)
	}

	return nil
}

func (d *ADS1x15Driver) writeWordBigEndian(reg uint8, val uint16) error {
	return d.connection.WriteWordData(reg, swapBytes(val))
}

func (d *ADS1x15Driver) readWordBigEndian(reg uint8) (uint16, error) {
	data, err := d.connection.ReadWordData(reg)
	if err != nil {
		return 0, err
	}
	return swapBytes(data), err
}

func (d *ADS1x15Driver) setChannelDataRates(ddr int) {
	for i := 0; i <= 3; i++ {
		d.channelCfgs[i].dataRate = ddr
	}
}

func (d *ADS1x15Driver) setChannelGains(gain int) {
	for i := 0; i <= 3; i++ {
		d.channelCfgs[i].gain = gain
	}
}

func ads1x15GetFullScaleRange(gain int) (float64, error) {
	if fsr, ok := ads1x15FullScaleRange[gain]; ok {
		return fsr, nil
	}

	keys := []int{}
	for k := range ads1x15FullScaleRange {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	return 0, fmt.Errorf("Gain (%d) must be one of: %d", gain, keys)
}

func ads1x15GetDataRateBits(dataRates map[int]uint16, dataRate int) (uint16, error) {
	if bits, ok := dataRates[dataRate]; ok {
		return bits, nil
	}

	keys := []int{}
	for k := range dataRates {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	return 0, fmt.Errorf("Invalid data rate (%d). Accepted values: %d", dataRate, keys)
}

// ads1x15BestGainForVoltage returns the gain the most adapted to read up to the specified difference of potential.
func ads1x15BestGainForVoltage(voltage float64) (int, error) {
	var maximum float64
	difference := math.MaxFloat64
	currentBestGain := -1

	for key, fsr := range ads1x15FullScaleRange {
		maximum = math.Max(maximum, fsr)
		newDiff := fsr - voltage
		if newDiff >= 0 && newDiff < difference {
			difference = newDiff
			currentBestGain = key
		}
	}

	if currentBestGain < 0 {
		return 0, fmt.Errorf("The maximum voltage which can be read is %f", maximum)
	}

	return currentBestGain, nil
}
