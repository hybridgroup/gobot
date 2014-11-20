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
	m := &MPU6050Driver{
		Driver: *gobot.NewDriver(
			name,
			"MPU6050Driver",
			a.(gobot.AdaptorInterface),
		),
	}
	m.AddEvent("error")
	return m
}

// adaptor returns MPU6050 adaptor
func (h *MPU6050Driver) adaptor() I2cInterface {
	return h.Adaptor().(I2cInterface)
}

// Start writes initialization bytes and reads from adaptor
// using specified interval to accelerometer andtemperature data
func (h *MPU6050Driver) Start() (errs []error) {
	if err := h.initialize(); err != nil {
		return []error{err}
	}

	gobot.Every(h.Interval(), func() {
		if err := h.adaptor().I2cWrite([]byte{MPU6050_RA_ACCEL_XOUT_H}); err != nil {
			gobot.Publish(h.Event("error"), err)
			return
		}

		ret, err := h.adaptor().I2cRead(14)
		if err != nil {
			gobot.Publish(h.Event("error"), err)
			return
		}
		buf := bytes.NewBuffer(ret)
		binary.Read(buf, binary.BigEndian, &h.Accelerometer)
		binary.Read(buf, binary.BigEndian, &h.Gyroscope)
		binary.Read(buf, binary.BigEndian, &h.Temperature)
	})
	return
}

// Halt returns true if devices is halted successfully
func (h *MPU6050Driver) Halt() (errs []error) { return }

func (h *MPU6050Driver) initialize() (err error) {
	if err = h.adaptor().I2cStart(0x68); err != nil {
		return
	}

	// setClockSource
	if err = h.adaptor().I2cWrite([]byte{MPU6050_RA_PWR_MGMT_1,
		MPU6050_PWR1_CLKSEL_BIT,
		MPU6050_PWR1_CLKSEL_LENGTH,
		MPU6050_CLOCK_PLL_XGYRO}); err != nil {
		return
	}

	// setFullScaleGyroRange
	if err = h.adaptor().I2cWrite([]byte{MPU6050_GYRO_FS_250,
		MPU6050_RA_GYRO_CONFIG,
		MPU6050_GCONFIG_FS_SEL_LENGTH,
		MPU6050_GCONFIG_FS_SEL_BIT}); err != nil {
		return
	}

	// setFullScaleAccelRange
	if err = h.adaptor().I2cWrite([]byte{MPU6050_RA_ACCEL_CONFIG,
		MPU6050_ACONFIG_AFS_SEL_BIT,
		MPU6050_ACONFIG_AFS_SEL_LENGTH,
		MPU6050_ACCEL_FS_2}); err != nil {
		return
	}

	// setSleepEnabled
	if err = h.adaptor().I2cWrite([]byte{MPU6050_RA_PWR_MGMT_1,
		MPU6050_PWR1_SLEEP_BIT,
		0}); err != nil {
		return
	}

	return nil
}
