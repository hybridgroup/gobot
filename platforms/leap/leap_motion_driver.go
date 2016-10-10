package leap

import (
	"encoding/json"
	"io"

	"github.com/hybridgroup/gobot"
	"golang.org/x/net/websocket"
)

const (
	// Message event
	MessageEvent = "message"
	// Hand event
	HandEvent = "hand"
	// Gesture event
	GestureEvent = "gesture"
)

type Driver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}

var receive = func(ws io.ReadWriteCloser, msg *[]byte) {
	websocket.Message.Receive(ws.(*websocket.Conn), msg)
}

// NewLDriver creates a new leap motion driver
//
// Adds the following events:
//		"message" - Gets triggered when receiving a message from leap motion
//		"hand" - Gets triggered per-message when leap motion detects a hand
//		"gesture" - Gets triggered per-message when leap motion detects a hand
func NewDriver(a *Adaptor) *Driver {
	l := &Driver{
		name:       "LeapMotion",
		connection: a,
		Eventer:    gobot.NewEventer(),
	}

	l.AddEvent(MessageEvent)
	l.AddEvent(HandEvent)
	l.AddEvent(GestureEvent)
	return l
}
func (l *Driver) Name() string                 { return l.name }
func (l *Driver) SetName(n string)             { l.name = n }
func (l *Driver) Connection() gobot.Connection { return l.connection }

// adaptor returns leap motion adaptor
func (l *Driver) adaptor() *Adaptor {
	return l.Connection().(*Adaptor)
}

// Start inits leap motion driver by enabling gestures
// and listening from incoming messages.
//
// Publishes the following events:
//		"message" - Emits Frame on new message received from Leap.
//		"hand" - Emits Hand when detected in message from Leap.
//		"gesture" - Emits Gesture when detected in message from Leap.
func (l *Driver) Start() (errs []error) {
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
			l.Publish(MessageEvent, frame)

			for _, hand := range frame.Hands {
				l.Publish(HandEvent, hand)
			}

			for _, gesture := range frame.Gestures {
				l.Publish(GestureEvent, gesture)
			}
		}
	}()

	return
}

// Halt returns true if driver is halted successfully
func (l *Driver) Halt() (errs []error) { return }
