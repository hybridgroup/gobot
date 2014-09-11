package i2c

import (
	"github.com/hybridgroup/gobot"
)

type HMC6352Driver struct {
	gobot.Driver
	Heading uint16
}

func NewHMC6352Driver(a I2cInterface, name string) *HMC6352Driver {
	return &HMC6352Driver{
		Driver: *gobot.NewDriver(
			name,
			"HMC6352Driver",
			a.(gobot.AdaptorInterface),
		),
	}
}

func (h *HMC6352Driver) adaptor() I2cInterface {
	return h.Adaptor().(I2cInterface)
}

func (h *HMC6352Driver) Start() bool {
	h.adaptor().I2cStart(0x21)
	h.adaptor().I2cWrite([]byte("A"))

	gobot.Every(h.Interval(), func() {
		h.adaptor().I2cWrite([]byte("A"))
		ret := h.adaptor().I2cRead(2)
		if len(ret) == 2 {
			h.Heading = (uint16(ret[1]) + uint16(ret[0])*256) / 10
		}
	})
	return true
}

func (h *HMC6352Driver) Init() bool { return true }
func (h *HMC6352Driver) Halt() bool { return true }
