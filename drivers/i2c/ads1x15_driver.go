package i2c

import (
	"errors"
	"math"
	"strconv"
	"time"

	"fmt"

	"gobot.io/x/gobot"
)

const (

	// ADS1x15DefaultAddress is the default I2C address for the ADS1x15 components
	ADS1x15DefaultAddress = 0x48

	ads1x15PointerConversion    = 0x00
	ads1x15PointerConfig        = 0x01
	ads1x15PointerLowThreshold  = 0x02
	ads1x15PointerHighThreshold = 0x03
	// Write: Set to start a single-conversion
	ads1x15ConfigOsSingle       = 0x8000
	ads1x15ConfigMuxOffset      = 12
	ads1x15ConfigModeContinuous = 0x0000
	//Single shoot mode
	ads1x15ConfigModeSingle = 0x0100

	ads1x15ConfigCompWindow      = 0x0010
	ads1x15ConfigCompAactiveHigh = 0x0008
	ads1x15ConfigCompLatching    = 0x0004
	ads1x15ConfigCompQueDisable  = 0x0003
)

// ADS1x15Driver is the Gobot driver for the ADS1015/ADS1115 ADC
type ADS1x15Driver struct {
	name            string
	connector       Connector
	connection      Connection
	gainConfig      map[int]uint16
	dataRates       map[int]uint16
	gainVoltage     map[int]float64
	converter       func([]byte) float64
	DefaultGain     int
	DefaultDataRate int
	Config
}

// NewADS1015Driver creates a new driver for the ADS1015 (12-bit ADC)
// Largely inspired by: https://github.com/adafruit/Adafruit_Python_ADS1x15
func NewADS1015Driver(a Connector, options ...func(Config)) *ADS1x15Driver {
	l := newADS1x15Driver(a, options...)

	l.dataRates = map[int]uint16{
		128:  0x0000,
		250:  0x0020,
		490:  0x0040,
		920:  0x0060,
		1600: 0x0080,
		2400: 0x00A0,
		3300: 0x00C0,
	}
	if l.DefaultDataRate == 0 {
		l.DefaultDataRate = 1600
	}

	l.converter = func(data []byte) (value float64) {
		result := (int(data[0]) << 8) | int(data[1])

		if result&0x8000 != 0 {
			result -= 1 << 16
		}

		return float64(result) / float64(1<<15)
	}

	return l
}

// NewADS1115Driver creates a new driver for the ADS1115 (16-bit ADC)
func NewADS1115Driver(a Connector, options ...func(Config)) *ADS1x15Driver {
	l := newADS1x15Driver(a, options...)

	l.dataRates = map[int]uint16{
		8:   0x0000,
		16:  0x0020,
		32:  0x0040,
		64:  0x0060,
		128: 0x0080,
		250: 0x00A0,
		475: 0x00C0,
		860: 0x00E0,
	}

	if l.DefaultDataRate == 0 {
		l.DefaultDataRate = 128
	}

	l.converter = func(data []byte) (value float64) {
		result := (int(data[0]) << 8) | int(data[1])

		if result&0x8000 != 0 {
			result -= 1 << 16
		}

		return float64(result) / float64(1<<15)
	}

	return l
}

func newADS1x15Driver(a Connector, options ...func(Config)) *ADS1x15Driver {
	l := &ADS1x15Driver{
		name:      gobot.DefaultName("ADS1x15"),
		connector: a,
		// Mapping of gain values to config register values.
		gainConfig: map[int]uint16{
			2 / 3: 0x0000,
			1:     0x0200,
			2:     0x0400,
			4:     0x0600,
			8:     0x0800,
			16:    0x0A00,
		},
		gainVoltage: map[int]float64{
			2 / 3: 6.144,
			1:     4.096,
			2:     2.048,
			4:     1.024,
			8:     0.512,
			16:    0.256,
		},
		DefaultGain: 1,

		Config: NewConfig(),
	}

	for _, option := range options {
		option(l)
	}

	// TODO: add commands to API
	return l
}

// Start initializes the sensor
func (d *ADS1x15Driver) Start() (err error) {
	bus := d.GetBusOrDefault(d.connector.GetDefaultBus())
	address := d.GetAddressOrDefault(ADS1x15DefaultAddress)

	if d.connection, err = d.connector.GetConnection(address, bus); err != nil {
		return err
	}

	return
}

// Name returns the Name for the Driver
func (d *ADS1x15Driver) Name() string { return d.name }

// SetName sets the Name for the Driver
func (d *ADS1x15Driver) SetName(n string) { d.name = n }

// Connection returns the connection for the Driver
func (d *ADS1x15Driver) Connection() gobot.Connection { return d.connector.(gobot.Connection) }

// Halt returns true if devices is halted successfully
func (d *ADS1x15Driver) Halt() (err error) { return }

// WithADS1x15Gain option sets the ADS1x15Driver gain option.
// Valid gain settings are any of the ADS1x15RegConfigPga* values
func WithADS1x15Gain(val int) func(Config) {
	return func(c Config) {
		d, ok := c.(*ADS1x15Driver)
		if ok {
			d.DefaultGain = val
		} else {
			// TODO: return error for trying to set Gain for non-ADS1015Driver
			return
		}
	}
}

// WithADS1x15DataRate option sets the ADS1x15Driver data rate option.
// Valid gain settings are any of the ADS1x15RegConfigPga* values
func WithADS1x15DataRate(val int) func(Config) {
	return func(c Config) {
		d, ok := c.(*ADS1x15Driver)
		if ok {
			d.DefaultDataRate = val
		} else {
			// TODO: return error for trying to set data rate for non-ADS1015Driver
			return
		}
	}
}

// BestGainForVoltage returns the gain the most adapted to read up to the specified difference of potential.
func (d *ADS1x15Driver) BestGainForVoltage(voltage float64) (bestGain int, err error) {
	var max float64
	difference := math.MaxFloat64
	currentBestGain := -1

	for key, value := range d.gainVoltage {
		max = math.Max(max, value)
		newDiff := value - voltage
		if newDiff >= 0 && newDiff < difference {
			difference = newDiff
			currentBestGain = key
		}
	}

	if currentBestGain < 0 {
		err = fmt.Errorf("The maximum voltage which can be read is %f", max)
		return
	}

	bestGain = currentBestGain
	return
}

// ReadDifferenceWithDefaults reads the difference in V between 2 inputs. It uses the default gain and data rate
// diff can be:
// * 0: Channel 0 - channel 1
// * 1: Channel 0 - channel 3
// * 2: Channel 1 - channel 3
// * 3: Channel 2 - channel 3
func (d *ADS1x15Driver) ReadDifferenceWithDefaults(diff int) (value float64, err error) {
	return d.ReadDifference(diff, d.DefaultGain, d.DefaultDataRate)
}

// ReadDifference reads the difference in V between 2 inputs.
// diff can be:
// * 0: Channel 0 - channel 1
// * 1: Channel 0 - channel 3
// * 2: Channel 1 - channel 3
// * 3: Channel 2 - channel 3
func (d *ADS1x15Driver) ReadDifference(diff int, gain int, dataRate int) (value float64, err error) {
	if err = d.checkChannel(diff); err != nil {
		return
	}

	return d.rawRead(diff, gain, dataRate)
}

// ReadWithDefaults reads the voltage at the specified channel (between 0 and 3).
// Default values are used for the gain and data rate. The result is in V.
func (d *ADS1x15Driver) ReadWithDefaults(channel int) (value float64, err error) {
	return d.Read(channel, d.DefaultGain, d.DefaultDataRate)
}

// Read reads the voltage at the specified channel (between 0 and 3). The result is in V.
func (d *ADS1x15Driver) Read(channel int, gain int, dataRate int) (value float64, err error) {
	if err = d.checkChannel(channel); err != nil {
		return
	}
	mux := channel + 0x04

	return d.rawRead(mux, gain, dataRate)
}

// AnalogRead returns value from analog reading of specified pin
func (d *ADS1x15Driver) AnalogRead(pin string) (value int, err error) {
	var useDifference = false
	var channel int
	var read float64

	// First case: the ADC is used in difference mode
	switch pin {
	case "0-1":
		useDifference = true
		channel = 0
		break
	case "0-3":
		useDifference = true
		channel = 1
		break
	case "1-3":
		useDifference = true
		channel = 2
		break
	case "2-3":
		useDifference = true
		channel = 3
		break
	}

	if useDifference {
		read, err = d.ReadDifferenceWithDefaults(channel)
	} else {
		// Second case: read the voltage at a specific pin, compared to the ground
		channel, err = strconv.Atoi(pin)
		if err != nil {
			return
		}

		read, err = d.ReadWithDefaults(channel)
	}

	if err == nil {
		value = int(gobot.ToScale(gobot.FromScale(read, 0, d.gainVoltage[d.DefaultGain]), 0, 1023))
	}

	return
}

func (d *ADS1x15Driver) rawRead(mux int, gain int, dataRate int) (value float64, err error) {
	var config uint16
	config = ads1x15ConfigOsSingle // Go out of power-down mode for conversion.
	// Specify mux value.
	config |= uint16((mux & 0x07) << ads1x15ConfigMuxOffset)
	// Validate the passed in gain and then set it in the config.

	gainConf, ok := d.gainConfig[gain]

	if !ok {
		err = errors.New("Gain must be one of: 2/3, 1, 2, 4, 8, 16")
		return
	}
	config |= gainConf
	// Set the mode (continuous or single shot).
	config |= ads1x15ConfigModeSingle
	// Get the default data rate if none is specified (default differs between
	// ADS1015 and ADS1115).
	dataRateConf, ok := d.dataRates[dataRate]

	if !ok {
		keys := []int{}
		for k := range d.dataRates {
			keys = append(keys, k)
		}

		err = fmt.Errorf("Invalid data rate. Accepted values: %d", keys)
		return
	}
	// Set the data rate (this is controlled by the subclass as it differs
	// between ADS1015 and ADS1115).
	config |= dataRateConf
	config |= ads1x15ConfigCompQueDisable // Disable comparator mode.

	// Send the config value to start the ADC conversion.
	// Explicitly break the 16-bit value down to a big endian pair of bytes.
	if _, err = d.connection.Write([]byte{ads1x15PointerConfig, byte((config >> 8) & 0xFF), byte(config & 0xFF)}); err != nil {
		return
	}

	// Wait for the ADC sample to finish based on the sample rate plus a
	// small offset to be sure (0.1 millisecond).
	time.Sleep(time.Duration(1000000/dataRate+100) * time.Microsecond)

	// Retrieve the result.
	if _, err = d.connection.Write([]byte{ads1x15PointerConversion}); err != nil {
		return
	}

	data := make([]byte, 2)
	_, err = d.connection.Read(data)
	if err != nil {
		return
	}

	voltageMultiplier, ok := d.gainVoltage[gain]

	if !ok {
		err = errors.New("Gain must be one of: 2/3, 1, 2, 4, 8, 16")
		return
	}

	value = d.converter(data) * voltageMultiplier

	return
}

func (d *ADS1x15Driver) checkChannel(channel int) (err error) {
	if channel < 0 || channel > 3 {
		err = errors.New("Invalid channel, must be between 0 and 3")
	}
	return
}
