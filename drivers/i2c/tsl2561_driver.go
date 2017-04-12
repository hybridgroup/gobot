package i2c

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
)

const (
	// TSL2561AddressLow - the address of the device when address pin is low
	TSL2561AddressLow = 0x29
	// TSL2561AddressFloat - the address of the device when address pin is floating
	TSL2561AddressFloat = 0x39
	// TSL2561AddressHigh - the address of the device when address pin is high
	TSL2561AddressHigh = 0x49

	tsl2561CommandBit = 0x80 // Must be 1
	tsl2561ClearBit   = 0x40 // Clears any pending interrupt (write 1 to clear)
	tsl2561WordBit    = 0x20 // 1 = read/write word (rather than byte)
	tsl2561BlockBit   = 0x10 // 1 = using block read/write

	tsl2561ControlPowerOn  = 0x03
	tsl2561ControlPowerOff = 0x00

	tsl2561LuxLuxScale     = 14     // Scale by 2^14
	tsl2561LuxRatioScale   = 9      // Scale ratio by 2^9
	tsl2561LuxChScale      = 10     // Scale channel values by 2^10
	tsl2561LuxCHScaleTInt0 = 0x7517 // 322/11 * 2^tsl2561LUXCHSCALE
	tsl2561LuxChScaleTInt1 = 0x0FE7 // 322/81 * 2^tsl2561LUXCHSCALE

	// T, FN and CL package values
	tsl2561LuxK1T = 0x0040 // 0.125 * 2^RATIO_SCALE
	tsl2561LuxB1T = 0x01f2 // 0.0304 * 2^LUX_SCALE
	tsl2561LuxM1T = 0x01be // 0.0272 * 2^LUX_SCALE
	tsl2561LuxK2T = 0x0080 // 0.250 * 2^RATIO_SCALE
	tsl2561LuxB2T = 0x0214 // 0.0325 * 2^LUX_SCALE
	tsl2561LuxM2T = 0x02d1 // 0.0440 * 2^LUX_SCALE
	tsl2561LuxK3T = 0x00c0 // 0.375 * 2^RATIO_SCALE
	tsl2561LuxB3T = 0x023f // 0.0351 * 2^LUX_SCALE
	tsl2561LuxM3T = 0x037b // 0.0544 * 2^LUX_SCALE
	tsl2561LuxK4T = 0x0100 // 0.50 * 2^RATIO_SCALE
	tsl2561LuxB4T = 0x0270 // 0.0381 * 2^LUX_SCALE
	tsl2561LuxM4T = 0x03fe // 0.0624 * 2^LUX_SCALE
	tsl2561LuxK5T = 0x0138 // 0.61 * 2^RATIO_SCALE
	tsl2561LuxB5T = 0x016f // 0.0224 * 2^LUX_SCALE
	tsl2561LuxM5T = 0x01fc // 0.0310 * 2^LUX_SCALE
	tsl2561LuxK6T = 0x019a // 0.80 * 2^RATIO_SCALE
	tsl2561LuxB6T = 0x00d2 // 0.0128 * 2^LUX_SCALE
	tsl2561LuxM6T = 0x00fb // 0.0153 * 2^LUX_SCALE
	tsl2561LuxK7T = 0x029a // 1.3 * 2^RATIO_SCALE
	tsl2561LuxB7T = 0x0018 // 0.00146 * 2^LUX_SCALE
	tsl2561LuxM7T = 0x0012 // 0.00112 * 2^LUX_SCALE
	tsl2561LuxK8T = 0x029a // 1.3 * 2^RATIO_SCALE
	tsl2561LuxB8T = 0x0000 // 0.000 * 2^LUX_SCALE
	tsl2561LuxM8T = 0x0000 // 0.000 * 2^LUX_SCALE

	// Auto-gain thresholds
	tsl2561AgcTHi13MS  = 4850 // Max value at Ti 13ms = 5047
	tsl2561AgcTLo13MS  = 100
	tsl2561AgcTHi101MS = 36000 // Max value at Ti 101ms = 37177
	tsl2561AgcTLo101MS = 200
	tsl2561AgcTHi402MS = 63000 // Max value at Ti 402ms = 65535
	tsl2561AgcTLo402MS = 500

	// Clipping thresholds
	tsl2561Clipping13MS  = 4900
	tsl2561Clipping101MS = 37000
	tsl2561Clipping402MS = 65000
)

const (
	tsl2561RegisterControl         = 0x00
	tsl2561RegisterTiming          = 0x01
	tsl2561RegisterThreshholdLLow  = 0x02
	tsl2561RegisterThreshholdLHigh = 0x03
	tsl2561RegisterThreshholdHLow  = 0x04
	tsl2561RegisterThreshholdHHigh = 0x05
	tsl2561RegisterInterrupt       = 0x06
	tsl2561RegisterCRC             = 0x08
	tsl2561RegisterID              = 0x0A
	tsl2561RegisterChan0Low        = 0x0C
	tsl2561RegisterChan0High       = 0x0D
	tsl2561RegisterChan1Low        = 0x0E
	tsl2561RegisterChan1High       = 0x0F
)

// TSL2561IntegrationTime is the type of all valid integration time settings
type TSL2561IntegrationTime int

const (
	// TSL2561IntegrationTime13MS integration time 13ms
	TSL2561IntegrationTime13MS TSL2561IntegrationTime = iota // 13.7ms
	// TSL2561IntegrationTime101MS integration time 101ms
	TSL2561IntegrationTime101MS // 101ms
	// TSL2561IntegrationTime402MS integration time 402ms
	TSL2561IntegrationTime402MS // 402ms
)

// TSL2561Gain is the type of all valid gain settings
type TSL2561Gain int

const (
	// TSL2561Gain1X gain == 1x
	TSL2561Gain1X TSL2561Gain = 0x00 // No gain
	// TSL2561Gain16X gain == 16x
	TSL2561Gain16X = 0x10 // 16x gain
)

// TSL2561Driver is the gobot driver for the Adafruit Digital Luminosity/Lux/Light Sensor
//
// Datasheet: http://www.adafruit.com/datasheets/TSL2561.pdf
//
// Ported from the Adafruit driver at https://github.com/adafruit/Adafruit_TSL2561 by
// K. Townsend
type TSL2561Driver struct {
	name       string
	connector  Connector
	connection Connection
	Config
	autoGain        bool
	gain            TSL2561Gain
	integrationTime TSL2561IntegrationTime
}

// NewTSL2561Driver creates a new driver for the TSL2561 device.
//
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):		bus to use with this driver
//		i2c.WithAddress(int):		address to use with this driver
//		i2c.WithTSL2561Gain1X:		sets the gain to 1X
//		i2c.WithTSL2561Gain16X:		sets the gain to 16X
//		i2c.WithTSL2561AutoGain:	turns on auto gain
//		i2c.WithTSL2561IntegrationTime13MS:	sets integration time to 13ms
//		i2c.WithTSL2561IntegrationTime101MS: 	sets integration time to 101ms
//		i2c.WithTSL2561IntegrationTime402MS: 	sets integration time to 402ms
//
func NewTSL2561Driver(conn Connector, options ...func(Config)) *TSL2561Driver {
	driver := &TSL2561Driver{
		name:            gobot.DefaultName("TSL2561"),
		connector:       conn,
		Config:          NewConfig(),
		integrationTime: TSL2561IntegrationTime402MS,
		gain:            TSL2561Gain1X,
		autoGain:        false,
	}

	for _, option := range options {
		option(driver)
	}

	return driver
}

// WithTSL2561Gain1X option sets the TSL2561Driver gain to 1X
func WithTSL2561Gain1X(c Config) {
	d, ok := c.(*TSL2561Driver)
	if ok {
		d.gain = TSL2561Gain1X
		return
	}
	// TODO: return errors.New("Trying to set Gain for non-TSL2561Driver")
}

// WithTSL2561Gain16X option sets the TSL2561Driver gain to 16X
func WithTSL2561Gain16X(c Config) {
	d, ok := c.(*TSL2561Driver)
	if ok {
		d.gain = TSL2561Gain16X
		return
	}
	// TODO: return errors.New("Trying to set Gain for non-TSL2561Driver")
}

// WithTSL2561AutoGain option turns on TSL2561Driver auto gain
func WithTSL2561AutoGain(c Config) {
	d, ok := c.(*TSL2561Driver)
	if ok {
		d.autoGain = true
		return
	}
	// TODO: return errors.New("Trying to set Auto Gain for non-TSL2561Driver")
}

func withTSL2561IntegrationTime(iTime TSL2561IntegrationTime) func(Config) {
	return func(c Config) {
		d, ok := c.(*TSL2561Driver)
		if ok {
			d.integrationTime = iTime
			return
		}
		// TODO: return errors.New("Trying to set integration time for non-TSL2561Driver")
	}
}

// WithTSL2561IntegrationTime13MS option sets the TSL2561Driver integration time
// to 13ms
func WithTSL2561IntegrationTime13MS(c Config) {
	withTSL2561IntegrationTime(TSL2561IntegrationTime13MS)(c)
}

// WithTSL2561IntegrationTime101MS option sets the TSL2561Driver integration time
// to 101ms
func WithTSL2561IntegrationTime101MS(c Config) {
	withTSL2561IntegrationTime(TSL2561IntegrationTime101MS)(c)
}

// WithTSL2561IntegrationTime402MS option sets the TSL2561Driver integration time
// to 402ms
func WithTSL2561IntegrationTime402MS(c Config) {
	withTSL2561IntegrationTime(TSL2561IntegrationTime402MS)(c)
}

// Name returns the name of the device.
func (d *TSL2561Driver) Name() string {
	return d.name
}

// SetName sets the name of the device.
func (d *TSL2561Driver) SetName(name string) {
	d.name = name
}

// Connection returns the connection of the device.
func (d *TSL2561Driver) Connection() gobot.Connection {
	return d.connector.(gobot.Connection)
}

// Start initializes the device.
func (d *TSL2561Driver) Start() (err error) {
	bus := d.GetBusOrDefault(d.connector.GetDefaultBus())
	address := d.GetAddressOrDefault(TSL2561AddressFloat)

	if d.connection, err = d.connector.GetConnection(address, bus); err != nil {
		return err
	}

	if err = d.enable(); err != nil {
		return err
	}

	var initialized byte
	if initialized, err = d.connection.ReadByteData(tsl2561RegisterID); err != nil {
		return err
	} else if (initialized & 0x0A) == 0 {
		return fmt.Errorf("TSL2561 device not found (0x%X)", initialized)
	}

	if err = d.SetIntegrationTime(d.integrationTime); err != nil {
		return err
	}

	if err = d.SetGain(d.gain); err != nil {
		return err
	}

	if err = d.disable(); err != nil {
		return err
	}

	return nil
}

// Halt stops the device
func (d *TSL2561Driver) Halt() error {
	return nil
}

// SetIntegrationTime sets integrations time for the TSL2561
func (d *TSL2561Driver) SetIntegrationTime(time TSL2561IntegrationTime) error {
	if err := d.enable(); err != nil {
		return err
	}

	timeGainVal := uint8(time) | uint8(d.gain)
	if err := d.connection.WriteByteData(tsl2561CommandBit|tsl2561RegisterTiming, timeGainVal); err != nil {
		return err
	}
	d.integrationTime = time

	return d.disable()
}

// SetGain adjusts the TSL2561 gain (sensitivity to light)
func (d *TSL2561Driver) SetGain(gain TSL2561Gain) error {
	if err := d.enable(); err != nil {
		return err
	}

	timeGainVal := uint8(d.integrationTime) | uint8(gain)
	if err := d.connection.WriteByteData(tsl2561CommandBit|tsl2561RegisterTiming, timeGainVal); err != nil {
		return err
	}
	d.gain = gain

	return d.disable()
}

// GetLuminocity gets the broadband and IR only values from the TSL2561,
// adjusting gain if auto-gain is enabled
func (d *TSL2561Driver) GetLuminocity() (broadband uint16, ir uint16, err error) {
	// if auto gain disabled get a single reading and continue
	if !d.autoGain {
		broadband, ir, err = d.getData()
		return
	}

	agcCheck := false
	hi, lo := d.getHiLo()

	// Read data until we find a valid range
	valid := false
	for {
		broadband, ir, err = d.getData()
		if err != nil {
			return
		}

		// Run an auto-gain check if we haven't already done so
		if !agcCheck {
			if (broadband < lo) && (d.gain == TSL2561Gain1X) {
				// increase gain and try again
				err = d.SetGain(TSL2561Gain16X)
				if err != nil {
					return
				}
				agcCheck = true
			} else if (broadband > hi) && (d.gain == TSL2561Gain16X) {
				// drop gain and try again
				err = d.SetGain(TSL2561Gain1X)
				if err != nil {
					return
				}
				agcCheck = true
			} else {
				// Reading is either valid, or we're already at the chips
				// limits
				valid = true
			}
		} else {
			// If we've already adjusted the gain once, just return the new results.
			// This avoids endless loops where a value is at one extreme pre-gain,
			// and the the other extreme post-gain
			valid = true
		}

		if valid {
			break
		}
	}

	return
}

// CalculateLux converts raw sensor values to the standard SI Lux equivalent.
// Returns 65536 if the sensor is saturated.
func (d *TSL2561Driver) CalculateLux(broadband uint16, ir uint16) (lux uint32) {
	var channel1 uint32
	var channel0 uint32

	// Set cliplevel and scaling based on integration time
	clipThreshold, chScale := d.getClipScaling()

	// Saturated sensor
	if (broadband > clipThreshold) || (ir > clipThreshold) {
		return 65536
	}

	// Adjust scale for gain
	if d.gain == TSL2561Gain1X {
		chScale = chScale * 16
	}

	channel0 = (uint32(broadband) * chScale) >> tsl2561LuxChScale
	channel1 = (uint32(ir) * chScale) >> tsl2561LuxChScale

	// Find the ratio of the channel values (channel1 / channel0)
	var ratio1 uint32
	if channel0 != 0 {
		ratio1 = (channel1 << (tsl2561LuxRatioScale + 1)) / channel0
	}

	// Round the ratio value
	ratio := (ratio1 + 1) / 2

	b, m := d.getBM(ratio)
	temp := (channel0 * b) - (channel1 * m)

	// Negative lux not allowed
	if temp < 0 {
		temp = 0
	}

	// Round lsb (2^(LUX_SCALE+1))
	temp += (1 << (tsl2561LuxLuxScale - 1))

	// Strip off fractional portion
	lux = temp >> tsl2561LuxLuxScale

	return lux
}

func (d *TSL2561Driver) enable() (err error) {
	err = d.connection.WriteByteData(uint8(tsl2561CommandBit|tsl2561RegisterControl), tsl2561ControlPowerOn)
	return err
}

func (d *TSL2561Driver) disable() (err error) {
	err = d.connection.WriteByteData(uint8(tsl2561CommandBit|tsl2561RegisterControl), tsl2561ControlPowerOff)
	return err
}

func (d *TSL2561Driver) getData() (broadband uint16, ir uint16, err error) {
	if err = d.enable(); err != nil {
		return
	}

	d.waitForADC()

	// Reads a two byte value from channel 0 (visible + infrared)
	broadband, err = d.connection.ReadWordData(tsl2561CommandBit | tsl2561WordBit | tsl2561RegisterChan0Low)
	if err != nil {
		return
	}

	// Reads a two byte value from channel 1 (infrared)
	ir, err = d.connection.ReadWordData(tsl2561CommandBit | tsl2561WordBit | tsl2561RegisterChan1Low)
	if err != nil {
		return
	}

	err = d.disable()

	return
}

func (d *TSL2561Driver) getHiLo() (hi, lo uint16) {
	switch d.integrationTime {
	case TSL2561IntegrationTime13MS:
		hi = tsl2561AgcTHi13MS
		lo = tsl2561AgcTLo13MS
	case TSL2561IntegrationTime101MS:
		hi = tsl2561AgcTHi101MS
		lo = tsl2561AgcTLo101MS
	case TSL2561IntegrationTime402MS:
		hi = tsl2561AgcTHi402MS
		lo = tsl2561AgcTLo402MS
	}
	return
}

func (d *TSL2561Driver) getClipScaling() (clipThreshold uint16, chScale uint32) {
	switch d.integrationTime {
	case TSL2561IntegrationTime13MS:
		clipThreshold = tsl2561Clipping13MS
		chScale = tsl2561LuxCHScaleTInt0
	case TSL2561IntegrationTime101MS:
		clipThreshold = tsl2561Clipping101MS
		chScale = tsl2561LuxChScaleTInt1
	case TSL2561IntegrationTime402MS:
		clipThreshold = tsl2561Clipping402MS
		chScale = (1 << tsl2561LuxChScale)
	}
	return
}

func (d *TSL2561Driver) getBM(ratio uint32) (b uint32, m uint32) {
	switch {
	case (ratio >= 0) && (ratio <= tsl2561LuxK1T):
		b = tsl2561LuxB1T
		m = tsl2561LuxM1T
	case (ratio <= tsl2561LuxK2T):
		b = tsl2561LuxB2T
		m = tsl2561LuxM2T
	case (ratio <= tsl2561LuxK3T):
		b = tsl2561LuxB3T
		m = tsl2561LuxM3T
	case (ratio <= tsl2561LuxK4T):
		b = tsl2561LuxB4T
		m = tsl2561LuxM4T
	case (ratio <= tsl2561LuxK5T):
		b = tsl2561LuxB5T
		m = tsl2561LuxM5T
	case (ratio <= tsl2561LuxK6T):
		b = tsl2561LuxB6T
		m = tsl2561LuxM6T
	case (ratio <= tsl2561LuxK7T):
		b = tsl2561LuxB7T
		m = tsl2561LuxM7T
	case (ratio > tsl2561LuxK8T): // TODO: there is a gap here...
		b = tsl2561LuxB8T
		m = tsl2561LuxM8T
	}
	return
}

func (d *TSL2561Driver) waitForADC() {
	switch d.integrationTime {
	case TSL2561IntegrationTime13MS:
		time.Sleep(15 * time.Millisecond)
	case TSL2561IntegrationTime101MS:
		time.Sleep(120 * time.Millisecond)
	case TSL2561IntegrationTime402MS:
		time.Sleep(450 * time.Millisecond)
	}
	return
}
