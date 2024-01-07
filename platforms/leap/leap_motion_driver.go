package leap

import (
	"encoding/json"
	"io"
	"log"

	"golang.org/x/net/websocket"

	"gobot.io/x/gobot/v2"
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
//
//	"message" - Gets triggered when receiving a message from leap motion
//	"hand" - Gets triggered per-message when leap motion detects a hand
//	"gesture" - Gets triggered per-message when leap motion detects a gesture
func NewDriver(a *Adaptor) *Driver {
	l := &Driver{
		name:       gobot.DefaultName("LeapMotion"),
		connection: a,
		Eventer:    gobot.NewEventer(),
		receive: func(ws io.ReadWriteCloser, msg *[]byte) {
			//nolint:forcetypeassert // ok here
			if err := websocket.Message.Receive(ws.(*websocket.Conn), msg); err != nil {
				panic(err)
			}
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
	//nolint:forcetypeassert // ok here
	return l.Connection().(*Adaptor)
}

func enableFeature(l *Driver, feature string) error {
	command := map[string]bool{feature: true}
	b, err := json.Marshal(command)
	if err != nil {
		return err
	}
	if _, err = l.adaptor().ws.Write(b); err != nil {
		return err
	}

	return nil
}

// Start inits leap motion driver by enabling gestures
// and listening from incoming messages.
//
// Publishes the following events:
//
//	"message" - Emits Frame on new message received from Leap.
//	"hand" - Emits Hand when detected in message from Leap.
//	"gesture" - Emits Gesture when detected in message from Leap.
func (l *Driver) Start() error {
	if err := enableFeature(l, "enableGestures"); err != nil {
		return err
	}
	if err := enableFeature(l, "background"); err != nil {
		return err
	}

	go func() {
		var msg []byte
		var frame Frame
		var err error
		for {
			l.receive(l.adaptor().ws, &msg)
			frame, err = l.ParseFrame(msg)
			if err != nil {
				log.Println(err)
				continue
			}

			l.Publish(MessageEvent, frame)

			for _, hand := range frame.Hands {
				l.Publish(HandEvent, hand)
			}

			for _, gesture := range frame.Gestures {
				l.Publish(GestureEvent, gesture)
			}
		}
	}()

	return nil
}

// Halt returns nil if driver is halted successfully
func (l *Driver) Halt() error { return nil }
