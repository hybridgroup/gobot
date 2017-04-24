package i2c

import (
	"errors"
	"strconv"
	"time"

	"gobot.io/x/gobot"
)

const ads1015Address = 0x48

// ADS1015Driver is a Driver for a ADS1015 analog to digital converter.
// Information used to create this driver came from the Adafruit C++ code
// for the ADS1015 located here:
// https://github.com/adafruit/Adafruit_ADS1X15
//
// It been tested using the Adafruit breakout board:
// https://www.adafruit.com/product/1083
//
type ADS1015Driver struct {
	name            string
	conversionDelay int
	gain            uint16
	connector       Connector
	connection      Connection
	Config
}

// NewADS1015Driver creates a new driver with specified i2c interface
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//		i2c.WithADS1015Gain(int):	input gain to use, valid options are
// 										any of the ADS1015RegConfigPga* const values
//
func NewADS1015Driver(a Connector, options ...func(Config)) *ADS1015Driver {
	d := &ADS1015Driver{
		name:            gobot.DefaultName("ADS1015"),
		connector:       a,
		conversionDelay: ADS1015ConversionDelay,
		gain:            ADS1015RegConfigPga6144V,
		Config:          NewConfig(),
	}

	for _, option := range options {
		option(d)
	}

	return d
}

// Name returns the name for this Driver
func (d *ADS1015Driver) Name() string { return d.name }

// SetName sets the name for this Driver
func (d *ADS1015Driver) SetName(n string) { d.name = n }

// Connection returns the connection for this Driver
func (d *ADS1015Driver) Connection() gobot.Connection { return d.connector.(gobot.Connection) }

// Start initializes the driver
func (d *ADS1015Driver) Start() (err error) {
	bus := d.GetBusOrDefault(d.connector.GetDefaultBus())
	address := d.GetAddressOrDefault(ads1015Address)

	d.connection, err = d.connector.GetConnection(address, bus)
	if err != nil {
		return err
	}

	return
}

// Halt returns true if device is halted successfully
func (d *ADS1015Driver) Halt() (err error) { return }

// WithADS1015Gain option sets the ADS1015Driver gain option.
// Valid gain settings are any of the ADS1015RegConfigPga* values
func WithADS1015Gain(val int) func(Config) {
	return func(c Config) {
		d, ok := c.(*ADS1015Driver)
		if ok {
			d.gain = uint16(val)
		} else {
			// TODO: return error for trying to set Gain for non-ADS1015Driver
			return
		}
	}
}

// ReadADC gets a single ADC reading from the specified channel
func (d *ADS1015Driver) ReadADC(c uint8) (val uint16, err error) {
	cfg := d.getDefaultConfig()

	switch c {
	case 0:
		cfg |= ADS1015RegConfigMuxSingle0
	case 1:
		cfg |= ADS1015RegConfigMuxSingle1
	case 2:
		cfg |= ADS1015RegConfigMuxSingle2
	case 3:
		cfg |= ADS1015RegConfigMuxSingle3
	default:
		err = errors.New("Invalid channel.")
		return
	}

	if _, err = d.connection.Write([]byte{ADS1015RegPointerConfig, byte(cfg >> 8), byte(cfg & 0xff)}); err != nil {
		return
	}

	time.Sleep(time.Duration(d.conversionDelay) * time.Millisecond)

	if _, err = d.connection.Write([]byte{ADS1015RegPointerConvert}); err != nil {
		return
	}

	b := []byte{0, 0}
	if _, err = d.connection.Read(b); err != nil {
		return
	}

	// Convert to 12-bit value
	val = uint16(b[0]&0xff) << 4
	val |= uint16(b[1]&0xff) >> 4

	return
}

// ReadADCDifference01 reads the conversion results, measuring the voltage
// difference between the P (AIN0) and N (AIN1) input.  Returns signed value
// since the difference can be either positive or negative.
//
func (d *ADS1015Driver) ReadADCDifference01() (val int16, err error) {
	cfg := d.getDefaultConfig()

	cfg |= ADS1015RegConfigMuxDiff01

	if _, err = d.connection.Write([]byte{ADS1015RegPointerConfig, byte(cfg >> 8), byte(cfg & 0xff)}); err != nil {
		return
	}

	time.Sleep(time.Duration(d.conversionDelay) * time.Millisecond)

	return d.getConversionResult()
}

// ReadADCDifference01 reads the conversion results, measuring the voltage
// difference between the P (AIN2) and N (AIN3) input. Returns signed value
// since the difference can be either positive or negative.
//
func (d *ADS1015Driver) ReadADCDifference23() (val int16, err error) {
	cfg := d.getDefaultConfig()

	cfg |= ADS1015RegConfigMuxDiff23

	if _, err = d.connection.Write([]byte{ADS1015RegPointerConfig, byte(cfg >> 8), byte(cfg & 0xff)}); err != nil {
		return
	}

	time.Sleep(time.Duration(d.conversionDelay) * time.Millisecond)

	return d.getConversionResult()
}

// AnalogRead returns value from analog reading of specified pin
func (d *ADS1015Driver) AnalogRead(pin string) (int, error) {
	switch pin {
	case "0-1":
		val, e := d.ReadADCDifference01()
		return int(val), e
	case "2-3":
		val, e := d.ReadADCDifference23()
		return int(val), e
	}

	p, e := strconv.Atoi(pin)
	if e != nil {
		return 0, e
	}

	val, e2 := d.ReadADC(uint8(p))
	return int(val), e2
}

func (d *ADS1015Driver) getDefaultConfig() uint16 {
	var cfg uint16
	cfg = ADS1015RegConfigCqueNone |
		ADS1015RegConfigClatNonLat |
		ADS1015RegConfigCpolActvLow |
		ADS1015RegConfigCmodeTrad |
		ADS1015RegConfigDr1600sps |
		ADS1015RegConfigModeSingle |
		ADS1015RegConfigOsSingle

	cfg |= d.gain

	return cfg
}

func (d *ADS1015Driver) getConversionResult() (val int16, err error) {
	d.connection.Write([]byte{ADS1015RegPointerConvert})
	b := []byte{0, 0}
	d.connection.Read(b)
	// Convert to 12-bit value
	val = int16(b[0]&0xff) << 4
	val |= int16(b[1]&0xff) >> 4

	if val&0x800 != 0 {
		val -= 1 << 12
	}

	return val, nil
}

const (
	// ADS1015ConversionDelay is the conversion delay in ms
	ADS1015ConversionDelay = 1

	// pointer register
	ADS1015RegPointerMask      = 0x03
	ADS1015RegPointerConvert   = 0x00
	ADS1015RegPointerConfig    = 0x01
	ADS1015RegPointerLowThresh = 0x02
	ADS1015RegPointerHiThresh  = 0x03

	// config register
	ADS1015RegConfigOsMask = 0x8000
	// Write: Set to start a single-conversion
	ADS1015RegConfigOsSingle = 0x8000
	// Read: Bit = 0 when conversion is in progress
	ADS1015RegConfigOsBusy = 0x0000
	// Read: Bit = 1 when device is not performing a conversion
	ADS1015RegConfigOsNotBusy = 0x8000

	ADS1015RegConfigMuxMask = 0x7000
	// Differential P = AIN0, N = AIN1 (default)
	ADS1015RegConfigMuxDiff01 = 0x0000
	// Differential P = AIN0, N = AIN3
	ADS1015RegConfigMuxDiff03 = 0x1000
	// Differential P = AIN1, N = AIN3
	ADS1015RegConfigMuxDiff13 = 0x2000
	// Differential P = AIN2, N = AIN3
	ADS1015RegConfigMuxDiff23 = 0x3000
	// Single-ended AIN0
	ADS1015RegConfigMuxSingle0 = 0x4000
	// Single-ended AIN1
	ADS1015RegConfigMuxSingle1 = 0x5000
	// Single-ended AIN2
	ADS1015RegConfigMuxSingle2 = 0x6000
	// Single-ended AIN3
	ADS1015RegConfigMuxSingle3 = 0x7000

	ADS1015RegConfigPgaMask = 0x0e00
	// +/-6.144V range = Gain 2/3
	ADS1015RegConfigPga6144V = 0x0000
	// +/-4.096V range = Gain 1
	ADS1015RegConfigPga4096V = 0x0200
	// +/-2.048V range = Gain 2 (default)
	ADS1015RegConfigPga2048V = 0x0400
	// +/-1.024V range = Gain 4
	ADS1015RegConfigPga1024V = 0x0600
	// +/-0.512V range = Gain 8
	ADS1015RegConfigPga0512V = 0x0800
	// +/-0.256V range = Gain 16
	ADS1015RegConfigPga0256V = 0x0800

	ADS1015RegConfigModeMask = 0x0100
	// Continuous conversion mode
	ADS1015RegConfigModeContin = 0x0000
	// Power-down single-shot mode (default)
	ADS1015RegConfigModeSingle = 0x0100

	ADS1015RegConfigDrMask = 0x00e0
	// 128 samples per second
	ADS1015RegConfigDr128sps = 0x0000
	// 250 samples per second
	ADS1015RegConfigDr250sps = 0x0020
	// 490 samples per second
	ADS1015RegConfigDr490sps = 0x0040
	// 960 samples per second
	ADS1015RegConfigDr960sps = 0x0060
	// 1600 samples per second
	ADS1015RegConfigDr1600sps = 0x0080
	// 2400 samples per second
	ADS1015RegConfigDr2400sps = 0x00a0
	// 3300 samples per second
	ADS1015RegConfigDr3300sps = 0x00c0

	ADS1015RegConfigCmodeMask = 0x0010
	// Traditional comparator with hysteresis (default)
	ADS1015RegConfigCmodeTrad = 0x0000
	// Window comparator
	ADS1015RegConfigCmodeWindow = 0x0010

	ADS1015RegConfigCpolMask = 0x0009
	// ALERT/RDY pin is low when active (default)
	ADS1015RegConfigCpolActvLow = 0x0000
	// ALERT/RDY pin is high when active
	ADS1015RegConfigCpolActvHi = 0x0008

	// Determines if ALERT/RDY pin latches once asserted
	ADS1015RegConfigClatMask = 0x0004
	// Non-latching comparator (default)
	ADS1015RegConfigClatNonLat = 0x0000
	// Latching comparator
	ADS1015RegConfigClatLat = 0x0004

	ADS1015RegConfigCqueMask = 0x0003
	// Assert ALERT/RDY after one conversions
	ADS1015RegConfigCque1Conv = 0x0000
	// Assert ALERT/RDY after two conversions
	ADS1015RegConfigCque2Conv = 0x0001
	// Assert ALERT/RDY after four conversions
	ADS1015RegConfigCque4Conv = 0x0002
	// Disable the comparator and put ALERT/RDY in high state (default)
	ADS1015RegConfigCqueNone = 0x0003
)
