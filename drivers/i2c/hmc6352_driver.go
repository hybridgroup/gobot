package i2c

import "gobot.io/x/gobot"

const hmc6352Address = 0x21

// HMC6352Driver is a Driver for a HMC6352 digital compass
type HMC6352Driver struct {
	name       string
	connection I2c
}

// NewHMC6352Driver creates a new driver with specified i2c interface
func NewHMC6352Driver(a I2c) *HMC6352Driver {
	return &HMC6352Driver{
		name:       "HMC6352",
		connection: a,
	}
}

// Name returns the name for this Driver
func (h *HMC6352Driver) Name() string                 { return h.name }

// SetName sets the name for this Driver
func (h *HMC6352Driver) SetName(n string)             { h.name = n }

// Connection returns the connection for this Driver
func (h *HMC6352Driver) Connection() gobot.Connection { return h.connection.(gobot.Connection) }

// Start initializes the hmc6352
func (h *HMC6352Driver) Start() (err error) {
	if err := h.connection.I2cStart(hmc6352Address); err != nil {
		return err
	}
	if err := h.connection.I2cWrite(hmc6352Address, []byte("A")); err != nil {
		return err
	}
	return
}

// Halt returns true if devices is halted successfully
func (h *HMC6352Driver) Halt() (err error) { return }

// Heading returns the current heading
func (h *HMC6352Driver) Heading() (heading uint16, err error) {
	if err = h.connection.I2cWrite(hmc6352Address, []byte("A")); err != nil {
		return
	}
	ret, err := h.connection.I2cRead(hmc6352Address, 2)
	if err != nil {
		return
	}
	if len(ret) == 2 {
		heading = (uint16(ret[1]) + uint16(ret[0])*256) / 10
		return
	} else {
		err = ErrNotEnoughBytes
	}
	return
}
