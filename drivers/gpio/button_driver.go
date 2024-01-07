package gpio

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
)

// buttonOptionApplier needs to be implemented by each configurable option type
type buttonOptionApplier interface {
	apply(cfg *buttonConfiguration)
}

// buttonConfiguration contains all changeable attributes of the driver.
type buttonConfiguration struct {
	readInterval time.Duration
	defaultState int
}

// buttonReadIntervalOption is the type for applying another read interval to the configuration
type buttonReadIntervalOption time.Duration

// buttonDefaultStateOption is the type for applying another default state to the configuration
type buttonDefaultStateOption int

// ButtonDriver Represents a digital Button
type ButtonDriver struct {
	*driver
	buttonCfg *buttonConfiguration
	gobot.Eventer
	active bool
	halt   chan struct{}
}

// NewButtonDriver returns a driver for a button with a polling interval for changed state of 10 milliseconds,
// given a DigitalReader and pin.
//
// Supported options:
//
//	"WithName"
//	"WithButtonPollInterval"
func NewButtonDriver(a DigitalReader, pin string, opts ...interface{}) *ButtonDriver {
	//nolint:forcetypeassert // no error return value, so there is no better way
	d := &ButtonDriver{
		driver:    newDriver(a.(gobot.Connection), "Button", withPin(pin)),
		buttonCfg: &buttonConfiguration{readInterval: 10 * time.Millisecond, defaultState: 0},
	}
	d.afterStart = d.initialize
	d.beforeHalt = d.shutdown

	for _, opt := range opts {
		switch o := opt.(type) {
		case optionApplier:
			o.apply(d.driverCfg)
		case buttonOptionApplier:
			o.apply(d.buttonCfg)
		case time.Duration:
			// TODO this is only for backward compatibility and will be removed after version 2.x
			d.buttonCfg.readInterval = o
		default:
			panic(fmt.Sprintf("'%s' can not be applied on '%s'", opt, d.driverCfg.name))
		}
	}

	return d
}

// WithButtonPollInterval change the asynchronous cyclic reading interval from default 10ms to the given value.
func WithButtonPollInterval(interval time.Duration) buttonOptionApplier {
	return buttonReadIntervalOption(interval)
}

// WithButtonDefaultState change the default state from default 0 to the given value.
func WithButtonDefaultState(s int) buttonOptionApplier {
	return buttonDefaultStateOption(s)
}

// Active gets the current state
func (d *ButtonDriver) Active() bool {
	// ensure that read and write can not interfere
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.active
}

// SetDefaultState for the next start.
// Deprecated: Please use option [gpio.WithButtonDefaultState] instead.
func (d *ButtonDriver) SetDefaultState(s int) {
	// ensure that read and write can not interfere
	d.mutex.Lock()
	defer d.mutex.Unlock()

	WithButtonDefaultState(s).apply(d.buttonCfg)
}

// initialize the ButtonDriver and polls the state of the button at the given interval.
//
// Emits the Events:
//
//	Push int - On button push
//	Release int - On button release
//	Error error - On button error
func (d *ButtonDriver) initialize() error {
	if d.buttonCfg.readInterval == 0 {
		return fmt.Errorf("the read interval for button needs to be greater than zero")
	}

	d.Eventer = gobot.NewEventer()
	d.AddEvent(ButtonPush)
	d.AddEvent(ButtonRelease)
	d.AddEvent(Error)

	d.halt = make(chan struct{})

	state := d.buttonCfg.defaultState

	go func() {
		for {
			select {
			case <-time.After(d.buttonCfg.readInterval):
				newValue, err := d.digitalRead(d.driverCfg.pin)
				if err != nil {
					d.Publish(Error, err)
				} else if newValue != state && newValue != -1 {
					state = newValue
					d.update(newValue)
				}
			case <-d.halt:
				return
			}
		}
	}()
	return nil
}

func (d *ButtonDriver) shutdown() error {
	if d.buttonCfg.readInterval == 0 || d.halt == nil {
		// cyclic reading deactivated
		return nil
	}

	close(d.halt) // broadcast halt, also to the test
	return nil
}

func (d *ButtonDriver) update(newValue int) {
	// ensure that read and write can not interfere
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if newValue != d.buttonCfg.defaultState {
		d.active = true
		d.Publish(ButtonPush, newValue)
	} else {
		d.active = false
		d.Publish(ButtonRelease, newValue)
	}
}

func (o buttonReadIntervalOption) String() string {
	return "read interval option for buttons"
}

func (o buttonDefaultStateOption) String() string {
	return "default state option for buttons"
}

func (o buttonReadIntervalOption) apply(cfg *buttonConfiguration) {
	cfg.readInterval = time.Duration(o)
}

func (o buttonDefaultStateOption) apply(cfg *buttonConfiguration) {
	cfg.defaultState = int(o)
}
