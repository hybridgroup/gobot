package ble

// SerialPort is a implementation of serial over Bluetooth LE
type SerialPort struct {
	address string
	rid     string
	tid     string
	client  *ClientAdaptor
}

// NewSerialPort returns a new serial over Bluetooth LE connection
func NewSerialPort(address string, rid string, wid string) *SerialPort {
	return &SerialPort{address: address, rid: rid, tid: wid}
}

// Open opens a connection to a BLE serial device
func (p *SerialPort) Open() (err error) {
	p.client = NewClientAdaptor(p.address)
	err = p.client.Connect()
	return
}

// Read reads bytes from BLE serial port connection
func (p *SerialPort) Read(b []byte) (n int, err error) {
	data, err := p.client.ReadCharacteristic(p.rid)
	if err != nil {
		return
	}
	copy(data, b)
	n = len(b)
	return
}

// Write writes to the BLE serial port connection
func (p *SerialPort) Write(b []byte) (n int, err error) {
	err = p.client.WriteCharacteristic(p.tid, b)
	return 0, err
}

// Close closes the BLE serial port connection
func (p *SerialPort) Close() (err error) {
	p.client.Disconnect()
	return
}

// Address returns the BLE address
func (p *SerialPort) Address() string {
	return p.address
}
