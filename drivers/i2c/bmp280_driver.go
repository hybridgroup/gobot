package i2c

import (
	"bytes"
	"encoding/binary"
	"math"

	"gobot.io/x/gobot"
)

const (
	bmp280RegisterControl      = 0xf4
	bmp280RegisterConfig       = 0xf5
	bmp280RegisterPressureData = 0xf7
	bmp280RegisterTempData     = 0xfa
	bmp280RegisterCalib00      = 0x88
	bmp280SeaLevelPressure     = 1013.25
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
	name       string
	connector  Connector
	connection Connection
	Config

	tpc *bmp280CalibrationCoefficients
}

// NewBMP280Driver creates a new driver with specified i2c interface.
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewBMP280Driver(c Connector, options ...func(Config)) *BMP280Driver {
	b := &BMP280Driver{
		name:      gobot.DefaultName("BMP280"),
		connector: c,
		Config:    NewConfig(),
		tpc:       &bmp280CalibrationCoefficients{},
	}

	for _, option := range options {
		option(b)
	}

	// TODO: expose commands to API
	return b
}

// Name returns the name of the device.
func (d *BMP280Driver) Name() string {
	return d.name
}

// SetName sets the name of the device.
func (d *BMP280Driver) SetName(n string) {
	d.name = n
}

// Connection returns the connection of the device.
func (d *BMP280Driver) Connection() gobot.Connection {
	return d.connector.(gobot.Connection)
}

// Start initializes the BMP280 and loads the calibration coefficients.
func (d *BMP280Driver) Start() (err error) {
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

// Halt halts the device.
func (d *BMP280Driver) Halt() (err error) {
	return nil
}

// Temperature returns the current temperature, in celsius degrees.
func (d *BMP280Driver) Temperature() (temp float32, err error) {
	var rawT int32
	if rawT, err = d.rawTemp(); err != nil {
		return 0.0, err
	}
	temp, _ = d.calculateTemp(rawT)
	return
}

// Pressure returns the current barometric pressure, in Pa
func (d *BMP280Driver) Pressure() (press float32, err error) {
	var rawT, rawP int32
	if rawT, err = d.rawTemp(); err != nil {
		return 0.0, err
	}

	if rawP, err = d.rawPressure(); err != nil {
		return 0.0, err
	}
	_, tFine := d.calculateTemp(rawT)
	return d.calculatePress(rawP, tFine), nil
}

// Altitude returns the current altitude in meters based on the
// current barometric pressure and estimated pressure at sea level.
// Calculation is based on code from Adafruit BME280 library
// 	https://github.com/adafruit/Adafruit_BME280_Library
func (d *BMP280Driver) Altitude() (alt float32, err error) {
	atmP, _ := d.Pressure()
	atmP /= 100.0
	alt = float32(44330.0 * (1.0 - math.Pow(float64(atmP/bmp280SeaLevelPressure), 0.1903)))

	return
}

// initialization reads the calibration coefficients.
func (d *BMP280Driver) initialization() (err error) {
	var coefficients []byte
	if coefficients, err = d.read(bmp280RegisterCalib00, 24); err != nil {
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

	d.connection.WriteByteData(bmp280RegisterControl, 0x3F)

	return nil
}

func (d *BMP280Driver) rawTemp() (temp int32, err error) {
	var data []byte
	var tp0, tp1, tp2 byte

	if data, err = d.read(bmp280RegisterTempData, 3); err != nil {
		return 0, err
	}
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &tp0)
	binary.Read(buf, binary.LittleEndian, &tp1)
	binary.Read(buf, binary.LittleEndian, &tp2)

	temp = ((int32(tp2) >> 4) | (int32(tp1) << 4) | (int32(tp0) << 12))

	return
}

func (d *BMP280Driver) rawPressure() (press int32, err error) {
	var data []byte
	var tp0, tp1, tp2 byte

	if data, err = d.read(bmp280RegisterPressureData, 3); err != nil {
		return 0, err
	}
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &tp0)
	binary.Read(buf, binary.LittleEndian, &tp1)
	binary.Read(buf, binary.LittleEndian, &tp2)

	press = ((int32(tp2) >> 4) | (int32(tp1) << 4) | (int32(tp0) << 12))

	return
}

func (d *BMP280Driver) calculateTemp(rawTemp int32) (float32, int32) {
	tcvar1 := ((float32(rawTemp) / 16384.0) - (float32(d.tpc.t1) / 1024.0)) * float32(d.tpc.t2)
	tcvar2 := (((float32(rawTemp) / 131072.0) - (float32(d.tpc.t1) / 8192.0)) * ((float32(rawTemp) / 131072.0) - float32(d.tpc.t1)/8192.0)) * float32(d.tpc.t3)
	temperatureComp := (tcvar1 + tcvar2) / 5120.0

	tFine := int32(tcvar1 + tcvar2)
	return temperatureComp, tFine
}

func (d *BMP280Driver) calculatePress(rawPress int32, tFine int32) float32 {
	var var1, var2, p int64

	var1 = int64(tFine) - 128000
	var2 = var1 * var1 * int64(d.tpc.p6)
	var2 = var2 + ((var1 * int64(d.tpc.p5)) << 17)
	var2 = var2 + (int64(d.tpc.p4) << 35)
	var1 = (var1 * var1 * int64(d.tpc.p3) >> 8) +
		((var1 * int64(d.tpc.p2)) << 12)
	var1 = ((int64(1) << 47) + var1) * (int64(d.tpc.p1)) >> 33

	if var1 == 0 {
		return 0 // avoid exception caused by division by zero
	}
	p = 1048576 - int64(rawPress)
	p = (((p << 31) - var2) * 3125) / var1
	var1 = (int64(d.tpc.p9) * (p >> 13) * (p >> 13)) >> 25
	var2 = (int64(d.tpc.p8) * p) >> 19

	p = ((p + var1 + var2) >> 8) + (int64(d.tpc.p7) << 4)
	return float32(p) / 256
}

func (d *BMP280Driver) read(address byte, n int) ([]byte, error) {
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
