package curie

import (
	"errors"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/firmata"
)

const (
	CURIE_IMU              = 0x11
	CURIE_IMU_READ_ACCEL   = 0x00
	CURIE_IMU_READ_GYRO    = 0x01
	CURIE_IMU_READ_TEMP    = 0x02
	CURIE_IMU_SHOCK_DETECT = 0x03
	CURIE_IMU_STEP_COUNTER = 0x04
	CURIE_IMU_TAP_DETECT   = 0x05
	CURIE_IMU_READ_MOTION  = 0x06
)

// AccelerometerData is what gets returned with the "Accelerometer" event.
type AccelerometerData struct {
	X int16
	Y int16
	Z int16
}

// GyroscopeData is what gets returned with the "Gyroscope" event.
type GyroscopeData struct {
	X int16
	Y int16
	Z int16
}

// ShockData is what gets returned with the "Shock" event.
type ShockData struct {
	Axis      byte
	Direction byte
}

// TapData is what gets returned with the "Tap" event.
type TapData struct {
	Axis      byte
	Direction byte
}

// MotionData is what gets returned with the "Motion" event.
type MotionData struct {
	AX int16
	AY int16
	AZ int16
	GX int16
	GY int16
	GZ int16
}

// IMUDriver represents the IMU that is built-in to the Curie
type IMUDriver struct {
	name       string
	connection *firmata.Adaptor
	gobot.Eventer
}

// NewIMUDriver returns a new IMUDriver
func NewIMUDriver(a *firmata.Adaptor) *IMUDriver {
	imu := &IMUDriver{
		name:       gobot.DefaultName("CurieIMU"),
		connection: a,
		Eventer:    gobot.NewEventer(),
	}

	return imu
}

// Start starts up the IMUDriver
func (imu *IMUDriver) Start() (err error) {
	imu.connection.On("SysexResponse", func(res interface{}) {
		data := res.([]byte)
		imu.handleEvent(data)
	})
	return
}

// Halt stops the IMUDriver
func (imu *IMUDriver) Halt() (err error) {
	return
}

// Name returns the IMUDriver's name
func (imu *IMUDriver) Name() string { return imu.name }

// SetName sets the IMUDriver'ss name
func (imu *IMUDriver) SetName(n string) { imu.name = n }

// Connection returns the IMUDriver's Connection
func (imu *IMUDriver) Connection() gobot.Connection { return imu.connection }

// ReadAccelerometer calls the Curie's built-in accelerometer. The result will
// be returned by the Sysex response message
func (imu *IMUDriver) ReadAccelerometer() error {
	return imu.connection.WriteSysex([]byte{CURIE_IMU, CURIE_IMU_READ_ACCEL})
}

// ReadGyroscope calls the Curie's built-in gyroscope. The result will
// be returned by the Sysex response message
func (imu *IMUDriver) ReadGyroscope() error {
	return imu.connection.WriteSysex([]byte{CURIE_IMU, CURIE_IMU_READ_GYRO})
}

// ReadTemperature calls the Curie's built-in temperature sensor.
// The result will be returned by the Sysex response message
func (imu *IMUDriver) ReadTemperature() error {
	return imu.connection.WriteSysex([]byte{CURIE_IMU, CURIE_IMU_READ_TEMP})
}

// EnableShockDetection turns on/off the Curie's built-in shock detection.
// The result will be returned by the Sysex response message
func (imu *IMUDriver) EnableShockDetection(detect bool) error {
	var d byte
	if detect {
		d = 1
	}
	return imu.connection.WriteSysex([]byte{CURIE_IMU, CURIE_IMU_SHOCK_DETECT, d})
}

// EnableStepCounter turns on/off the Curie's built-in step counter.
// The result will be returned by the Sysex response message
func (imu *IMUDriver) EnableStepCounter(count bool) error {
	var c byte
	if count {
		c = 1
	}
	return imu.connection.WriteSysex([]byte{CURIE_IMU, CURIE_IMU_STEP_COUNTER, c})
}

// EnableTapDetection turns on/off the Curie's built-in tap detection.
// The result will be returned by the Sysex response message
func (imu *IMUDriver) EnableTapDetection(detect bool) error {
	var d byte
	if detect {
		d = 1
	}
	return imu.connection.WriteSysex([]byte{CURIE_IMU, CURIE_IMU_TAP_DETECT, d})
}

// ReadMotion calls the Curie's built-in accelerometer & gyroscope.
// The result will be returned by the Sysex response message
func (imu *IMUDriver) ReadMotion() error {
	return imu.connection.WriteSysex([]byte{CURIE_IMU, CURIE_IMU_READ_MOTION})
}

func (imu *IMUDriver) handleEvent(data []byte) (err error) {
	if data[1] == CURIE_IMU {
		switch data[2] {
		case CURIE_IMU_READ_ACCEL:
			val, err := parseAccelerometerData(data)
			if err == nil {
				imu.Publish("Accelerometer", val)
			}
		case CURIE_IMU_READ_GYRO:
			val, err := parseGyroscopeData(data)
			if err == nil {
				imu.Publish("Gyroscope", val)
			}
		case CURIE_IMU_READ_TEMP:
			val, err := parseTemperatureData(data)
			if err == nil {
				imu.Publish("Temperature", val)
			}
		case CURIE_IMU_SHOCK_DETECT:
			val, err := parseShockData(data)
			if err == nil {
				imu.Publish("Shock", val)
			}
		case CURIE_IMU_STEP_COUNTER:
			val, err := parseStepData(data)
			if err == nil {
				imu.Publish("Steps", val)
			}
		case CURIE_IMU_TAP_DETECT:
			val, err := parseTapData(data)
			if err == nil {
				imu.Publish("Tap", val)
			}
		case CURIE_IMU_READ_MOTION:
			val, err := parseMotionData(data)
			if err == nil {
				imu.Publish("Motion", val)
			}
		}
	}
	return
}

func parseAccelerometerData(data []byte) (*AccelerometerData, error) {
	if len(data) < 9 {
		return nil, errors.New("Invalid data")
	}
	x := int16(uint16(data[3]) | uint16(data[4])<<7)
	y := int16(uint16(data[5]) | uint16(data[6])<<7)
	z := int16(uint16(data[7]) | uint16(data[8])<<7)

	res := &AccelerometerData{X: x, Y: y, Z: z}
	return res, nil
}

func parseGyroscopeData(data []byte) (*GyroscopeData, error) {
	if len(data) < 9 {
		return nil, errors.New("Invalid data")
	}
	x := int16(uint16(data[3]) | uint16(data[4])<<7)
	y := int16(uint16(data[5]) | uint16(data[6])<<7)
	z := int16(uint16(data[7]) | uint16(data[8])<<7)

	res := &GyroscopeData{X: x, Y: y, Z: z}
	return res, nil
}

func parseTemperatureData(data []byte) (float32, error) {
	if len(data) < 8 {
		return 0, errors.New("Invalid data")
	}
	t1 := int16(uint16(data[3]) | uint16(data[4])<<7)
	t2 := int16(uint16(data[5]) | uint16(data[6])<<7)

	res := (float32(t1+(t2*8)) / 512.0) + 23.0
	return res, nil
}

func parseShockData(data []byte) (*ShockData, error) {
	if len(data) < 6 {
		return nil, errors.New("Invalid data")
	}

	res := &ShockData{Axis: data[3], Direction: data[4]}
	return res, nil
}

func parseStepData(data []byte) (int16, error) {
	if len(data) < 6 {
		return 0, errors.New("Invalid data")
	}

	res := int16(uint16(data[3]) | uint16(data[4])<<7)
	return res, nil
}

func parseTapData(data []byte) (*TapData, error) {
	if len(data) < 6 {
		return nil, errors.New("Invalid data")
	}

	res := &TapData{Axis: data[3], Direction: data[4]}
	return res, nil
}

func parseMotionData(data []byte) (*MotionData, error) {
	if len(data) < 16 {
		return nil, errors.New("Invalid data")
	}
	ax := int16(uint16(data[3]) | uint16(data[4])<<7)
	ay := int16(uint16(data[5]) | uint16(data[6])<<7)
	az := int16(uint16(data[7]) | uint16(data[8])<<7)

	gx := int16(uint16(data[9]) | uint16(data[10])<<7)
	gy := int16(uint16(data[11]) | uint16(data[12])<<7)
	gz := int16(uint16(data[13]) | uint16(data[14])<<7)

	res := &MotionData{AX: ax, AY: ay, AZ: az, GX: gx, GY: gy, GZ: gz}
	return res, nil
}
