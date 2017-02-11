package i2c

import (
	"bytes"
	"encoding/binary"
	"time"

	"gobot.io/x/gobot"
)

const mpu6050Address = 0x68

const MPU6050_RA_ACCEL_XOUT_H = 0x3B
const MPU6050_RA_PWR_MGMT_1 = 0x6B
const MPU6050_PWR1_CLKSEL_BIT = 2
const MPU6050_PWR1_CLKSEL_LENGTH = 3
const MPU6050_CLOCK_PLL_XGYRO = 0x01
const MPU6050_GYRO_FS_250 = 0x00
const MPU6050_RA_GYRO_CONFIG = 0x1B
const MPU6050_GCONFIG_FS_SEL_LENGTH = 2
const MPU6050_GCONFIG_FS_SEL_BIT = 4
const MPU6050_RA_ACCEL_CONFIG = 0x1C
const MPU6050_ACONFIG_AFS_SEL_BIT = 4
const MPU6050_ACONFIG_AFS_SEL_LENGTH = 2
const MPU6050_ACCEL_FS_2 = 0x00
const MPU6050_PWR1_SLEEP_BIT = 6
const MPU6050_PWR1_ENABLE_BIT = 0

type ThreeDData struct {
	X int16
	Y int16
	Z int16
}

// MPU6050Driver is a new Gobot Driver for an MPU6050 I2C Accelerometer/Gyroscope.
type MPU6050Driver struct {
	name       string
	connector  Connector
	connection Connection
	Config
	interval      time.Duration
	Accelerometer ThreeDData
	Gyroscope     ThreeDData
	Temperature   int16
	gobot.Eventer
}

// NewMPU6050Driver creates a new Gobot Driver for an MPU6050 I2C Accelerometer/Gyroscope.
//
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewMPU6050Driver(a Connector, options ...func(Config)) *MPU6050Driver {
	m := &MPU6050Driver{
		name:      gobot.DefaultName("MPU6050"),
		connector: a,
		Config:    NewConfig(),
		Eventer:   gobot.NewEventer(),
	}

	for _, option := range options {
		option(m)
	}

	// TODO: add commands to API
	return m
}

// Name returns the name of the device.
func (h *MPU6050Driver) Name() string { return h.name }

// SetName sets the name of the device.
func (h *MPU6050Driver) SetName(n string) { h.name = n }

// Connection returns the connection for the device.
func (h *MPU6050Driver) Connection() gobot.Connection { return h.connector.(gobot.Connection) }

// Start writes initialization bytes to sensor
func (h *MPU6050Driver) Start() (err error) {
	if err := h.initialize(); err != nil {
		return err
	}

	return
}

// Halt returns true if devices is halted successfully
func (h *MPU6050Driver) Halt() (err error) { return }

// GetData fetches the latest data from the MPU6050
func (h *MPU6050Driver) GetData() (err error) {
	if _, err = h.connection.Write([]byte{MPU6050_RA_ACCEL_XOUT_H}); err != nil {
		return
	}

	data := make([]byte, 14)
	_, err = h.connection.Read(data)
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.BigEndian, &h.Accelerometer)
	binary.Read(buf, binary.BigEndian, &h.Temperature)
	binary.Read(buf, binary.BigEndian, &h.Gyroscope)
	h.convertToCelsius()
	return
}

func (h *MPU6050Driver) initialize() (err error) {
	bus := h.GetBusOrDefault(h.connector.GetDefaultBus())
	address := h.GetAddressOrDefault(mpu6050Address)

	h.connection, err = h.connector.GetConnection(address, bus)
	if err != nil {
		return err
	}

	// setClockSource
	if _, err = h.connection.Write([]byte{MPU6050_RA_PWR_MGMT_1,
		MPU6050_PWR1_CLKSEL_BIT,
		MPU6050_PWR1_CLKSEL_LENGTH,
		MPU6050_CLOCK_PLL_XGYRO}); err != nil {
		return
	}

	// setFullScaleGyroRange
	if _, err = h.connection.Write([]byte{MPU6050_RA_GYRO_CONFIG,
		MPU6050_GCONFIG_FS_SEL_BIT,
		MPU6050_GCONFIG_FS_SEL_LENGTH,
		MPU6050_GYRO_FS_250}); err != nil {
		return
	}

	// setFullScaleAccelRange
	if _, err = h.connection.Write([]byte{MPU6050_RA_ACCEL_CONFIG,
		MPU6050_ACONFIG_AFS_SEL_BIT,
		MPU6050_ACONFIG_AFS_SEL_LENGTH,
		MPU6050_ACCEL_FS_2}); err != nil {
		return
	}

	// setSleepEnabled
	if _, err = h.connection.Write([]byte{MPU6050_RA_PWR_MGMT_1,
		MPU6050_PWR1_ENABLE_BIT,
		0}); err != nil {
		return
	}

	return nil
}

// The temperature sensor is -40 to +85 degrees Celsius.
// It is a signed integer.
// According to the datasheet:
//   340 per degrees Celsius, -512 at 35 degrees.
// At 0 degrees: -512 - (340 * 35) = -12412
func (h *MPU6050Driver) convertToCelsius() {
	h.Temperature = (h.Temperature + 12412) / 340
}
