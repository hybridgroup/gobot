package gobotLeap

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/hybridgroup/gobot"
)

type LeapAdaptor struct {
	gobot.Adaptor
	Leap *websocket.Conn
}

func (me *LeapAdaptor) Connect() bool {
	origin := fmt.Sprintf("http://%v", me.Port)
	url := fmt.Sprintf("ws://%v/v3.json", me.Port)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		panic(err)
	}
	me.Leap = ws
	me.Connected = true
	return true
}
func (me *LeapAdaptor) Reconnect() bool  { return me.Connect() }
func (me *LeapAdaptor) Disconnect() bool { return false }
func (me *LeapAdaptor) Finalize() bool   { return false }
