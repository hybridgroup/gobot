package spi

import (
	xspi "golang.org/x/exp/io/spi"
	"time"
)

const (
	// BusNotInitialized is the initial value for a bus
	BusNotInitialized = -1

	// AddressNotInitialized is the initial value for an address
	AddressNotInitialized = -1
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
	//	SetAddress(int) error
}

// Connector lets Adaptors provide the interface for Drivers
// to get access to the SPI buses on platforms that support SPI.
type Connector interface {
	// GetConnection returns a connection to device at the specified address
	// and bus. Bus numbering starts at index 0, the range of valid buses is
	// platform specific.
	GetSpiConnection(busNum, address, mode int, maxSpeed int64) (device Connection, err error)

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
	address  int
}

// NewConnection creates and returns a new connection to a specific
// spi device on a bus and address
func NewConnection(bus SPIDevice, address int) (connection *SpiConnection) {
	return &SpiConnection{bus: bus, address: address}
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

func (c *SpiConnection) ReadBytes(address byte, msg byte, numBytes int) (val []byte, err error) {
	w := make([]byte, numBytes)
	w[0] = address
	w[1] = msg
	r := make([]byte, len(w))
	err = c.Tx(w, r)
	if err != nil {
		return val, err
	}
	return r, nil
}

func (c *SpiConnection) ReadUint8(address, msg byte) (val uint8, err error) {
	r, err := c.ReadBytes(address, msg, 8)
	if err != nil {
		return val, err
	}
	return uint8(r[4]) << 8, nil
}

func (c *SpiConnection) ReadUint16(address, msg byte) (val uint16, err error) {
	r, err := c.ReadBytes(address, msg, 8)
	if err != nil {
		return val, err
	}
	return uint16(r[4])<<8 | uint16(r[5]), nil
}

func (c *SpiConnection) ReadUint32(address, msg byte) (val uint32, err error) {
	r, err := c.ReadBytes(address, msg, 8)
	if err != nil {
		return val, err
	}
	return uint32(r[4])<<24 | uint32(r[5])<<16 | uint32(r[6])<<8 | uint32(r[7]), nil
}

func (c *SpiConnection) WriteBytes(w []byte) (err error) {
	return c.Tx(w, nil)
}
