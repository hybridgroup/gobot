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
	binary.Read(buf, binary.LittleEndian, &d.hc.h2)
	binary.Read(buf, binary.BigEndian, &d.hc.h3)
	binary.Read(buf, binary.BigEndian, &d.hc.h4)
	binary.Read(buf, binary.BigEndian, &d.hc.h5)
	binary.Read(buf, binary.BigEndian, &d.hc.h6)
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
	// TODO: real adjustment based on hc coefficients
	return 0.0 / 1024.0
}
