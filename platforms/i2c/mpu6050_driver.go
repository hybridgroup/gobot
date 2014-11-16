package i2c

import (
	"bytes"
	"encoding/binary"
	"github.com/hybridgroup/gobot"
)

var _ gobot.DriverInterface = (*MPU6050Driver)(nil)

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

type ThreeDData struct {
	X int16
	Y int16
	Z int16
}

type MPU6050Driver struct {
	gobot.Driver
	Accelerometer ThreeDData
	Gyroscope     ThreeDData
	Temperature   int16
}

// NewMPU6050Driver creates a new driver with specified name and i2c interface
func NewMPU6050Driver(a I2cInterface, name string) *MPU6050Driver {
	return &MPU6050Driver{
		Driver: *gobot.NewDriver(
			name,
			"MPU6050Driver",
			a.(gobot.AdaptorInterface),
		),
	}
}

// adaptor returns MPU6050 adaptor
func (h *MPU6050Driver) adaptor() I2cInterface {
	return h.Adaptor().(I2cInterface)
}

// Start writes initialization bytes and reads from adaptor
// using specified interval to accelerometer andtemperature data
func (h *MPU6050Driver) Start() error {
	h.initialize()

	gobot.Every(h.Interval(), func() {
		h.adaptor().I2cWrite([]byte{MPU6050_RA_ACCEL_XOUT_H})

		ret := h.adaptor().I2cRead(14)
		buf := bytes.NewBuffer(ret)
		binary.Read(buf, binary.BigEndian, &h.Accelerometer)
		binary.Read(buf, binary.BigEndian, &h.Gyroscope)
		binary.Read(buf, binary.BigEndian, &h.Temperature)
	})
	return nil
}

// Halt returns true if devices is halted successfully
func (h *MPU6050Driver) Halt() error { return nil }

func (h *MPU6050Driver) initialize() bool {
	h.adaptor().I2cStart(0x68)

	// setClockSource
	h.adaptor().I2cWrite([]byte{MPU6050_RA_PWR_MGMT_1,
		MPU6050_PWR1_CLKSEL_BIT,
		MPU6050_PWR1_CLKSEL_LENGTH,
		MPU6050_CLOCK_PLL_XGYRO})

	// setFullScaleGyroRange
	h.adaptor().I2cWrite([]byte{MPU6050_GYRO_FS_250,
		MPU6050_RA_GYRO_CONFIG,
		MPU6050_GCONFIG_FS_SEL_LENGTH,
		MPU6050_GCONFIG_FS_SEL_BIT})

	// setFullScaleAccelRange
	h.adaptor().I2cWrite([]byte{MPU6050_RA_ACCEL_CONFIG,
		MPU6050_ACONFIG_AFS_SEL_BIT,
		MPU6050_ACONFIG_AFS_SEL_LENGTH,
		MPU6050_ACCEL_FS_2})

	// setSleepEnabled
	h.adaptor().I2cWrite([]byte{MPU6050_RA_PWR_MGMT_1,
		MPU6050_PWR1_SLEEP_BIT,
		0})

	return true
}
