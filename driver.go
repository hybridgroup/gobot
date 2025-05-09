package gobot

// Driver is the interface that describes a driver in gobot
type Driver interface {
	// Name returns the label for the Driver
	Name() string
	// SetName sets the label for the Driver (deprecated, use WithName() instead).
	// Please use options [aio.WithName, ble.WithName, gpio.WithName, onewire.WithName or serial.WithName] instead.
	SetName(s string)
	// Start initiates the Driver
	Start() error
	// Halt terminates the Driver
	Halt() error
	// Connection returns the Connection associated with the Driver
	Connection() Connection
}

// Pinner is the interface that describes a driver's pin
type Pinner interface {
	Pin() string
}
