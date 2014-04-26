package gobotSphero

import (
	"github.com/hybridgroup/gobot"
	"github.com/tarm/goserial"
	"io"
)

type SpheroAdaptor struct {
	gobot.Adaptor
	sp io.ReadWriteCloser
}

var connect = func(me *SpheroAdaptor) {
	c := &serial.Config{Name: me.Adaptor.Port, Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		panic(err)
	}
	me.sp = s
}

func (me *SpheroAdaptor) Connect() bool {
	connect(me)
	me.Connected = true
	return true
}

func (me *SpheroAdaptor) Reconnect() bool {
	if me.Connected == true {
		me.Disconnect()
	}
	return me.Connect()
}

func (me *SpheroAdaptor) Disconnect() bool {
	me.sp.Close()
	me.Connected = false
	return true
}

func (me *SpheroAdaptor) Finalize() bool {
	return true
}
