//go:build !windows
// +build !windows

package firmata

import (
	"io"
	"net"

	"gobot.io/x/gobot/v2"
)

// TCPAdaptor represents a TCP based connection to a microcontroller running
// WiFiFirmata
type TCPAdaptor struct {
	*Adaptor
}

// NewTCPAdaptor opens and uses a TCP connection to a microcontroller running
// WiFiFirmata
func NewTCPAdaptor(args ...interface{}) *TCPAdaptor {
	address := args[0].(string) //nolint:forcetypeassert // ok here

	a := NewAdaptor(address)
	a.SetName(gobot.DefaultName("TCPFirmata"))
	a.PortOpener = connect

	return &TCPAdaptor{
		Adaptor: a,
	}
}

func connect(address string) (io.ReadWriteCloser, error) {
	return net.Dial("tcp", address)
}
