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

type MPL115A2Driver struct {
	name       string
	connection I2c
	interval   time.Duration
	gobot.Eventer
	A0          float32
	B1          float32
	B2          float32
	C12         float32
	Pressure    float32
	Temperature float32
}

// NewMPL115A2Driver creates a new driver with specified i2c interface
func NewMPL115A2Driver(a I2c, v ...time.Duration) *MPL115A2Driver {
	m := &MPL115A2Driver{
		name:       "MPL115A2",
		connection: a,
		Eventer:    gobot.NewEventer(),
		interval:   10 * time.Millisecond,
	}

	if len(v) > 0 {
		m.interval = v[0]
	}
	m.AddEvent(Error)
	return m
}

func (h *MPL115A2Driver) Name() string                 { return h.name }
func (h *MPL115A2Driver) SetName(n string)             { h.name = n }
func (h *MPL115A2Driver) Connection() gobot.Connection { return h.connection.(gobot.Connection) }

// Start writes initialization bytes and reads from adaptor
// using specified interval to accelerometer and temperature data
func (h *MPL115A2Driver) Start() (err error) {
	var temperature uint16
	var pressure uint16
	var pressureComp float32

	if err := h.initialization(); err != nil {
		return err
	}

	go func() {
		for {
			if err := h.connection.I2cWrite(mpl115a2Address, []byte{MPL115A2_REGISTER_STARTCONVERSION, 0}); err != nil {
				h.Publish(h.Event(Error), err)
				continue

			}
			time.Sleep(5 * time.Millisecond)

			if err := h.connection.I2cWrite(mpl115a2Address, []byte{MPL115A2_REGISTER_PRESSURE_MSB}); err != nil {
				h.Publish(h.Event(Error), err)
				continue
			}

			ret, err := h.connection.I2cRead(mpl115a2Address, 4)
			if err != nil {
				h.Publish(h.Event(Error), err)
				continue
			}
			if len(ret) == 4 {
				buf := bytes.NewBuffer(ret)
				binary.Read(buf, binary.BigEndian, &pressure)
				binary.Read(buf, binary.BigEndian, &temperature)

				temperature = temperature >> 6
				pressure = pressure >> 6

				pressureComp = float32(h.A0) + (float32(h.B1)+float32(h.C12)*float32(temperature))*float32(pressure) + float32(h.B2)*float32(temperature)
				h.Pressure = (65.0/1023.0)*pressureComp + 50.0
				h.Temperature = ((float32(temperature) - 498.0) / -5.35) + 25.0
			}
			time.Sleep(h.interval)
		}
	}()
	return
}

// Halt returns true if devices is halted successfully
func (h *MPL115A2Driver) Halt() (err error) { return }

func (h *MPL115A2Driver) initialization() (err error) {
	var coA0 int16
	var coB1 int16
	var coB2 int16
	var coC12 int16

	if err = h.connection.I2cStart(mpl115a2Address); err != nil {
		return
	}
	if err = h.connection.I2cWrite(mpl115a2Address, []byte{MPL115A2_REGISTER_A0_COEFF_MSB}); err != nil {
		return
	}
	ret, err := h.connection.I2cRead(mpl115a2Address, 8)
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(ret)

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
