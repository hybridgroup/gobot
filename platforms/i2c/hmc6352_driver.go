package i2c

import (
	"github.com/hybridgroup/gobot"
)

var _ gobot.DriverInterface = (*HMC6352Driver)(nil)

type HMC6352Driver struct {
	gobot.Driver
	Heading uint16
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
func (h *HMC6352Driver) Start() error {
	h.adaptor().I2cStart(0x21)
	h.adaptor().I2cWrite([]byte("A"))

	gobot.Every(h.Interval(), func() {
		h.adaptor().I2cWrite([]byte("A"))
		ret := h.adaptor().I2cRead(2)
		if len(ret) == 2 {
			h.Heading = (uint16(ret[1]) + uint16(ret[0])*256) / 10
		}
	})
	return nil
}

// Halt returns true if devices is halted successfully
func (h *HMC6352Driver) Halt() error { return nil }
