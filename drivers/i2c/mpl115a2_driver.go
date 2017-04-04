package i2c

import (
	"gobot.io/x/gobot"

	"bytes"
	"encoding/binary"
	"time"
)

const mpl115a2Address = 0x60

const MPL115A2_REGISTER_PRESSURE_MSB = 0x00
const MPL115A2_REGISTER_PRESSURE_LSB = 0x01
const MPL115A2_REGISTER_TEMP_MSB = 0x02
const MPL115A2_REGISTER_TEMP_LSB = 0x03
const MPL115A2_REGISTER_A0_COEFF_MSB = 0x04
const MPL115A2_REGISTER_A0_COEFF_LSB = 0x05
const MPL115A2_REGISTER_B1_COEFF_MSB = 0x06
const MPL115A2_REGISTER_B1_COEFF_LSB = 0x07
const MPL115A2_REGISTER_B2_COEFF_MSB = 0x08
const MPL115A2_REGISTER_B2_COEFF_LSB = 0x09
const MPL115A2_REGISTER_C12_COEFF_MSB = 0x0A
const MPL115A2_REGISTER_C12_COEFF_LSB = 0x0B
const MPL115A2_REGISTER_STARTCONVERSION = 0x12

// MPL115A2Driver is a Gobot Driver for the MPL115A2 I2C digitial pressure/temperature sensor.
type MPL115A2Driver struct {
	name       string
	connector  Connector
	connection Connection
	Config
	gobot.Eventer
	A0  float32
	B1  float32
	B2  float32
	C12 float32
}

// NewMPL115A2Driver creates a new Gobot Driver for an MPL115A2
// I2C Pressure/Temperature sensor.
//
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewMPL115A2Driver(a Connector, options ...func(Config)) *MPL115A2Driver {
	m := &MPL115A2Driver{
		name:      gobot.DefaultName("MPL115A2"),
		connector: a,
		Config:    NewConfig(),
		Eventer:   gobot.NewEventer(),
	}

	for _, option := range options {
		option(m)
	}

	// TODO: add commands to API
	m.AddEvent(Error)

	return m
}

// Name returns the name of the device.
func (h *MPL115A2Driver) Name() string { return h.name }

// SetName sets the name of the device.
func (h *MPL115A2Driver) SetName(n string) { h.name = n }

// Connection returns the Connection of the device.
func (h *MPL115A2Driver) Connection() gobot.Connection { return h.connector.(gobot.Connection) }

// Start writes initialization bytes and reads from adaptor
// using specified interval to accelerometer andtemperature data
func (h *MPL115A2Driver) Start() (err error) {
	if err := h.initialization(); err != nil {
		return err
	}

	return
}

// Halt returns true if devices is halted successfully
func (h *MPL115A2Driver) Halt() (err error) { return }

// Pressure fetches the latest data from the MPL115A2, and returns the pressure
func (h *MPL115A2Driver) Pressure() (p float32, err error) {
	p, _, err = h.getData()
	return
}

// Temperature fetches the latest data from the MPL115A2, and returns the temperature
func (h *MPL115A2Driver) Temperature() (t float32, err error) {
	_, t, err = h.getData()
	return
}

func (h *MPL115A2Driver) initialization() (err error) {
	var coA0 int16
	var coB1 int16
	var coB2 int16
	var coC12 int16

	bus := h.GetBusOrDefault(h.connector.GetDefaultBus())
	address := h.GetAddressOrDefault(mpl115a2Address)

	h.connection, err = h.connector.GetConnection(address, bus)
	if err != nil {
		return err
	}

	if _, err = h.connection.Write([]byte{MPL115A2_REGISTER_A0_COEFF_MSB}); err != nil {
		return
	}

	data := make([]byte, 8)
	if _, err = h.connection.Read(data); err != nil {
		return
	}
	buf := bytes.NewBuffer(data)

	binary.Read(buf, binary.BigEndian, &coA0)
	binary.Read(buf, binary.BigEndian, &coB1)
	binary.Read(buf, binary.BigEndian, &coB2)
	binary.Read(buf, binary.BigEndian, &coC12)

	coC12 = coC12 >> 2

	h.A0 = float32(coA0) / 8.0
	h.B1 = float32(coB1) / 8192.0
	h.B2 = float32(coB2) / 16384.0
	h.C12 = float32(coC12) / 4194304.0

	return
}

// getData fetches the latest data from the MPL115A2
func (h *MPL115A2Driver) getData() (p, t float32, err error) {
	var temperature uint16
	var pressure uint16
	var pressureComp float32

	if _, err = h.connection.Write([]byte{MPL115A2_REGISTER_STARTCONVERSION, 0}); err != nil {
		return
	}
	time.Sleep(5 * time.Millisecond)

	if _, err = h.connection.Write([]byte{MPL115A2_REGISTER_PRESSURE_MSB}); err != nil {
		return
	}

	data := []byte{0, 0, 0, 0}
	bytesRead, err1 := h.connection.Read(data)
	if err1 != nil {
		err = err1
		return
	}

	if bytesRead == 4 {
		buf := bytes.NewBuffer(data)
		binary.Read(buf, binary.BigEndian, &pressure)
		binary.Read(buf, binary.BigEndian, &temperature)

		temperature = temperature >> 6
		pressure = pressure >> 6

		pressureComp = float32(h.A0) + (float32(h.B1)+float32(h.C12)*float32(temperature))*float32(pressure) + float32(h.B2)*float32(temperature)
		p = (65.0/1023.0)*pressureComp + 50.0
		t = ((float32(temperature) - 498.0) / -5.35) + 25.0
	}

	return
}
