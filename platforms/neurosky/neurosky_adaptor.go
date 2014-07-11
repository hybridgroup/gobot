package neurosky

import (
	"github.com/hybridgroup/gobot"
	"github.com/tarm/goserial"
	"io"
)

type NeuroskyAdaptor struct {
	gobot.Adaptor
	sp      io.ReadWriteCloser
	connect func(string) io.ReadWriteCloser
}

func NewNeuroskyAdaptor(name string, port string) *NeuroskyAdaptor {
	return &NeuroskyAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"NeuroskyAdaptor",
			port,
		),
		connect: func(port string) io.ReadWriteCloser {
			sp, err := serial.OpenPort(&serial.Config{Name: port, Baud: 57600})
			if err != nil {
				panic(err)
			}
			return sp
		},
	}
}

func (n *NeuroskyAdaptor) Connect() bool {
	n.sp = n.connect(n.Adaptor.Port())
	n.SetConnected(true)
	return true
}

func (n *NeuroskyAdaptor) Finalize() bool {
	n.sp.Close()
	n.SetConnected(false)
	return true
}
