package serialport

import (
	"io"

	"go.bug.st/serial"

	"gobot.io/x/gobot/v2"
)

// Adaptor represents a Gobot Adaptor for the Serial Communication
type Adaptor struct {
	name      string
	port      string
	sp        io.ReadWriteCloser
	connected bool
	connect   func(string) (io.ReadWriteCloser, error)
}

// NewAdaptor returns a new adaptor given a port for the serial communication
func NewAdaptor(port string) *Adaptor {
	return &Adaptor{
		name: gobot.DefaultName("Serial"),
		port: port,
		connect: func(port string) (io.ReadWriteCloser, error) {
			return serial.Open(port, &serial.Mode{BaudRate: 115200})
		},
	}
}

// Name returns the Adaptor's name
func (a *Adaptor) Name() string { return a.name }

// SetName sets the Adaptor's name
func (a *Adaptor) SetName(n string) { a.name = n }

// Connect initiates a connection to the serial port.
func (a *Adaptor) Connect() error {
	sp, err := a.connect(a.Port())
	if err != nil {
		return err
	}

	a.sp = sp
	a.connected = true
	return nil
}

// Finalize finalizes the adaptor by disconnect
func (a *Adaptor) Finalize() error {
	return a.Disconnect()
}

// Disconnect terminates the connection to the port.
func (a *Adaptor) Disconnect() error {
	if a.connected {
		if err := a.sp.Close(); err != nil {
			return err
		}
		a.connected = false
	}
	return nil
}

// Reconnect attempts to reconnect to the port. If the port is connected it will first close
// that connection and then establish a new connection.
func (a *Adaptor) Reconnect() error {
	if a.connected {
		if err := a.Disconnect(); err != nil {
			return err
		}
	}
	return a.Connect()
}

// Port returns the Adaptor's port
func (a *Adaptor) Port() string { return a.port }

// IsConnected returns the connection state
func (a *Adaptor) IsConnected() bool {
	return a.connected
}

// SerialRead reads from the port
func (a *Adaptor) SerialRead(p []byte) (int, error) {
	return a.sp.Read(p)
}

// SerialWrite writes to the port
func (a *Adaptor) SerialWrite(p []byte) (int, error) {
	return a.sp.Write(p)
}
