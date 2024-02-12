package megapi

import (
	"bytes"
	"encoding/binary"
	"log"
	"sync"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/serial"
)

var _ gobot.Driver = (*MotorDriver)(nil)

type megapiMotorSerialAdaptor interface {
	gobot.Adaptor
	serial.SerialWriter
}

// MotorDriver represents a motor
type MotorDriver struct {
	*serial.Driver
	port              byte
	halted            bool
	writeBytesChannel chan []byte
	finalizeChannel   chan struct{}
	syncRoot          *sync.Mutex
}

// NewMotorDriver creates a new MotorDriver at the given port
func NewMotorDriver(a megapiMotorSerialAdaptor, port byte, opts ...serial.OptionApplier) *MotorDriver {
	d := &MotorDriver{
		port:              port,
		halted:            true,
		syncRoot:          &sync.Mutex{},
		writeBytesChannel: make(chan []byte),
		finalizeChannel:   make(chan struct{}),
	}
	d.Driver = serial.NewDriver(a, "MegaPiMotor", d.initialize, d.shutdown, opts...)

	return d
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

// initialize implements the Driver interface
func (d *MotorDriver) initialize() error {
	d.syncRoot.Lock()
	defer d.syncRoot.Unlock()

	// sleeping is required to give the board a chance to reset after connection is done
	time.Sleep(2 * time.Second)

	// kick off thread to send bytes to the board
	go func() {
		for {
			select {
			case bytes := <-d.writeBytesChannel:
				if _, err := d.adaptor().SerialWrite(bytes); err != nil {
					panic(err)
				}
				time.Sleep(10 * time.Millisecond)
			case <-d.finalizeChannel:
				d.finalizeChannel <- struct{}{}
				return
			default:
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	d.halted = false
	return d.speedHelper(0)
}

// Halt terminates the Driver interface
func (d *MotorDriver) shutdown() error {
	d.syncRoot.Lock()
	defer d.syncRoot.Unlock()

	d.finalizeChannel <- struct{}{}
	<-d.finalizeChannel

	d.halted = true
	return d.speedHelper(0)
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
	d.writeBytesChannel <- bufOut.Bytes()

	return nil
}

func (d *MotorDriver) adaptor() megapiMotorSerialAdaptor {
	if a, ok := d.Connection().(megapiMotorSerialAdaptor); ok {
		return a
	}

	log.Printf("%s has no MegaPi serial connector\n", d.Name())
	return nil
}
