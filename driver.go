package gobot

// Driver is the interface that describes a driver in gobot
type Driver interface {
	// Name returns the label for the Driver
	Name() string
	// Start initiates the Driver
	Start() []error
	// Halt terminates the Driver
	Halt() []error
	// Connection returns the Connection assiciated with the Driver
	Connection() Connection
}

// Pinner is the interface that describes a driver's pin
type Pinner interface {
	Pin() string
}
