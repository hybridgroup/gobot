package leap

import (
	"encoding/json"
	"io"

	"code.google.com/p/go.net/websocket"
	"github.com/hybridgroup/gobot"
)

type LeapMotionDriver struct {
	gobot.Driver
}

var receive = func(ws io.ReadWriteCloser) []byte {
	var msg []byte
	websocket.Message.Receive(ws.(*websocket.Conn), &msg)
	return msg
}

// NewLeapMotionDriver creates a new leap motion driver with specified name
//
// Adds the following events:
//		"message" - Gets triggered when receiving a message from leap motion
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

// adaptor returns leap motion adaptor
func (l *LeapMotionDriver) adaptor() *LeapMotionAdaptor {
	return l.Adaptor().(*LeapMotionAdaptor)
}

// Start inits leap motion driver by enabling gestures
// and listening from incoming messages.
//
// Publishes the following events:
//		"message" - Emits Frame on new message received from Leap.
func (l *LeapMotionDriver) Start() bool {
	enableGestures := map[string]bool{"enableGestures": true}
	b, _ := json.Marshal(enableGestures)
	_, err := l.adaptor().ws.Write(b)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			gobot.Publish(l.Event("message"), l.ParseFrame(receive(l.adaptor().ws)))
		}
	}()

	return true
}

// Halt returns true if driver is halted succesfully
func (l *LeapMotionDriver) Halt() bool { return true }
