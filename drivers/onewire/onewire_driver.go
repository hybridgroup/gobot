package onewire

import (
	"fmt"
	"log"
	"sync"

	"gobot.io/x/gobot/v2"
)

// connector lets adaptors provide the drivers to get access to the 1-wire devices on platforms.
type connector interface {
	// GetOneWireConnection returns a connection to a 1-wire device with family code and serial number.
	GetOneWireConnection(familyCode byte, serialNumber uint64) (Connection, error)
}

// Connection is a connection to a 1-wire device with family code and serial number on a specific bus, provided by
// an adaptor, usually just by calling the onewire package's GetOneWireConnection() function.
type Connection gobot.OneWireOperations

// optionApplier needs to be implemented by each configurable option type
type optionApplier interface {
	apply(cfg *configuration)
}

// configuration contains all changeable attributes of the driver.
type configuration struct {
	name         string
	familyCode   byte
	serialNumber uint64
}

// nameOption is the type for applying another name to the configuration
type nameOption string

// Driver implements the interface gobot.Driver.
type driver struct {
	driverCfg  *configuration
	connector  connector
	connection Connection
	afterStart func() error
	beforeHalt func() error
	gobot.Commander
	mutex *sync.Mutex // mutex often needed to ensure that write-read sequences are not interrupted
}

// newDriver creates a new generic and basic 1-wire gobot driver.
//
// Supported options:
//
//	"WithName"
func newDriver(a connector, name string, familyCode byte, serialNumber uint64, opts ...interface{}) *driver {
	d := &driver{
		driverCfg:  &configuration{name: gobot.DefaultName(name), familyCode: familyCode, serialNumber: serialNumber},
		connector:  a,
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

// Name returns the name of the device.
func (d *driver) Name() string {
	return d.driverCfg.name
}

// SetName sets the name of the device (deprecated, use WithName() instead).
func (d *driver) SetName(name string) {
	d.driverCfg.name = name
}

// Connection returns the connection of the device.
func (d *driver) Connection() gobot.Connection {
	if conn, ok := d.connection.(gobot.Connection); ok {
		return conn
	}

	log.Printf("%s has no gobot connection\n", d.driverCfg.name)
	return nil
}

// Start initializes the device.
func (d *driver) Start() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	var err error
	d.connection, err = d.connector.GetOneWireConnection(d.driverCfg.familyCode, d.driverCfg.serialNumber)
	if err != nil {
		return err
	}

	return d.afterStart()
}

// Halt halts the device.
func (d *driver) Halt() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// currently there is nothing to do here for the driver, the connection is cached on adaptor side
	// and will be closed on adaptor Finalize()

	return d.beforeHalt()
}

func (o nameOption) String() string {
	return "name option for 1-wire drivers"
}

// apply change the name in the configuration.
func (o nameOption) apply(c *configuration) {
	c.name = string(o)
}
