package gobot

type Adaptor struct {
  Name string
  Robot Robot
  Connected bool 
  Port Port
  Params map[string]string
}

func (Adaptor) NewAdaptor(a Adaptor) Adaptor {
  return a
}

// Closes connection with device if connected
// @return [Boolean]
func (a *Adaptor) Finalize() bool{
  if a.IsConnected() {
    a.Disconnect()
  }
  return true
}

// Makes connected flag true
// @return [Boolean]
func (a *Adaptor) Connect() bool {
  a.Connected = true
  return true
}

// Makes connected flag false
// @return [Boolean]
func (a *Adaptor) Disconnect() bool {
  a.Connected = false
  return true
}

// Makes connected flag true
// @return [Boolean] true unless connected
func (a *Adaptor) Reconnect() bool {
  if !a.IsConnected(){
    return a.Connect()
  }
  return true
}

// @return [Boolean] connected flag status
func (a *Adaptor) IsConnected() bool {
  return a.Connected
}

/*
# Connects to configured port
# @return [TCPSocket] tcp socket of tcp port
# @return [String] port configured
def connect_to
  if port.is_tcp?
   connect_to_tcp
  else
   port.port
  end
end

# @return [TCPSocket] TCP socket connection
def connect_to_tcp
  @socket ||= TCPSocket.new(port.host, port.port)
end
      
# @return [UDPSocket] UDP socket connection
def connect_to_udp
  @udp_socket ||= UDPSocket.new
end

# Creates serial connection
# @param speed [int]
# @param data_bits [int]
# @param stop_bits [int]
# @param parity
# @return [SerialPort] new connection
def connect_to_serial(speed=57600, data_bits=8, stop_bits=1, parity=nil)
  require 'serialport'
  parity = ::SerialPort::NONE unless parity
  @sp = ::SerialPort.new(port.port, speed, data_bits, stop_bits, parity)
rescue LoadError
  Logger.error "Please 'gem install hybridgroup-serialport' for serial port support."
end
*/