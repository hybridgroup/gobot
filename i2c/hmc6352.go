package gobotI2C

import (
	"github.com/hybridgroup/gobot"
)

type HMC6352 struct {
	gobot.Driver
	Adaptor I2cInterface
	Heading uint16
}

func NewHMC6352(a I2cInterface) *HMC6352 {
	d := new(HMC6352)
	d.Adaptor = a
	d.Events = make(map[string]chan interface{})
	d.Commands = []string{}
	return d
}

func (self *HMC6352) Start() bool {
	self.Adaptor.I2cStart(0x21)
	self.Adaptor.I2cWrite([]uint16{uint16([]byte("A")[0])})

	gobot.Every(self.Interval, func() {
		self.Adaptor.I2cWrite([]uint16{uint16([]byte("A")[0])})
		ret := self.Adaptor.I2cRead(2)
		if len(ret) == 2 {
			self.Heading = (ret[1] + ret[0]*256) / 10
		}
	})
	return true
}
func (self *HMC6352) Init() bool { return true }
func (self *HMC6352) Halt() bool { return true }
