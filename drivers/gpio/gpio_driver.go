package gpio

import (
	"errors"
	"sync"

	"gobot.io/x/gobot/v2"
)

var (
	// ErrServoWriteUnsupported is the error resulting when a driver attempts to use
	// hardware capabilities which a connection does not support
	ErrServoWriteUnsupported = errors.New("ServoWrite is not supported by this platform")
	// ErrPwmWriteUnsupported is the error resulting when a driver attempts to use
	// hardware capabilities which a connection does not support
	ErrPwmWriteUnsupported = errors.New("PwmWrite is not supported by this platform")
	// ErrAnalogReadUnsupported is error resulting when a driver attempts to use
	// hardware capabilities which a connection does not support
	ErrAnalogReadUnsupported = errors.New("AnalogRead is not supported by this platform")
	// ErrDigitalWriteUnsupported is the error resulting when a driver attempts to use
	// hardware capabilities which a connection does not support
	ErrDigitalWriteUnsupported = errors.New("DigitalWrite is not supported by this platform")
	// ErrDigitalReadUnsupported is the error resulting when a driver attempts to use
	// hardware capabilities which a connection does not support
	ErrDigitalReadUnsupported = errors.New("DigitalRead is not supported by this platform")
	// ErrServoOutOfRange is the error resulting when a driver attempts to use
	// hardware capabilities which a connection does not support
	ErrServoOutOfRange = errors.New("servo angle must be between 0-180")
)

const (
	// Error event
	Error = "error"
	// ButtonRelease event
	ButtonRelease = "release"
	// ButtonPush event
	ButtonPush = "push"
	// Data event
	Data = "data"
	// Vibration event
	Vibration = "vibration"
	// MotionDetected event
	MotionDetected = "motion-detected"
	// MotionStopped event
	MotionStopped = "motion-stopped"
)

// PwmWriter interface represents an Adaptor which has Pwm capabilities
type PwmWriter interface {
	PwmWrite(string, byte) (err error)
}

// ServoWriter interface represents an Adaptor which has Servo capabilities
type ServoWriter interface {
	ServoWrite(string, byte) (err error)
}

// DigitalWriter interface represents an Adaptor which has DigitalWrite capabilities
type DigitalWriter interface {
	DigitalWrite(string, byte) (err error)
}

// DigitalReader interface represents an Adaptor which has DigitalRead capabilities
type DigitalReader interface {
	DigitalRead(string) (val int, err error)
}

// Driver implements the interface gobot.Driver.
type Driver struct {
	name       string
	connection gobot.Adaptor
	afterStart func() error
	beforeHalt func() error
	gobot.Commander
	mutex *sync.Mutex // mutex often needed to ensure that write-read sequences are not interrupted
}

// NewDriver creates a new generic and basic gpio gobot driver.
func NewDriver(a gobot.Adaptor, name string) *Driver {
	d := &Driver{
		name:       gobot.DefaultName(name),
		connection: a,
		afterStart: func() error { return nil },
		beforeHalt: func() error { return nil },
		Commander:  gobot.NewCommander(),
		mutex:      &sync.Mutex{},
	}

	return d
}

// Name returns the name of the gpio device.
func (d *Driver) Name() string {
	return d.name
}

// SetName sets the name of the gpio device.
func (d *Driver) SetName(name string) {
	d.name = name
}

// Connection returns the connection of the gpio device.
func (d *Driver) Connection() gobot.Connection {
	return d.connection.(gobot.Connection)
}

// Start initializes the gpio device.
func (d *Driver) Start() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// currently there is nothing to do here for the driver

	return d.afterStart()
}

// Halt halts the gpio device.
func (d *Driver) Halt() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// currently there is nothing to do after halt for the driver

	return d.beforeHalt()
}
