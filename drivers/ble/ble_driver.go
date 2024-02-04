package ble

import (
	"log"
	"sync"

	"gobot.io/x/gobot/v2"
)

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
type Driver struct {
	gobot.Commander
	connection interface{}
	driverCfg  *configuration
	afterStart func() error
	beforeHalt func() error
	mutex      *sync.Mutex
}

// NewDriver creates a new basic BLE gobot driver.
func NewDriver(a interface{}, name string, afterStart func() error, beforeHalt func() error) *Driver {
	if afterStart == nil {
		afterStart = func() error { return nil }
	}

	if beforeHalt == nil {
		beforeHalt = func() error { return nil }
	}

	d := Driver{
		driverCfg:  &configuration{name: gobot.DefaultName(name)},
		connection: a,
		afterStart: afterStart,
		beforeHalt: beforeHalt,
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
func (d *Driver) Name() string {
	return d.driverCfg.name
}

// SetName sets the name of the driver.
// Deprecated: Please use option [aio.WithName] instead.
func (d *Driver) SetName(name string) {
	WithName(name).apply(d.driverCfg)
}

// Connection returns the connection of the driver.
func (d *Driver) Connection() gobot.Connection {
	if conn, ok := d.connection.(gobot.Connection); ok {
		return conn
	}

	log.Printf("%s has no gobot connection\n", d.driverCfg.name)
	return nil
}

// Start initializes the driver.
func (d *Driver) Start() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// currently there is nothing to do here for the driver

	return d.afterStart()
}

// Halt halts the driver.
func (d *Driver) Halt() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// currently there is nothing to do after halt for the driver

	return d.beforeHalt()
}

// Adaptor returns the BLE adaptor
func (d *Driver) Adaptor() gobot.BLEConnector {
	if a, ok := d.connection.(gobot.BLEConnector); ok {
		return a
	}

	log.Printf("%s has no BLE connector\n", d.driverCfg.name)
	return nil
}

func (d *Driver) Mutex() *sync.Mutex {
	return d.mutex
}

// apply change the name in the configuration.
func (o nameOption) apply(c *configuration) {
	c.name = string(o)
}
