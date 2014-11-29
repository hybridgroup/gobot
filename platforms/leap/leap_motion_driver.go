package leap

import (
	"encoding/json"
	"io"

	"code.google.com/p/go.net/websocket"
	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*LeapMotionDriver)(nil)

type LeapMotionDriver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
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
		name:       name,
		connection: a,
		Eventer:    gobot.NewEventer(),
	}

	l.AddEvent("message")
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
		for {
			gobot.Publish(l.Event("message"), l.ParseFrame(receive(l.adaptor().ws)))
		}
	}()

	return
}

// Halt returns true if driver is halted succesfully
func (l *LeapMotionDriver) Halt() (errs []error) { return }
