package gpio

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
)

// pirMotionOptionApplier needs to be implemented by each configurable option type
type pirMotionOptionApplier interface {
	apply(cfg *pirMotionConfiguration)
}

// pirMotionConfiguration contains all changeable attributes of the driver.
type pirMotionConfiguration struct {
	readInterval time.Duration
}

// pirMotionReadIntervalOption is the type for applying another read interval to the configuration
type pirMotionReadIntervalOption time.Duration

// PIRMotionDriver represents a digital Proximity Infra Red (PIR) motion detecter
//
// Supported options:
//
//	"WithName"
type PIRMotionDriver struct {
	*driver
	pirMotionCfg *pirMotionConfiguration
	gobot.Eventer
	active bool
	halt   chan struct{}
}

// NewPIRMotionDriver returns a new driver for  PIR motion sensor with a polling interval of 10 Milliseconds,
// given a DigitalReader and pin.
//
// Supported options:
//
//	"WithName"
//	"WithButtonPollInterval"
func NewPIRMotionDriver(a DigitalReader, pin string, opts ...interface{}) *PIRMotionDriver {
	//nolint:forcetypeassert // no error return value, so there is no better way
	d := &PIRMotionDriver{
		driver:       newDriver(a.(gobot.Connection), "PIRMotion", withPin(pin)),
		pirMotionCfg: &pirMotionConfiguration{readInterval: 10 * time.Millisecond},
	}
	d.afterStart = d.initialize
	d.beforeHalt = d.shutdown

	for _, opt := range opts {
		switch o := opt.(type) {
		case optionApplier:
			o.apply(d.driverCfg)
		case pirMotionOptionApplier:
			o.apply(d.pirMotionCfg)
		case time.Duration:
			// TODO this is only for backward compatibility and will be removed after version 2.x
			d.pirMotionCfg.readInterval = o
		default:
			panic(fmt.Sprintf("'%s' can not be applied on '%s'", opt, d.driverCfg.name))
		}
	}

	return d
}

// WithPIRMotionPollInterval change the asynchronous cyclic reading interval from default 10ms to the given value.
func WithPIRMotionPollInterval(interval time.Duration) pirMotionOptionApplier {
	return pirMotionReadIntervalOption(interval)
}

// Active gets the current state
func (d *PIRMotionDriver) Active() bool {
	// ensure that read and write can not interfere
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.active
}

// initialize the PIRMotionDriver and polls the state of the sensor at the given interval.
//
// Emits the Events:
//
//	MotionDetected - On motion detected
//	MotionStopped int - On motion stopped
//	Error error - On pirMotion error
//
// The PIRMotionDriver will send the MotionDetected event over and over,
// just as long as motion is still being detected.
// It will only send the MotionStopped event once, however, until
// motion starts being detected again
func (d *PIRMotionDriver) initialize() error {
	if d.pirMotionCfg.readInterval == 0 {
		return fmt.Errorf("the read interval for pirMotion needs to be greater than zero")
	}

	d.Eventer = gobot.NewEventer()
	d.AddEvent(MotionDetected)
	d.AddEvent(MotionStopped)
	d.AddEvent(Error)

	d.halt = make(chan struct{})

	go func() {
		for {
			select {
			case <-time.After(d.pirMotionCfg.readInterval):
				newValue, err := d.digitalRead(d.driverCfg.pin)
				if err != nil {
					d.Publish(Error, err)
				}
				d.update(newValue)
			case <-d.halt:
				return
			}
		}
	}()
	return nil
}

// shutdown stops polling
func (d *PIRMotionDriver) shutdown() error {
	if d.pirMotionCfg.readInterval == 0 || d.halt == nil {
		// cyclic reading deactivated
		return nil
	}

	close(d.halt) // broadcast halt, also to the test
	return nil
}

func (d *PIRMotionDriver) update(newValue int) {
	// ensure that read and write can not interfere
	d.mutex.Lock()
	defer d.mutex.Unlock()

	switch newValue {
	case 1:
		if !d.active {
			d.active = true
			d.Publish(MotionDetected, newValue)
		}
	case 0:
		if d.active {
			d.active = false
			d.Publish(MotionStopped, newValue)
		}
	}
}

func (o pirMotionReadIntervalOption) String() string {
	return "read interval option for PIR motion sensor"
}

func (o pirMotionReadIntervalOption) apply(cfg *pirMotionConfiguration) {
	cfg.readInterval = time.Duration(o)
}
