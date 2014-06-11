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

func NewLeapMotionDriver(a *LeapMotionAdaptor, name string) *LeapMotionDriver {
	return &LeapMotionDriver{
		Driver: gobot.Driver{
			Name: name,
			Events: map[string]*gobot.Event{
				"Message": gobot.NewEvent(),
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
func (l *LeapMotionDriver) Init() bool { return true }
func (l *LeapMotionDriver) Halt() bool { return true }
