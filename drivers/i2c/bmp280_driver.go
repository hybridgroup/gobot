package i2c

import (
	"bytes"
	"encoding/binary"
	"log"
	"math"
)

const bmp280Debug = true

// the default address is applicable for SDO to VDD, for SDO to GND it will be 0x76
// this is also true for bme280 (which using this address as well)
const bmp280DefaultAddress = 0x77

type (
	BMP280PressureOversampling    uint8
	BMP280TemperatureOversampling uint8
	BMP280IIRFilter               uint8
)

const (
	bmp280RegCalib00      = 0x88 // 12 x 16 bit calibration data (T1..T3, P1..P9)
	bmp280RegCtrl         = 0xF4 // data acquisition options (oversampling of temperature and pressure, power mode)
	bmp280RegConf         = 0xF5 // rate, IIR-filter and interface options (SPI)
	bmp280RegPressureData = 0xF7
	bmp280RegTempData     = 0xFA

	// bits 0, 1 of control register
	bmp280CtrlPwrSleepMode   = 0x00
	bmp280CtrlPwrForcedMode  = 0x01
	bmp280CtrlPwrForcedMode2 = 0x02 // same function as 0x01
	bmp280CtrlPwrNormalMode  = 0x03

	// bits 2, 3, 4 of control register (will be shifted on write)
	BMP280CtrlPressNoMeasurement  BMP280PressureOversampling = 0x00 // no measurement (value will be 0x08 0x00 0x00)
	BMP280CtrlPressOversampling1  BMP280PressureOversampling = 0x01 // resolution 16 bit
	BMP280CtrlPressOversampling2  BMP280PressureOversampling = 0x02 // resolution 17 bit
	BMP280CtrlPressOversampling4  BMP280PressureOversampling = 0x03 // resolution 18 bit
	BMP280CtrlPressOversampling8  BMP280PressureOversampling = 0x04 // resolution 19 bit
	BMP280CtrlPressOversampling16 BMP280PressureOversampling = 0x05 // resolution 20 bit (same as 0x06, 0x07)

	// bits 5, 6, 7 of control register (will be shifted on write)
	BMP280CtrlTempNoMeasurement  BMP280TemperatureOversampling = 0x00 // no measurement (value will be 0x08 0x00 0x00)
	BMP280CtrlTempOversampling1  BMP280TemperatureOversampling = 0x01 // resolution 16 bit
	BMP280CtrlTempOversampling2  BMP280TemperatureOversampling = 0x02 // resolution 17 bit
	BMP280CtrlTempOversampling4  BMP280TemperatureOversampling = 0x03 // resolution 18 bit
	BMP280CtrlTempOversampling8  BMP280TemperatureOversampling = 0x04 // resolution 19 bit
	BMP280CtrlTempOversampling16 BMP280TemperatureOversampling = 0x05 // resolution 20 bit

	// bit 0 of config register
	bmp280ConfSPIBit = 0x01 // if set, SPI is used

	// bits 2, 3, 4 of config register (bit 1 is unused, will be shifted on write)
	bmp280ConfStandBy0005 = 0x00 //	0.5 ms
	bmp280ConfStandBy0625 = 0x01 //	62.5 ms
	bmp280ConfStandBy0125 = 0x02 //	125 ms
	bmp280ConfStandBy0250 = 0x03 //	250 ms
	bmp280ConfStandBy0500 = 0x04 //	500 ms
	bmp280ConfStandBy1000 = 0x05 //	1000 ms
	bmp280ConfStandBy2000 = 0x06 //	2000 ms
	bmp280ConfStandBy4000 = 0x07 //	4000 ms

	// bits 5, 6, 7 of config register
	BMP280ConfFilterOff BMP280IIRFilter = 0x00
	BMP280ConfFilter2   BMP280IIRFilter = 0x01
	BMP280ConfFilter4   BMP280IIRFilter = 0x02
	BMP280ConfFilter8   BMP280IIRFilter = 0x03
	BMP280ConfFilter16  BMP280IIRFilter = 0x04

	bmp280SeaLevelPressure = 1013.25
)

type bmp280CalibrationCoefficients struct {
	t1 uint16
	t2 int16
	t3 int16
	p1 uint16
	p2 int16
	p3 int16
	p4 int16
	p5 int16
	p6 int16
	p7 int16
	p8 int16
	p9 int16
}

// BMP280Driver is a driver for the BMP280 temperature/pressure sensor
type BMP280Driver struct {
	*Driver
	calCoeffs         *bmp280CalibrationCoefficients
	ctrlPwrMode       uint8
	ctrlPressOversamp BMP280PressureOversampling
	ctrlTempOversamp  BMP280TemperatureOversampling
	confFilter        BMP280IIRFilter
}

// NewBMP280Driver creates a new driver with specified i2c interface.
// Params:
//
//	c Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
func NewBMP280Driver(c Connector, options ...func(Config)) *BMP280Driver {
	d := &BMP280Driver{
		Driver:            NewDriver(c, "BMP280", bmp280DefaultAddress),
		calCoeffs:         &bmp280CalibrationCoefficients{},
		ctrlPwrMode:       bmp280CtrlPwrNormalMode,
		ctrlPressOversamp: BMP280CtrlPressOversampling16,
		ctrlTempOversamp:  BMP280CtrlTempOversampling1,
		confFilter:        BMP280ConfFilterOff,
	}
	d.afterStart = d.initialization

	for _, option := range options {
		option(d)
	}

	// TODO: expose commands to API
	return d
}

// WithBMP280PressureOversampling option sets the oversampling for pressure.
// Valid settings are of type "BMP280PressureOversampling"
func WithBMP280PressureOversampling(val BMP280PressureOversampling) func(Config) {
	return func(c Config) {
		if d, ok := c.(*BMP280Driver); ok {
			d.ctrlPressOversamp = val
		} else if bmp280Debug {
			log.Printf("Trying to set pressure oversampling for non-BMP280Driver %v", c)
		}
	}
}

// WithBMP280TemperatureOversampling option sets oversampling for temperature.
// Valid settings are of type "BMP280TemperatureOversampling"
func WithBMP280TemperatureOversampling(val BMP280TemperatureOversampling) func(Config) {
	return func(c Config) {
		if d, ok := c.(*BMP280Driver); ok {
			d.ctrlTempOversamp = val
		} else if bmp280Debug {
			log.Printf("Trying to set temperature oversampling for non-BMP280Driver %v", c)
		}
	}
}

// WithBMP280IIRFilter option sets the count of IIR filter coefficients.
// Valid settings are of type "BMP280IIRFilter"
func WithBMP280IIRFilter(val BMP280IIRFilter) func(Config) {
	return func(c Config) {
		if d, ok := c.(*BMP280Driver); ok {
			d.confFilter = val
		} else if bmp280Debug {
			log.Printf("Trying to set IIR filter for non-BMP280Driver %v", c)
		}
	}
}

// Temperature returns the current temperature, in celsius degrees.
func (d *BMP280Driver) Temperature() (float32, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	rawT, err := d.rawTemp()
	if err != nil {
		return 0.0, err
	}
	temp, _ := d.calculateTemp(rawT)
	return temp, nil
}

// Pressure returns the current barometric pressure, in Pa
func (d *BMP280Driver) Pressure() (float32, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	rawT, err := d.rawTemp()
	if err != nil {
		return 0.0, err
	}

	rawP, err := d.rawPressure()
	if err != nil {
		return 0.0, err
	}
	_, tFine := d.calculateTemp(rawT)
	return d.calculatePress(rawP, tFine), nil
}

// Altitude returns the current altitude in meters based on the
// current barometric pressure and estimated pressure at sea level.
// Calculation is based on code from Adafruit BME280 library
//
//	https://github.com/adafruit/Adafruit_BME280_Library
func (d *BMP280Driver) Altitude() (float32, error) {
	atmP, err := d.Pressure()
	if err != nil {
		return 0, err
	}
	atmP /= 100.0
	alt := float32(44330.0 * (1.0 - math.Pow(float64(atmP/bmp280SeaLevelPressure), 0.1903)))

	return alt, nil
}

// initialization reads the calibration coefficients.
func (d *BMP280Driver) initialization() error {
	coefficients := make([]byte, 24)
	if err := d.connection.ReadBlockData(bmp280RegCalib00, coefficients); err != nil {
		return err
	}
	buf := bytes.NewBuffer(coefficients)
	if err := binary.Read(buf, binary.LittleEndian, &d.calCoeffs.t1); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &d.calCoeffs.t2); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &d.calCoeffs.t3); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &d.calCoeffs.p1); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &d.calCoeffs.p2); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &d.calCoeffs.p3); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &d.calCoeffs.p4); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &d.calCoeffs.p5); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &d.calCoeffs.p6); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &d.calCoeffs.p7); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &d.calCoeffs.p8); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &d.calCoeffs.p9); err != nil {
		return err
	}

	ctrlReg := d.ctrlPwrMode | uint8(d.ctrlPressOversamp)<<2 | uint8(d.ctrlTempOversamp)<<5
	if err := d.connection.WriteByteData(bmp280RegCtrl, ctrlReg); err != nil {
		return err
	}

	confReg := uint8(bmp280ConfStandBy0005)<<2 | uint8(d.confFilter)<<5
	return d.connection.WriteByteData(bmp280RegConf, confReg & ^uint8(bmp280ConfSPIBit))
}

func (d *BMP280Driver) rawTemp() (int32, error) {
	var tp0, tp1, tp2 byte

	data := make([]byte, 3)
	if err := d.connection.ReadBlockData(bmp280RegTempData, data); err != nil {
		return 0, err
	}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.LittleEndian, &tp0); err != nil {
		return 0, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &tp1); err != nil {
		return 0, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &tp2); err != nil {
		return 0, err
	}

	return ((int32(tp2) >> 4) | (int32(tp1) << 4) | (int32(tp0) << 12)), nil
}

func (d *BMP280Driver) rawPressure() (int32, error) {
	var tp0, tp1, tp2 byte

	data := make([]byte, 3)
	if err := d.connection.ReadBlockData(bmp280RegPressureData, data); err != nil {
		return 0, err
	}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.LittleEndian, &tp0); err != nil {
		return 0, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &tp1); err != nil {
		return 0, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &tp2); err != nil {
		return 0, err
	}

	return ((int32(tp2) >> 4) | (int32(tp1) << 4) | (int32(tp0) << 12)), nil
}

func (d *BMP280Driver) calculateTemp(rawTemp int32) (float32, int32) {
	tcvar1 := ((float32(rawTemp) / 16384.0) - (float32(d.calCoeffs.t1) / 1024.0)) * float32(d.calCoeffs.t2)
	tcvar2 := (((float32(rawTemp) / 131072.0) - (float32(d.calCoeffs.t1) / 8192.0)) * ((float32(rawTemp) / 131072.0) -
		float32(d.calCoeffs.t1)/8192.0)) * float32(d.calCoeffs.t3)
	temperatureComp := (tcvar1 + tcvar2) / 5120.0

	tFine := int32(tcvar1 + tcvar2)
	return temperatureComp, tFine
}

func (d *BMP280Driver) calculatePress(rawPress int32, tFine int32) float32 {
	var var1, var2, p int64

	var1 = int64(tFine) - 128000
	var2 = var1 * var1 * int64(d.calCoeffs.p6)
	var2 = var2 + ((var1 * int64(d.calCoeffs.p5)) << 17)
	var2 = var2 + (int64(d.calCoeffs.p4) << 35)
	var1 = (var1 * var1 * int64(d.calCoeffs.p3) >> 8) +
		((var1 * int64(d.calCoeffs.p2)) << 12)
	var1 = ((int64(1) << 47) + var1) * (int64(d.calCoeffs.p1)) >> 33

	if var1 == 0 {
		return 0 // avoid exception caused by division by zero
	}
	p = 1048576 - int64(rawPress)
	p = (((p << 31) - var2) * 3125) / var1
	var1 = (int64(d.calCoeffs.p9) * (p >> 13) * (p >> 13)) >> 25
	var2 = (int64(d.calCoeffs.p8) * p) >> 19

	p = ((p + var1 + var2) >> 8) + (int64(d.calCoeffs.p7) << 4)
	return float32(p) / 256
}
