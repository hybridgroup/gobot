package leap

import (
	"fmt"
	"io"

	"code.google.com/p/go.net/websocket"
	"github.com/hybridgroup/gobot"
)

var _ gobot.Adaptor = (*LeapMotionAdaptor)(nil)

type LeapMotionAdaptor struct {
	name    string
	port    string
	ws      io.ReadWriteCloser
	connect func(*LeapMotionAdaptor) (err error)
}

// NewLeapMotionAdaptor creates a new leap motion adaptor using specified name and port
func NewLeapMotionAdaptor(name string, port string) *LeapMotionAdaptor {
	return &LeapMotionAdaptor{
		name: name,
		port: port,
		connect: func(l *LeapMotionAdaptor) (err error) {
			ws, err := websocket.Dial(
				fmt.Sprintf("ws://%v/v3.json", l.Port()),
				"",
				fmt.Sprintf("http://%v", l.Port()),
			)
			if err != nil {
				return err
			}
			l.ws = ws
			return
		},
	}
}
func (l *LeapMotionAdaptor) Name() string { return l.name }
func (l *LeapMotionAdaptor) Port() string { return l.port }

// Connect returns true if connection to leap motion is established succesfully
func (l *LeapMotionAdaptor) Connect() (errs []error) {
	if err := l.connect(l); err != nil {
		return []error{err}
	}
	return
}

// Finalize ends connection to leap motion
func (l *LeapMotionAdaptor) Finalize() (errs []error) { return }
