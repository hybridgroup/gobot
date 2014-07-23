package neurosky

import (
	"io"

	"github.com/hybridgroup/gobot"
	"github.com/tarm/goserial"
)

type NeuroskyAdaptor struct {
	gobot.Adaptor
	sp      io.ReadWriteCloser
	connect func(*NeuroskyAdaptor)
}

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

func (n *NeuroskyAdaptor) Connect() bool {
	n.connect(n)
	n.SetConnected(true)
	return true
}

func (n *NeuroskyAdaptor) Finalize() bool {
	n.sp.Close()
	n.SetConnected(false)
	return true
}
