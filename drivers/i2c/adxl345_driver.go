package i2c

import (
	"encoding/binary"
	"fmt"
	"log"
)

const adxl345Debug = false

// ADXL345 supports 2 addresses, which can be changed by the address pin, there is no internal pull-up/down resistor!
// pin to GND: 0x53, pin to VDD: 0x1D
const (
	ADXL345AddressPullUp  = 0x1D // can be used by WithAddress()
	adxl345DefaultAddress = 0x53
)

type (
	ADXL345RateConfig    uint8
	ADXL345FsRangeConfig uint8
)

const (
	// registers are named according to the datasheet
	adxl345Reg_DEVID          = 0x00 // R,     11100101,   Device ID
	adxl345Reg_THRESH_TAP     = 0x1D // R/W,   00000000,   Tap threshold
	adxl345Reg_OFSX           = 0x1E // R/W,   00000000,   X-axis offset
	adxl345Reg_OFSY           = 0x1F // R/W,   00000000,   Y-axis offset
	adxl345Reg_OFSZ           = 0x20 // R/W,   00000000,   Z-axis offset
	adxl345Reg_DUR            = 0x21 // R/W,   00000000,   Tap duration
	adxl345Reg_LATENT         = 0x22 // R/W,   00000000,   Tap latency
	adxl345Reg_WINDOW         = 0x23 // R/W,   00000000,   Tap window
	adxl345Reg_THRESH_ACT     = 0x24 // R/W,   00000000,   Activity threshold
	adxl345Reg_THRESH_INACT   = 0x25 // R/W,   00000000,   Inactivity threshold
	adxl345Reg_TIME_INACT     = 0x26 // R/W,   00000000,   Inactivity time
	adxl345Reg_ACT_INACT_CTL  = 0x27 // R/W,   00000000,   Axis enable control for activity and inactivity detection
	adxl345Reg_THRESH_FF      = 0x28 // R/W,   00000000,   Free-fall threshold
	adxl345Reg_TIME_FF        = 0x29 // R/W,   00000000,   Free-fall time
	adxl345Reg_TAP_AXES       = 0x2A // R/W,   00000000,   Axis control for single tap/double tap
	adxl345Reg_ACT_TAP_STATUS = 0x2B // R,     00000000,   Source of single tap/double tap
	adxl345Reg_BW_RATE        = 0x2C // R/W,   00001010,   Data rate and power mode control
	adxl345Reg_POWER_CTL      = 0x2D // R/W,   00000000,   Power-saving features control
	adxl345Reg_INT_ENABLE     = 0x2E // R/W,   00000000,   Interrupt enable control
	adxl345Reg_INT_MAP        = 0x2F // R/W,   00000000,   Interrupt mapping control
	adxl345Reg_INT_SOUCE      = 0x30 // R,     00000010,   Source of interrupts
	adxl345Reg_DATA_FORMAT    = 0x31 // R/W,   00000000,   Data format control (FS range, justify, full resolution)
	adxl345Reg_DATAX0         = 0x32 // R,     00000000,   X-Axis Data 0 (LSByte)
	adxl345Reg_DATAX1         = 0x33 // R,     00000000,   X-Axis Data 1 (MSByte)
	adxl345Reg_DATAY0         = 0x34 // R,     00000000,   Y-Axis Data 0
	adxl345Reg_DATAY1         = 0x35 // R,     00000000,   Y-Axis Data 1
	adxl345Reg_DATAZ0         = 0x36 // R,     00000000,   Z-Axis Data 0
	adxl345Reg_DATAZ1         = 0x37 // R,     00000000,   Z-Axis Data 1
	adxl345Reg_FIFO_CTL       = 0x38 // R/W,   00000000,   FIFO control
	adxl345Reg_FIFO_STATUS    = 0x39 // R,     00000000,   FIFO status

	adxl345Rate_LowPowerBit = 0x10 // set the device to low power, but increase the noise by ~2.5x

	ADXL345Rate_100mHZ   ADXL345RateConfig = 0x00 // 0.10 Hz
	ADXL345Rate_200mHZ   ADXL345RateConfig = 0x01 // 0.20 Hz
	ADXL345Rate_390mHZ   ADXL345RateConfig = 0x02 // 0.39 Hz
	ADXL345Rate_780mHZ   ADXL345RateConfig = 0x03 // 0.78 Hz
	ADXL345Rate_1560mHZ  ADXL345RateConfig = 0x04 // 1.56 Hz
	ADXL345Rate_3130mHZ  ADXL345RateConfig = 0x05 // 3.13 Hz
	ADXL345Rate_6250mHZ  ADXL345RateConfig = 0x06 // 6.25 Hz
	ADXL345Rate_12500mHZ ADXL345RateConfig = 0x07 // 12.5 Hz
	ADXL345Rate_25HZ     ADXL345RateConfig = 0x08 // 25 Hz
	ADXL345Rate_50HZ     ADXL345RateConfig = 0x09 // 50 Hz
	ADXL345Rate_100HZ    ADXL345RateConfig = 0x0A // 100 Hz
	ADXL345Rate_200HZ    ADXL345RateConfig = 0x0B // 200 Hz
	ADXL345Rate_400HZ    ADXL345RateConfig = 0x0C // 400 Hz
	ADXL345Rate_800HZ    ADXL345RateConfig = 0x0D // 800 Hz
	ADXL345Rate_1600HZ   ADXL345RateConfig = 0x0E // 1600 Hz
	ADXL345Rate_3200HZ   ADXL345RateConfig = 0x0F // 3200 Hz

	ADXL345FsRange_2G  ADXL345FsRangeConfig = 0x00 // +-2 g
	ADXL345FsRange_4G  ADXL345FsRangeConfig = 0x01 // +-4 g
	ADXL345FsRange_8G  ADXL345FsRangeConfig = 0x02 // +-8 g
	ADXL345FsRange_16G ADXL345FsRangeConfig = 0x03 // +-16 g)
)

// ADXL345Driver is the gobot driver for the digital accelerometer ADXL345
//
// Datasheet EN: http://www.analog.com/media/en/technical-documentation/data-sheets/ADXL345.pdf
// Datasheet JP: http://www.analog.com/media/jp/technical-documentation/data-sheets/ADXL345_jp.pdf
//
// Ported from the Arduino driver https://github.com/jakalada/Arduino-ADXL345
type ADXL345Driver struct {
	*Driver
	powerCtl   adxl345PowerCtl
	dataFormat adxl345DataFormat
	bwRate     adxl345BwRate
}

// Internal structure for the power configuration
type adxl345PowerCtl struct {
	link      uint8
	autoSleep uint8
	measure   uint8
	sleep     uint8
	wakeUp    uint8
}

// Internal structure for the sensor's data format configuration
type adxl345DataFormat struct {
	selfTest       uint8
	spi            uint8
	intInvert      uint8
	fullRes        uint8
	justify        uint8
	fullScaleRange ADXL345FsRangeConfig
}

// Internal structure for the sampling rate configuration
type adxl345BwRate struct {
	lowPower bool
	rate     ADXL345RateConfig
}

// NewADXL345Driver creates a new driver with specified i2c interface
// Params:
//
//	c Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
func NewADXL345Driver(c Connector, options ...func(Config)) *ADXL345Driver {
	d := &ADXL345Driver{
		Driver: NewDriver(c, "ADXL345", adxl345DefaultAddress),
		powerCtl: adxl345PowerCtl{
			measure: 1,
		},
		dataFormat: adxl345DataFormat{
			fullScaleRange: ADXL345FsRange_2G,
		},
		bwRate: adxl345BwRate{
			lowPower: true,
			rate:     ADXL345Rate_100HZ,
		},
	}
	d.afterStart = d.initialize
	d.beforeHalt = d.shutdown

	for _, option := range options {
		option(d)
	}

	// TODO: add commands for API
	return d
}

// WithADXL345LowPowerMode option modifies the low power mode.
func WithADXL345LowPowerMode(val bool) func(Config) {
	return func(c Config) {
		if d, ok := c.(*ADXL345Driver); ok {
			d.bwRate.lowPower = val
		} else if adxl345Debug {
			log.Printf("Trying to modify low power mode for non-ADXL345Driver %v", c)
		}
	}
}

// WithADXL345DataOutputRate option sets the data output rate.
// Valid settings are of type "ADXL345RateConfig"
func WithADXL345DataOutputRate(val ADXL345RateConfig) func(Config) {
	return func(c Config) {
		if d, ok := c.(*ADXL345Driver); ok {
			d.bwRate.rate = val
		} else if adxl345Debug {
			log.Printf("Trying to set data output rate for non-ADXL345Driver %v", c)
		}
	}
}

// WithADXL345FullScaleRange option sets the full scale range.
// Valid settings are of type "ADXL345FsRangeConfig"
func WithADXL345FullScaleRange(val ADXL345FsRangeConfig) func(Config) {
	return func(c Config) {
		if d, ok := c.(*ADXL345Driver); ok {
			d.dataFormat.fullScaleRange = val
		} else if adxl345Debug {
			log.Printf("Trying to set full scale range for non-ADXL345Driver %v", c)
		}
	}
}

// UseLowPower change the current rate of the sensor
func (d *ADXL345Driver) UseLowPower(lowPower bool) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.bwRate.lowPower = lowPower
	return d.connection.WriteByteData(adxl345Reg_BW_RATE, d.bwRate.toByte())
}

// SetRate change the current rate of the sensor immediately
func (d *ADXL345Driver) SetRate(rate ADXL345RateConfig) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.bwRate.rate = rate
	return d.connection.WriteByteData(adxl345Reg_BW_RATE, d.bwRate.toByte())
}

// SetRange change the current range of the sensor immediately
func (d *ADXL345Driver) SetRange(fullScaleRange ADXL345FsRangeConfig) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.dataFormat.fullScaleRange = fullScaleRange
	return d.connection.WriteByteData(adxl345Reg_DATA_FORMAT, d.dataFormat.toByte())
}

// XYZ returns the adjusted x, y and z axis, unit [g]
func (d *ADXL345Driver) XYZ() (float64, float64, float64, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	xr, yr, zr, err := d.readRawData()
	if err != nil {
		return 0, 0, 0, err
	}

	return d.dataFormat.convertToG(xr), d.dataFormat.convertToG(yr), d.dataFormat.convertToG(zr), nil
}

// RawXYZ returns the raw x,y and z axis
func (d *ADXL345Driver) RawXYZ() (int16, int16, int16, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.readRawData()
}

func (d *ADXL345Driver) readRawData() (int16, int16, int16, error) {
	buf := []byte{0, 0, 0, 0, 0, 0}
	if err := d.connection.ReadBlockData(adxl345Reg_DATAX0, buf); err != nil {
		return 0, 0, 0, err
	}

	rx := int16(binary.LittleEndian.Uint16(buf[0:2])) //nolint:gosec // TODO: fix later
	ry := int16(binary.LittleEndian.Uint16(buf[2:4])) //nolint:gosec // TODO: fix later
	rz := int16(binary.LittleEndian.Uint16(buf[4:6])) //nolint:gosec // TODO: fix later
	return rx, ry, rz, nil
}

func (d *ADXL345Driver) initialize() error {
	if err := d.connection.WriteByteData(adxl345Reg_BW_RATE, d.bwRate.toByte()); err != nil {
		return err
	}
	if err := d.connection.WriteByteData(adxl345Reg_POWER_CTL, d.powerCtl.toByte()); err != nil {
		return err
	}
	if err := d.connection.WriteByteData(adxl345Reg_DATA_FORMAT, d.dataFormat.toByte()); err != nil {
		return err
	}

	return nil
}

func (d *ADXL345Driver) shutdown() error {
	d.powerCtl.measure = 0
	if d.connection == nil {
		return fmt.Errorf("connection not available")
	}
	return d.connection.WriteByteData(adxl345Reg_POWER_CTL, d.powerCtl.toByte())
}

// convertToG converts the given raw value by range configuration to the unit [g]
func (d *adxl345DataFormat) convertToG(rawValue int16) float64 {
	switch d.fullScaleRange {
	case ADXL345FsRange_2G:
		return float64(rawValue) * 2 / 512
	case ADXL345FsRange_4G:
		return float64(rawValue) * 4 / 512
	case ADXL345FsRange_8G:
		return float64(rawValue) * 8 / 512
	case ADXL345FsRange_16G:
		return float64(rawValue) * 16 / 512
	default:
		return 0
	}
}

// toByte returns a byte from the powerCtl configuration
func (p *adxl345PowerCtl) toByte() uint8 {
	bits := p.wakeUp
	bits = bits | (p.sleep << 2)
	bits = bits | (p.measure << 3)
	bits = bits | (p.autoSleep << 4)
	return bits | (p.link << 5)
}

// toByte returns a byte from the dataFormat configuration
func (d *adxl345DataFormat) toByte() uint8 {
	bits := uint8(d.fullScaleRange)
	bits = bits | (d.justify << 2)
	bits = bits | (d.fullRes << 3)
	bits = bits | (d.intInvert << 5)
	bits = bits | (d.spi << 6)
	return bits | (d.selfTest << 7)
}

// toByte returns a byte from the bwRate configuration
func (b *adxl345BwRate) toByte() uint8 {
	bits := uint8(b.rate)
	if b.lowPower {
		bits = bits | adxl345Rate_LowPowerBit
	}
	return bits
}
