package leap

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"github.com/hybridgroup/gobot"
)

type LeapMotionDriver struct {
	gobot.Driver
	Adaptor *LeapMotionAdaptor
}

func NewLeapMotionDriver(a *LeapMotionAdaptor) *LeapMotionDriver {
	return &LeapMotionDriver{
		Driver: gobot.Driver{
			Events: map[string]chan interface{}{
				"Message": make(chan interface{}),
			},
		},
		Adaptor: a,
	}
}

func (l *LeapMotionDriver) Start() bool {
	enableGestures := map[string]bool{"enableGestures": true}
	b, _ := json.Marshal(enableGestures)
	_, err := l.Adaptor.ws.Write(b)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			var msg []byte
			websocket.Message.Receive(l.Adaptor.ws, &msg)
			gobot.Publish(l.Events["Message"], l.ParseFrame(msg))
		}
	}()

	return true
}
func (me *LeapMotionDriver) Init() bool { return true }
func (me *LeapMotionDriver) Halt() bool { return true }
