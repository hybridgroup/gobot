package i2c

import (
	"bytes"
	"encoding/binary"

	"gobot.io/x/gobot"
)

const l3gd20hAddress = 0x6B

// Control Register 1
const l3gd20hRegisterCtl1 = 0x20
const l3gd20hNormalMode = 0x8
const l3gd20hEnableZ = 0x04
const l3gd20hEnableY = 0x02
const l3gd20hEnableX = 0x01

// Control Register 4
const l3gd20hRegisterCtl4 = 0x23

const l3gd20hRegisterOutXLSB = 0x28 | 0x80 // set auto-increment bit.

const (
	// L3GD20HScale250dps is the 250 degress-per-second scale.
	L3GD20HScale250dps L3GD20HScale = 0x00
	// L3GD20HScale500dps is the 500 degress-per-second scale.
	L3GD20HScale500dps L3GD20HScale = 0x10
	// L3GD20HScale2000dps is the 2000 degress-per-second scale.
	L3GD20HScale2000dps L3GD20HScale = 0x30
)

// L3GD20HScale is the scale sensitivity of degrees-per-second.
type L3GD20HScale byte

// L3GD20HDriver is the gobot driver for the Adafruit Triple-Axis Gyroscope L3GD20H.
// Device datasheet: http://www.st.com/internet/com/TECHNICAL_RESOURCES/TECHNICAL_LITERATURE/DATASHEET/DM00036465.pdf
type L3GD20HDriver struct {
	name       string
	connector  Connector
	connection Connection
	Config
	scale L3GD20HScale
}

// NewL3GD20HDriver creates a new Gobot driver for the
// L3GD20H I2C Triple-Axis Gyroscope.
//
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewL3GD20HDriver(c Connector, options ...func(Config)) *L3GD20HDriver {
	l := &L3GD20HDriver{
		name:      gobot.DefaultName("L3GD20H"),
		connector: c,
		Config:    NewConfig(),
		scale:     L3GD20HScale250dps,
	}

	for _, option := range options {
		option(l)
	}

	// TODO: add commands to API
	return l
}

// Name returns the name of the device.
func (d *L3GD20HDriver) Name() string {
	return d.name
}

// SetName sets the name of the device.
func (d *L3GD20HDriver) SetName(name string) {
	d.name = name
}

// Connection returns the connection of the device.
func (d *L3GD20HDriver) Connection() gobot.Connection {
	return d.connector.(gobot.Connection)
}

// Scale returns the scale sensitivity of the device.
func (d *L3GD20HDriver) Scale() L3GD20HScale {
	return d.scale
}

// SetScale sets the scale sensitivity of the device.
func (d *L3GD20HDriver) SetScale(s L3GD20HScale) {
	d.scale = s
}

// Start initializes the device.
func (d *L3GD20HDriver) Start() (err error) {
	if err := d.initialization(); err != nil {
		return err
	}
	return nil
}

// Halt halts the device.
func (d *L3GD20HDriver) Halt() (err error) {
	return nil
}

// XYZ returns the current change in degrees per second, for the 3 axis.
func (d *L3GD20HDriver) XYZ() (x float32, y float32, z float32, err error) {
	if _, err = d.connection.Write([]byte{l3gd20hRegisterOutXLSB}); err != nil {
		return 0, 0, 0, err
	}
	measurements := make([]byte, 6)
	if _, err = d.connection.Read(measurements); err != nil {
		return 0, 0, 0, err
	}

	var rawX int16
	var rawY int16
	var rawZ int16
	buf := bytes.NewBuffer(measurements)
	binary.Read(buf, binary.LittleEndian, &rawX)
	binary.Read(buf, binary.LittleEndian, &rawY)
	binary.Read(buf, binary.LittleEndian, &rawZ)

	// Sensitivity values from the mechanical characteristics in the datasheet.
	sensitivity := d.getSensitivity()

	return float32(rawX) * sensitivity, float32(rawY) * sensitivity, float32(rawZ) * sensitivity, nil
}

func (d *L3GD20HDriver) initialization() (err error) {
	bus := d.GetBusOrDefault(d.connector.GetDefaultBus())
	address := d.GetAddressOrDefault(l3gd20hAddress)

	d.connection, err = d.connector.GetConnection(address, bus)
	if err != nil {
		return err
	}
	// reset the gyroscope.
	if _, err := d.connection.Write([]byte{l3gd20hRegisterCtl1, 0x00}); err != nil {
		return err
	}
	// Enable Z, Y and X axis.
	if _, err := d.connection.Write([]byte{l3gd20hRegisterCtl1, l3gd20hNormalMode | l3gd20hEnableZ | l3gd20hEnableY | l3gd20hEnableX}); err != nil {
		return err
	}
	// Set the sensitivity scale.
	if _, err := d.connection.Write([]byte{l3gd20hRegisterCtl4, byte(d.scale)}); err != nil {
		return err
	}
	return nil
}

func (d *L3GD20HDriver) getSensitivity() float32 {
	switch d.scale {
	case L3GD20HScale250dps:
		return 0.00875
	case L3GD20HScale500dps:
		return 0.0175
	case L3GD20HScale2000dps:
		return 0.07
	}
	return 0
}
