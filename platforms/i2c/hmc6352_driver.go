package i2c

import (
	"github.com/hybridgroup/gobot"
	"time"
)

type HMC6352Driver struct {
	gobot.Driver
	Adaptor I2cInterface
	Heading uint16
}

func NewHMC6352Driver(a I2cInterface) *HMC6352Driver {
	return &HMC6352Driver{
		Adaptor: a,
	}
}

func (h *HMC6352Driver) Start() bool {
	h.Adaptor.I2cStart(0x21)
	h.Adaptor.I2cWrite([]uint16{uint16([]byte("A")[0])})

	gobot.Every(1*time.Second, func() {
		h.Adaptor.I2cWrite([]uint16{uint16([]byte("A")[0])})
		ret := h.Adaptor.I2cRead(2)
		if len(ret) == 2 {
			h.Heading = (ret[1] + ret[0]*256) / 10
		}
	})
	return true
}
func (self *HMC6352Driver) Init() bool { return true }
func (self *HMC6352Driver) Halt() bool { return true }
