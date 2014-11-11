package i2c

import (
	"github.com/hybridgroup/gobot"

	"bytes"
	"encoding/binary"
	"time"
)

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
	gobot.Driver
	A0          float32
	B1          float32
	B2          float32
	C12         float32
	Pressure    float32
	Temperature float32
}

// NewMPL115A2Driver creates a new driver with specified name and i2c interface
func NewMPL115A2Driver(a I2cInterface, name string) *MPL115A2Driver {
	return &MPL115A2Driver{
		Driver: *gobot.NewDriver(
			name,
			"MPL115A2Driver",
			a.(gobot.AdaptorInterface),
		),
	}
}

// adaptor returns MPL115A2 adaptor
func (h *MPL115A2Driver) adaptor() I2cInterface {
	return h.Adaptor().(I2cInterface)
}

// Start writes initialization bytes and reads from adaptor
// using specified interval to accelerometer andtemperature data
func (h *MPL115A2Driver) Start() bool {
	var temperature uint16
	var pressure uint16
	var pressureComp float32

	h.initialization()

	gobot.Every(h.Interval(), func() {
		h.adaptor().I2cWrite([]byte{MPL115A2_REGISTER_STARTCONVERSION, 0})
		<-time.After(5 * time.Millisecond)

		h.adaptor().I2cWrite([]byte{MPL115A2_REGISTER_PRESSURE_MSB})
		ret := h.adaptor().I2cRead(4)
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
	})
	return true
}

// Init returns true if device is initialized correctly
func (h *MPL115A2Driver) Init() bool { return true }

// Halt returns true if devices is halted successfully
func (h *MPL115A2Driver) Halt() bool { return true }

func (h *MPL115A2Driver) initialization() bool {
	var coA0 int16
	var coB1 int16
	var coB2 int16
	var coC12 int16

	h.adaptor().I2cStart(0x60)
	h.adaptor().I2cWrite([]byte{MPL115A2_REGISTER_A0_COEFF_MSB})
	ret := h.adaptor().I2cRead(8)
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

	return true
}
