package ollie

import (
	"bytes"
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
)

var _ gobot.Driver = (*Driver)(nil)

type Driver struct {
	name          string
	connection    gobot.Connection
	seq           uint8
	packetChannel chan *packet
	gobot.Eventer
}

const (
	// service IDs
	SpheroBLEService    = "22bb746f2bb075542d6f726568705327"
	RobotControlService = "22bb746f2ba075542d6f726568705327"

	// characteristic IDs
	WakeCharacteristic    = "22bb746f2bbf75542d6f726568705327"
	TXPowerCharacteristic = "22bb746f2bb275542d6f726568705327"
	AntiDosCharacteristic = "22bb746f2bbd75542d6f726568705327"

	CommandsCharacteristic = "22bb746f2ba175542d6f726568705327"
	ResponseCharacteristic = "22bb746f2ba675542d6f726568705327"

	// gobot events
	SensorData = "sensordata"
	Collision  = "collision"
	Error      = "error"
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
func (b *Driver) Connection() gobot.Connection { return b.connection }

func (b *Driver) Name() string { return b.name }

func (b *Driver) SetName(n string) { b.name = n }

// adaptor returns BLE adaptor
func (b *Driver) adaptor() *ble.ClientAdaptor {
	return b.Connection().(*ble.ClientAdaptor)
}

// Start tells driver to get ready to do work
func (s *Driver) Start() (err error) {
	s.Init()

	// send commands
	go func() {
		for {
			packet := <-s.packetChannel
			err := s.write(packet)
			if err != nil {
				s.Publish(s.Event(Error), err)
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

func (b *Driver) Init() (err error) {
	b.AntiDOSOff()
	b.SetTXPower(7)
	b.Wake()

	// subscribe to Sphero response notifications
	b.adaptor().Subscribe(RobotControlService, ResponseCharacteristic, b.HandleResponses)

	return
}

// Turns off Anti-DOS code so we can control Ollie
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

// Wakes Ollie up so we can play
func (b *Driver) Wake() (err error) {
	buf := []byte{0x01}

	err = b.adaptor().WriteCharacteristic(SpheroBLEService, WakeCharacteristic, buf)
	if err != nil {
		fmt.Println("Wake error:", err)
		return err
	}

	return
}

// Sets transmit level
func (b *Driver) SetTXPower(level int) (err error) {
	buf := []byte{byte(level)}

	err = b.adaptor().WriteCharacteristic(SpheroBLEService, TXPowerCharacteristic, buf)
	if err != nil {
		fmt.Println("SetTXLevel error:", err)
		return err
	}

	return
}

// Handle responses returned from Ollie
func (b *Driver) HandleResponses(data []byte, e error) {
	fmt.Println("response data:", data)

	return
}

// SetRGB sets the Ollie to the given r, g, and b values
func (s *Driver) SetRGB(r uint8, g uint8, b uint8) {
	s.packetChannel <- s.craftPacket([]uint8{r, g, b, 0x01}, 0x02, 0x20)
}

// Tells the Ollie to roll
func (s *Driver) Roll(speed uint8, heading uint16) {
	s.packetChannel <- s.craftPacket([]uint8{speed, uint8(heading >> 8), uint8(heading & 0xFF), 0x01}, 0x02, 0x30)
}

// Tells the Ollie to stop
func (s *Driver) Stop() {
	s.Roll(0, 0)
}

// Go to sleep
func (s *Driver) Sleep() {
	s.packetChannel <- s.craftPacket([]uint8{0x00, 0x00, 0x00, 0x00, 0x00}, 0x00, 0x22)
}

func (s *Driver) EnableStopOnDisconnect() {
	s.packetChannel <- s.craftPacket([]uint8{0x00, 0x00, 0x00, 0x01}, 0x02, 0x37)
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
