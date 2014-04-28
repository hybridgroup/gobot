package leap

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/hybridgroup/gobot"
)

type LeapMotionAdaptor struct {
	gobot.Adaptor
	ws      *websocket.Conn
	connect func(*LeapMotionAdaptor)
}

func NewLeapMotionAdaptor() *LeapMotionAdaptor {
	return &LeapMotionAdaptor{
		connect: func(l *LeapMotionAdaptor) {
			origin := fmt.Sprintf("http://%v", l.Port)
			url := fmt.Sprintf("ws://%v/v3.json", l.Port)
			ws, err := websocket.Dial(url, "", origin)
			if err != nil {
				panic(err)
			}
			l.ws = ws
		},
	}
}

func (l *LeapMotionAdaptor) Connect() bool {
	l.connect(l)
	l.Connected = true
	return true
}
func (me *LeapMotionAdaptor) Finalize() bool { return true }
