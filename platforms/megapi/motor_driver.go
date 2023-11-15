package megapi

import (
	"bytes"
	"encoding/binary"
	"sync"

	"gobot.io/x/gobot/v2"
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
func (d *MotorDriver) Name() string {
	return d.name
}

// SetName sets the name of this motor
func (d *MotorDriver) SetName(n string) {
	d.name = n
}

// Start implements the Driver interface
func (d *MotorDriver) Start() error {
	d.syncRoot.Lock()
	defer d.syncRoot.Unlock()
	d.halted = false
	return d.speedHelper(0)
}

// Halt terminates the Driver interface
func (d *MotorDriver) Halt() error {
	d.syncRoot.Lock()
	defer d.syncRoot.Unlock()
	d.halted = true
	return d.speedHelper(0)
}

// Connection returns the Connection associated with the Driver
func (d *MotorDriver) Connection() gobot.Connection {
	return gobot.Connection(d.megaPi)
}

// Speed sets the motors speed to the specified value
func (d *MotorDriver) Speed(speed int16) error {
	d.syncRoot.Lock()
	defer d.syncRoot.Unlock()
	if d.halted {
		return nil
	}
	return d.speedHelper(speed)
}

// there is some sort of bug on the hardware such that you cannot
// send the exact same speed to 2 different motors consecutively
// hence we ensure we always alternate speeds
func (d *MotorDriver) speedHelper(speed int16) error {
	if err := d.sendSpeed(speed - 1); err != nil {
		return err
	}
	return d.sendSpeed(speed)
}

// sendSpeed sets the motors speed to the specified value
func (d *MotorDriver) sendSpeed(speed int16) error {
	bufOut := new(bytes.Buffer)

	// byte sequence: 0xff, 0x55, id, action, device, port
	bufOut.Write([]byte{0xff, 0x55, 0x6, 0x0, 0x2, 0xa, d.port})
	if err := binary.Write(bufOut, binary.LittleEndian, speed); err != nil {
		return err
	}
	bufOut.Write([]byte{0xa})
	d.megaPi.writeBytesChannel <- bufOut.Bytes()

	return nil
}
