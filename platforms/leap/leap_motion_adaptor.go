package leap

import (
	"io"

	"gobot.io/x/gobot"

	"golang.org/x/net/websocket"
)

// Adaptor is the Gobot Adaptor connection to the Leap Motion
type Adaptor struct {
	name    string
	port    string
	ws      io.ReadWriteCloser
	connect func(string) (io.ReadWriteCloser, error)
}

// NewAdaptor creates a new leap motion adaptor using specified port,
// which is this case is the host IP or name of the Leap Motion daemon
func NewAdaptor(port string) *Adaptor {
	return &Adaptor{
		name: gobot.DefaultName("LeapMotion"),
		port: port,
		connect: func(host string) (io.ReadWriteCloser, error) {
			return websocket.Dial("ws://"+host+"/v3.json", "", "http://"+host)
		},
	}
}

// Name returns the Adaptor Name
func (l *Adaptor) Name() string { return l.name }

// SetName sets the Adaptor Name
func (l *Adaptor) SetName(n string) { l.name = n }

// Port returns the Adaptor Port which is this case is the host IP or name
func (l *Adaptor) Port() string { return l.port }

// Connect returns true if connection to leap motion is established successfully
func (l *Adaptor) Connect() (err error) {
	ws, e := l.connect(l.Port())
	if e != nil {
		return e
	}

	l.ws = ws
	return
}

// Finalize ends connection to leap motion
func (l *Adaptor) Finalize() (err error) { return }
