package i2c

import (
	"fmt"
	"strconv"
	"sync"

	"gobot.io/x/gobot"
)

// Driver implements the interface gobot.Driver.
type Driver struct {
	name           string
	defaultAddress int
	connector      Connector
	connection     Connection
	afterStart     func() error
	beforeHalt     func() error
	Config
	gobot.Commander
	mutex *sync.Mutex // mutex often needed to ensure that write-read sequences are not interrupted
}

// NewDriver creates a new generic and basic i2c gobot driver.
func NewDriver(c Connector, name string, address int, options ...func(Config)) *Driver {
	d := &Driver{
		name:           gobot.DefaultName(name),
		defaultAddress: address,
		connector:      c,
		afterStart:     func() error { return nil },
		beforeHalt:     func() error { return nil },
		Config:         NewConfig(),
		Commander:      gobot.NewCommander(),
		mutex:          &sync.Mutex{},
	}

	for _, option := range options {
		option(d)
	}

	return d
}

// Name returns the name of the i2c device.
func (d *Driver) Name() string {
	return d.name
}

// SetName sets the name of the i2c device.
func (d *Driver) SetName(name string) {
	d.name = name
}

// Connection returns the connection of the i2c device.
func (d *Driver) Connection() gobot.Connection {
	return d.connector.(gobot.Connection)
}

// Start initializes the i2c device.
func (d *Driver) Start() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	var err error
	bus := d.GetBusOrDefault(d.connector.GetDefaultBus())
	address := d.GetAddressOrDefault(int(d.defaultAddress))

	if d.connection, err = d.connector.GetConnection(address, bus); err != nil {
		return err
	}

	return d.afterStart()
}

// Halt halts the i2c device.
func (d *Driver) Halt() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.beforeHalt(); err != nil {
		return err
	}

	// currently there is nothing to do here for the driver
	return nil
}

// Write implements a simple write mechanism to the given register of an i2c device.
func (d *Driver) Write(pin string, val int) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	register, err := driverParseRegister(pin)
	if err != nil {
		return err
	}

	// TODO: create buffer from size
	// currently only one byte value is supported
	b := []byte{uint8(val)}
	return d.connection.WriteBlockData(uint8(register), b)
}

// Read implements a simple read mechanism from the given register of an i2c device.
func (d *Driver) Read(pin string) (int, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	register, err := driverParseRegister(pin)
	if err != nil {
		return 0, err
	}

	// TODO: create buffer from size
	// currently only one byte value is supported
	b := []byte{0}
	if err := d.connection.ReadBlockData(register, b); err != nil {
		return 0, err
	}

	return int(b[0]), nil
}

func driverParseRegister(pin string) (uint8, error) {
	register, err := strconv.ParseUint(pin, 10, 8)
	if err != nil {
		return 0, fmt.Errorf("Could not parse the register from given pin '%s'", pin)
	}
	return uint8(register), nil
}
