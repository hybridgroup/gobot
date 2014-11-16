package neurosky

import (
	"io"

	"github.com/hybridgroup/gobot"
	"github.com/tarm/goserial"
)

var _ gobot.AdaptorInterface = (*NeuroskyAdaptor)(nil)

type NeuroskyAdaptor struct {
	gobot.Adaptor
	sp      io.ReadWriteCloser
	connect func(*NeuroskyAdaptor)
}

// NewNeuroskyAdaptor creates a neurosky adaptor with specified name
func NewNeuroskyAdaptor(name string, port string) *NeuroskyAdaptor {
	return &NeuroskyAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"NeuroskyAdaptor",
			port,
		),
		connect: func(n *NeuroskyAdaptor) {
			sp, err := serial.OpenPort(&serial.Config{Name: n.Port(), Baud: 57600})
			if err != nil {
				panic(err)
			}
			n.sp = sp
		},
	}
}

// Connect returns true if connection to device is successful
func (n *NeuroskyAdaptor) Connect() error {
	n.connect(n)
	n.SetConnected(true)
	return nil
}

// Finalize returns true if device finalization is successful
func (n *NeuroskyAdaptor) Finalize() error {
	n.sp.Close()
	n.SetConnected(false)
	return nil
}
