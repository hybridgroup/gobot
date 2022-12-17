package gobot

import "io"

// DigitalPinOptioner is the interface to provide the possibility to change pin behavior for the next usage
type DigitalPinOptioner interface {
	// SetLabel change the pins label
	SetLabel(string) (changed bool)
	// SetDirectionOutput sets the pins direction to output with the given initial value
	SetDirectionOutput(initialState int) (changed bool)
	// SetDirectionInput sets the pins direction to input
	SetDirectionInput() (changed bool)
}

// DigitalPinOptionApplier is the interface to apply options to change pin behavior immediately
type DigitalPinOptionApplier interface {
	// ApplyOptions apply all given options to the pin immediately
	ApplyOptions(...func(DigitalPinOptioner) bool) error
}

// DigitalPinner is the interface for system gpio interactions
type DigitalPinner interface {
	// Export exports the pin for use by the adaptor
	Export() error
	// Unexport releases the pin from the adaptor, so it is free for the operating system
	Unexport() error
	// Read reads the current value of the pin
	Read() (int, error)
	// Write writes to the pin
	Write(int) error
	// DigitalPinOptionApplier is the interface to change pin behavior immediately
	DigitalPinOptionApplier
}

// DigitalPinValuer is the interface to get pin behavior for the next usage. The interface is and should be rarely used.
type DigitalPinValuer interface {
	// DirectionBehavior gets the direction behavior when the pin is used the next time.
	// This means its possibly not in this direction type at the moment.
	DirectionBehavior() string
}

// DigitalPinnerProvider is the interface that an Adaptor should implement to allow clients to obtain
// access to any DigitalPin's available on that board. If the pin is initially acquired, it is an input.
// Pin direction and other options can be changed afterwards by pin.ApplyOptions() at any time.
type DigitalPinnerProvider interface {
	DigitalPin(id string) (DigitalPinner, error)
}

// PWMPinner is the interface for system PWM interactions
type PWMPinner interface {
	// Export exports the PWM pin for use by the operating system
	Export() error
	// Unexport releases the PWM pin from the operating system
	Unexport() error
	// Enabled returns the enabled state of the PWM pin
	Enabled() (bool, error)
	// SetEnabled enables/disables the PWM pin
	SetEnabled(bool) error
	// Polarity returns true if the polarity of the PWM pin is normal, otherwise false
	Polarity() (bool, error)
	// SetPolarity sets the polarity of the PWM pin to normal if called with true and to inverted if called with false
	SetPolarity(normal bool) error
	// Period returns the current PWM period in nanoseconds for pin
	Period() (uint32, error)
	// SetPeriod sets the current PWM period in nanoseconds for pin
	SetPeriod(uint32) error
	// DutyCycle returns the duty cycle in nanoseconds for the PWM pin
	DutyCycle() (uint32, error)
	// SetDutyCycle writes the duty cycle in nanoseconds to the PWM pin
	SetDutyCycle(uint32) error
}

// PWMPinnerProvider is the interface that an Adaptor should implement to allow
// clients to obtain access to any PWMPin's available on that board.
type PWMPinnerProvider interface {
	PWMPin(id string) (PWMPinner, error)
}

// BusOperations are functions provided by a bus device, e.g. SPI, i2c.
type BusOperations interface {
	// WriteBlockData writes the given data starting from the given register of bus device.
	WriteBlockData(reg uint8, data []byte) error
	// WriteByte writes the given byte value to the current register of a bus device.
	WriteByte(val byte) error
	// WriteBytes writes the given data starting from the current register of an bus device.
	WriteBytes(data []byte) error
	// Close the connection.
	Close() error
}

// I2cOperations represents the i2c methods according to I2C/SMBus specification.
// Some functions are not in the interface yet:
// * Process Call (WriteWordDataReadWordData)
// * Block Write - Block Read (WriteBlockDataReadBlockData)
// * Host Notify - WriteWordData() can be used instead
//
// see: https://docs.kernel.org/i2c/smbus-protocol.html#key-to-symbols
//
// S: Start condition; Sr: Repeated start condition, used to switch from write to read mode.
// P: Stop condition; Rd/Wr (1 bit): Read/Write bit. Rd equals 1, Wr equals 0.
// A, NA (1 bit): Acknowledge (ACK) and Not Acknowledge (NACK) bit
// Addr (7 bits): I2C 7 bit address. (10 bit I2C address not yet supported by gobot).
// Comm (8 bits): Command byte, a data byte which often selects a register on the device.
// Data (8 bits): A plain data byte. DataLow and DataHigh represent the low and high byte of a 16 bit word.
// Count (8 bits): A data byte containing the length of a block operation.
// [..]: Data sent by I2C device, as opposed to data sent by the host adapter.
//
// WriteByte must be implemented as the sequence:
// "S Addr Wr [A] Data [A] P"
// WriteBlockData must be implemented as the sequence:
// "S Addr Wr [A] Comm [A] Count [A] Data [A] Data [A] ... [A] Data [A] P"
type I2cOperations interface {
	io.ReadWriteCloser
	BusOperations

	// ReadByte must be implemented as the sequence:
	// "S Addr Rd [A] [Data] NA P"
	ReadByte() (byte, error)

	// ReadByteData must be implemented as the sequence:
	// "S Addr Wr [A] Comm [A] Sr Addr Rd [A] [Data] NA P"
	ReadByteData(reg uint8) (uint8, error)

	// ReadWordData must be implemented as the sequence:
	// "S Addr Wr [A] Comm [A] Sr Addr Rd [A] [DataLow] A [DataHigh] NA P"
	ReadWordData(reg uint8) (uint16, error)

	// ReadBlockData must be implemented as the sequence:
	// "S Addr Wr [A] Comm [A] Sr Addr Rd [A] [Count] A [Data] A [Data] A ... A [Data] NA P"
	ReadBlockData(reg uint8, b []byte) error

	// WriteByteData must be implemented as the sequence:
	// "S Addr Wr [A] Comm [A] Data [A] P"
	WriteByteData(reg uint8, val uint8) error

	// WriteWordData must be implemented as the sequence:
	// "S Addr Wr [A] Comm [A] DataLow [A] DataHigh [A] P"
	WriteWordData(reg uint8, val uint16) error
}

// SpiOperations are the wrappers around the actual functions used by the SPI device interface
type SpiOperations interface {
	BusOperations
	// ReadCommandData uses the SPI device TX to send/receive data.
	ReadCommandData(command []byte, data []byte) error
}

// Adaptor is the interface that describes an adaptor in gobot
type Adaptor interface {
	// Name returns the label for the Adaptor
	Name() string
	// SetName sets the label for the Adaptor
	SetName(string)
	// Connect initiates the Adaptor
	Connect() error
	// Finalize terminates the Adaptor
	Finalize() error
}

// Porter is the interface that describes an adaptor's port
type Porter interface {
	Port() string
}
