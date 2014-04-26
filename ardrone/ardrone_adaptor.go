package ardrone

import (
	"github.com/hybridgroup/go-ardrone/client"
	"github.com/hybridgroup/gobot"
)

type drone interface{}

type ArdroneAdaptor struct {
	gobot.Adaptor
	ardrone drone
}

var connect = func(me *ArdroneAdaptor) {
	ardrone, err := ardrone.Connect(ardrone.DefaultConfig())
	if err != nil {
		panic(err)
	}
	me.ardrone = ardrone
}

func (me *ArdroneAdaptor) Connect() bool {
	connect(me)
	return true
}

func (me *ArdroneAdaptor) Reconnect() bool {
	return true
}

func (me *ArdroneAdaptor) Disconnect() bool {
	return true
}

func (me *ArdroneAdaptor) Finalize() bool {
	return true
}

func (me *ArdroneAdaptor) Drone() drone {
	return me.ardrone
}
