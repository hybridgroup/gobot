package curie

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/firmata"
)

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
