package leap

import (
	"encoding/json"

	"code.google.com/p/go.net/websocket"
	"github.com/hybridgroup/gobot"
)

type LeapMotionDriver struct {
	gobot.Driver
}

func NewLeapMotionDriver(a *LeapMotionAdaptor, name string) *LeapMotionDriver {
	l := &LeapMotionDriver{
		Driver: *gobot.NewDriver(
			name,
			"LeapMotionDriver",
			a,
		),
	}

	l.AddEvent("message")
	return l
}

func (l *LeapMotionDriver) adaptor() *LeapMotionAdaptor {
	return l.Adaptor().(*LeapMotionAdaptor)
}
func (l *LeapMotionDriver) Start() bool {
	enableGestures := map[string]bool{"enableGestures": true}
	b, _ := json.Marshal(enableGestures)
	_, err := l.adaptor().ws.Write(b)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			var msg []byte
			websocket.Message.Receive(l.adaptor().ws, &msg)
			gobot.Publish(l.Event("message"), l.ParseFrame(msg))
		}
	}()

	return true
}
func (l *LeapMotionDriver) Init() bool { return true }
func (l *LeapMotionDriver) Halt() bool { return true }
