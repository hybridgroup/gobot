package serial

import (
	"log"
	"sync"

	"gobot.io/x/gobot/v2"
)

type SerialReader interface {
	SerialRead(b []byte) (n int, err error)
}

type SerialWriter interface {
	SerialWrite(b []byte) (n int, err error)
}

// OptionApplier needs to be implemented by each configurable option type
type OptionApplier interface {
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

// NewDriver creates a new basic serial gobot driver.
func NewDriver(a interface{}, name string,
	afterStart func() error, beforeHalt func() error,
	opts ...OptionApplier,
) *Driver {
	if afterStart == nil {
		afterStart = func() error { return nil }
	}

	if beforeHalt == nil {
		beforeHalt = func() error { return nil }
	}
	d := Driver{
		Commander:  gobot.NewCommander(),
		connection: a,
		driverCfg:  &configuration{name: gobot.DefaultName(name)},
		afterStart: afterStart,
		beforeHalt: beforeHalt,
		mutex:      &sync.Mutex{},
	}

	for _, o := range opts {
		o.apply(d.driverCfg)
	}

	return &d
}

// WithName is used to replace the default name of the driver.
func WithName(name string) OptionApplier {
	return nameOption(name)
}

// Name returns the name of the driver.
func (d *Driver) Name() string {
	return d.driverCfg.name
}

// SetName sets the name of the driver.
// Deprecated: Please use option [serial.WithName] instead.
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

func (d *Driver) Mutex() *sync.Mutex {
	return d.mutex
}

// apply change the name in the configuration.
func (o nameOption) apply(c *configuration) {
	c.name = string(o)
}
