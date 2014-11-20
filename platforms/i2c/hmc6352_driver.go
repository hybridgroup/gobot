package i2c

import (
	"errors"

	"github.com/hybridgroup/gobot"
)

var _ gobot.DriverInterface = (*HMC6352Driver)(nil)

type HMC6352Driver struct {
	gobot.Driver
}

// NewHMC6352Driver creates a new driver with specified name and i2c interface
func NewHMC6352Driver(a I2cInterface, name string) *HMC6352Driver {
	return &HMC6352Driver{
		Driver: *gobot.NewDriver(
			name,
			"HMC6352Driver",
			a.(gobot.AdaptorInterface),
		),
	}
}

// adaptor returns HMC6352 adaptor
func (h *HMC6352Driver) adaptor() I2cInterface {
	return h.Adaptor().(I2cInterface)
}

// Start writes initialization bytes and reads from adaptor
// using specified interval to update Heading
func (h *HMC6352Driver) Start() (err error) {
	if err = h.adaptor().I2cStart(0x21); err != nil {
		return
	}
	return h.adaptor().I2cWrite([]byte("A"))
}

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

// Halt returns true if devices is halted successfully
func (h *HMC6352Driver) Halt() error { return nil }
