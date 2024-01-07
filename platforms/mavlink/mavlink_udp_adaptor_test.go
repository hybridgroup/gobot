package mavlink

import (
	"bytes"
	"errors"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	mavlink "gobot.io/x/gobot/v2/platforms/mavlink/common"
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
	assert.Equal(t, ":14550", a.Port())
}

func TestMavlinkUDPAdaptorName(t *testing.T) {
	a := initTestMavlinkUDPAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "Mavlink"))
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestMavlinkUDPAdaptorConnectAndFinalize(t *testing.T) {
	a := initTestMavlinkUDPAdaptor()
	require.NoError(t, a.Connect())
	require.NoError(t, a.Finalize())
}

func TestMavlinkUDPAdaptorWrite(t *testing.T) {
	a := initTestMavlinkUDPAdaptor()
	_ = a.Connect()
	defer func() { _ = a.Finalize() }()

	m := NewMockUDPConnection()
	m.TestWriteTo = func([]byte, net.Addr) (int, error) {
		return 3, nil
	}
	a.sock = m

	i, err := a.Write([]byte{0x01, 0x02, 0x03})
	assert.Equal(t, 3, i)
	require.NoError(t, err)
}

func TestMavlinkReadMAVLinkReadDefaultPacket(t *testing.T) {
	a := initTestMavlinkUDPAdaptor()
	_ = a.Connect()
	defer func() { _ = a.Finalize() }()

	m := NewMockUDPConnection()

	m.TestReadFromUDP = func(b []byte) (int, *net.UDPAddr, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{mavlink.MAVLINK_10_STX, 0x02, 0x03})
		copy(b, buf.Bytes())
		return buf.Len(), nil, nil
	}
	a.sock = m

	p, _ := a.ReadMAVLinkPacket()
	assert.Equal(t, uint8(254), p.Protocol)
}

func TestMavlinkReadMAVLinkPacketReadError(t *testing.T) {
	a := initTestMavlinkUDPAdaptor()
	_ = a.Connect()
	defer func() { _ = a.Finalize() }()

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
	require.ErrorContains(t, err, "read error")
}
