package neurosky

import (
	"io"

	"github.com/tarm/serial"
)

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
			return serial.OpenPort(&serial.Config{Name: n.Port(), Baud: 57600})
		},
	}
}

func (n *Adaptor) Name() string        { return n.name }
func (n *Adaptor) SetName(name string) { n.name = name }
func (n *Adaptor) Port() string        { return n.port }

// Connect returns true if connection to device is successful
func (n *Adaptor) Connect() error {
	if sp, err := n.connect(n); err != nil {
		return err
	} else {
		n.sp = sp
	}
	return nil
}

// Finalize returns true if device finalization is successful
func (n *Adaptor) Finalize() (err error) {
	err = n.sp.Close()
	return
}
