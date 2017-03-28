package i2c

import (
	"bytes"
	"encoding/binary"

	"gobot.io/x/gobot"
)

const bme280RegisterHumidityMSB = 0xFD
const bme280RegisterCalibDigH1 = 0xa1
const bme280RegisterCalibDigH2LSB = 0xe1
const bmp280RegisterCalib00 = 0x88

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

type bmeHumidityCalibrationCoefficients struct {
	h1 uint8
	h2 int16
	h3 uint8
	h4 int16
	h5 int16
	h6 int8
}

// BME280Driver is a driver for the BME280 temperature/humidity sensor
type BME280Driver struct {
	name       string
	connector  Connector
	connection Connection
	Config

	tpc *bmp280CalibrationCoefficients
	hc  *bmeHumidityCalibrationCoefficients
}

// NewBME280Driver creates a new driver with specified i2c interface.
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewBME280Driver(c Connector, options ...func(Config)) *BME280Driver {
	b := &BME280Driver{
		name:      gobot.DefaultName("BME280"),
		connector: c,
		Config:    NewConfig(),
		tpc:       &bmp280CalibrationCoefficients{},
		hc:        &bmeHumidityCalibrationCoefficients{},
	}

	for _, option := range options {
		option(b)
	}

	// TODO: expose commands to API
	return b
}

// Name returns the name of the device.
func (d *BME280Driver) Name() string {
	return d.name
}

// SetName sets the name of the device.
func (d *BME280Driver) SetName(n string) {
	d.name = n
}

// Connection returns the connection of the device.
func (d *BME280Driver) Connection() gobot.Connection {
	return d.connector.(gobot.Connection)
}

// Start initializes the BME280 and loads the calibration coefficients.
func (d *BME280Driver) Start() (err error) {
	bus := d.GetBusOrDefault(d.connector.GetDefaultBus())
	address := d.GetAddressOrDefault(bmp180Address)

	if d.connection, err = d.connector.GetConnection(address, bus); err != nil {
		return err
	}

	// TODO: set sleep mode here...

	if err := d.initialization(); err != nil {
		return err
	}
	if err := d.initHumidity(); err != nil {
		return err
	}

	// TODO: set usage mode here...

	// TODO: set default sea level here

	return nil
}

// Halt halts the device.
func (d *BME280Driver) Halt() (err error) {
	return nil
}

// Temperature returns the current temperature, in celsius degrees.
func (d *BME280Driver) Temperature() (temp float32, err error) {
	// TODO: implement this
	return 0, nil
}

// Pressure returns the current barometric pressure, in Pa
func (d *BME280Driver) Pressure() (press float32, err error) {
	// TODO: implement this
	return 0, nil
}

// Humidity returns the current humidity in percentage of relative humidity
func (d *BME280Driver) Humidity() (humidity float32, err error) {
	// TODO: implement this
	return 0, nil
}

// initialization reads the calibration coefficients.
func (d *BME280Driver) initialization() (err error) {
	var coefficients []byte
	if coefficients, err = d.read(bmp280RegisterCalib00, 26); err != nil {
		return err
	}
	buf := bytes.NewBuffer(coefficients)
	binary.Read(buf, binary.LittleEndian, &d.tpc.t1)
	binary.Read(buf, binary.LittleEndian, &d.tpc.t2)
	binary.Read(buf, binary.LittleEndian, &d.tpc.t3)
	binary.Read(buf, binary.LittleEndian, &d.tpc.p1)
	binary.Read(buf, binary.LittleEndian, &d.tpc.p2)
	binary.Read(buf, binary.LittleEndian, &d.tpc.p3)
	binary.Read(buf, binary.LittleEndian, &d.tpc.p4)
	binary.Read(buf, binary.LittleEndian, &d.tpc.p5)
	binary.Read(buf, binary.LittleEndian, &d.tpc.p6)
	binary.Read(buf, binary.LittleEndian, &d.tpc.p7)
	binary.Read(buf, binary.LittleEndian, &d.tpc.p8)
	binary.Read(buf, binary.LittleEndian, &d.tpc.p9)

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
	var addrE4 byte
	var addrE5 byte
	var addrE6 byte

	binary.Read(buf, binary.LittleEndian, &d.hc.h2) // E1 ...
	binary.Read(buf, binary.BigEndian, &d.hc.h3)    // E3
	binary.Read(buf, binary.BigEndian, &addrE4)     // E4
	binary.Read(buf, binary.BigEndian, &addrE5)     // E5
	binary.Read(buf, binary.BigEndian, &addrE6)     // E6
	binary.Read(buf, binary.BigEndian, &d.hc.h6)    // ... E7

	d.hc.h4 = 0 + (int16(addrE4) << 4) | (int16(addrE5 & 0x0F))
	d.hc.h5 = 0 + (int16(addrE6) << 4) | (int16(addrE5) >> 4)

	return nil
}

// TODO: implement
func (d *BME280Driver) rawTempPress() (temp int16, press int16, err error) {
	return 0, 0, nil
}

// TODO: implement
func (d *BME280Driver) calculateTemp(rawTemp int16) float32 {
	return 0
}

// TODO: implement
func (d *BME280Driver) calculatePress(rawPress int16) float32 {
	return 0
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

// Adapted from https://github.com/BoschSensortec/BME280_driver/blob/master/bme280.c
// function bme280_compensate_humidity_double(s32 v_uncom_humidity_s32)
func (d *BME280Driver) calculateHumidity(rawH int16) float32 {
	var rawT int16
	var err error
	var h float32
	var fine int32

	rawT, _, err = d.rawTempPress()
	if err != nil {
		return 0
	}

	// TODO: calculate fine temp adjust
	fine = int32(rawT)
	h = float32(fine - 76800)

	if h == 0 {
		return 0 // TODO err is 'invalid data' from Bosch - include errors or not?
	}

	h = (float32(rawH) - (float32(d.hc.h4) * 64.0) + // H4 double * 64.0 double
		(float32(d.hc.h5) / 16384.0 * h)) // H5 double / 16384.0 * var_h

	y :=
		(float32(d.hc.h2) / 65536.0) * // H2 double / 65536.0
			(1.0 + float32(d.hc.h6)/67108864.0*h) * // 1.0 + (H6 double / 67108664.0 * var_h
			(1.0 + float32(d.hc.h3)/67108864.0*h) // 1.0 + H3 double / 67108864.0 * var_h

	return h * y
}

func (d *BME280Driver) read(address byte, n int) ([]byte, error) {
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
