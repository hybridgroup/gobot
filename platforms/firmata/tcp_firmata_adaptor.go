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
	conn, err := net.Dial("tcp", args[0].(string))
	if err != nil {
		// TODO: handle error
	}

	a := NewAdaptor(conn)
	a.SetName("TCPFirmata")

	return &TCPAdaptor{
		Adaptor: a,
	}
}
