package i2c

import (
	"errors"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*HMC6352Driver)(nil)

type HMC6352Driver struct {
	name       string
	connection gobot.Connection
}

// NewHMC6352Driver creates a new driver with specified name and i2c interface
func NewHMC6352Driver(a I2cInterface, name string) *HMC6352Driver {
	return &HMC6352Driver{
		name:       name,
		connection: a.(gobot.Connection),
	}
}

func (h *HMC6352Driver) Name() string                 { return h.name }
func (h *HMC6352Driver) Connection() gobot.Connection { return h.connection }

// adaptor returns HMC6352 adaptor
func (h *HMC6352Driver) adaptor() I2cInterface {
	return h.Connection().(I2cInterface)
}

// Start writes initialization bytes and reads from adaptor
// using specified interval to update Heading
func (h *HMC6352Driver) Start() (errs []error) {
	if err := h.adaptor().I2cStart(0x21); err != nil {
		return []error{err}
	}
	if err := h.adaptor().I2cWrite([]byte("A")); err != nil {
		return []error{err}
	}
	return
}

// Halt returns true if devices is halted successfully
func (h *HMC6352Driver) Halt() (errs []error) { return }

// Heading returns the current heading
func (h *HMC6352Driver) Heading() (heading uint16, err error) {
	if err = h.adaptor().I2cWrite([]byte("A")); err != nil {
		return
	}
	ret, err := h.adaptor().I2cRead(2)
	if err != nil {
		return
	}
	if len(ret) == 2 {
		heading = (uint16(ret[1]) + uint16(ret[0])*256) / 10
		return
	} else {
		err = errors.New("Not enough bytes read")
	}
	return
}
