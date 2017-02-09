package i2c

import "gobot.io/x/gobot"

const hmc6352Address = 0x21

// HMC6352Driver is a Driver for a HMC6352 digital compass
type HMC6352Driver struct {
	name       string
	connector  I2cConnector
	connection I2cConnection
	I2cConfig
}

// NewHMC6352Driver creates a new driver with specified i2c interface
func NewHMC6352Driver(a I2cConnector, options ...func(I2cConfig)) *HMC6352Driver {
	hmc := &HMC6352Driver{
		name:      gobot.DefaultName("HMC6352"),
		connector: a,
		I2cConfig: NewI2cConfig(),
	}

	for _, option := range options {
		option(hmc)
	}

	return hmc
}

// Name returns the name for this Driver
func (h *HMC6352Driver) Name() string { return h.name }

// SetName sets the name for this Driver
func (h *HMC6352Driver) SetName(n string) { h.name = n }

// Connection returns the connection for this Driver
func (h *HMC6352Driver) Connection() gobot.Connection { return h.connector.(gobot.Connection) }

// Start initializes the hmc6352
func (h *HMC6352Driver) Start() (err error) {
	bus := h.GetBus(h.connector.I2cGetDefaultBus())
	address := h.GetAddress(hmc6352Address)

	h.connection, err = h.connector.I2cGetConnection(address, bus)
	if err != nil {
		return err
	}

	if _, err := h.connection.Write([]byte("A")); err != nil {
		return err
	}
	return
}

// Halt returns true if devices is halted successfully
func (h *HMC6352Driver) Halt() (err error) { return }

// Heading returns the current heading
func (h *HMC6352Driver) Heading() (heading uint16, err error) {
	if _, err = h.connection.Write([]byte("A")); err != nil {
		return
	}
	buf := []byte{0, 0}
	bytesRead, err := h.connection.Read(buf)
	if err != nil {
		return
	}
	if bytesRead == 2 {
		heading = (uint16(buf[1]) + uint16(buf[0])*256) / 10
		return
	} else {
		err = ErrNotEnoughBytes
	}
	return
}
