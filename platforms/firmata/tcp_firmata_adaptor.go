package firmata

import (
	"io"
	"net"

	"gobot.io/x/gobot"
)

// TCPAdaptor represents a TCP based connection to a microcontroller running
// WiFiFirmata
type TCPAdaptor struct {
	*Adaptor
}

func connect(address string) (io.ReadWriteCloser, error) {
	return net.Dial("tcp", address)
}

// NewTCPAdaptor opens and uses a TCP connection to a microcontroller running
// WiFiFirmata
func NewTCPAdaptor(args ...interface{}) *TCPAdaptor {
	address := args[0].(string)

	a := NewAdaptor(address)
	a.SetName(gobot.DefaultName("TCPFirmata"))
	a.PortOpener = func(port string) (io.ReadWriteCloser, error) {
		return connect(port)
	}

	return &TCPAdaptor{
		Adaptor: a,
	}
}
