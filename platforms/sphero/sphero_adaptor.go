package sphero

import (
	"io"

	"github.com/hybridgroup/gobot"
	"github.com/tarm/goserial"
)

var _ gobot.Adaptor = (*SpheroAdaptor)(nil)

// Represents a Connection to a Sphero
type SpheroAdaptor struct {
	name      string
	port      string
	sp        io.ReadWriteCloser
	connected bool
	connect   func(string) (io.ReadWriteCloser, error)
}

// NewSpheroAdaptor returns a new SpheroAdaptor given a name and port
func NewSpheroAdaptor(name string, port string) *SpheroAdaptor {
	return &SpheroAdaptor{
		name: name,
		port: port,
		connect: func(port string) (io.ReadWriteCloser, error) {
			return serial.OpenPort(&serial.Config{Name: port, Baud: 115200})
		},
	}
}
func (a *SpheroAdaptor) Name() string { return a.name }
func (a *SpheroAdaptor) Port() string { return a.port }

// Connect initiates a connection to the Sphero. Returns true on successful connection.
func (a *SpheroAdaptor) Connect() (errs []error) {
	if sp, err := a.connect(a.Port()); err != nil {
		return []error{err}
	} else {
		a.sp = sp
		a.connected = true
	}
	return
}

// Reconnect attempts to reconnect to the Sphero. If the Sphero has an active connection
// it will first close that connection and then establish a new connection.
// Returns true on Successful reconnection
func (a *SpheroAdaptor) Reconnect() (errs []error) {
	if a.connected {
		a.Disconnect()
	}
	return a.Connect()
}

// Disconnect terminates the connection to the Sphero. Returns true on successful disconnect.
func (a *SpheroAdaptor) Disconnect() (errs []error) {
	if a.connected {
		if err := a.sp.Close(); err != nil {
			return []error{err}
		}
		a.connected = false
	}
	return
}

// Finalize finalizes the SpheroAdaptor
func (a *SpheroAdaptor) Finalize() (errs []error) {
	return a.Disconnect()
}
