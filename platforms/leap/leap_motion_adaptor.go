package leap

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/hybridgroup/gobot"
	"io"
)

type LeapMotionAdaptor struct {
	gobot.Adaptor
	ws      io.ReadWriteCloser
	connect func(*LeapMotionAdaptor)
}

// NewLeapMotionAdaptor creates a new leap motion adaptor using specified name and port
func NewLeapMotionAdaptor(name string, port string) *LeapMotionAdaptor {
	return &LeapMotionAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"LeapMotionAdaptor",
			port,
		),
		connect: func(l *LeapMotionAdaptor) {
			ws, err := websocket.Dial(
				fmt.Sprintf("ws://%v/v3.json", l.Port()),
				"",
				fmt.Sprintf("http://%v", l.Port()),
			)
			if err != nil {
				panic(err)
			}
			l.ws = ws
		},
	}
}

// Connect returns true if connection to leap motion is established succesfully
func (l *LeapMotionAdaptor) Connect() bool {
	l.connect(l)
	l.SetConnected(true)
	return true
}

// Finalize ends connection to leap motion
func (l *LeapMotionAdaptor) Finalize() bool { return true }
