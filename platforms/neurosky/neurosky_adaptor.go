// Package neurosky is the Gobot platform for the Neurosky Mindwave EEG
package neurosky

import (
	"io"

	serial "go.bug.st/serial.v1"
)

// Adaptor is the Gobot Adaptor for the Neurosky Mindwave
type Adaptor struct {
	name    string
	port    string
	sp      io.ReadWriteCloser
	connect func(*Adaptor) (io.ReadWriteCloser, error)
}

// NewAdaptor creates a neurosky adaptor with specified port
func NewAdaptor(port string) *Adaptor {
	return &Adaptor{
		name: "Neurosky",
		port: port,
		connect: func(n *Adaptor) (io.ReadWriteCloser, error) {
			return serial.Open(n.Port(), &serial.Mode{BaudRate: 57600})
		},
	}
}

// Name returns the Adaptor Name
func (n *Adaptor) Name() string { return n.name }

// SetName sets the Adaptor Name
func (n *Adaptor) SetName(name string) { n.name = name }

// Port returns the Adaptor port
func (n *Adaptor) Port() string { return n.port }

// Connect returns true if connection to device is successful
func (n *Adaptor) Connect() error {
	sp, err := n.connect(n)
	if err != nil {
		return err
	}

	n.sp = sp
	return nil
}

// Finalize returns true if device finalization is successful
func (n *Adaptor) Finalize() (err error) {
	err = n.sp.Close()
	return
}
