package mavlink

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	common "github.com/hybridgroup/gobot/platforms/mavlink/common"
)

type MavlinkDriver struct {
	gobot.Driver
}

type MavlinkInterface interface {
}

func NewMavlinkDriver(a *MavlinkAdaptor, name string) *MavlinkDriver {
	m := &MavlinkDriver{
		Driver: *gobot.NewDriver(
			name,
			"mavlink.MavlinkDriver",
			a,
		),
	}

	m.AddEvent("packet")
	m.AddEvent("message")

	return m
}

func (m *MavlinkDriver) adaptor() *MavlinkAdaptor {
	return m.Driver.Adaptor().(*MavlinkAdaptor)
}

func (m *MavlinkDriver) Start() bool {
	go func() {
		for {
			packet, err := common.ReadMAVLinkPacket(m.adaptor().sp)
			if err != nil {
				fmt.Println(err)
				continue
			}
			gobot.Publish(m.Event("packet"), packet)
			message, err := packet.MAVLinkMessage()
			if err != nil {
				fmt.Println(err)
				continue
			}
			gobot.Publish(m.Event("message"), message)
			<-time.After(m.Interval())
		}
	}()
	return true
}

func (m *MavlinkDriver) SendPacket(packet *common.MAVLinkPacket) {
	m.adaptor().sp.Write(packet.Pack())
}

func (m *MavlinkDriver) Halt() bool { return true }
