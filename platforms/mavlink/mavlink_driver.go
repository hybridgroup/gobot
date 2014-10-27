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

// NewMavlinkDriver creates a new mavlink driver with specified name.
//
// It add the following events:
//	"packet" - triggered when a new packet is read
//	"message" - triggered when a new valid message is processed
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

// adaptor returns driver associated adaptor
func (m *MavlinkDriver) adaptor() *MavlinkAdaptor {
	return m.Driver.Adaptor().(*MavlinkAdaptor)
}

// Start begins process to read mavlink packets every m.Interval
// and process them
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

// SendPacket sends a packet to mavlink device
func (m *MavlinkDriver) SendPacket(packet *common.MAVLinkPacket) {
	m.adaptor().sp.Write(packet.Pack())
}

// Halt returns true if device is halted successfully
func (m *MavlinkDriver) Halt() bool { return true }
