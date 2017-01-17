package ble

// SerialPort is a implementation of serial over Bluetooth LE
// Inspired by https://github.com/monteslu/ble-serial by @monteslu
type SerialPort struct {
	address string
	rid     string
	tid     string
	client  *ClientAdaptor
}

// NewSerialPort returns a new serial over Bluetooth LE connection
func NewSerialPort(address string, rid string, tid string) *SerialPort {
	return &SerialPort{address: address, rid: rid, tid: tid}
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
	n = len(data)
	//log.Println("reading", p.rid, "data:", data)
	return
}

// Write writes to the BLE serial port connection
func (p *SerialPort) Write(b []byte) (n int, err error) {
	err = p.client.WriteCharacteristic(p.tid, b)
	n = len(b)
	//log.Println("writing", p.tid, "data:", b)
	return
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
