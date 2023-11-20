package aio

import (
	"log"
	"sync"

	"gobot.io/x/gobot/v2"
)

const (
	// Error event
	Error = "error"
	// Data event
	Data = "data"
	// Value event
	Value = "value"
	// Vibration event
	Vibration = "vibration"
)

// AnalogReader interface represents an Adaptor which has AnalogRead capabilities
type AnalogReader interface {
	// gobot.Adaptor
	AnalogRead(pin string) (val int, err error)
}

// AnalogWriter interface represents an Adaptor which has AnalogWrite capabilities
type AnalogWriter interface {
	// gobot.Adaptor
	AnalogWrite(pin string, val int) error
}

// optionApplier needs to be implemented by each configurable option type
type optionApplier interface {
	apply(cfg *configuration)
}

// configuration contains all changeable attributes of the driver.
type configuration struct {
	name string
}

// nameOption is the type for applying another name to the configuration
type nameOption string

// Driver implements the interface gobot.Driver.
type driver struct {
	driverCfg  *configuration
	connection interface{}
	afterStart func() error
	beforeHalt func() error
	gobot.Commander
	mutex *sync.Mutex // e.g. used to prevent data race between cyclic and single shot write/read to values and scaler
}

// newDriver creates a new basic analog gobot driver.
func newDriver(a interface{}, name string) *driver {
	d := driver{
		driverCfg:  &configuration{name: gobot.DefaultName(name)},
		connection: a,
		afterStart: func() error { return nil },
		beforeHalt: func() error { return nil },
		Commander:  gobot.NewCommander(),
		mutex:      &sync.Mutex{},
	}

	return &d
}

// WithName is used to replace the default name of the driver.
func WithName(name string) optionApplier {
	return nameOption(name)
}

// Name returns the name of the driver.
func (d *driver) Name() string {
	return d.driverCfg.name
}

// SetName sets the name of the driver.
// Deprecated: Please use option [aio.WithName] instead.
func (d *driver) SetName(name string) {
	WithName(name).apply(d.driverCfg)
}

// Connection returns the connection of the driver.
func (d *driver) Connection() gobot.Connection {
	if conn, ok := d.connection.(gobot.Connection); ok {
		return conn
	}

	log.Printf("%s has no gobot connection\n", d.driverCfg.name)
	return nil
}

// Start initializes the driver.
func (d *driver) Start() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// currently there is nothing to do here for the driver

	return d.afterStart()
}

// Halt halts the driver.
func (d *driver) Halt() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// currently there is nothing to do after halt for the driver

	return d.beforeHalt()
}

// apply change the name in the configuration.
func (o nameOption) apply(c *configuration) {
	c.name = string(o)
}
