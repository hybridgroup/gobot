package leap

import (
	"encoding/json"
	"io"

	"gobot.io/x/gobot"
	"golang.org/x/net/websocket"
)

const (
	// MessageEvent event
	MessageEvent = "message"
	// HandEvent event
	HandEvent = "hand"
	// GestureEvent event
	GestureEvent = "gesture"
)

// Driver the Gobot software device to the Leap Motion
type Driver struct {
	name       string
	connection gobot.Connection
	receive    func(ws io.ReadWriteCloser, msg *[]byte)
	gobot.Eventer
}

// NewDriver creates a new leap motion driver
//
// Adds the following events:
//		"message" - Gets triggered when receiving a message from leap motion
//		"hand" - Gets triggered per-message when leap motion detects a hand
//		"gesture" - Gets triggered per-message when leap motion detects a hand
func NewDriver(a *Adaptor) *Driver {
	l := &Driver{
		name:       gobot.DefaultName("LeapMotion"),
		connection: a,
		Eventer:    gobot.NewEventer(),
		receive: func(ws io.ReadWriteCloser, msg *[]byte) {
			websocket.Message.Receive(ws.(*websocket.Conn), msg)
		},
	}

	l.AddEvent(MessageEvent)
	l.AddEvent(HandEvent)
	l.AddEvent(GestureEvent)
	return l
}

// Name returns the Driver Name
func (l *Driver) Name() string { return l.name }

// SetName sets the Driver Name
func (l *Driver) SetName(n string) { l.name = n }

// Connection returns the Driver's Connection
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
func (l *Driver) Start() (err error) {
	enableGestures := map[string]bool{"enableGestures": true}
	b, e := json.Marshal(enableGestures)
	if e != nil {
		return e
	}
	_, e = l.adaptor().ws.Write(b)
	if e != nil {
		return e
	}

	go func() {
		var msg []byte
		var frame Frame
		for {
			l.receive(l.adaptor().ws, &msg)
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

// Halt returns nil if driver is halted successfully
func (l *Driver) Halt() (errs error) { return }
