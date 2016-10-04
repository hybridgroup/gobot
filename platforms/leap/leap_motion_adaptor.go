package leap

import (
	"io"

	"golang.org/x/net/websocket"
)

type Adaptor struct {
	name    string
	port    string
	ws      io.ReadWriteCloser
	connect func(string) (io.ReadWriteCloser, error)
}

// NewAdaptor creates a new leap motion adaptor using specified port
func NewAdaptor(port string) *Adaptor {
	return &Adaptor{
		name: "LeapMotion",
		port: port,
		connect: func(port string) (io.ReadWriteCloser, error) {
			return websocket.Dial("ws://"+port+"/v3.json", "", "http://"+port)
		},
	}
}
func (l *Adaptor) Name() string     { return l.name }
func (l *Adaptor) SetName(n string) { l.name = n }
func (l *Adaptor) Port() string     { return l.port }

// Connect returns true if connection to leap motion is established successfully
func (l *Adaptor) Connect() (errs []error) {
	if ws, err := l.connect(l.Port()); err != nil {
		return []error{err}
	} else {
		l.ws = ws
	}
	return
}

// Finalize ends connection to leap motion
func (l *Adaptor) Finalize() (errs []error) { return }
