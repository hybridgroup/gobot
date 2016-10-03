package mavlink

import (
	"io"

	"github.com/tarm/serial"
)

type Adaptor struct {
	name    string
	port    string
	sp      io.ReadWriteCloser
	connect func(string) (io.ReadWriteCloser, error)
}

// NewAdaptor creates a new mavlink adaptor with specified port
func NewAdaptor(port string) *Adaptor {
	return &Adaptor{
		port: port,
		connect: func(port string) (io.ReadWriteCloser, error) {
			return serial.OpenPort(&serial.Config{Name: port, Baud: 57600})
		},
	}
}

func (m *Adaptor) Name() string     { return m.name }
func (m *Adaptor) SetName(n string) { m.name = n }
func (m *Adaptor) Port() string     { return m.port }

// Connect returns true if connection to device is successful
func (m *Adaptor) Connect() (errs []error) {
	if sp, err := m.connect(m.Port()); err != nil {
		return []error{err}
	} else {
		m.sp = sp
	}
	return
}

// Finalize returns true if connection to devices is closed successfully
func (m *Adaptor) Finalize() (errs []error) {
	if err := m.sp.Close(); err != nil {
		return []error{err}
	}
	return
}
