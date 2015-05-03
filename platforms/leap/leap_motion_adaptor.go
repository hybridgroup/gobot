package leap

import (
	"io"

	"github.com/hybridgroup/gobot"
	"golang.org/x/net/websocket"
)

var _ gobot.Adaptor = (*LeapMotionAdaptor)(nil)

type LeapMotionAdaptor struct {
	name    string
	port    string
	ws      io.ReadWriteCloser
	connect func(string) (io.ReadWriteCloser, error)
}

// NewLeapMotionAdaptor creates a new leap motion adaptor using specified name and port
func NewLeapMotionAdaptor(name string, port string) *LeapMotionAdaptor {
	return &LeapMotionAdaptor{
		name: name,
		port: port,
		connect: func(port string) (io.ReadWriteCloser, error) {
			return websocket.Dial("ws://"+port+"/v3.json", "", "http://"+port)
		},
	}
}
func (l *LeapMotionAdaptor) Name() string { return l.name }
func (l *LeapMotionAdaptor) Port() string { return l.port }

// Connect returns true if connection to leap motion is established succesfully
func (l *LeapMotionAdaptor) Connect() (errs []error) {
	if ws, err := l.connect(l.Port()); err != nil {
		return []error{err}
	} else {
		l.ws = ws
	}
	return
}

// Finalize ends connection to leap motion
func (l *LeapMotionAdaptor) Finalize() (errs []error) { return }
