package firmata

import "net"

// TCPAdaptor represents a TCP based connection to a microcontroller running
// WiFiFirmata
type TCPAdaptor struct {
	*Adaptor
}

// NewTCPAdaptor opens and uses a TCP connection to a microcontroller running
// WiFiFirmata
func NewTCPAdaptor(args ...interface{}) *TCPAdaptor {
	address := args[0].(string)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		// TODO: handle error
	}

	a := NewAdaptor(conn, address)
	a.SetName("TCPFirmata")

	return &TCPAdaptor{
		Adaptor: a,
	}
}
