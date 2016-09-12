package mavlink

import (
	"time"

	"github.com/hybridgroup/gobot"
	common "github.com/hybridgroup/gobot/platforms/mavlink/common"
)

const (
	// PacketEvent event
	PacketEvent = "packet"
	// MessageEvent event
	MessageEvent = "message"
	// ErrorIOE event
	ErrorIOEvent = "errorIO"
	// ErrorMAVLinkEvent event
	ErrorMAVLinkEvent = "errorMAVLink"
)

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

	m.AddEvent(PacketEvent)
	m.AddEvent(MessageEvent)
	m.AddEvent(ErrorIOEvent)
	m.AddEvent(ErrorMAVLinkEvent)

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
				m.Publish(ErrorIOEvent, err)
				continue
			}
			m.Publish(PacketEvent, packet)
			message, err := packet.MAVLinkMessage()
			if err != nil {
				m.Publish(ErrorMAVLinkEvent, err)
				continue
			}
			m.Publish(MessageEvent, message)
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
