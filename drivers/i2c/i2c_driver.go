package i2c

import (
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
	"sync"

	"gobot.io/x/gobot/v2"
)

// Config is the interface to set and get I2C device related parameters.
type Config interface {
	// SetBus sets which bus to use
	SetBus(bus int)

	// GetBusOrDefault gets which bus to use
	GetBusOrDefault(def int) int

	// SetAddress sets which address to use
	SetAddress(address int)

	// GetAddressOrDefault gets which address to use
	GetAddressOrDefault(def int) int
}

// Connector lets adaptors (platforms) provide the interface for Drivers to get access to the I2C buses on platforms
// that support I2C. The "I2C" specifier is part of the name to differentiate to SPI at platform level.
type Connector interface {
	// GetI2cConnection creates and returns a connection to device at the specified address
	// and bus. Bus numbering starts at index 0, the range of valid buses is
	// platform specific.
	GetI2cConnection(address int, busNr int) (device Connection, err error)

	// DefaultI2cBus returns the default I2C bus index
	DefaultI2cBus() int
}

// Driver implements the interface gobot.Driver.
type Driver struct {
	Config
	gobot.Commander

	name           string
	defaultAddress int
	connector      Connector
	connection     Connection
	afterStart     func() error
	beforeHalt     func() error
	mutex          *sync.Mutex // mutex often needed to ensure that write-read sequences are not interrupted
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

// Connection returns the gobot connection of the i2c device.
func (d *Driver) Connection() gobot.Connection {
	if d.connector == nil {
		log.Printf("%s has no connector\n", d.name)
		return nil
	}

	if conn, ok := d.connector.(gobot.Connection); ok {
		return conn
	}

	log.Printf("%s has no gobot connection\n", d.name)
	return nil
}

// Start initializes the i2c device.
func (d *Driver) Start() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.connector == nil {
		return fmt.Errorf("%s has no connector", d.name)
	}

	var err error
	bus := d.GetBusOrDefault(d.connector.DefaultI2cBus())
	address := d.GetAddressOrDefault(d.defaultAddress)

	if d.connection, err = d.connector.GetI2cConnection(address, bus); err != nil {
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

	d.connection = nil

	return nil
}

// Write implements a simple write mechanism, starting from the given register of an i2c device.
func (d *Driver) Write(pin string, val int) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.connection == nil {
		return fmt.Errorf("i2c driver not started for '%s'", d.name)
	}

	register, err := driverParseRegister(pin)
	if err != nil {
		return err
	}

	if val > 0xFFFF {
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(val)) //nolint:gosec // ok here
		return d.connection.WriteBlockData(register, buf)
	}
	if val > 0xFF {
		return d.connection.WriteWordData(register, uint16(val))
	}
	return d.connection.WriteByteData(register, uint8(val)) //nolint:gosec // ok here
}

// Read implements a simple read mechanism from the given register of an i2c device.
func (d *Driver) Read(pin string) (int, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.connection == nil {
		return 0, fmt.Errorf("i2c driver not started for '%s'", d.name)
	}

	register, err := driverParseRegister(pin)
	if err != nil {
		return 0, err
	}

	val, err := d.connection.ReadByteData(register)
	if err != nil {
		return 0, err
	}

	return int(val), nil
}

func (d *Driver) write(data []byte) (int, error) {
	if d.connection == nil {
		return 0, fmt.Errorf("i2c driver not started for '%s'", d.name)
	}

	return d.connection.Write(data)
}

func (d *Driver) writeByte(val byte) error {
	if d.connection == nil {
		return fmt.Errorf("i2c driver not started for '%s'", d.name)
	}

	return d.connection.WriteByte(val)
}

func (d *Driver) writeByteData(reg uint8, val byte) error {
	if d.connection == nil {
		return fmt.Errorf("i2c driver not started for '%s'", d.name)
	}

	return d.connection.WriteByteData(reg, val)
}

func (d *Driver) writeWordData(reg uint8, val uint16) error {
	if d.connection == nil {
		return fmt.Errorf("i2c driver not started for '%s'", d.name)
	}

	return d.connection.WriteWordData(reg, val)
}

func (d *Driver) writeBlockData(reg uint8, data []byte) error {
	if d.connection == nil {
		return fmt.Errorf("i2c driver not started for '%s'", d.name)
	}

	return d.connection.WriteBlockData(reg, data)
}

func (d *Driver) read(data []byte) (int, error) {
	if d.connection == nil {
		return 0, fmt.Errorf("i2c driver not started for '%s'", d.name)
	}

	return d.connection.Read(data)
}

func (d *Driver) readByte() (byte, error) {
	if d.connection == nil {
		return 0, fmt.Errorf("i2c driver not started for '%s'", d.name)
	}

	return d.connection.ReadByte()
}

func (d *Driver) readByteData(reg uint8) (byte, error) {
	if d.connection == nil {
		return 0, fmt.Errorf("i2c driver not started for '%s'", d.name)
	}

	return d.connection.ReadByteData(reg)
}

func (d *Driver) readWordData(reg uint8) (uint16, error) {
	if d.connection == nil {
		return 0, fmt.Errorf("i2c driver not started for '%s'", d.name)
	}

	return d.connection.ReadWordData(reg)
}

func (d *Driver) readBlockData(reg uint8, data []byte) error {
	if d.connection == nil {
		return fmt.Errorf("i2c driver not started for '%s'", d.name)
	}

	return d.connection.ReadBlockData(reg, data)
}

func driverParseRegister(pin string) (uint8, error) {
	register, err := strconv.ParseUint(pin, 10, 8)
	if err != nil {
		return 0, fmt.Errorf("could not parse the register from given pin '%s'", pin)
	}
	return uint8(register), nil
}
