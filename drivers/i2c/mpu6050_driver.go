package i2c

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"time"
)

const (
	mpu6050Debug                = false
	mpu6050DefaultAddress       = 0x68
	mpu6050EarthStandardGravity = 9.80665 // [m/s²] standard gravity (pole: 9.834,  equator: 9.764)
)

type (
	MPU6050DlpfConfig      uint8
	MPU6050FrameSyncConfig uint8
	MPU6050GyroFsConfig    uint8
	MPU6050AccelFsConfig   uint8
	MPU6050Pwr1ClockConfig uint8
)

const (
	mpu6050Reg_GeneralConfig   = 0x1A // external frame synchronization and digital low pass filter
	mpu6050Reg_GyroConfig      = 0x1B // self test and full scale range
	mpu6050Reg_AccelConfig     = 0x1C // self test and full scale range
	mpu6050Reg_AccelXoutH      = 0x3B // first data register
	mpu6050Reg_SignalPathReset = 0x68
	mpu6050Reg_PwrMgmt1        = 0x6B

	MPU6050General_Dlpf260Hz MPU6050DlpfConfig = 0x00
	MPU6050General_Dlpf184Hz MPU6050DlpfConfig = 0x01
	MPU6050General_Dlpf94Hz  MPU6050DlpfConfig = 0x02
	MPU6050General_Dlpf44Hz  MPU6050DlpfConfig = 0x03
	MPU6050General_Dlpf21Hz  MPU6050DlpfConfig = 0x04
	MPU6050General_Dlpf10Hz  MPU6050DlpfConfig = 0x05
	MPU6050General_Dlpf5Hz   MPU6050DlpfConfig = 0x06

	MPU6050General_FrameSyncDisabled MPU6050FrameSyncConfig = 0x00
	MPU6050General_FrameSyncTemp     MPU6050FrameSyncConfig = 0x01
	MPU6050General_FrameSyncGyroX    MPU6050FrameSyncConfig = 0x02
	MPU6050General_FrameSyncGyroY    MPU6050FrameSyncConfig = 0x03
	MPU6050General_FrameSyncGyroZ    MPU6050FrameSyncConfig = 0x04
	MPU6050General_FrameSyncAccelX   MPU6050FrameSyncConfig = 0x05
	MPU6050General_FrameSyncAccelY   MPU6050FrameSyncConfig = 0x06
	MPU6050General_FrameSyncAccelZ   MPU6050FrameSyncConfig = 0x07

	MPU6050Gyro_FsSel250dps  MPU6050GyroFsConfig = 0x00 // +/- 250 °/s
	MPU6050Gyro_FsSel500dps  MPU6050GyroFsConfig = 0x01 // +/- 500 °/s
	MPU6050Gyro_FsSel1000dps MPU6050GyroFsConfig = 0x02 // +/- 1000 °/s
	MPU6050Gyro_FsSel2000dps MPU6050GyroFsConfig = 0x03 // +/- 2000 °/s

	MPU6050Accel_AFsSel2g  MPU6050AccelFsConfig = 0x00 // +/- 2 g
	MPU6050Accel_AFsSel4g  MPU6050AccelFsConfig = 0x01 // +/- 4 g
	MPU6050Accel_AFsSel8g  MPU6050AccelFsConfig = 0x02 // +/- 8 g
	MPU6050Accel_AFsSel16g MPU6050AccelFsConfig = 0x03 // +/- 16 g

	mpu6050SignalReset_TempBit  = 0x01
	mpu6050SignalReset_AccelBit = 0x02
	mpu6050SignalReset_GyroBit  = 0x04

	MPU6050Pwr1_ClockIntern8G  MPU6050Pwr1ClockConfig = 0x00 // internal 8GHz
	MPU6050Pwr1_ClockPllXGyro  MPU6050Pwr1ClockConfig = 0x01 // PLL with X axis gyroscope reference
	MPU6050Pwr1_ClockPllYGyro  MPU6050Pwr1ClockConfig = 0x02 // PLL with Y axis gyroscope reference
	MPU6050Pwr1_ClockPllZGyro  MPU6050Pwr1ClockConfig = 0x03 // PLL with Z axis gyroscope reference
	MPU6050Pwr1_ClockPllExt32K MPU6050Pwr1ClockConfig = 0x04 // PLL with external 32.768kHz reference
	MPU6050Pwr1_ClockPllExt19M MPU6050Pwr1ClockConfig = 0x05 // PLL with external 19.2MHz reference
	MPU6050Pwr1_ClockStop      MPU6050Pwr1ClockConfig = 0x07 // Stops the clock and keeps the timing generator in reset

	mpu6050Pwr1_SleepOnBit     = 0x40 // put into low power sleep mode
	mpu6050Pwr1_DeviceResetBit = 0x80
)

type MPU6050ThreeDData struct {
	X float64
	Y float64
	Z float64
}

// MPU6050Driver is a Gobot Driver for an MPU6050 I2C Accelerometer/Gyroscope/Temperature sensor.
//
// This driver was tested with Tinkerboard & Digispark adaptor and a MPU6050 breakout board GY-521,
// available from various distributors.
//
// datasheet:
// https://product.tdk.com/system/files/dam/doc/product/sensor/mortion-inertial/imu/data_sheet/mpu-6000-datasheet1.pdf
//
// reference implementations:
// * https://github.com/adafruit/Adafruit_CircuitPython_MPU6050
// * https://github.com/ElectronicCats/mpu6050
type MPU6050Driver struct {
	*Driver
	Accelerometer MPU6050ThreeDData
	Gyroscope     MPU6050ThreeDData
	Temperature   float64
	dlpf          MPU6050DlpfConfig
	frameSync     MPU6050FrameSyncConfig
	accelFs       MPU6050AccelFsConfig
	gyroFs        MPU6050GyroFsConfig
	clock         MPU6050Pwr1ClockConfig
	gravity       float64 // set to 1.0 leads to [g]
}

// mpu6050AccelGain in 1/g
var mpu6050AccelGain = map[MPU6050AccelFsConfig]uint16{
	MPU6050Accel_AFsSel2g:  16384,
	MPU6050Accel_AFsSel4g:  8192,
	MPU6050Accel_AFsSel8g:  4096,
	MPU6050Accel_AFsSel16g: 2028,
}

// mpu6050GyroGain in s/°
var mpu6050GyroGain = map[MPU6050GyroFsConfig]float64{
	MPU6050Gyro_FsSel250dps:  131.0,
	MPU6050Gyro_FsSel500dps:  65.5,
	MPU6050Gyro_FsSel1000dps: 32.8,
	MPU6050Gyro_FsSel2000dps: 16.4,
}

// NewMPU6050Driver creates a new Gobot Driver for an MPU6050 I2C Accelerometer/Gyroscope/Temperature sensor.
//
// Params:
//
//	conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
func NewMPU6050Driver(a Connector, options ...func(Config)) *MPU6050Driver {
	m := &MPU6050Driver{
		Driver:    NewDriver(a, "MPU6050", mpu6050DefaultAddress),
		dlpf:      MPU6050General_Dlpf260Hz,
		frameSync: MPU6050General_FrameSyncDisabled,
		accelFs:   MPU6050Accel_AFsSel2g,
		gyroFs:    MPU6050Gyro_FsSel250dps,
		clock:     MPU6050Pwr1_ClockPllXGyro,
		gravity:   mpu6050EarthStandardGravity,
	}
	m.afterStart = m.initialize

	for _, option := range options {
		option(m)
	}

	// TODO: add commands to API
	return m
}

// WithMPU6050DigitalFilter option sets the digital low pass filter bandwidth frequency.
// Valid settings are of type "MPU6050DlpfConfig"
func WithMPU6050DigitalFilter(val MPU6050DlpfConfig) func(Config) {
	return func(c Config) {
		if d, ok := c.(*MPU6050Driver); ok {
			d.dlpf = val
		} else if mpu6050Debug {
			log.Printf("Trying to set digital low pass filter for non-MPU6050Driver %v", c)
		}
	}
}

// WithMPU6050FrameSync option sets the external frame synchronization.
// Valid settings are of type "MPU6050FrameSyncConfig"
func WithMPU6050FrameSync(val MPU6050FrameSyncConfig) func(Config) {
	return func(c Config) {
		if d, ok := c.(*MPU6050Driver); ok {
			d.frameSync = val
		} else if mpu6050Debug {
			log.Printf("Trying to set external frame synchronization for non-MPU6050Driver %v", c)
		}
	}
}

// WithMPU6050AccelFullScaleRange option sets the full scale range for the accelerometer.
// Valid settings are of type "MPU6050AccelFsConfig"
func WithMPU6050AccelFullScaleRange(val MPU6050AccelFsConfig) func(Config) {
	return func(c Config) {
		if d, ok := c.(*MPU6050Driver); ok {
			d.accelFs = val
		} else if mpu6050Debug {
			log.Printf("Trying to set full scale range of accelerometer for non-MPU6050Driver %v", c)
		}
	}
}

// WithMPU6050GyroFullScaleRange option sets the full scale range for the gyroscope.
// Valid settings are of type "MPU6050GyroFsConfig"
func WithMPU6050GyroFullScaleRange(val MPU6050GyroFsConfig) func(Config) {
	return func(c Config) {
		if d, ok := c.(*MPU6050Driver); ok {
			d.gyroFs = val
		} else if mpu6050Debug {
			log.Printf("Trying to set full scale range of gyroscope for non-MPU6050Driver %v", c)
		}
	}
}

// WithMPU6050ClockSource option sets the clock source.
// Valid settings are of type "MPU6050Pwr1ClockConfig"
func WithMPU6050ClockSource(val MPU6050Pwr1ClockConfig) func(Config) {
	return func(c Config) {
		if d, ok := c.(*MPU6050Driver); ok {
			d.clock = val
		} else if mpu6050Debug {
			log.Printf("Trying to set clock source for non-MPU6050Driver %v", c)
		}
	}
}

// WithMPU6050Gravity option sets the gravity in [m/s²/g].
// Useful settings are "1.0" (to use unit "g") or a value between 9.834 (pole) and 9.764 (equator)
func WithMPU6050Gravity(val float64) func(Config) {
	return func(c Config) {
		if d, ok := c.(*MPU6050Driver); ok {
			d.gravity = val
		} else if mpu6050Debug {
			log.Printf("Trying to set gravity for non-MPU6050Driver %v", c)
		}
	}
}

// GetData fetches the latest data from the MPU6050
func (m *MPU6050Driver) GetData() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	data := make([]byte, 14)
	if err := m.connection.ReadBlockData(mpu6050Reg_AccelXoutH, data); err != nil {
		return err
	}

	var accel struct {
		X int16
		Y int16
		Z int16
	}
	var temp int16
	var gyro struct {
		X int16
		Y int16
		Z int16
	}

	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &accel); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &temp); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &gyro); err != nil {
		return err
	}

	ag := float64(mpu6050AccelGain[m.accelFs]) / m.gravity
	m.Accelerometer.X = float64(accel.X) / ag
	m.Accelerometer.Y = float64(accel.Y) / ag
	m.Accelerometer.Z = float64(accel.Z) / ag

	m.Temperature = float64(temp)/340 + 36.53

	gg := mpu6050GyroGain[m.gyroFs]
	m.Gyroscope.X = float64(gyro.X) / gg
	m.Gyroscope.Y = float64(gyro.Y) / gg
	m.Gyroscope.Z = float64(gyro.Z) / gg

	return nil
}

func (m *MPU6050Driver) waitForReset() error {
	wait := 100 * time.Millisecond
	start := time.Now()
	for {
		if time.Since(start) > wait {
			return fmt.Errorf("timeout on wait for reset is done")
		}
		if val, err := m.connection.ReadByteData(mpu6050Reg_PwrMgmt1); (val&mpu6050Pwr1_DeviceResetBit == 0) && (err == nil) {
			return nil
		}
		time.Sleep(wait / 10)
	}
}

func (m *MPU6050Driver) initialize() error {
	// reset device and wait for reset is finished
	if err := m.connection.WriteByteData(mpu6050Reg_PwrMgmt1, mpu6050Pwr1_DeviceResetBit); err != nil {
		return err
	}
	if err := m.waitForReset(); err != nil {
		return err
	}

	// reset signal path register
	reset := uint8(mpu6050SignalReset_TempBit | mpu6050SignalReset_AccelBit | mpu6050SignalReset_GyroBit)
	if err := m.connection.WriteByteData(mpu6050Reg_SignalPathReset, reset); err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)

	// configure digital filter bandwidth and external frame synchronization (bits 3...5 are used)
	generalConf := uint8(m.dlpf) | uint8(m.frameSync)<<3
	if err := m.connection.WriteByteData(mpu6050Reg_GeneralConfig, generalConf); err != nil {
		return err
	}

	// set full scale range of gyroscope (bits 3 and 4 are used)
	if err := m.connection.WriteByteData(mpu6050Reg_GyroConfig, uint8(m.gyroFs)<<3); err != nil {
		return err
	}

	// set full scale range of accelerometer (bits 3 and 4 are used)
	if err := m.connection.WriteByteData(mpu6050Reg_AccelConfig, uint8(m.accelFs)<<3); err != nil {
		return err
	}

	// set clock source and reset sleep
	pwr1 := uint8(m.clock) & ^uint8(mpu6050Pwr1_SleepOnBit)

	return m.connection.WriteByteData(mpu6050Reg_PwrMgmt1, pwr1)
}
