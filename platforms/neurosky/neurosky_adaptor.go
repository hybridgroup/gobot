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

func NewNeuroskyAdaptor() *NeuroskyAdaptor {
	return &NeuroskyAdaptor{
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
	n.sp = n.connect(n.Adaptor.Port)
	n.Connected = true
	return true
}

func (n *NeuroskyAdaptor) Finalize() bool {
	n.sp.Close()
	n.Connected = false
	return true
}
