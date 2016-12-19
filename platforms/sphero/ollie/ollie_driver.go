package ollie

import (
	"bytes"
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
)

// Driver is the Gobot driver for the Sphero Ollie robot
type Driver struct {
	name          string
	connection    gobot.Connection
	seq           uint8
	packetChannel chan *packet
	gobot.Eventer
}

const (
	// SpheroBLEService is the primary service ID
	SpheroBLEService = "22bb746f2bb075542d6f726568705327"

	// RobotControlService is the service ID for the Sphero command API
	RobotControlService = "22bb746f2ba075542d6f726568705327"

	// WakeCharacteristic characteristic ID
	WakeCharacteristic = "22bb746f2bbf75542d6f726568705327"

	// TXPowerCharacteristic characteristic ID
	TXPowerCharacteristic = "22bb746f2bb275542d6f726568705327"

	// AntiDosCharacteristic characteristic ID
	AntiDosCharacteristic = "22bb746f2bbd75542d6f726568705327"

	// CommandsCharacteristic characteristic ID
	CommandsCharacteristic = "22bb746f2ba175542d6f726568705327"

	// ResponseCharacteristic characteristic ID
	ResponseCharacteristic = "22bb746f2ba675542d6f726568705327"

	// SensorData event
	SensorData = "sensordata"

	// Collision event
	Collision = "collision"

	// Error event
	Error = "error"
)

type packet struct {
	header   []uint8
	body     []uint8
	checksum uint8
}

// NewDriver creates a Driver for a Sphero Ollie
func NewDriver(a *ble.ClientAdaptor) *Driver {
	n := &Driver{
		name:          "Ollie",
		connection:    a,
		Eventer:       gobot.NewEventer(),
		packetChannel: make(chan *packet, 1024),
	}

	return n
}

// Connection returns the connection to this Ollie
func (b *Driver) Connection() gobot.Connection { return b.connection }

// Name returns the name for the Driver
func (b *Driver) Name() string { return b.name }

// SetName sets the Name for the Driver
func (b *Driver) SetName(n string) { b.name = n }

// adaptor returns BLE adaptor
func (b *Driver) adaptor() *ble.ClientAdaptor {
	return b.Connection().(*ble.ClientAdaptor)
}

// Start tells driver to get ready to do work
func (b *Driver) Start() (err error) {
	b.Init()

	// send commands
	go func() {
		for {
			packet := <-b.packetChannel
			err := b.write(packet)
			if err != nil {
				b.Publish(b.Event(Error), err)
			}
		}
	}()

	return
}

// Halt stops Ollie driver (void)
func (b *Driver) Halt() (err error) {
	b.Sleep()
	time.Sleep(750 * time.Microsecond)
	return
}

// Init is used to initialize the Ollie
func (b *Driver) Init() (err error) {
	b.AntiDOSOff()
	b.SetTXPower(7)
	b.Wake()

	// subscribe to Sphero response notifications
	b.adaptor().Subscribe(RobotControlService, ResponseCharacteristic, b.HandleResponses)

	return
}

// AntiDOSOff turns off Anti-DOS code so we can control Ollie
func (b *Driver) AntiDOSOff() (err error) {
	str := "011i3"
	buf := &bytes.Buffer{}
	buf.WriteString(str)

	err = b.adaptor().WriteCharacteristic(SpheroBLEService, AntiDosCharacteristic, buf.Bytes())
	if err != nil {
		fmt.Println("AntiDOSOff error:", err)
		return err
	}

	return
}

// Wake wakes Ollie up so we can play
func (b *Driver) Wake() (err error) {
	buf := []byte{0x01}

	err = b.adaptor().WriteCharacteristic(SpheroBLEService, WakeCharacteristic, buf)
	if err != nil {
		fmt.Println("Wake error:", err)
		return err
	}

	return
}

// SetTXPower sets transmit level
func (b *Driver) SetTXPower(level int) (err error) {
	buf := []byte{byte(level)}

	err = b.adaptor().WriteCharacteristic(SpheroBLEService, TXPowerCharacteristic, buf)
	if err != nil {
		fmt.Println("SetTXLevel error:", err)
		return err
	}

	return
}

// HandleResponses handles responses returned from Ollie
func (b *Driver) HandleResponses(data []byte, e error) {
	fmt.Println("response data:", data)

	return
}

// SetRGB sets the Ollie to the given r, g, and b values
func (b *Driver) SetRGB(r uint8, g uint8, bl uint8) {
	b.packetChannel <- b.craftPacket([]uint8{r, g, bl, 0x01}, 0x02, 0x20)
}

// Roll tells the Ollie to roll
func (b *Driver) Roll(speed uint8, heading uint16) {
	b.packetChannel <- b.craftPacket([]uint8{speed, uint8(heading >> 8), uint8(heading & 0xFF), 0x01}, 0x02, 0x30)
}

// Stop tells the Ollie to stop
func (b *Driver) Stop() {
	b.Roll(0, 0)
}

// Sleep says Go to sleep
func (b *Driver) Sleep() {
	b.packetChannel <- b.craftPacket([]uint8{0x00, 0x00, 0x00, 0x00, 0x00}, 0x00, 0x22)
}

// EnableStopOnDisconnect auto-sends a Stop command after losing the connection
func (b *Driver) EnableStopOnDisconnect() {
	b.packetChannel <- b.craftPacket([]uint8{0x00, 0x00, 0x00, 0x01}, 0x02, 0x37)
}

func (s *Driver) write(packet *packet) (err error) {
	buf := append(packet.header, packet.body...)
	buf = append(buf, packet.checksum)
	err = s.adaptor().WriteCharacteristic(RobotControlService, CommandsCharacteristic, buf)
	if err != nil {
		fmt.Println("send command error:", err)
		return err
	}

	s.seq++
	return
}

func (s *Driver) craftPacket(body []uint8, did byte, cid byte) *packet {
	packet := new(packet)
	packet.body = body
	dlen := len(packet.body) + 1
	packet.header = []uint8{0xFF, 0xFF, did, cid, s.seq, uint8(dlen)}
	packet.checksum = s.calculateChecksum(packet)
	return packet
}

func (s *Driver) calculateChecksum(packet *packet) uint8 {
	buf := append(packet.header, packet.body...)
	return calculateChecksum(buf[2:])
}

func calculateChecksum(buf []byte) byte {
	var calculatedChecksum uint16
	for i := range buf {
		calculatedChecksum += uint16(buf[i])
	}
	return uint8(^(calculatedChecksum % 256))
}
