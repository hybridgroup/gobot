package ble

import (
	"sync"

	"gobot.io/x/gobot/v2"
)

// SerialPortDriver is a implementation of serial over Bluetooth LE
// Inspired by https://github.com/monteslu/ble-serial by @monteslu
type SerialPortDriver struct {
	*Driver
	rid string
	tid string
	// buffer of responseData and mutex to protect it
	responseData  []byte
	responseMutex sync.Mutex
}

// NewSerialPortDriver returns a new serial over Bluetooth LE connection
func NewSerialPortDriver(a gobot.BLEConnector, rid string, tid string, opts ...OptionApplier) *SerialPortDriver {
	d := &SerialPortDriver{
		Driver: NewDriver(a, "BleSerial", nil, nil, opts...),
		rid:    rid,
		tid:    tid,
	}

	return d
}

// Open opens a connection to a BLE serial device
func (p *SerialPortDriver) Open() error {
	if err := p.Adaptor().Connect(); err != nil {
		return err
	}

	// subscribe to response notifications
	return p.Adaptor().Subscribe(p.rid, func(data []byte) {
		p.responseMutex.Lock()
		defer p.responseMutex.Unlock()
		p.responseData = append(p.responseData, data...)
	})
}

// Read reads bytes from BLE serial port connection
func (p *SerialPortDriver) Read(b []byte) (int, error) {
	p.responseMutex.Lock()
	defer p.responseMutex.Unlock()

	if len(p.responseData) == 0 {
		return 0, nil
	}

	n := len(b)
	if len(p.responseData) < n {
		n = len(p.responseData)
	}
	copy(b, p.responseData[:n])

	if len(p.responseData) > n {
		p.responseData = p.responseData[n:]
	} else {
		p.responseData = nil
	}

	return n, nil
}

// Write writes to the BLE serial port connection
func (p *SerialPortDriver) Write(b []byte) (int, error) {
	err := p.Adaptor().WriteCharacteristic(p.tid, b)
	n := len(b)
	return n, err
}

// Close closes the BLE serial port connection
func (p *SerialPortDriver) Close() error {
	return p.Adaptor().Disconnect()
}

// Address returns the BLE address
func (p *SerialPortDriver) Address() string {
	return p.Adaptor().Address()
}
