package mavlink

import (
	"bytes"
	"errors"
	"net"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
	mavlink "gobot.io/x/gobot/platforms/mavlink/common"
)

var _ gobot.Adaptor = (*UDPAdaptor)(nil)

type MockUDPConnection struct {
	TestClose       func() error
	TestReadFromUDP func([]byte) (int, *net.UDPAddr, error)
	TestWriteTo     func([]byte, net.Addr) (int, error)
}

func (m *MockUDPConnection) Close() error {
	return m.TestClose()
}

func (m *MockUDPConnection) ReadFromUDP(b []byte) (int, *net.UDPAddr, error) {
	return m.TestReadFromUDP(b)
}

func (m *MockUDPConnection) WriteTo(b []byte, a net.Addr) (int, error) {
	return m.TestWriteTo(b, a)
}

func NewMockUDPConnection() *MockUDPConnection {
	return &MockUDPConnection{
		TestClose: func() error {
			return nil
		},
		TestReadFromUDP: func([]byte) (int, *net.UDPAddr, error) {
			return 0, nil, nil
		},
		TestWriteTo: func([]byte, net.Addr) (int, error) {
			return 0, nil
		},
	}
}

func initTestMavlinkUDPAdaptor() *UDPAdaptor {
	m := NewUDPAdaptor(":14550")
	return m
}

func TestMavlinkUDPAdaptor(t *testing.T) {
	a := initTestMavlinkUDPAdaptor()
	gobottest.Assert(t, a.Port(), ":14550")
}

func TestMavlinkUDPAdaptorName(t *testing.T) {
	a := initTestMavlinkUDPAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "Mavlink"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestMavlinkUDPAdaptorConnectAndFinalize(t *testing.T) {
	a := initTestMavlinkUDPAdaptor()
	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestMavlinkUDPAdaptorWrite(t *testing.T) {
	a := initTestMavlinkUDPAdaptor()
	a.Connect()
	defer a.Finalize()

	m := NewMockUDPConnection()
	m.TestWriteTo = func([]byte, net.Addr) (int, error) {
		return 3, nil
	}
	a.sock = m

	i, err := a.Write([]byte{0x01, 0x02, 0x03})
	gobottest.Assert(t, i, 3)
	gobottest.Assert(t, err, nil)
}

func TestMavlinkReadMAVLinkReadDefaultPacket(t *testing.T) {
	a := initTestMavlinkUDPAdaptor()
	a.Connect()
	defer a.Finalize()

	m := NewMockUDPConnection()

	m.TestReadFromUDP = func(b []byte) (int, *net.UDPAddr, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{mavlink.MAVLINK_10_STX, 0x02, 0x03})
		copy(b, buf.Bytes())
		return buf.Len(), nil, nil
	}
	a.sock = m

	p, _ := a.ReadMAVLinkPacket()
	gobottest.Assert(t, p.Protocol, uint8(254))
}

func TestMavlinkReadMAVLinkPacketReadError(t *testing.T) {
	a := initTestMavlinkUDPAdaptor()
	a.Connect()
	defer a.Finalize()

	m := NewMockUDPConnection()

	i := 0
	m.TestReadFromUDP = func(b []byte) (int, *net.UDPAddr, error) {
		switch i {
		case 0:
			i = 1
			return 1, nil, nil
		case 1:
			i = 2
			buf := new(bytes.Buffer)
			buf.Write([]byte{0x01, 0x02, 0x03})
			copy(b, buf.Bytes())
			return buf.Len(), nil, nil
		case 2:
			i = 3
			buf := new(bytes.Buffer)
			buf.Write([]byte{mavlink.MAVLINK_10_STX, 255})
			copy(b, buf.Bytes())
			return buf.Len(), nil, nil
		}

		return 0, nil, errors.New("read error")
	}
	a.sock = m

	_, err := a.ReadMAVLinkPacket()
	gobottest.Assert(t, err, errors.New("read error"))
}
