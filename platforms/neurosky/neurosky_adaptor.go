package neurosky

import (
	"io"

	"github.com/hybridgroup/gobot"
	"github.com/tarm/goserial"
)

var _ gobot.Adaptor = (*NeuroskyAdaptor)(nil)

type NeuroskyAdaptor struct {
	name    string
	port    string
	sp      io.ReadWriteCloser
	connect func(*NeuroskyAdaptor) (err error)
}

// NewNeuroskyAdaptor creates a neurosky adaptor with specified name
func NewNeuroskyAdaptor(name string, port string) *NeuroskyAdaptor {
	return &NeuroskyAdaptor{
		name: name,
		port: port,
		connect: func(n *NeuroskyAdaptor) (err error) {
			sp, err := serial.OpenPort(&serial.Config{Name: n.Port(), Baud: 57600})
			if err != nil {
				return err
			}
			n.sp = sp
			return
		},
	}
}
func (n *NeuroskyAdaptor) Name() string { return n.name }
func (n *NeuroskyAdaptor) Port() string { return n.port }

// Connect returns true if connection to device is successful
func (n *NeuroskyAdaptor) Connect() (errs []error) {
	if err := n.connect(n); err != nil {
		return []error{err}
	}
	return
}

// Finalize returns true if device finalization is successful
func (n *NeuroskyAdaptor) Finalize() (errs []error) {
	if err := n.sp.Close(); err != nil {
		return []error{err}
	}
	return
}
