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
func (a *Adaptor) Connect() (err error) {
	if sp, e := a.connect(a.Port()); e != nil {
		return e
	} else {
		a.sp = sp
		a.connected = true
	}
	return
}

// Reconnect attempts to reconnect to the Sphero. If the Sphero has an active connection
// it will first close that connection and then establish a new connection.
// Returns true on Successful reconnection
func (a *Adaptor) Reconnect() (err error) {
	if a.connected {
		a.Disconnect()
	}
	return a.Connect()
}

// Disconnect terminates the connection to the Sphero. Returns true on successful disconnect.
func (a *Adaptor) Disconnect() error {
	if a.connected {
		if e := a.sp.Close(); e != nil {
			return e
		}
		a.connected = false
	}
	return nil
}

// Finalize finalizes the Sphero Adaptor
func (a *Adaptor) Finalize() error {
	return a.Disconnect()
}
