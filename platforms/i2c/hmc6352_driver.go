package i2c

import (
	"github.com/hybridgroup/gobot"
)

type HMC6352Driver struct {
	gobot.Driver
	Adaptor I2cInterface
	Heading uint16
}

func NewHMC6352Driver(a I2cInterface, name string) *HMC6352Driver {
	return &HMC6352Driver{
		Driver: gobot.Driver{
			Name: name,
		},
		Adaptor: a,
	}
}

func (h *HMC6352Driver) Start() bool {
	h.Adaptor.I2cStart(0x21)
	h.Adaptor.I2cWrite([]byte("A"))

	gobot.Every(h.Interval, func() {
		h.Adaptor.I2cWrite([]byte("A"))
		ret := h.Adaptor.I2cRead(2)
		if len(ret) == 2 {
			h.Heading = (uint16(ret[1]) + uint16(ret[0])*256) / 10
		}
	})
	return true
}
func (self *HMC6352Driver) Init() bool { return true }
func (self *HMC6352Driver) Halt() bool { return true }
