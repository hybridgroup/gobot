package ardrone

import (
	client "github.com/hybridgroup/go-ardrone/client"
	"github.com/hybridgroup/gobot"
)

type drone interface{}

type ArdroneAdaptor struct {
	gobot.Adaptor
	drone   drone
	connect func(*ArdroneAdaptor)
}

func NewArdroneAdaptor(name string) *ArdroneAdaptor {
	return &ArdroneAdaptor{
		Adaptor: gobot.Adaptor{
			Name: name,
		},
		connect: func(a *ArdroneAdaptor) {
			d, err := client.Connect(client.DefaultConfig())
			if err != nil {
				panic(err)
			}
			a.drone = d
		},
	}
}

func (a *ArdroneAdaptor) Connect() bool {
	a.connect(a)
	return true
}

func (a *ArdroneAdaptor) Reconnect() bool {
	return true
}

func (a *ArdroneAdaptor) Disconnect() bool {
	return true
}

func (a *ArdroneAdaptor) Finalize() bool {
	return true
}

func (a *ArdroneAdaptor) Drone() drone {
	return a.drone
}
