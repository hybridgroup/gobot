package leap

import (
	"encoding/json"
	"io"

	"github.com/hybridgroup/gobot"
	"golang.org/x/net/websocket"
)

var _ gobot.Driver = (*LeapMotionDriver)(nil)

type LeapMotionDriver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}

var receive = func(ws io.ReadWriteCloser, msg *[]byte) {
	websocket.Message.Receive(ws.(*websocket.Conn), msg)
}

// NewLeapMotionDriver creates a new leap motion driver with specified name
//
// Adds the following events:
//		"message" - Gets triggered when receiving a message from leap motion
//		"hand" - Gets triggered per-message when leap motion detects a hand
//		"gesture" - Gets triggered per-message when leap motion detects a hand
func NewLeapMotionDriver(a *LeapMotionAdaptor, name string) *LeapMotionDriver {
	l := &LeapMotionDriver{
		name:       name,
		connection: a,
		Eventer:    gobot.NewEventer(),
	}

	l.AddEvent("message")
	l.AddEvent("hand")
	l.AddEvent("gesture")
	return l
}
func (l *LeapMotionDriver) Name() string                 { return l.name }
func (l *LeapMotionDriver) Connection() gobot.Connection { return l.connection }

// adaptor returns leap motion adaptor
func (l *LeapMotionDriver) adaptor() *LeapMotionAdaptor {
	return l.Connection().(*LeapMotionAdaptor)
}

// Start inits leap motion driver by enabling gestures
// and listening from incoming messages.
//
// Publishes the following events:
//		"message" - Emits Frame on new message received from Leap.
//		"hand" - Emits Hand when detected in message from Leap.
//		"gesture" - Emits Gesture when detected in message from Leap.
func (l *LeapMotionDriver) Start() (errs []error) {
	enableGestures := map[string]bool{"enableGestures": true}
	b, err := json.Marshal(enableGestures)
	if err != nil {
		return []error{err}
	}
	_, err = l.adaptor().ws.Write(b)
	if err != nil {
		return []error{err}
	}

	go func() {
		var msg []byte
		var frame Frame
		for {
			receive(l.adaptor().ws, &msg)
			frame = l.ParseFrame(msg)
			gobot.Publish(l.Event("message"), frame)

			for _, hand := range frame.Hands {
				gobot.Publish(l.Event("hand"), hand)
			}

			for _, gesture := range frame.Gestures {
				gobot.Publish(l.Event("gesture"), gesture)
			}
		}
	}()

	return
}

// Halt returns true if driver is halted succesfully
func (l *LeapMotionDriver) Halt() (errs []error) { return }
