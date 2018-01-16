package i2c

import (
	"encoding/binary"

	"github.com/pkg/errors"
	"gobot.io/x/gobot"
)

const ADXL345AddressLow = 0x53
const ADXL345AddressHigh = 0x1D

const (
	// Data rate
	ADXL345_RATE_3200HZ = 0x0F // 3200 Hz
	ADXL345_RATE_1600HZ = 0x0E // 1600 Hz
	ADXL345_RATE_800HZ  = 0x0D // 800 Hz
	ADXL345_RATE_400HZ  = 0x0C // 400 Hz
	ADXL345_RATE_200HZ  = 0x0B // 200 Hz
	ADXL345_RATE_100HZ  = 0x0A // 100 Hz
	ADXL345_RATE_50HZ   = 0x09 // 50 Hz
	ADXL345_RATE_25HZ   = 0x08 // 25 Hz
	ADXL345_RATE_12_5HZ = 0x07 // 12.5 Hz
	ADXL345_RATE_6_25HZ = 0x06 // 6.25 Hz
	ADXL345_RATE_3_13HZ = 0x05 // 3.13 Hz
	ADXL345_RATE_1_56HZ = 0x04 // 1.56 Hz
	ADXL345_RATE_0_78HZ = 0x03 // 0.78 Hz
	ADXL345_RATE_0_39HZ = 0x02 // 0.39 Hz
	ADXL345_RATE_0_20HZ = 0x01 // 0.20 Hz
	ADXL345_RATE_0_10HZ = 0x00 // 0.10 Hz

	// Data range
	ADXL345_RANGE_2G  = 0x00 // +-2 g
	ADXL345_RANGE_4G  = 0x01 // +-4 g
	ADXL345_RANGE_8G  = 0x02 // +-8 g
	ADXL345_RANGE_16G = 0x03 // +-16 g)

	ADXL345_REG_DEVID          = 0x00 // R,     11100101,   Device ID
	ADXL345_REG_THRESH_TAP     = 0x1D // R/W,   00000000,   Tap threshold
	ADXL345_REG_OFSX           = 0x1E // R/W,   00000000,   X-axis offset
	ADXL345_REG_OFSY           = 0x1F // R/W,   00000000,   Y-axis offset
	ADXL345_REG_OFSZ           = 0x20 // R/W,   00000000,   Z-axis offset
	ADXL345_REG_DUR            = 0x21 // R/W,   00000000,   Tap duration
	ADXL345_REG_LATENT         = 0x22 // R/W,   00000000,   Tap latency
	ADXL345_REG_WINDOW         = 0x23 // R/W,   00000000,   Tap window
	ADXL345_REG_THRESH_ACT     = 0x24 // R/W,   00000000,   Activity threshold
	ADXL345_REG_THRESH_INACT   = 0x25 // R/W,   00000000,   Inactivity threshold
	ADXL345_REG_TIME_INACT     = 0x26 // R/W,   00000000,   Inactivity time
	ADXL345_REG_ACT_INACT_CTL  = 0x27 // R/W,   00000000,   Axis enable control for activity and inactiv ity detection
	ADXL345_REG_THRESH_FF      = 0x28 // R/W,   00000000,   Free-fall threshold
	ADXL345_REG_TIME_FF        = 0x29 // R/W,   00000000,   Free-fall time
	ADXL345_REG_TAP_AXES       = 0x2A // R/W,   00000000,   Axis control for single tap/double tap
	ADXL345_REG_ACT_TAP_STATUS = 0x2B // R,     00000000,   Source of single tap/double tap
	ADXL345_REG_BW_RATE        = 0x2C // R/W,   00001010,   Data rate and power mode control
	ADXL345_REG_POWER_CTL      = 0x2D // R/W,   00000000,   Power-saving features control
	ADXL345_REG_INT_ENABLE     = 0x2E // R/W,   00000000,   Interrupt enable control
	ADXL345_REG_INT_MAP        = 0x2F // R/W,   00000000,   Interrupt mapping control
	ADXL345_REG_INT_SOUCE      = 0x30 // R,     00000010,   Source of interrupts
	ADXL345_REG_DATA_FORMAT    = 0x31 // R/W,   00000000,   Data format control
	ADXL345_REG_DATAX0         = 0x32 // R,     00000000,   X-Axis Data 0
	ADXL345_REG_DATAX1         = 0x33 // R,     00000000,   X-Axis Data 1
	ADXL345_REG_DATAY0         = 0x34 // R,     00000000,   Y-Axis Data 0
	ADXL345_REG_DATAY1         = 0x35 // R,     00000000,   Y-Axis Data 1
	ADXL345_REG_DATAZ0         = 0x36 // R,     00000000,   Z-Axis Data 0
	ADXL345_REG_DATAZ1         = 0x37 // R,     00000000,   Z-Axis Data 1
	ADXL345_REG_FIFO_CTL       = 0x38 // R/W,   00000000,   FIFO control
	ADXL345_REG_FIFO_STATUS    = 0x39 // R,     00000000,   FIFO status
)

type ADXL345Driver struct {
	name       string
	connector  Connector
	connection Connection

	powerCtl   adxl345PowerCtl
	dataFormat adxl345DataFormat
	bwRate     adxl345BwRate

	x, y, z          float64
	rawX, rawY, rawZ int16

	Config
}

type adxl345PowerCtl struct {
	Link      uint8
	AutoSleep uint8
	Measure   uint8
	Sleep     uint8
	WakeUp    uint8
}

type adxl345DataFormat struct {
	SelfTest  uint8
	SPI       uint8
	IntInvert uint8
	FullRes   uint8
	Justify   uint8
	Range     uint8
}

type adxl345BwRate struct {
	LowPower uint8
	Rate     uint8
}

// NewADXL345Driver creates a new driver with specified i2c interface
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewADXL345Driver(a Connector, options ...func(Config)) *ADXL345Driver {
	m := &ADXL345Driver{
		name:      gobot.DefaultName("ADXL345"),
		connector: a,
		powerCtl: adxl345PowerCtl{
			Measure: 1,
		},
		dataFormat: adxl345DataFormat{
			Range: ADXL345_RANGE_2G,
		},
		bwRate: adxl345BwRate{
			LowPower: 1,
			Rate:     ADXL345_RATE_100HZ,
		},
		Config: NewConfig(),
	}

	for _, option := range options {
		option(m)
	}

	// TODO: add commands for API
	return m
}

// Name returns the Name for the Driver
func (h *ADXL345Driver) Name() string { return h.name }

// SetName sets the Name for the Driver
func (h *ADXL345Driver) SetName(n string) { h.name = n }

// Connection returns the connection for the Driver
func (h *ADXL345Driver) Connection() gobot.Connection { return h.connector.(gobot.Connection) }

// Start initialized the adxl345
func (h *ADXL345Driver) Start() (err error) {
	bus := h.GetBusOrDefault(h.connector.GetDefaultBus())
	address := h.GetAddressOrDefault(ADXL345AddressLow)

	h.connection, err = h.connector.GetConnection(address, bus)
	if err != nil {
		return err
	}

	if _, err := h.connection.Write([]byte{ADXL345_REG_BW_RATE, h.bwRate.toByte()}); err != nil {
		return err
	}

	if _, err := h.connection.Write([]byte{ADXL345_REG_POWER_CTL, h.powerCtl.toByte()}); err != nil {
		return err
	}

	if _, err := h.connection.Write([]byte{ADXL345_REG_DATA_FORMAT, h.dataFormat.toByte()}); err != nil {
		return err
	}

	return
}

// Stop adxl345
func (h *ADXL345Driver) Stop() (err error) {
	if _, err := h.connection.Write([]byte{ADXL345_REG_POWER_CTL, 0}); err != nil {
		return err
	}

	return
}

// Halt returns true if devices is halted successfully
func (h *ADXL345Driver) Halt() (err error) {
	h.Stop()
	return
}

// XYZ returns the adjusted x, y and z axis from the adxl345
func (h *ADXL345Driver) XYZ() (x float64, y float64, z float64) {
	h.update()
	return h.x, h.y, h.z
}

// XYZ returns the raw x,y and z axis from the adxl345
func (h *ADXL345Driver) RawXYZ() (x int16, y int16, z int16) {
	h.update()
	return h.rawX, h.rawY, h.rawZ
}

// update the cached values for the axis to avoid errors if the connection is not available (polling too frequently)
func (h *ADXL345Driver) update() (err error) {

	if h.connection == nil {
		return errors.New("connection not available")
	}

	h.connection.Write([]byte{ADXL345_REG_DATAX0})
	buf := []byte{0, 0, 0, 0, 0, 0}

	_, err = h.connection.Read(buf)
	if err != nil {
		return
	}

	h.rawX = int16(binary.LittleEndian.Uint16(buf[0:2]))
	h.rawY = int16(binary.LittleEndian.Uint16(buf[2:4]))
	h.rawZ = int16(binary.LittleEndian.Uint16(buf[4:6]))

	h.x = h.dataFormat.ConvertToSI(h.rawX)
	h.y = h.dataFormat.ConvertToSI(h.rawY)
	h.z = h.dataFormat.ConvertToSI(h.rawZ)

	return
}

// ConvertToSI adjusts the raw values from the adxl345 with the range configuration
func (d *adxl345DataFormat) ConvertToSI(rawValue int16) float64 {
	switch d.Range {
	case ADXL345_RANGE_2G:
		return float64(rawValue) * 2 / 512
	case ADXL345_RANGE_4G:
		return float64(rawValue) * 4 / 512
	case ADXL345_RANGE_8G:
		return float64(rawValue) * 8 / 512
	case ADXL345_RANGE_16G:
		return float64(rawValue) * 16 / 512
	default:
		return 0
	}
}

// toByte returns a byte from the powerCtl configuration
func (p *adxl345PowerCtl) toByte() (bits uint8) {
	bits = 0x00
	bits = bits | (p.Link << 5)
	bits = bits | (p.AutoSleep << 4)
	bits = bits | (p.Measure << 3)
	bits = bits | (p.Sleep << 2)
	bits = bits | p.WakeUp

	return bits
}

// toByte returns a byte from the dataFormat configuration
func (d *adxl345DataFormat) toByte() (bits uint8) {
	bits = 0x00
	bits = bits | (d.SelfTest << 7)
	bits = bits | (d.SPI << 6)
	bits = bits | (d.IntInvert << 5)
	bits = bits | (d.FullRes << 3)
	bits = bits | (d.Justify << 2)
	bits = bits | d.Range

	return bits
}

// toByte returns a byte from the bwRate configuration
func (b *adxl345BwRate) toByte() (bits uint8) {
	bits = 0x00
	bits = bits | (b.LowPower << 4)
	bits = bits | b.Rate

	return bits
}
