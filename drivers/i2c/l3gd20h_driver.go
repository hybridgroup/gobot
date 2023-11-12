//nolint:lll // ok here
package i2c

import (
	"bytes"
	"encoding/binary"
	"log"
)

const (
	l3gd20hDebug          = false
	l3gd20hDefaultAddress = 0x6B
)

const (
	l3gd20hReg_Ctl1    = 0x20 // output data rate selection, bandwidth selection, power mode, axis X/Y/Z enable
	l3gd20hReg_Ctl4    = 0x23 // block data update, big/little-endian, full scale, level sensitive latch, self test, serial interface mode
	l3gd20hReg_OutXLSB = 0x28 // X-axis angular rate data, LSB

	l3gd20hCtl1_NormalModeBit = 0x08
	l3gd20hCtl1_EnableZBit    = 0x04
	l3gd20hCtl1_EnableYBit    = 0x02
	l3gd20hCtl1_EnableXBit    = 0x01

	l3gd20hCtl4_FullScaleRangeBits = 0x30
)

// L3GD20HScale is for configurable full scale range.
type L3GD20HScale byte

const (
	// L3GD20HScale250dps is the +/-250 degrees-per-second full scale range (+/-245 from datasheet, but can hold around +/-286).
	L3GD20HScale250dps L3GD20HScale = 0x00
	// L3GD20HScale500dps is the +/-500 degrees-per-second full scale range.
	L3GD20HScale500dps L3GD20HScale = 0x10
	// L3GD20HScale2001dps is the +/-2000 degrees-per-second full scale range by using 0x20 setting.
	L3GD20HScale2001dps L3GD20HScale = 0x20
	// L3GD20HScale2000dps is the +/-2000 degrees-per-second full scale range.
	L3GD20HScale2000dps L3GD20HScale = 0x30
)

// l3gdhSensibility in °/s, see the mechanical characteristics in the datasheet
var l3gdhSensibility = map[L3GD20HScale]float32{
	L3GD20HScale250dps:  0.00875,
	L3GD20HScale500dps:  0.0175,
	L3GD20HScale2001dps: 0.07,
	L3GD20HScale2000dps: 0.07,
}

// L3GD20HDriver is the gobot driver for the Adafruit Triple-Axis Gyroscope L3GD20H.
// Device datasheet: http://www.st.com/internet/com/TECHNICAL_RESOURCES/TECHNICAL_LITERATURE/DATASHEET/DM00036465.pdf
type L3GD20HDriver struct {
	*Driver
	scale L3GD20HScale
}

// NewL3GD20HDriver creates a new Gobot driver for the
// L3GD20H I2C Triple-Axis Gyroscope.
//
// Params:
//
//	c Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
func NewL3GD20HDriver(c Connector, options ...func(Config)) *L3GD20HDriver {
	l := &L3GD20HDriver{
		Driver: NewDriver(c, "L3GD20H", l3gd20hDefaultAddress, options...),
		scale:  L3GD20HScale250dps,
	}
	l.afterStart = l.initialize

	// TODO: add commands to API
	return l
}

// WithL3GD20HFullScaleRange option sets the full scale range for the gyroscope.
// Valid settings are of type "L3GD20HScale"
func WithL3GD20HFullScaleRange(val L3GD20HScale) func(Config) {
	return func(c Config) {
		d, ok := c.(*L3GD20HDriver)
		if ok {
			d.scale = val
		} else if l3gd20hDebug {
			log.Printf("Trying to set full scale range of gyroscope for non-L3GD20HDriver %v", c)
		}
	}
}

// SetScale sets the full scale range of the device (deprecated, use WithL3GD20HFullScaleRange() instead).
func (d *L3GD20HDriver) SetScale(s L3GD20HScale) {
	d.scale = s
}

// Scale returns the full scale range (deprecated, use FullScaleRange() instead).
func (d *L3GD20HDriver) Scale() L3GD20HScale {
	return d.scale
}

// FullScaleRange returns the full scale range of the device.
func (d *L3GD20HDriver) FullScaleRange() (uint8, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	val, err := d.connection.ReadByteData(l3gd20hReg_Ctl4)
	if err != nil {
		return 0, err
	}
	return val & l3gd20hCtl4_FullScaleRangeBits, nil
}

// XYZ returns the current change in degrees per second, for the 3 axis.
//
//nolint:nonamedreturns // is sufficient here
func (d *L3GD20HDriver) XYZ() (x float32, y float32, z float32, err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	measurements := make([]byte, 6)
	reg := l3gd20hReg_OutXLSB | 0x80 // set auto-increment bit
	if err := d.connection.ReadBlockData(uint8(reg), measurements); err != nil {
		return 0, 0, 0, err
	}

	var rawX int16
	var rawY int16
	var rawZ int16
	buf := bytes.NewBuffer(measurements)
	if err := binary.Read(buf, binary.LittleEndian, &rawX); err != nil {
		return 0, 0, 0, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &rawY); err != nil {
		return 0, 0, 0, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &rawZ); err != nil {
		return 0, 0, 0, err
	}

	sensitivity := l3gdhSensibility[d.scale]

	return float32(rawX) * sensitivity, float32(rawY) * sensitivity, float32(rawZ) * sensitivity, nil
}

func (d *L3GD20HDriver) initialize() error {
	// reset the gyroscope.
	if err := d.connection.WriteByteData(l3gd20hReg_Ctl1, 0x00); err != nil {
		return err
	}
	// Enable Z, Y and X axis.
	ctl1 := l3gd20hCtl1_NormalModeBit | l3gd20hCtl1_EnableZBit | l3gd20hCtl1_EnableYBit | l3gd20hCtl1_EnableXBit
	if err := d.connection.WriteByteData(l3gd20hReg_Ctl1, uint8(ctl1)); err != nil {
		return err
	}
	// Set the sensitivity scale.
	if err := d.connection.WriteByteData(l3gd20hReg_Ctl4, byte(d.scale)); err != nil {
		return err
	}
	return nil
}
