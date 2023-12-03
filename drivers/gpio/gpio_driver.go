package gpio

import (
	"errors"
	"fmt"
	"log"
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
	// ErrDigitalWriteUnsupported is the error resulting when a driver attempts to use
	// hardware capabilities which a connection does not support
	ErrDigitalWriteUnsupported = errors.New("DigitalWrite is not supported by this platform")
	// ErrDigitalReadUnsupported is the error resulting when a driver attempts to use
	// hardware capabilities which a connection does not support
	ErrDigitalReadUnsupported = errors.New("DigitalRead is not supported by this platform")
)

const (
	// Error event
	Error = "error"
	// ButtonRelease event
	ButtonRelease = "release"
	// ButtonPush event
	ButtonPush = "push"
	// MotionDetected event
	MotionDetected = "motion-detected"
	// MotionStopped event
	MotionStopped = "motion-stopped"
)

// PwmWriter interface represents an Adaptor which has Pwm capabilities
type PwmWriter interface {
	PwmWrite(pin string, val byte) error
}

// ServoWriter interface represents an Adaptor which has Servo capabilities
type ServoWriter interface {
	ServoWrite(pin string, val byte) error
}

// DigitalWriter interface represents an Adaptor which has DigitalWrite capabilities
type DigitalWriter interface {
	DigitalWrite(pin string, val byte) error
}

// DigitalReader interface represents an Adaptor which has DigitalRead capabilities
type DigitalReader interface {
	DigitalRead(pin string) (val int, err error)
}

// optionApplier needs to be implemented by each configurable option type
type optionApplier interface {
	apply(cfg *configuration)
}

// configuration contains all changeable attributes of the driver.
type configuration struct {
	name string
	pin  string
}

// nameOption is the type for applying another name to the configuration
type nameOption string

// pinOption is the type for applying a pin to the configuration
type pinOption string

// Driver implements the interface gobot.Driver.
type driver struct {
	driverCfg  *configuration
	connection gobot.Adaptor
	afterStart func() error
	beforeHalt func() error
	gobot.Commander
	mutex *sync.Mutex // mutex often needed to ensure that write-read sequences are not interrupted
}

// newDriver creates a new generic and basic gpio gobot driver.
//
// Supported options:
//
//	"WithName"
//	"withPin"
func newDriver(a gobot.Adaptor, name string, opts ...interface{}) *driver {
	d := &driver{
		driverCfg:  &configuration{name: gobot.DefaultName(name)},
		connection: a,
		afterStart: func() error { return nil },
		beforeHalt: func() error { return nil },
		Commander:  gobot.NewCommander(),
		mutex:      &sync.Mutex{},
	}

	for _, opt := range opts {
		switch o := opt.(type) {
		case optionApplier:
			o.apply(d.driverCfg)
		default:
			panic(fmt.Sprintf("'%s' can not be applied on '%s'", opt, d.driverCfg.name))
		}
	}

	return d
}

// WithName is used to replace the default name of the driver.
func WithName(name string) optionApplier {
	return nameOption(name)
}

// withPin is used to add a pin to the driver. Only one pin can be linked.
// This option is not available outside gpio package.
func withPin(pin string) optionApplier {
	return pinOption(pin)
}

// Name returns the name of the gpio device.
func (d *driver) Name() string {
	return d.driverCfg.name
}

// SetName sets the name of the gpio device.
// Deprecated: Please use option [gpio.WithName] instead.
func (d *driver) SetName(name string) {
	WithName(name).apply(d.driverCfg)
}

// Pin returns the pin associated with the driver.
func (d *driver) Pin() string {
	return d.driverCfg.pin
}

// Connection returns the connection of the gpio device.
func (d *driver) Connection() gobot.Connection {
	if conn, ok := d.connection.(gobot.Connection); ok {
		return conn
	}

	log.Printf("%s has no gobot connection\n", d.driverCfg.name)
	return nil
}

// Start initializes the gpio device.
func (d *driver) Start() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// currently there is nothing to do here for the driver

	return d.afterStart()
}

// Halt halts the gpio device.
func (d *driver) Halt() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// currently there is nothing to do after halt for the driver

	return d.beforeHalt()
}

// digitalRead is a helper function with check that the connection implements DigitalReader
func (d *driver) digitalRead(pin string) (int, error) {
	if reader, ok := d.connection.(DigitalReader); ok {
		return reader.DigitalRead(pin)
	}

	return 0, ErrDigitalReadUnsupported
}

// digitalWrite is a helper function with check that the connection implements DigitalWriter
func (d *driver) digitalWrite(pin string, val byte) error {
	if writer, ok := d.connection.(DigitalWriter); ok {
		return writer.DigitalWrite(pin, val)
	}

	return ErrDigitalWriteUnsupported
}

// pwmWrite is a helper function with check that the connection implements PwmWriter
func (d *driver) pwmWrite(pin string, level byte) error {
	if writer, ok := d.connection.(PwmWriter); ok {
		return writer.PwmWrite(pin, level)
	}

	return ErrPwmWriteUnsupported
}

// servoWrite is a helper function with check that the connection implements ServoWriter
func (d *driver) servoWrite(pin string, level byte) error {
	if writer, ok := d.connection.(ServoWriter); ok {
		return writer.ServoWrite(pin, level)
	}

	return ErrServoWriteUnsupported
}

func (o nameOption) String() string {
	return "name option for digital drivers"
}

func (o pinOption) String() string {
	return "pin option for digital drivers"
}

// apply change the name in the configuration.
func (o nameOption) apply(c *configuration) {
	c.name = string(o)
}

// apply change the pins list of the configuration.
func (o pinOption) apply(c *configuration) {
	c.pin = string(o)
}
