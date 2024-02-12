package spi

import (
	"log"
	"sync"

	"gobot.io/x/gobot/v2"
)

const (
	// NotInitialized is the initial value for a bus/chip
	NotInitialized = -1
)

// Connector lets adaptors provide the interface for Drivers
// to get access to the SPI buses on platforms that support SPI.
type Connector interface {
	// GetSpiConnection returns a connection to a SPI device at the specified bus and chip.
	// Bus numbering starts at index 0, the range of valid buses is
	// platform specific. Same with chip numbering.
	GetSpiConnection(busNum, chip, mode, bits int, maxSpeed int64) (device Connection, err error)

	// SpiDefaultBusNumber returns the default SPI bus index
	SpiDefaultBusNumber() int

	// SpiDefaultChipNumber returns the default SPI chip index
	SpiDefaultChipNumber() int

	// DefaultMode returns the default SPI mode (0/1/2/3)
	SpiDefaultMode() int

	// SpiDefaultBitCount returns the default SPI number of bits (8)
	SpiDefaultBitCount() int

	// SpiDefaultMaxSpeed returns the max SPI speed
	SpiDefaultMaxSpeed() int64
}

// Connection is a connection to a SPI device with a specific bus/chip.
// Provided by an Adaptor, usually just by calling the spi package's GetSpiConnection() function.
type Connection gobot.SpiOperations

// Config is the interface which describes how a Driver can specify
// optional SPI params such as which SPI bus it wants to use.
type Config interface {
	// SetBusNumber sets which bus to use
	SetBusNumber(bus int)

	// GetBusNumberOrDefault gets which bus to use
	GetBusNumberOrDefault(def int) int

	// SetChipNumber sets which chip to use
	SetChipNumber(chip int)

	// GetChipNumberOrDefault gets which chip to use
	GetChipNumberOrDefault(def int) int

	// SetMode sets which mode to use
	SetMode(mode int)

	// GetModeOrDefault gets which mode to use
	GetModeOrDefault(def int) int

	// SetUsedBits sets how many bits to use
	SetBitCount(count int)

	// GetBitCountOrDefault gets how many bits to use
	GetBitCountOrDefault(def int) int

	// SetSpeed sets which speed to use (in Hz)
	SetSpeed(speed int64)

	// GetSpeedOrDefault gets which speed to use (in Hz)
	GetSpeedOrDefault(def int64) int64
}

// Driver implements the interface gobot.Driver for SPI devices.
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

// NewDriver creates a new generic and basic SPI gobot driver.
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
func (d *Driver) Connection() gobot.Connection {
	if conn, ok := d.connector.(gobot.Connection); ok {
		return conn
	}

	log.Printf("%s has no gobot connection\n", d.name)
	return nil
}

// Start initializes the driver.
func (d *Driver) Start() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	bus := d.GetBusNumberOrDefault(d.connector.SpiDefaultBusNumber())
	chip := d.GetChipNumberOrDefault(d.connector.SpiDefaultChipNumber())
	mode := d.GetModeOrDefault(d.connector.SpiDefaultMode())
	bits := d.GetBitCountOrDefault(d.connector.SpiDefaultBitCount())
	maxSpeed := d.GetSpeedOrDefault(d.connector.SpiDefaultMaxSpeed())

	var err error
	d.connection, err = d.connector.GetSpiConnection(bus, chip, mode, bits, maxSpeed)
	if err != nil {
		return err
	}
	return d.afterStart()
}

// Halt stops the driver.
func (d *Driver) Halt() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.beforeHalt(); err != nil {
		return err
	}

	// currently there is nothing to do here for the driver, the connection is cached on adaptor side
	// and will be closed on adaptor Finalize()
	return nil
}
