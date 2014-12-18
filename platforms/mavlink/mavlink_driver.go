package mavlink

import (
	"time"

	"github.com/hybridgroup/gobot"
	common "github.com/hybridgroup/gobot/platforms/mavlink/common"
)

var _ gobot.Driver = (*MavlinkDriver)(nil)

type MavlinkDriver struct {
	name       string
	connection gobot.Connection
	interval   time.Duration
	gobot.Eventer
}

type MavlinkInterface interface {
}

// NewMavlinkDriver creates a new mavlink driver with specified name.
//
// It add the following events:
//	"packet" - triggered when a new packet is read
//	"message" - triggered when a new valid message is processed
func NewMavlinkDriver(a *MavlinkAdaptor, name string, v ...time.Duration) *MavlinkDriver {
	m := &MavlinkDriver{
		name:       name,
		connection: a,
		Eventer:    gobot.NewEventer(),
		interval:   10 * time.Millisecond,
	}

	if len(v) > 0 {
		m.interval = v[0]
	}

	m.AddEvent("packet")
	m.AddEvent("message")
	m.AddEvent("errorIO")
	m.AddEvent("errorMAVLink")

	return m
}

func (m *MavlinkDriver) Connection() gobot.Connection { return m.connection }
func (m *MavlinkDriver) Name() string                 { return m.name }

// adaptor returns driver associated adaptor
func (m *MavlinkDriver) adaptor() *MavlinkAdaptor {
	return m.Connection().(*MavlinkAdaptor)
}

// Start begins process to read mavlink packets every m.Interval
// and process them
func (m *MavlinkDriver) Start() (errs []error) {
	go func() {
		for {
			packet, err := common.ReadMAVLinkPacket(m.adaptor().sp)
			if err != nil {
				gobot.Publish(m.Event("errorIO"), err)
				continue
			}
			gobot.Publish(m.Event("packet"), packet)
			message, err := packet.MAVLinkMessage()
			if err != nil {
				gobot.Publish(m.Event("errorMAVLink"), err)
				continue
			}
			gobot.Publish(m.Event("message"), message)
			<-time.After(m.interval)
		}
	}()
	return
}

// Halt returns true if device is halted successfully
func (m *MavlinkDriver) Halt() (errs []error) { return }

// SendPacket sends a packet to mavlink device
func (m *MavlinkDriver) SendPacket(packet *common.MAVLinkPacket) (err error) {
	_, err = m.adaptor().sp.Write(packet.Pack())
	return err
}
