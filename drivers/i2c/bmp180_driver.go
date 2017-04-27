package i2c

import (
	"bytes"
	"encoding/binary"
	"time"

	"gobot.io/x/gobot"
)

const bmp180Address = 0x77

const bmp180RegisterAC1MSB = 0xAA

const bmp180RegisterCtl = 0xF4
const bmp180CmdTemp = 0x2E
const bmp180RegisterTempMSB = 0xF6
const bmp180CmdPressure = 0x34
const bmp180RegisterPressureMSB = 0xF6

const (
	// BMP180UltraLowPower is the lowest oversampling mode of the pressure measurement.
	BMP180UltraLowPower BMP180OversamplingMode = iota
	// BMP180Standard is the standard oversampling mode of the pressure measurement.
	BMP180Standard
	// BMP180HighResolution is a high oversampling mode of the pressure measurement.
	BMP180HighResolution
	// BMP180UltraHighResolution is the highest oversampling mode of the pressure measurement.
	BMP180UltraHighResolution
)

// BMP180OversamplingMode is the oversampling ratio of the pressure measurement.
type BMP180OversamplingMode uint

type calibrationCoefficients struct {
	ac1 int16
	ac2 int16
	ac3 int16
	ac4 uint16
	ac5 uint16
	ac6 uint16
	b1  int16
	b2  int16
	mb  int16
	mc  int16
	md  int16
}

// BMP180Driver is the gobot driver for the Bosch pressure sensor BMP180.
// Device datasheet: https://cdn-shop.adafruit.com/datasheets/BST-BMP180-DS000-09.pdf
type BMP180Driver struct {
	name       string
	Mode       BMP180OversamplingMode
	connector  Connector
	connection Connection
	Config
	calibrationCoefficients *calibrationCoefficients
}

// NewBMP180Driver creates a new driver with the i2c interface for the BMP180 device.
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewBMP180Driver(c Connector, options ...func(Config)) *BMP180Driver {
	b := &BMP180Driver{
		name:                    gobot.DefaultName("BMP180"),
		connector:               c,
		Mode:                    BMP180UltraLowPower,
		Config:                  NewConfig(),
		calibrationCoefficients: &calibrationCoefficients{},
	}

	for _, option := range options {
		option(b)
	}

	// TODO: expose commands to API
	return b
}

// Name returns the name of the device.
func (d *BMP180Driver) Name() string {
	return d.name
}

// SetName sets the name of the device.
func (d *BMP180Driver) SetName(n string) {
	d.name = n
}

// Connection returns the connection of the device.
func (d *BMP180Driver) Connection() gobot.Connection {
	return d.connector.(gobot.Connection)
}

// Start initializes the BMP180 and loads the calibration coefficients.
func (d *BMP180Driver) Start() (err error) {
	bus := d.GetBusOrDefault(d.connector.GetDefaultBus())
	address := d.GetAddressOrDefault(bmp180Address)

	if d.connection, err = d.connector.GetConnection(address, bus); err != nil {
		return err
	}
	if err := d.initialization(); err != nil {
		return err
	}
	return nil
}

func (d *BMP180Driver) initialization() (err error) {
	var coefficients []byte
	// read the 11 calibration coefficients.
	if coefficients, err = d.read(bmp180RegisterAC1MSB, 22); err != nil {
		return err
	}
	buf := bytes.NewBuffer(coefficients)
	binary.Read(buf, binary.BigEndian, &d.calibrationCoefficients.ac1)
	binary.Read(buf, binary.BigEndian, &d.calibrationCoefficients.ac2)
	binary.Read(buf, binary.BigEndian, &d.calibrationCoefficients.ac3)
	binary.Read(buf, binary.BigEndian, &d.calibrationCoefficients.ac4)
	binary.Read(buf, binary.BigEndian, &d.calibrationCoefficients.ac5)
	binary.Read(buf, binary.BigEndian, &d.calibrationCoefficients.ac6)
	binary.Read(buf, binary.BigEndian, &d.calibrationCoefficients.b1)
	binary.Read(buf, binary.BigEndian, &d.calibrationCoefficients.b2)
	binary.Read(buf, binary.BigEndian, &d.calibrationCoefficients.mb)
	binary.Read(buf, binary.BigEndian, &d.calibrationCoefficients.mc)
	binary.Read(buf, binary.BigEndian, &d.calibrationCoefficients.md)

	return nil
}

// Halt halts the device.
func (d *BMP180Driver) Halt() (err error) {
	return nil
}

// Temperature returns the current temperature, in celsius degrees.
func (d *BMP180Driver) Temperature() (temp float32, err error) {
	var rawTemp int16
	if rawTemp, err = d.rawTemp(); err != nil {
		return 0, err
	}
	return d.calculateTemp(rawTemp), nil
}

// Pressure returns the current pressure, in pascals.
func (d *BMP180Driver) Pressure() (pressure float32, err error) {
	var rawTemp int16
	var rawPressure int32
	if rawTemp, err = d.rawTemp(); err != nil {
		return 0, err
	}
	if rawPressure, err = d.rawPressure(d.Mode); err != nil {
		return 0, err
	}
	return d.calculatePressure(rawTemp, rawPressure, d.Mode), nil
}

func (d *BMP180Driver) rawTemp() (int16, error) {
	if _, err := d.connection.Write([]byte{bmp180RegisterCtl, bmp180CmdTemp}); err != nil {
		return 0, err
	}
	time.Sleep(5 * time.Millisecond)
	ret, err := d.read(bmp180RegisterTempMSB, 2)
	if err != nil {
		return 0, err
	}
	buf := bytes.NewBuffer(ret)
	var rawTemp int16
	binary.Read(buf, binary.BigEndian, &rawTemp)
	return rawTemp, nil
}

func (d *BMP180Driver) read(address byte, n int) ([]byte, error) {
	if _, err := d.connection.Write([]byte{address}); err != nil {
		return nil, err
	}
	buf := make([]byte, n)
	bytesRead, err := d.connection.Read(buf)
	if bytesRead != n || err != nil {
		return nil, err
	}
	return buf, nil
}

func (d *BMP180Driver) calculateTemp(rawTemp int16) float32 {
	b5 := d.calculateB5(rawTemp)
	t := (b5 + 8) >> 4
	return float32(t) / 10
}

func (d *BMP180Driver) calculateB5(rawTemp int16) int32 {
	x1 := (int32(rawTemp) - int32(d.calibrationCoefficients.ac6)) * int32(d.calibrationCoefficients.ac5) >> 15
	x2 := int32(d.calibrationCoefficients.mc) << 11 / (x1 + int32(d.calibrationCoefficients.md))
	return x1 + x2
}

func (d *BMP180Driver) rawPressure(mode BMP180OversamplingMode) (rawPressure int32, err error) {
	if _, err = d.connection.Write([]byte{bmp180RegisterCtl, bmp180CmdPressure + byte(mode<<6)}); err != nil {
		return 0, err
	}
	time.Sleep(pauseForReading(mode))
	var ret []byte
	if ret, err = d.read(bmp180RegisterPressureMSB, 3); err != nil {
		return 0, err
	}
	rawPressure = (int32(ret[0])<<16 + int32(ret[1])<<8 + int32(ret[2])) >> (8 - uint(mode))
	return rawPressure, nil
}

func (d *BMP180Driver) calculatePressure(rawTemp int16, rawPressure int32, mode BMP180OversamplingMode) float32 {
	b5 := d.calculateB5(rawTemp)
	b6 := b5 - 4000
	x1 := (int32(d.calibrationCoefficients.b2) * (b6 * b6 >> 12)) >> 11
	x2 := (int32(d.calibrationCoefficients.ac2) * b6) >> 11
	x3 := x1 + x2
	b3 := (((int32(d.calibrationCoefficients.ac1)*4 + x3) << uint(mode)) + 2) >> 2
	x1 = (int32(d.calibrationCoefficients.ac3) * b6) >> 13
	x2 = (int32(d.calibrationCoefficients.b1) * ((b6 * b6) >> 12)) >> 16
	x3 = ((x1 + x2) + 2) >> 2
	b4 := (uint32(d.calibrationCoefficients.ac4) * uint32(x3+32768)) >> 15
	b7 := (uint32(rawPressure-b3) * (50000 >> uint(mode)))
	var p int32
	if b7 < 0x80000000 {
		p = int32((b7 << 1) / b4)
	} else {
		p = int32((b7 / b4) << 1)
	}
	x1 = (p >> 8) * (p >> 8)
	x1 = (x1 * 3038) >> 16
	x2 = (-7357 * p) >> 16
	return float32(p + ((x1 + x2 + 3791) >> 4))
}

func pauseForReading(mode BMP180OversamplingMode) time.Duration {
	var d time.Duration
	switch mode {
	case BMP180UltraLowPower:
		d = 5 * time.Millisecond
	case BMP180Standard:
		d = 8 * time.Millisecond
	case BMP180HighResolution:
		d = 14 * time.Millisecond
	case BMP180UltraHighResolution:
		d = 26 * time.Millisecond
	}
	return d
}
