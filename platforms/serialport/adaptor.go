package serialport

import (
	"fmt"
	"io"

	"go.bug.st/serial"

	"gobot.io/x/gobot/v2"
)

// configuration contains all changeable attributes of the driver.
type configuration struct {
	name     string
	baudRate int
}

// Adaptor represents a Gobot Adaptor for the Serial Communication
type Adaptor struct {
	port string
	cfg  *configuration

	sp          io.ReadWriteCloser
	connectFunc func(string, int) (io.ReadWriteCloser, error)
}

// NewAdaptor returns a new adaptor given a port for the serial communication
func NewAdaptor(port string, opts ...optionApplier) *Adaptor {
	cfg := configuration{
		name:     gobot.DefaultName("Serial"),
		baudRate: 115200,
	}

	a := Adaptor{
		cfg:  &cfg,
		port: port,
		connectFunc: func(port string, baudRate int) (io.ReadWriteCloser, error) {
			return serial.Open(port, &serial.Mode{BaudRate: baudRate})
		},
	}

	for _, o := range opts {
		o.apply(a.cfg)
	}

	return &a
}

// WithName is used to replace the default name of the driver.
func WithName(name string) optionApplier {
	return nameOption(name)
}

// WithName is used to replace the default name of the driver.
func WithBaudRate(baudRate int) optionApplier {
	return baudRateOption(baudRate)
}

// Name returns the adaptors name
func (a *Adaptor) Name() string {
	return a.cfg.name
}

// SetName sets the adaptors name
// Deprecated: Please use option [serialport.WithName] instead.
func (a *Adaptor) SetName(n string) {
	WithName(n).apply(a.cfg)
}

// Connect initiates a connection to the serial port.
func (a *Adaptor) Connect() error {
	if a.sp != nil {
		return fmt.Errorf("serial port is already connected, try reconnect or run disconnect first")
	}

	sp, err := a.connectFunc(a.port, a.cfg.baudRate)
	if err != nil {
		return err
	}

	a.sp = sp
	return nil
}

// Finalize finalizes the adaptor by disconnect
func (a *Adaptor) Finalize() error {
	return a.Disconnect()
}

// Disconnect terminates the connection to the port.
func (a *Adaptor) Disconnect() error {
	if a.sp != nil {
		if err := a.sp.Close(); err != nil {
			return err
		}
		a.sp = nil
	}
	return nil
}

// Reconnect attempts to reconnect to the port. If the port is connected it will first close
// that connection and then establish a new connection.
func (a *Adaptor) Reconnect() error {
	if a.sp != nil {
		if err := a.Disconnect(); err != nil {
			return err
		}
	}
	return a.Connect()
}

// Port returns the adaptors port
func (a *Adaptor) Port() string { return a.port }

// IsConnected returns the connection state
func (a *Adaptor) IsConnected() bool {
	return a.sp != nil
}

// SerialRead reads from the port to the given reference
func (a *Adaptor) SerialRead(pData []byte) (int, error) {
	return a.sp.Read(pData)
}

// SerialWrite writes to the port
func (a *Adaptor) SerialWrite(data []byte) (int, error) {
	return a.sp.Write(data)
}
