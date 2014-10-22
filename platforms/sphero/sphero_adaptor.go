package sphero

import (
	"io"

	"github.com/hybridgroup/gobot"
	"github.com/tarm/goserial"
)

// Represents a Connection to a Sphero
type SpheroAdaptor struct {
	gobot.Adaptor
	sp      io.ReadWriteCloser
	connect func(*SpheroAdaptor)
}

// NewSpheroAdaptor returns a new SpheroAdaptor given a name and port
func NewSpheroAdaptor(name string, port string) *SpheroAdaptor {
	return &SpheroAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"SpheroAdaptor",
			port,
		),
		connect: func(a *SpheroAdaptor) {
			c := &serial.Config{Name: a.Port(), Baud: 115200}
			s, err := serial.OpenPort(c)
			if err != nil {
				panic(err)
			}
			a.sp = s
		},
	}
}

// Connect initiates a connection to the Sphero. Returns true on successful connection.
func (a *SpheroAdaptor) Connect() bool {
	a.connect(a)
	a.SetConnected(true)
	return true
}

// Reconnect attempts to reconnect to the Sphero. If the Sphero has an active connection
// it will first close that connection and then establish a new connection.
// Returns true on Successful reconnection
func (a *SpheroAdaptor) Reconnect() bool {
	if a.Connected() == true {
		a.Disconnect()
	}
	return a.Connect()
}

// Disconnect terminates the connection to the Sphero. Returns true on successful disconnect.
func (a *SpheroAdaptor) Disconnect() bool {
	a.sp.Close()
	a.SetConnected(false)
	return true
}

// Finalize finalizes the SpheroAdaptor
func (a *SpheroAdaptor) Finalize() bool {
	return true
}
