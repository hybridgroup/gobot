package sphero

import (
	"github.com/hybridgroup/gobot"
	"github.com/tarm/goserial"
	"io"
)

type SpheroAdaptor struct {
	gobot.Adaptor
	sp      io.ReadWriteCloser
	connect func(*SpheroAdaptor)
}

func NewSpheroAdaptor(name string, port string) *SpheroAdaptor {
	return &SpheroAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"SpheroAdaptor",
			port,
		),
		connect: func(a *SpheroAdaptor) {
			c := &serial.Config{Name: a.Port(), Baud: 115200}
			s, err := serial.OpenPort(c)
			if err != nil {
				panic(err)
			}
			a.sp = s
		},
	}
}

func (a *SpheroAdaptor) Connect() bool {
	a.connect(a)
	a.SetConnected(true)
	return true
}

func (a *SpheroAdaptor) Reconnect() bool {
	if a.Connected() == true {
		a.Disconnect()
	}
	return a.Connect()
}

func (a *SpheroAdaptor) Disconnect() bool {
	a.sp.Close()
	a.SetConnected(false)
	return true
}

func (a *SpheroAdaptor) Finalize() bool {
	return true
}
