package spi

import (
	xspi "golang.org/x/exp/io/spi"
	"time"
)

const (
	// BusNotInitialized is the initial value for a bus
	BusNotInitialized = -1
)

type SPIOperations interface {
	Close() error
	SetBitOrder(o xspi.Order) error
	SetBitsPerWord(bits int) error
	SetCSChange(leaveEnabled bool) error
	SetDelay(t time.Duration) error
	SetMaxSpeed(speed int) error
	SetMode(mode xspi.Mode) error
	Tx(w, r []byte) error
}

// SPIDevice is the interface to a specific spi bus
type SPIDevice interface {
	SPIOperations
}

// Connector lets Adaptors provide the interface for Drivers
// to get access to the SPI buses on platforms that support SPI.
type Connector interface {
	// GetConnection returns a connection to device at the specified bus.
	// Bus numbering starts at index 0, the range of valid buses is
	// platform specific.
	GetSpiConnection(busNum, mode int, maxSpeed int64) (device Connection, err error)

	// GetDefaultBus returns the default SPI bus index
	GetSpiDefaultBus() int

	// GetDefaultMode returns the default SPI mode (0/1/2/3)
	GetSpiDefaultMode() int

	// GetSpiDefaultMaxSpeed returns the max SPI speed
	GetSpiDefaultMaxSpeed() int64
}

// Connection is a connection to an SPI device with a specified address
// on a specific bus. Used as an alternative to the SPI interface.
// Implements SPIOperations to talk to the device, wrapping the
// calls in SetAddress to always target the specified device.
// Provided by an Adaptor by implementing the SPIConnector interface.
type Connection SPIOperations

type SpiConnection struct {
	bus      SPIDevice
	mode     int
	maxSpeed int64
}

// NewConnection creates and returns a new connection to a specific
// spi device on a bus and address
func NewConnection(bus SPIDevice) (connection *SpiConnection) {
	return &SpiConnection{bus: bus}
}

func (c *SpiConnection) Close() error {
	return c.bus.Close()
}

func (c *SpiConnection) SetBitOrder(o xspi.Order) error {
	return c.bus.SetBitOrder(o)
}

func (c *SpiConnection) SetBitsPerWord(bits int) error {
	return c.bus.SetBitsPerWord(bits)
}

func (c *SpiConnection) SetCSChange(leaveEnabled bool) error {
	return c.bus.SetCSChange(leaveEnabled)
}

func (c *SpiConnection) SetDelay(t time.Duration) error {
	return c.bus.SetDelay(t)
}

func (c *SpiConnection) SetMaxSpeed(speed int) error {
	return c.bus.SetMaxSpeed(speed)
}

func (c *SpiConnection) SetMode(mode xspi.Mode) error {
	return c.bus.SetMode(mode)
}

func (c *SpiConnection) Tx(w, r []byte) error {
	return c.bus.Tx(w, r)
}
