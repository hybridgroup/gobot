package spi

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
