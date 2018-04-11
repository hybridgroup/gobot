package spi

import (
	"fmt"
	"log"

	xspi "periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/host"
)

const (
	// BusNotInitialized is the initial value for a bus
	BusNotInitialized = -1
)

// Operations are the wrappers around the actual functions used by the SPI Device interface
type Operations interface {
	Close() error
	Tx(w, r []byte) error
}

// Device is the interface to a specific spi bus/chip
type Device interface {
	Operations
}

// Connector lets Adaptors provide the interface for Drivers
// to get access to the SPI buses on platforms that support SPI.
type Connector interface {
	// GetSpiConnection returns a connection to a SPI device at the specified bus and chip.
	// Bus numbering starts at index 0, the range of valid buses is
	// platform specific. Same with chip numbering.
	GetSpiConnection(busNum, chip, mode, bits int, maxSpeed int64) (device Device, err error)

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

// Connection is a connection to an SPI device with a specified bus
// on a specific chip.
// Implements SPIOperations to talk to the device, wrapping the
// calls in SetAddress to always target the specified device.
// Provided by an Adaptor by implementing the SPIConnector interface.
type Connection Operations

type SpiConnection struct {
	Connection
	port     xspi.PortCloser
	dev      xspi.Conn
	bus      int
	chip     int
	bits     int
	mode     int
	maxSpeed int64
}

// NewConnection creates and returns a new connection to a specific
// spi device on a bus/chip using the periph.io interface
func NewConnection(port xspi.PortCloser, conn xspi.Conn) (connection *SpiConnection) {
	return &SpiConnection{port: port, dev: conn}
}

// Close the SPI connection
func (c *SpiConnection) Close() error {
	return c.port.Close()
}

// Tx uses the SPI device to send/receive data.
func (c *SpiConnection) Tx(w, r []byte) error {
	return c.dev.Tx(w, r)
}

// GetSpiConnection is a helper to return a SPI device
func GetSpiConnection(busNum, chipNum, mode, bits int, maxSpeed int64) (spiDevice Device, err error) {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	var spiMode xspi.Mode
	switch mode {
	case 0:
		spiMode = xspi.Mode0
	case 1:
		spiMode = xspi.Mode1
	case 2:
		spiMode = xspi.Mode2
	case 3:
		spiMode = xspi.Mode3
	default:
		spiMode = xspi.Mode0
	}
	devName := fmt.Sprintf("/dev/spidev%d.%d", busNum, chipNum)
	p, err := spireg.Open(devName)
	if err != nil {
		return nil, err
	}

	c, err := p.Connect(maxSpeed, spiMode, bits)
	if err != nil {
		return nil, err
	}
	spiDevice = NewConnection(p, c)

	return
}
