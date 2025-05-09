package i2c

import (
	"bytes"
	"encoding/binary"
	"log"
	"time"
)

const bmp180Debug = false

// the default address is applicable for SDO to VDD, for SDO to GND it will be 0x76
const bmp180DefaultAddress = 0x77

const (
	bmp180RegisterAC1MSB  = 0xAA // 11 x 16 bit calibration data (AC1..AC6, B1, B2, MB, MC, MD)
	bmp180RegisterCtl     = 0xF4 // control the value to read
	bmp180RegisterDataMSB = 0xF6 // 16 bit data (temperature or pressure)

	bmp180CtlTemp     = 0x2E
	bmp180CtlPressure = 0x34
)

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

type bmp180CalibrationCoefficients struct {
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

// BMP180Driver is the gobot driver for the Bosch pressure and temperature sensor BMP180.
// Device datasheet: https://cdn-shop.adafruit.com/datasheets/BST-BMP180-DS000-09.pdf
type BMP180Driver struct {
	*Driver
	oversampling BMP180OversamplingMode
	calCoeffs    *bmp180CalibrationCoefficients
}

// NewBMP180Driver creates a new driver with the i2c interface for the BMP180 device.
// Params:
//
//	conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
func NewBMP180Driver(c Connector, options ...func(Config)) *BMP180Driver {
	d := &BMP180Driver{
		Driver:       NewDriver(c, "BMP180", bmp180DefaultAddress),
		oversampling: BMP180UltraLowPower,
		calCoeffs:    &bmp180CalibrationCoefficients{},
	}
	d.afterStart = d.initialization

	for _, option := range options {
		option(d)
	}

	// TODO: expose commands to API
	return d
}

// WithBMP180OversamplingMode option sets oversampling mode.
// Valid settings are of type "BMP180OversamplingMode"
func WithBMP180OversamplingMode(val BMP180OversamplingMode) func(Config) {
	return func(c Config) {
		if d, ok := c.(*BMP180Driver); ok {
			d.oversampling = val
		} else if bmp180Debug {
			log.Printf("Trying to set oversampling mode for non-BMP180Driver %v", c)
		}
	}
}

// Temperature returns the current temperature, in celsius degrees.
func (d *BMP180Driver) Temperature() (float32, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	rawTemp, err := d.rawTemp()
	if err != nil {
		return 0, err
	}
	return d.calculateTemp(rawTemp), nil
}

// Pressure returns the current pressure, in pascals.
func (d *BMP180Driver) Pressure() (float32, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	rawTemp, err := d.rawTemp()
	if err != nil {
		return 0, err
	}
	rawPressure, err := d.rawPressure(d.oversampling)
	if err != nil {
		return 0, err
	}
	return d.calculatePressure(rawTemp, rawPressure, d.oversampling), nil
}

func (d *BMP180Driver) initialization() error {
	// read the 11 calibration coefficients.
	coefficients := make([]byte, 22)
	if err := d.connection.ReadBlockData(bmp180RegisterAC1MSB, coefficients); err != nil {
		return err
	}
	buf := bytes.NewBuffer(coefficients)
	if err := binary.Read(buf, binary.BigEndian, &d.calCoeffs.ac1); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &d.calCoeffs.ac2); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &d.calCoeffs.ac3); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &d.calCoeffs.ac4); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &d.calCoeffs.ac5); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &d.calCoeffs.ac6); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &d.calCoeffs.b1); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &d.calCoeffs.b2); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &d.calCoeffs.mb); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &d.calCoeffs.mc); err != nil {
		return err
	}
	return binary.Read(buf, binary.BigEndian, &d.calCoeffs.md)
}

func (d *BMP180Driver) rawTemp() (int16, error) {
	if _, err := d.connection.Write([]byte{bmp180RegisterCtl, bmp180CtlTemp}); err != nil {
		return 0, err
	}
	time.Sleep(5 * time.Millisecond)
	ret := make([]byte, 2)
	err := d.connection.ReadBlockData(bmp180RegisterDataMSB, ret)
	if err != nil {
		return 0, err
	}
	buf := bytes.NewBuffer(ret)
	var rawTemp int16
	if err := binary.Read(buf, binary.BigEndian, &rawTemp); err != nil {
		return 0, err
	}
	return rawTemp, nil
}

func (d *BMP180Driver) calculateTemp(rawTemp int16) float32 {
	b5 := d.calculateB5(rawTemp)
	t := (b5 + 8) >> 4
	return float32(t) / 10
}

func (d *BMP180Driver) calculateB5(rawTemp int16) int32 {
	x1 := (int32(rawTemp) - int32(d.calCoeffs.ac6)) * int32(d.calCoeffs.ac5) >> 15
	x2 := int32(d.calCoeffs.mc) << 11 / (x1 + int32(d.calCoeffs.md))
	return x1 + x2
}

func (d *BMP180Driver) rawPressure(oversampling BMP180OversamplingMode) (int32, error) {
	if _, err := d.connection.Write([]byte{bmp180RegisterCtl, bmp180CtlPressure + byte(oversampling<<6)}); err != nil {
		return 0, err
	}
	time.Sleep(bmp180PauseForReading(oversampling))
	ret := make([]byte, 3)
	if err := d.connection.ReadBlockData(bmp180RegisterDataMSB, ret); err != nil {
		return 0, err
	}
	rawPressure := (int32(ret[0])<<16 + int32(ret[1])<<8 + int32(ret[2])) >> (8 - uint(oversampling))
	return rawPressure, nil
}

func (d *BMP180Driver) calculatePressure(
	rawTemp int16,
	rawPressure int32,
	oversampling BMP180OversamplingMode,
) float32 {
	b5 := d.calculateB5(rawTemp)
	b6 := b5 - 4000
	x1 := (int32(d.calCoeffs.b2) * (b6 * b6 >> 12)) >> 11
	x2 := (int32(d.calCoeffs.ac2) * b6) >> 11
	x3 := x1 + x2
	b3 := (((int32(d.calCoeffs.ac1)*4 + x3) << uint(oversampling)) + 2) >> 2
	x1 = (int32(d.calCoeffs.ac3) * b6) >> 13
	x2 = (int32(d.calCoeffs.b1) * ((b6 * b6) >> 12)) >> 16
	x3 = ((x1 + x2) + 2) >> 2
	b4 := (uint32(d.calCoeffs.ac4) * uint32(x3+32768)) >> 15       //nolint:gosec // TODO: fix later
	b7 := (uint32(rawPressure-b3) * (50000 >> uint(oversampling))) //nolint:gosec // TODO: fix later
	var p int32
	if b7 < 0x80000000 {
		p = int32((b7 << 1) / b4) //nolint:gosec // TODO: fix later
	} else {
		p = int32((b7 / b4) << 1) //nolint:gosec // TODO: fix later
	}
	x1 = (p >> 8) * (p >> 8)
	x1 = (x1 * 3038) >> 16
	x2 = (-7357 * p) >> 16
	return float32(p + ((x1 + x2 + 3791) >> 4))
}

func bmp180PauseForReading(oversampling BMP180OversamplingMode) time.Duration {
	var d time.Duration
	switch oversampling {
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
