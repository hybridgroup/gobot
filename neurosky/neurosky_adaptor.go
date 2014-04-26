package gobotNeurosky

import (
	"github.com/hybridgroup/gobot"
	"github.com/tarm/goserial"
	"io"
)

type NeuroskyAdaptor struct {
	gobot.Adaptor
	sp io.ReadWriteCloser
}

func (me *NeuroskyAdaptor) Connect() bool {
	c := &serial.Config{Name: me.Adaptor.Port, Baud: 57600}
	s, err := serial.OpenPort(c)
	if err != nil {
		panic(err)
	}
	me.sp = s
	me.Connected = true
	return true
}

func (me *NeuroskyAdaptor) Reconnect() bool {
	if me.Connected == true {
		me.Disconnect()
	}
	return me.Connect()
}

func (me *NeuroskyAdaptor) Disconnect() bool {
	me.sp.Close()
	me.Connected = false
	return true
}

func (me *NeuroskyAdaptor) Finalize() bool {
	return true
}
