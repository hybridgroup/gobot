package i2c

import (
	"bytes"
	"encoding/binary"
	"time"

	"gobot.io/x/gobot/v2"
)

const mpl115a2DefaultAddress = 0x60

const (
	mpl115A2Reg_PressureMSB = 0x00 // first ADC register
	mpl115A2Reg_PressureLSB = 0x01
	mpl115A2Reg_TempMSB     = 0x02
	mpl115A2Reg_TempLSB     = 0x03

	mpl115A2Reg_A0_MSB  = 0x04 // first coefficient register
	mpl115A2Reg_A0_LSB  = 0x05
	mpl115A2Reg_B1_MSB  = 0x06
	mpl115A2Reg_B1_LSB  = 0x07
	mpl115A2Reg_B2_MSB  = 0x08
	mpl115A2Reg_B2_LSB  = 0x09
	mpl115A2Reg_C12_MSB = 0x0A
	mpl115A2Reg_C12_LSB = 0x0B

	mpl115A2Reg_StartConversion = 0x12
)

// MPL115A2Driver is a Gobot Driver for the MPL115A2 I2C digital pressure/temperature sensor.
// datasheet:
// https://www.nxp.com/docs/en/data-sheet/MPL115A2.pdf
//
// reference implementations:
// * https://github.com/adafruit/Adafruit_MPL115A2
type MPL115A2Driver struct {
	*Driver
	gobot.Eventer
	a0  float32
	b1  float32
	b2  float32
	c12 float32
}

// NewMPL115A2Driver creates a new Gobot Driver for an MPL115A2
// I2C Pressure/Temperature sensor.
//
// Params:
//
//	c Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
func NewMPL115A2Driver(c Connector, options ...func(Config)) *MPL115A2Driver {
	d := &MPL115A2Driver{
		Driver:  NewDriver(c, "MPL115A2", mpl115a2DefaultAddress),
		Eventer: gobot.NewEventer(),
	}
	d.afterStart = d.initialization

	for _, option := range options {
		option(d)
	}

	// TODO: add commands to API
	d.AddEvent(Error)

	return d
}

// Pressure fetches the latest data from the MPL115A2, and returns the pressure in kPa
func (d *MPL115A2Driver) Pressure() (float32, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	p, _, err := d.getData()
	return p, err
}

// Temperature fetches the latest data from the MPL115A2, and returns the temperature in °C
func (d *MPL115A2Driver) Temperature() (float32, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, t, err := d.getData()
	return t, err
}

func (d *MPL115A2Driver) initialization() error {
	data := make([]byte, 8)
	if err := d.connection.ReadBlockData(mpl115A2Reg_A0_MSB, data); err != nil {
		return err
	}

	var coA0 int16
	var coB1 int16
	var coB2 int16
	var coC12 int16

	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &coA0); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &coB1); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &coB2); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &coC12); err != nil {
		return err
	}

	coC12 = coC12 >> 2

	d.a0 = float32(coA0) / 8.0
	d.b1 = float32(coB1) / 8192.0
	d.b2 = float32(coB2) / 16384.0
	d.c12 = float32(coC12) / 4194304.0

	return nil
}

// getData fetches the latest data from the MPL115A2
//
//nolint:nonamedreturns // is sufficient here
func (d *MPL115A2Driver) getData() (p, t float32, err error) {
	var temperature uint16
	var pressure uint16
	var pressureComp float32

	if err = d.connection.WriteByteData(mpl115A2Reg_StartConversion, 0); err != nil {
		return 0, 0, err
	}
	time.Sleep(5 * time.Millisecond)

	data := []byte{0, 0, 0, 0}
	if err = d.connection.ReadBlockData(mpl115A2Reg_PressureMSB, data); err != nil {
		return 0, 0, err
	}

	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &pressure); err != nil {
		return 0, 0, err
	}
	if err := binary.Read(buf, binary.BigEndian, &temperature); err != nil {
		return 0, 0, err
	}

	temperature = temperature >> 6
	pressure = pressure >> 6

	pressureComp = d.a0 + (d.b1+d.c12*float32(temperature))*float32(pressure) + d.b2*float32(temperature)
	p = (65.0/1023.0)*pressureComp + 50.0
	t = ((float32(temperature) - 498.0) / -5.35) + 25.0

	return p, t, err
}
