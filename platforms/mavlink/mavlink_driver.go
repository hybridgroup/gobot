package mavlink

import (
	"time"

	"github.com/hybridgroup/gobot"
	common "github.com/hybridgroup/gobot/platforms/mavlink/common"
)

var _ gobot.DriverInterface = (*MavlinkDriver)(nil)

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
	m.AddEvent("error")

	return m
}

// adaptor returns driver associated adaptor
func (m *MavlinkDriver) adaptor() *MavlinkAdaptor {
	return m.Driver.Adaptor().(*MavlinkAdaptor)
}

// Start begins process to read mavlink packets every m.Interval
// and process them
func (m *MavlinkDriver) Start() error {
	go func() {
		for {
			packet, err := common.ReadMAVLinkPacket(m.adaptor().sp)
			if err != nil {
				gobot.Publish(m.Event("error"), err)
				continue
			}
			gobot.Publish(m.Event("packet"), packet)
			message, err := packet.MAVLinkMessage()
			if err != nil {
				gobot.Publish(m.Event("error"), err)
				continue
			}
			gobot.Publish(m.Event("message"), message)
			<-time.After(m.Interval())
		}
	}()
	return nil
}

// SendPacket sends a packet to mavlink device
func (m *MavlinkDriver) SendPacket(packet *common.MAVLinkPacket) (err error) {
	_, err = m.adaptor().sp.Write(packet.Pack())
	return err
}

// Halt returns true if device is halted successfully
func (m *MavlinkDriver) Halt() error { return nil }
