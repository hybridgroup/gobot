package i2c

import (
	"bytes"
	"encoding/binary"

	"gobot.io/x/gobot"
)

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

	// TODO: set sleep mode here...

	if err := d.initialization(); err != nil {
		return err
	}

	// TODO: set usage mode here...

	// TODO: set default sea level here

	return nil
}

// Halt halts the device.
func (d *BMP280Driver) Halt() (err error) {
	return nil
}

// Temperature returns the current temperature, in celsius degrees.
func (d *BMP280Driver) Temperature() (temp float32, err error) {
	// TODO: implement this
	return 0, nil
}

// Pressure returns the current barometric pressure, in Pa
func (d *BMP280Driver) Pressure() (press float32, err error) {
	// TODO: implement this
	return 0, nil
}

// initialization reads the calibration coefficients.
func (d *BMP280Driver) initialization() (err error) {
	// TODO: set sleep mode here...

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

	// TODO: set usage mode here...
	// TODO: set default sea level here

	return nil
}

// TODO: implement
func (d *BMP280Driver) rawTempPress() (temp int16, press int16, err error) {
	return 0, 0, nil
}

// TODO: implement
func (d *BMP280Driver) calculateTemp(rawTemp int16) float32 {
	return 0
}

// TODO: implement
func (d *BMP280Driver) calculatePress(rawPress int16) float32 {
	return 0
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
