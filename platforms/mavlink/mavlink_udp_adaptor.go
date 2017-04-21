package mavlink

import (
	"net"

	common "gobot.io/x/gobot/platforms/mavlink/common"
)

type UDPConnection interface {
	Close() error
	ReadFromUDP([]byte) (int, *net.UDPAddr, error)
	WriteTo([]byte, net.Addr) (int, error)
}

type UDPAdaptor struct {
	name string
	port string
	sock UDPConnection
}

var _ BaseAdaptor = (*UDPAdaptor)(nil)

// NewAdaptor creates a new Mavlink-over-UDP adaptor with specified
// port.
func NewUDPAdaptor(port string) *UDPAdaptor {
	return &UDPAdaptor{
		name: "Mavlink",
		port: port,
	}
}

func (m *UDPAdaptor) Name() string     { return m.name }
func (m *UDPAdaptor) SetName(n string) { m.name = n }
func (m *UDPAdaptor) Port() string     { return m.port }

// Connect returns true if connection to device is successful
func (m *UDPAdaptor) Connect() error {
	m.close()

	addr, err := net.ResolveUDPAddr("udp", m.Port())
	if err != nil {
		return err
	}

	m.sock, err = net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	return nil
}

func (m *UDPAdaptor) close() error {
	sock := m.sock
	m.sock = nil

	if sock != nil {
		return sock.Close()
	} else {
		return nil
	}
}

// Finalize returns true if connection to devices is closed successfully
func (m *UDPAdaptor) Finalize() (err error) {
	return m.close()
}

func (m *UDPAdaptor) ReadMAVLinkPacket() (*common.MAVLinkPacket, error) {
	buf := make([]byte, 4096)

	for {
		got, _, err := m.sock.ReadFromUDP(buf)
		if err != nil {
			return nil, err
		}
		if got < 2 {
			continue
		}
		sof := buf[0]
		length := buf[1]

		if sof != common.MAVLINK_10_STX {
			continue
		}
		if length > 250 {
			continue
		}
		m := &common.MAVLinkPacket{}
		m.Decode(buf)
		return m, nil
	}
}

func (m *UDPAdaptor) Write(b []byte) (int, error) {
	addr, err := net.ResolveUDPAddr("udp", m.Port())
	if err != nil {
		return 0, err
	}

	return m.sock.WriteTo(b, addr)
}
