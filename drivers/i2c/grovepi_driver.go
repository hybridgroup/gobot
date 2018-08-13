package i2c

import (
	"strconv"
	"time"

	"gobot.io/x/gobot"
)

const grovePiAddress = 0x04

// Commands format
const (
	CommandReadDigital  = 1
	CommandWriteDigital = 2
	CommandReadAnalog   = 3
	CommandWriteAnalog  = 4
	CommandReadDHT      = 40
)

// GrovePiDriver is a driver for the GrovePi for I²C bus interface.
type GrovePiDriver struct {
	name       string
	connector  Connector
	connection Connection
	Config
}

// NewGrovePiDriver creates a new driver with specified i2c interface
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewGrovePiDriver(a Connector, options ...func(Config)) *GrovePiDriver {
	d := &GrovePiDriver{
		name:      gobot.DefaultName("GrovePi"),
		connector: a,
		Config:    NewConfig(),
	}

	for _, option := range options {
		option(d)
	}

	// TODO: add commands for API
	return d
}

// Name returns the Name for the Driver
func (d *GrovePiDriver) Name() string { return d.name }

// SetName sets the Name for the Driver
func (d *GrovePiDriver) SetName(n string) { d.name = n }

// Connection returns the connection for the Driver
func (d *GrovePiDriver) Connection() gobot.Connection { return d.connector.(gobot.Connection) }

// Start initialized the GrovePi
func (d *GrovePiDriver) Start() (err error) {
	bus := d.GetBusOrDefault(d.connector.GetDefaultBus())
	address := d.GetAddressOrDefault(grovePiAddress)

	d.connection, err = d.connector.GetConnection(address, bus)
	if err != nil {
		return err
	}

	return
}

// Halt returns true if devices is halted successfully
func (d *GrovePiDriver) Halt() (err error) { return }

// AnalogRead returns value from analog pin implementing the AnalogReader interface.
func (d *GrovePiDriver) AnalogRead(pin string) (value int, err error) {
	// TODO: strip off the leading "A"

	var channel int
	channel, err = strconv.Atoi(pin)
	if err != nil {
		return
	}
	value, err = d.ReadAnalog(byte(channel))

	return
}

// ReadAnalog reads analog value from the GrovePi
func (d *GrovePiDriver) ReadAnalog(pin byte) (int, error) {
	b := []byte{1, CommandReadAnalog, pin, 0, 0}
	_, err := d.connection.Write(b)
	if err != nil {
		return 0, err
	}

	time.Sleep(100 * time.Millisecond)

	d.connection.Write([]byte{1})
	d.connection.ReadByte()

	data := make([]byte, 4)
	d.connection.Write([]byte{1})
	_, err = d.connection.Read(data)
	if err != nil {
		return 0, err
	}

	v1 := int(data[1])
	v2 := int(data[2])
	return ((v1 * 256) + v2), nil
}

// ReadDigital reads digitally to the GrovePi
func (d *GrovePiDriver) ReadDigital(pin byte, size int) ([]byte, error) {
	buf := []byte{CommandReadDigital, pin, 0, 0}
	_, err := d.connection.Write(buf)
	if err != nil {
		return nil, err
	}

	time.Sleep(100 * time.Millisecond)

	d.connection.Write([]byte{1})
	data := make([]byte, size)
	_, err = d.connection.Read(data)
	if err != nil {
		return nil, err
	}

	return data, err
}

// WriteDigital writes digitally to the GrovePi
func (d *GrovePiDriver) WriteDigital(pin byte, val byte) error {
	buf := []byte{1, CommandWriteDigital, pin, val, 0}
	_, err := d.connection.Write(buf)
	time.Sleep(100 * time.Millisecond)
	return err
}

// DigitalWrite writes a value to a specific digital pin implementing the DigitalWriter interface.
func (d *GrovePiDriver) DigitalWrite(pin string, val byte) (err error) {
	// TODO: strip off the leading "D"
	var channel int
	channel, err = strconv.Atoi(pin)
	if err != nil {
		return
	}
	err = d.WriteDigital(byte(channel), val)

	return
}

// WriteAnalog writes analog to the GrovePi
func (d *GrovePiDriver) WriteAnalog(pin byte, val byte) error {
	buf := []byte{1, CommandWriteAnalog, pin, val, 0}
	_, err := d.connection.Write(buf)
	time.Sleep(100 * time.Millisecond)
	return err
}

// ReadDHT returns temperature and humidity from DHT sensor
func (d *GrovePiDriver) ReadDHT(pin byte, size int) ([]byte, error) {
	cmd := []byte{1, CommandReadDHT, pin, 0, 0}

	// prepare and read raw data
	_, err := d.connection.Write(cmd)
	if err != nil {
		return nil, err
	}
	time.Sleep(600 * time.Millisecond)
	d.connection.Write([]byte{1})
	d.connection.ReadByte()
	time.Sleep(100 * time.Millisecond)

	data := make([]byte, size)
	d.connection.Write([]byte{1})
	_, err = d.connection.Read(data)
	if err != nil {
		return nil, err
	}

	return data, err
}
