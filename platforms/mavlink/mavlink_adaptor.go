package mavlink

import (
	"io"

	"github.com/hybridgroup/gobot"
	"github.com/tarm/goserial"
)

var _ gobot.Adaptor = (*MavlinkAdaptor)(nil)

type MavlinkAdaptor struct {
	name    string
	port    string
	sp      io.ReadWriteCloser
	connect func(string) (io.ReadWriteCloser, error)
}

// NewMavLinkAdaptor creates a new mavlink adaptor with specified name and port
func NewMavlinkAdaptor(name string, port string) *MavlinkAdaptor {
	return &MavlinkAdaptor{
		name: name,
		port: port,
		connect: func(port string) (io.ReadWriteCloser, error) {
			return serial.OpenPort(&serial.Config{Name: port, Baud: 57600})
		},
	}
}

func (m *MavlinkAdaptor) Name() string { return m.name }
func (m *MavlinkAdaptor) Port() string { return m.port }

// Connect returns true if connection to device is successful
func (m *MavlinkAdaptor) Connect() (errs []error) {
	if sp, err := m.connect(m.Port()); err != nil {
		return []error{err}
	} else {
		m.sp = sp
	}
	return
}

// Finalize returns true if connection to devices is closed successfully
func (m *MavlinkAdaptor) Finalize() (errs []error) {
	if err := m.sp.Close(); err != nil {
		return []error{err}
	}
	return
}
