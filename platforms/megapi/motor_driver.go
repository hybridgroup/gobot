package megapi

import (
	"bytes"
	"encoding/binary"
	"sync"

	"gobot.io/x/gobot"
)

var _ gobot.Driver = (*MotorDriver)(nil)

// MotorDriver represents a motor
type MotorDriver struct {
	name     string
	megaPi   *Adaptor
	port     byte
	halted   bool
	syncRoot *sync.Mutex
}

// NewMotorDriver creates a new MotorDriver at the given port
func NewMotorDriver(megaPi *Adaptor, port byte) *MotorDriver {
	return &MotorDriver{
		name:     "MegaPiMotor",
		megaPi:   megaPi,
		port:     port,
		halted:   true,
		syncRoot: &sync.Mutex{},
	}
}

// Name returns the name of this motor
func (m *MotorDriver) Name() string {
	return m.name
}

// SetName sets the name of this motor
func (m *MotorDriver) SetName(n string) {
	m.name = n
}

// Start implements the Driver interface
func (m *MotorDriver) Start() error {
	m.syncRoot.Lock()
	defer m.syncRoot.Unlock()
	m.halted = false
	m.speedHelper(0)
	return nil
}

// Halt terminates the Driver interface
func (m *MotorDriver) Halt() error {
	m.syncRoot.Lock()
	defer m.syncRoot.Unlock()
	m.halted = true
	m.speedHelper(0)
	return nil
}

// Connection returns the Connection associated with the Driver
func (m *MotorDriver) Connection() gobot.Connection {
	return gobot.Connection(m.megaPi)
}

// Speed sets the motors speed to the specified value
func (m *MotorDriver) Speed(speed int16) error {
	m.syncRoot.Lock()
	defer m.syncRoot.Unlock()
	if m.halted {
		return nil
	}
	m.speedHelper(speed)
	return nil
}

// there is some sort of bug on the hardware such that you cannot
// send the exact same speed to 2 different motors consecutively
// hence we ensure we always alternate speeds
func (m *MotorDriver) speedHelper(speed int16) {
	m.sendSpeed(speed - 1)
	m.sendSpeed(speed)
}

// sendSpeed sets the motors speed to the specified value
func (m *MotorDriver) sendSpeed(speed int16) {
	bufOut := new(bytes.Buffer)

	// byte sequence: 0xff, 0x55, id, action, device, port
	bufOut.Write([]byte{0xff, 0x55, 0x6, 0x0, 0x2, 0xa, m.port})
	binary.Write(bufOut, binary.LittleEndian, speed)
	bufOut.Write([]byte{0xa})
	m.megaPi.writeBytesChannel <- bufOut.Bytes()
}
