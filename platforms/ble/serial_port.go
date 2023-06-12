package ble

import "sync"

// SerialPort is a implementation of serial over Bluetooth LE
// Inspired by https://github.com/monteslu/ble-serial by @monteslu
type SerialPort struct {
	address string
	rid     string
	tid     string
	client  *ClientAdaptor

	// buffer of responseData and mutex to protect it
	responseData  []byte
	responseMutex sync.Mutex
}

// NewSerialPort returns a new serial over Bluetooth LE connection
func NewSerialPort(address string, rid string, tid string) *SerialPort {
	return &SerialPort{address: address, rid: rid, tid: tid}
}

// Open opens a connection to a BLE serial device
func (p *SerialPort) Open() error {
	p.client = NewClientAdaptor(p.address)

	if err := p.client.Connect(); err != nil {
		return err
	}

	// subscribe to response notifications
	return p.client.Subscribe(p.rid, func(data []byte, e error) {
		p.responseMutex.Lock()
		p.responseData = append(p.responseData, data...)
		p.responseMutex.Unlock()
	})
}

// Read reads bytes from BLE serial port connection
func (p *SerialPort) Read(b []byte) (n int, err error) {
	if len(p.responseData) == 0 {
		return
	}

	p.responseMutex.Lock()
	n = len(b)
	if len(p.responseData) < n {
		n = len(p.responseData)
	}
	copy(b, p.responseData[:n])

	if len(p.responseData) > n {
		p.responseData = p.responseData[n:]
	} else {
		p.responseData = nil
	}
	p.responseMutex.Unlock()

	return
}

// Write writes to the BLE serial port connection
func (p *SerialPort) Write(b []byte) (n int, err error) {
	err = p.client.WriteCharacteristic(p.tid, b)
	n = len(b)
	return
}

// Close closes the BLE serial port connection
func (p *SerialPort) Close() error {
	return p.client.Disconnect()
}

// Address returns the BLE address
func (p *SerialPort) Address() string {
	return p.address
}
