package i2c

import (
	"bytes"
	"encoding/binary"
)

const bme280RegisterHumidityMSB = 0xFD
const bme280RegisterCalibDigH1 = 0xa1
const bme280RegisterCalibDigH2LSB = 0xe1

type humidityCalibrationCoefficients struct {
	h1 uint8
	h2 int16
	h3 uint8
	h4 int16
	h5 int16
	h6 int8
}

// BME280Driver is a driver for the BME280 temperature/humidity sensor
type BME280Driver struct {
	*BMP180Driver
	hc *humidityCalibrationCoefficients
}

// NewBME280Driver creates a new driver with specified i2c interface.
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewBME280Driver(a Connector, options ...func(Config)) *BME280Driver {
	b := &BME280Driver{
		BMP180Driver: NewBMP180Driver(a),
		hc:           &humidityCalibrationCoefficients{},
	}

	for _, option := range options {
		option(b)
	}

	return b
}

// Start initializes the BME280 and loads the calibration coefficients.
func (d *BME280Driver) Start() (err error) {
	bus := d.GetBusOrDefault(d.connector.GetDefaultBus())
	address := d.GetAddressOrDefault(bmp180Address)

	if d.connection, err = d.connector.GetConnection(address, bus); err != nil {
		return err
	}
	if err := d.initialization(); err != nil {
		return err
	}
	if err := d.initHumidity(); err != nil {
		return err
	}

	return nil
}

// read the humidity calibration coefficients.
func (d *BME280Driver) initHumidity() (err error) {
	var coefficients []byte
	if coefficients, err = d.read(bme280RegisterCalibDigH1, 1); err != nil {
		return err
	}
	buf := bytes.NewBuffer(coefficients)
	binary.Read(buf, binary.BigEndian, &d.hc.h1)

	if coefficients, err = d.read(bme280RegisterCalibDigH2LSB, 7); err != nil {
		return err
	}
	buf = bytes.NewBuffer(coefficients)

	// H4 and H5 laid out strangely on the bme280
	var addrE4	byte
	var addrE5	byte
	var addrE6	byte

	binary.Read(hu_buf, binary.LittleEndian, &d.hc.h2) // E1 ...
	binary.Read(hu_buf, binary.BigEndian, &d.hc.h3) // E3
	binary.Read(hu_buf, binary.BigEndian, &addrE4) // E4
	binary.Read(hu_buf, binary.BigEndian, &addrE5) // E5
	binary.Read(hu_buf, binary.BigEndian, &addrE6) // E6
	binary.Read(hu_buf, binary.BigEndian, &d.hc.h6) // ... E7

	d.hc.h4 = 0 + (int16(addrE4) << 4) | (int16(addrE5 & 0x0F))
	d.hc.h5 = 0 + (int16(addrE6) << 4) | (int16(addrE5) >> 4)

	return nil
}

// Humidity returns the current humidity in percentage of relative humidity
func (d *BME280Driver) Humidity() (humidity float32, err error) {
	var rawH int16
	if rawH, err = d.rawHumidity(); err != nil {
		return 0, nil
	}
	//TODO: return d.calculateHumidity(rawH), nil
	return float32(rawH / 1024.0), nil
}

func (d *BME280Driver) rawHumidity() (int16, error) {
	ret, err := d.read(bme280RegisterHumidityMSB, 2)
	if err != nil {
		return 0, err
	}
	buf := bytes.NewBuffer(ret)
	var rawH int16
	binary.Read(buf, binary.BigEndian, &rawH)
	return rawH, nil
}


func (d *BME280Driver) calculateHumidity(rawH int16) float32 {
	// Adapted from https://github.com/BoschSensortec/BME280_driver/blob/master/bme280.c
	// function bme280_compensate_humidity_double(s32 v_uncom_humidity_s32)
	var rawT int16
	var err error
	var var_h float32
	var t_fine int32

	rawT, err = d.rawTemp()
	if err != nil {
		return 0, err
	}

	t_fine := d.calculateB5(rawT)
	var_h = float32(t_fine - 76800)

	if var_h == 0 {
		return 0 // TODO err is 'invalid data' from Bosch - include errors or not?
	}

	var_h = (
		rawH -
			(float32(d.hc.h4) * 64.0) + // H4 double * 64.0 double
			(float32(d.hc.h5) / 16384.0 * var_h) // H5 double / 16384.0 * var_h
		) *
		(
			(float32(d.hc.h2) / 65536.0) * // H2 double / 65536.0
			(
				(1.0 + float32(d.hc.h6) / 67108864.0 * var_h) * // 1.0 + (H6 double / 67108664.0 * var_h
				(1.0 + float32(d.hc.h3) / 67108864.0 * var_h) // 1.0 + H3 double / 67108864.0 * var_h
			)
		)

	return var_h
}
