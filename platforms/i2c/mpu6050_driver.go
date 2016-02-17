package i2c

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*MPU6050Driver)(nil)

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

type MPU6050Driver struct {
	name          string
	connection    I2c
	interval      time.Duration
	Accelerometer ThreeDData
	Gyroscope     ThreeDData
	Temperature   int16
	gobot.Eventer
}

// NewMPU6050Driver creates a new driver with specified name and i2c interface
func NewMPU6050Driver(a I2c, name string, v ...time.Duration) *MPU6050Driver {
	m := &MPU6050Driver{
		name:       name,
		connection: a,
		interval:   10 * time.Millisecond,
		Eventer:    gobot.NewEventer(),
	}

	if len(v) > 0 {
		m.interval = v[0]
	}

	m.AddEvent(Error)
	return m
}

func (h *MPU6050Driver) Name() string                 { return h.name }
func (h *MPU6050Driver) Connection() gobot.Connection { return h.connection.(gobot.Connection) }

// Start writes initialization bytes and reads from adaptor
// using specified interval to accelerometer andtemperature data
func (h *MPU6050Driver) Start() (errs []error) {
	if err := h.initialize(); err != nil {
		return []error{err}
	}

	go func() {
		for {
			if err := h.connection.I2cWrite(mpu6050Address, []byte{MPU6050_RA_ACCEL_XOUT_H}); err != nil {
				gobot.Publish(h.Event(Error), err)
				continue
			}

			ret, err := h.connection.I2cRead(mpu6050Address, 14)
			if err != nil {
				gobot.Publish(h.Event(Error), err)
				continue
			}
			buf := bytes.NewBuffer(ret)
			binary.Read(buf, binary.BigEndian, &h.Accelerometer)
			binary.Read(buf, binary.BigEndian, &h.Temperature)
			binary.Read(buf, binary.BigEndian, &h.Gyroscope)
			h.convertToCelsius()
			<-time.After(h.interval)
		}
	}()
	return
}

// Halt returns true if devices is halted successfully
func (h *MPU6050Driver) Halt() (errs []error) { return }

func (h *MPU6050Driver) initialize() (err error) {
	if err = h.connection.I2cStart(mpu6050Address); err != nil {
		return
	}

	// setClockSource
	if err = h.connection.I2cWrite(mpu6050Address, []byte{MPU6050_RA_PWR_MGMT_1,
		MPU6050_PWR1_CLKSEL_BIT,
		MPU6050_PWR1_CLKSEL_LENGTH,
		MPU6050_CLOCK_PLL_XGYRO}); err != nil {
		return
	}

	// setFullScaleGyroRange
  if err = h.connection.I2cWrite(mpu6050Address, []byte{MPU6050_RA_GYRO_CONFIG,
    MPU6050_GCONFIG_FS_SEL_BIT,
    MPU6050_GCONFIG_FS_SEL_LENGTH,
    MPU6050_GYRO_FS_250}); err != nil {
    return
  }

	// setFullScaleAccelRange
	if err = h.connection.I2cWrite(mpu6050Address, []byte{MPU6050_RA_ACCEL_CONFIG,
		MPU6050_ACONFIG_AFS_SEL_BIT,
		MPU6050_ACONFIG_AFS_SEL_LENGTH,
		MPU6050_ACCEL_FS_2}); err != nil {
		return
	}

	// setSleepEnabled
	if err = h.connection.I2cWrite(mpu6050Address, []byte{MPU6050_RA_PWR_MGMT_1,
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
