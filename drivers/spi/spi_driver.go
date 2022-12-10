package spi

import (
	"sync"

	"gobot.io/x/gobot"
)

const (
	// NotInitialized is the initial value for a bus/chip
	NotInitialized = -1
)

// Operations are the wrappers around the actual functions used by the SPI device interface
type Operations interface {
	Close() error
	Tx(w, r []byte) error
}

// TODO: rename to golang getter spec (no prefix "Get" for simple getters)

// Connector lets Adaptors provide the interface for Drivers
// to get access to the SPI buses on platforms that support SPI.
type Connector interface {
	// GetSpiConnection returns a connection to a SPI device at the specified bus and chip.
	// Bus numbering starts at index 0, the range of valid buses is
	// platform specific. Same with chip numbering.
	GetSpiConnection(busNum, chip, mode, bits int, maxSpeed int64) (device Connection, err error)

	// GetSpiDefaultBus returns the default SPI bus index
	GetSpiDefaultBus() int

	// GetSpiDefaultChip returns the default SPI chip index
	GetSpiDefaultChip() int

	// GetDefaultMode returns the default SPI mode (0/1/2/3)
	GetSpiDefaultMode() int

	// GetDefaultMode returns the default SPI number of bits (8)
	GetSpiDefaultBits() int

	// GetSpiDefaultMaxSpeed returns the max SPI speed
	GetSpiDefaultMaxSpeed() int64
}

// Connection is a connection to a SPI device with a specific bus/chip.
// Provided by an Adaptor, usually just by calling the spi package's GetSpiConnection() function.
type Connection Operations

// Config is the interface which describes how a Driver can specify
// optional SPI params such as which SPI bus it wants to use.
type Config interface {
	// WithBus sets which bus to use
	WithBus(bus int)

	// GetBusOrDefault gets which bus to use
	GetBusOrDefault(def int) int

	// WithChip sets which chip to use
	WithChip(chip int)

	// GetChipOrDefault gets which chip to use
	GetChipOrDefault(def int) int

	// WithMode sets which mode to use
	WithMode(mode int)

	// GetModeOrDefault gets which mode to use
	GetModeOrDefault(def int) int

	// WithBIts sets how many bits to use
	WithBits(bits int)

	// GetBitsOrDefault gets how many bits to use
	GetBitsOrDefault(def int) int

	// WithSpeed sets which speed to use (in Hz)
	WithSpeed(speed int64)

	// GetSpeedOrDefault gets which speed to use (in Hz)
	GetSpeedOrDefault(def int64) int64
}

type Driver struct {
	name       string
	connector  Connector
	connection Connection
	afterStart func() error
	beforeHalt func() error
	Config
	gobot.Commander
	mutex sync.Mutex
}

func NewDriver(a Connector, name string, options ...func(Config)) *Driver {
	d := &Driver{
		name:       gobot.DefaultName(name),
		connector:  a,
		afterStart: func() error { return nil },
		beforeHalt: func() error { return nil },
		Config:     NewConfig(),
		Commander:  gobot.NewCommander(),
	}
	for _, option := range options {
		option(d)
	}
	return d
}

// Name returns the name of the device.
func (d *Driver) Name() string { return d.name }

// SetName sets the name of the device.
func (d *Driver) SetName(n string) { d.name = n }

// Connection returns the Connection of the device.
func (d *Driver) Connection() gobot.Connection { return d.connection.(gobot.Connection) }

// Start initializes the driver.
func (d *Driver) Start() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	bus := d.GetBusOrDefault(d.connector.GetSpiDefaultBus())
	chip := d.GetChipOrDefault(d.connector.GetSpiDefaultChip())
	mode := d.GetModeOrDefault(d.connector.GetSpiDefaultMode())
	bits := d.GetBitsOrDefault(d.connector.GetSpiDefaultBits())
	maxSpeed := d.GetSpeedOrDefault(d.connector.GetSpiDefaultMaxSpeed())

	var err error
	d.connection, err = d.connector.GetSpiConnection(bus, chip, mode, bits, maxSpeed)
	if err != nil {
		return err
	}
	return d.afterStart()
}

// Halt stops the driver.
func (d *Driver) Halt() (err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.beforeHalt(); err != nil {
		return err
	}

	return d.connection.Close()
}
