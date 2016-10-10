package sphero

import (
	"io"

	"github.com/tarm/serial"
)

// Represents a Connection to a Sphero
type Adaptor struct {
	name      string
	port      string
	sp        io.ReadWriteCloser
	connected bool
	connect   func(string) (io.ReadWriteCloser, error)
}

// NewAdaptor returns a new Sphero Adaptor given a port
func NewAdaptor(port string) *Adaptor {
	return &Adaptor{
		name: "Sphero",
		port: port,
		connect: func(port string) (io.ReadWriteCloser, error) {
			return serial.OpenPort(&serial.Config{Name: port, Baud: 115200})
		},
	}
}

func (a *Adaptor) Name() string     { return a.name }
func (a *Adaptor) SetName(n string) { a.name = n }
func (a *Adaptor) Port() string     { return a.port }
func (a *Adaptor) SetPort(p string) { a.port = p }

// Connect initiates a connection to the Sphero. Returns true on successful connection.
func (a *Adaptor) Connect() (errs []error) {
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
func (a *Adaptor) Reconnect() (errs []error) {
	if a.connected {
		a.Disconnect()
	}
	return a.Connect()
}

// Disconnect terminates the connection to the Sphero. Returns true on successful disconnect.
func (a *Adaptor) Disconnect() (errs []error) {
	if a.connected {
		if err := a.sp.Close(); err != nil {
			return []error{err}
		}
		a.connected = false
	}
	return
}

// Finalize finalizes the Sphero Adaptor
func (a *Adaptor) Finalize() (errs []error) {
	return a.Disconnect()
}
