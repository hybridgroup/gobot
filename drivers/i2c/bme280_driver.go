package i2c

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
)

const bme280Debug = true

type BME280HumidityOversampling uint8

const (
	bme280RegCalibDigH1      = 0xA1
	bme280RegCalibDigH2LSB   = 0xE1
	bme280RegControlHumidity = 0xF2
	bme280RegHumidityMSB     = 0xFD

	// bits 0, 1, 3 of control humidity register
	BME280CtrlHumidityNoMeasurement  BME280HumidityOversampling = 0x00 // no measurement (value will be 0x08 0x00 0x00)
	BME280CtrlHumidityOversampling1  BME280HumidityOversampling = 0x01
	BME280CtrlHumidityOversampling2  BME280HumidityOversampling = 0x02
	BME280CtrlHumidityOversampling4  BME280HumidityOversampling = 0x03
	BME280CtrlHumidityOversampling8  BME280HumidityOversampling = 0x04
	BME280CtrlHumidityOversampling16 BME280HumidityOversampling = 0x05 // same as 0x06, 0x07
)

type bmeHumidityCalibrationCoefficients struct {
	h1 uint8
	h2 int16
	h3 uint8
	h4 int16
	h5 int16
	h6 int8
}

// BME280Driver is a driver for the BME280 temperature/humidity sensor.
// It implements all of the same functions as the BMP280Driver, but also
// adds the Humidity() function by reading the BME280's humidity sensor.
// For details on the BMP280Driver please see:
//
//	https://godoc.org/gobot.io/x/gobot/v2/drivers/i2c#BMP280Driver
type BME280Driver struct {
	*BMP280Driver
	humCalCoeffs    *bmeHumidityCalibrationCoefficients
	ctrlHumOversamp BME280HumidityOversampling
}

// NewBME280Driver creates a new driver with specified i2c interface.
// Params:
//
//	conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
func NewBME280Driver(c Connector, options ...func(Config)) *BME280Driver {
	d := &BME280Driver{
		BMP280Driver:    NewBMP280Driver(c),
		humCalCoeffs:    &bmeHumidityCalibrationCoefficients{},
		ctrlHumOversamp: BME280CtrlHumidityOversampling16,
	}
	d.afterStart = d.initializationBME280

	// this loop is for options of this class, all options of base class BMP280Driver
	// must be added in this class for usage
	for _, option := range options {
		option(d)
	}

	// TODO: expose commands to API
	return d
}

// WithBME280PressureOversampling option sets the oversampling for pressure.
// Valid settings are of type "BMP280PressureOversampling"
func WithBME280PressureOversampling(val BMP280PressureOversampling) func(Config) {
	return func(c Config) {
		if d, ok := c.(*BME280Driver); ok {
			d.ctrlPressOversamp = val
		} else if bme280Debug {
			log.Printf("Trying to set pressure oversampling for non-BME280Driver %v", c)
		}
	}
}

// WithBME280TemperatureOversampling option sets oversampling for temperature.
// Valid settings are of type "BMP280TemperatureOversampling"
func WithBME280TemperatureOversampling(val BMP280TemperatureOversampling) func(Config) {
	return func(c Config) {
		if d, ok := c.(*BME280Driver); ok {
			d.ctrlTempOversamp = val
		} else if bme280Debug {
			log.Printf("Trying to set temperature oversampling for non-BME280Driver %v", c)
		}
	}
}

// WithBME280IIRFilter option sets the count of IIR filter coefficients.
// Valid settings are of type "BMP280IIRFilter"
func WithBME280IIRFilter(val BMP280IIRFilter) func(Config) {
	return func(c Config) {
		if d, ok := c.(*BME280Driver); ok {
			d.confFilter = val
		} else if bme280Debug {
			log.Printf("Trying to set IIR filter for non-BME280Driver %v", c)
		}
	}
}

// WithBME280HumidityOversampling option sets the oversampling for humidity.
// Valid settings are of type "BME280HumidityOversampling"
func WithBME280HumidityOversampling(val BME280HumidityOversampling) func(Config) {
	return func(c Config) {
		if d, ok := c.(*BME280Driver); ok {
			d.ctrlHumOversamp = val
		} else if bme280Debug {
			log.Printf("Trying to set humidity oversampling for non-BME280Driver %v", c)
		}
	}
}

// Humidity returns the current humidity in percentage of relative humidity
func (d *BME280Driver) Humidity() (float32, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	rawH, err := d.rawHumidity()
	if err != nil {
		return 0.0, err
	}
	humidity := d.calculateHumidity(rawH)
	return humidity, nil
}

func (d *BME280Driver) initializationBME280() error {
	// call the initialization routine of base class BMP280Driver, which do:
	// * initializes temperature and pressure calibration coefficients
	// * set the control register
	// * set the configuration register
	if err := d.initialization(); err != nil {
		return err
	}

	if err := d.initHumidity(); err != nil {
		return err
	}

	return nil
}

// read the humidity calibration coefficients.
func (d *BME280Driver) initHumidity() error {
	hch1, err := d.connection.ReadByteData(bme280RegCalibDigH1)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer([]byte{hch1})
	if err := binary.Read(buf, binary.BigEndian, &d.humCalCoeffs.h1); err != nil {
		return err
	}

	coefficients := make([]byte, 7)
	if err = d.connection.ReadBlockData(bme280RegCalibDigH2LSB, coefficients); err != nil {
		return err
	}
	buf = bytes.NewBuffer(coefficients)

	// H4 and H5 laid out strangely on the bme280
	var addrE4 byte
	var addrE5 byte
	var addrE6 byte
	// E1 ...
	if err := binary.Read(buf, binary.LittleEndian, &d.humCalCoeffs.h2); err != nil {
		return err
	}
	// E3
	if err := binary.Read(buf, binary.BigEndian, &d.humCalCoeffs.h3); err != nil {
		return err
	}
	// E4
	if err := binary.Read(buf, binary.BigEndian, &addrE4); err != nil {
		return err
	}
	// E5
	if err := binary.Read(buf, binary.BigEndian, &addrE5); err != nil {
		return err
	}
	// E6
	if err := binary.Read(buf, binary.BigEndian, &addrE6); err != nil {
		return err
	}
	// ... E7
	if err := binary.Read(buf, binary.BigEndian, &d.humCalCoeffs.h6); err != nil {
		return err
	}

	d.humCalCoeffs.h4 = 0 + (int16(addrE4) << 4) | (int16(addrE5 & 0x0F))
	d.humCalCoeffs.h5 = 0 + (int16(addrE6) << 4) | (int16(addrE5) >> 4)

	// The 'ctrl_hum' register (0xF2) sets the humidity data acquisition options of
	// the device. Changes to this register only become effective after a write
	// operation to 'ctrl_meas' (0xF4). So we read the current value in, then write it back
	if err := d.connection.WriteByteData(bme280RegControlHumidity, uint8(d.ctrlHumOversamp)); err != nil {
		return err
	}

	cmr, err := d.connection.ReadByteData(bmp280RegCtrl)
	if err != nil {
		return err
	}

	return d.connection.WriteByteData(bmp280RegCtrl, cmr)
}

func (d *BME280Driver) rawHumidity() (uint32, error) {
	ret := make([]byte, 2)
	if err := d.connection.ReadBlockData(bme280RegHumidityMSB, ret); err != nil {
		return 0, err
	}
	if ret[0] == 0x80 && ret[1] == 0x00 {
		return 0, errors.New("Humidity disabled")
	}
	buf := bytes.NewBuffer(ret)
	var rawH uint16
	if err := binary.Read(buf, binary.BigEndian, &rawH); err != nil {
		return 0, err
	}
	return uint32(rawH), nil
}

// Adapted from https://github.com/BoschSensortec/BME280_driver/blob/master/bme280.c
// function bme280_compensate_humidity_double(s32 v_uncom_humidity_s32)
func (d *BME280Driver) calculateHumidity(rawH uint32) float32 {
	var rawT int32
	var err error
	var h float32

	rawT, err = d.rawTemp()
	if err != nil {
		return 0
	}

	_, tFine := d.calculateTemp(rawT)
	h = float32(tFine) - 76800

	if h == 0 {
		return 0 // TODO err is 'invalid data' from Bosch - include errors or not?
	}

	x := float32(rawH) - (float32(d.humCalCoeffs.h4)*64.0 +
		(float32(d.humCalCoeffs.h5) / 16384.0 * h))

	y := float32(d.humCalCoeffs.h2) / 65536.0 *
		(1.0 + float32(d.humCalCoeffs.h6)/67108864.0*h*
			(1.0+float32(d.humCalCoeffs.h3)/67108864.0*h))

	h = x * y
	h = h * (1 - float32(d.humCalCoeffs.h1)*h/524288)
	return h
}
