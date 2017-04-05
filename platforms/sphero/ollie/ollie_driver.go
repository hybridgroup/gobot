package ollie

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sync"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/sphero"
)

// Driver is the Gobot driver for the Sphero Ollie robot
type Driver struct {
	name              string
	connection        gobot.Connection
	seq               uint8
	mtx               sync.Mutex
	collisionResponse []uint8
	packetChannel     chan *packet
	gobot.Eventer
}

const (
	// bluetooth service IDs
	spheroBLEService    = "22bb746f2bb075542d6f726568705327"
	robotControlService = "22bb746f2ba075542d6f726568705327"

	// BLE characteristic IDs
	wakeCharacteristic     = "22bb746f2bbf75542d6f726568705327"
	txPowerCharacteristic  = "22bb746f2bb275542d6f726568705327"
	antiDosCharacteristic  = "22bb746f2bbd75542d6f726568705327"
	commandsCharacteristic = "22bb746f2ba175542d6f726568705327"
	responseCharacteristic = "22bb746f2ba675542d6f726568705327"

	// SensorData event
	SensorData = "sensordata"

	// Collision event
	Collision = "collision"

	// Error event
	Error = "error"

	// Packet header size
	PacketHeaderSize = 5

	// Response packet max size
	ResponsePacketMaxSize = 20

	// Collision Packet data size: The number of bytes following the DLEN field through the end of the packet
	CollisionDataSize = 17

	// Full size of the collision response
	CollisionResponseSize = PacketHeaderSize + CollisionDataSize
)

type packet struct {
	header   []uint8
	body     []uint8
	checksum uint8
}

// NewDriver creates a Driver for a Sphero Ollie
func NewDriver(a ble.BLEConnector) *Driver {
	n := &Driver{
		name:          gobot.DefaultName("Ollie"),
		connection:    a,
		Eventer:       gobot.NewEventer(),
		packetChannel: make(chan *packet, 1024),
	}

	n.AddEvent(Collision)

	return n
}

// Connection returns the connection to this Ollie
func (b *Driver) Connection() gobot.Connection { return b.connection }

// Name returns the name for the Driver
func (b *Driver) Name() string { return b.name }

// SetName sets the Name for the Driver
func (b *Driver) SetName(n string) { b.name = n }

// adaptor returns BLE adaptor
func (b *Driver) adaptor() ble.BLEConnector {
	return b.Connection().(ble.BLEConnector)
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

	go func() {
		for {
			b.adaptor().ReadCharacteristic(responseCharacteristic)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	b.ConfigureCollisionDetection(DefaultCollisionConfig())

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
	b.adaptor().Subscribe(responseCharacteristic, b.HandleResponses)

	return
}

// AntiDOSOff turns off Anti-DOS code so we can control Ollie
func (b *Driver) AntiDOSOff() (err error) {
	str := "011i3"
	buf := &bytes.Buffer{}
	buf.WriteString(str)

	err = b.adaptor().WriteCharacteristic(antiDosCharacteristic, buf.Bytes())
	if err != nil {
		fmt.Println("AntiDOSOff error:", err)
		return err
	}

	return
}

// Wake wakes Ollie up so we can play
func (b *Driver) Wake() (err error) {
	buf := []byte{0x01}

	err = b.adaptor().WriteCharacteristic(wakeCharacteristic, buf)
	if err != nil {
		fmt.Println("Wake error:", err)
		return err
	}

	return
}

// SetTXPower sets transmit level
func (b *Driver) SetTXPower(level int) (err error) {
	buf := []byte{byte(level)}

	err = b.adaptor().WriteCharacteristic(txPowerCharacteristic, buf)
	if err != nil {
		fmt.Println("SetTXLevel error:", err)
		return err
	}

	return
}

// HandleResponses handles responses returned from Ollie
func (b *Driver) HandleResponses(data []byte, e error) {
	//fmt.Println("response data:", data, e)

	b.handleCollisionDetected(data)
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

// ConfigureCollisionDetection configures the sensitivity of the detection.
func (b *Driver) ConfigureCollisionDetection(cc sphero.CollisionConfig) {
	b.packetChannel <- b.craftPacket([]uint8{cc.Method, cc.Xt, cc.Yt, cc.Xs, cc.Ys, cc.Dead}, 0x02, 0x12)
}

func (b *Driver) write(packet *packet) (err error) {
	buf := append(packet.header, packet.body...)
	buf = append(buf, packet.checksum)
	err = b.adaptor().WriteCharacteristic(commandsCharacteristic, buf)
	if err != nil {
		fmt.Println("send command error:", err)
		return err
	}

	b.mtx.Lock()
	defer b.mtx.Unlock()
	b.seq++
	return
}

func (b *Driver) craftPacket(body []uint8, did byte, cid byte) *packet {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	packet := new(packet)
	packet.body = body
	dlen := len(packet.body) + 1
	packet.header = []uint8{0xFF, 0xFF, did, cid, b.seq, uint8(dlen)}
	packet.checksum = b.calculateChecksum(packet)
	return packet
}

func (b *Driver) handleCollisionDetected(data []uint8) {

	if len(data) == ResponsePacketMaxSize {
		// Check if this is the header of collision response. (i.e. first part of data)
		// Collision response is 22 bytes long. (individual packet size is maxed at 20)
		switch data[1] {
		case 0xFE:
			if data[2] == 0x07 {
				// response code 7 is for a detected collision
				if len(b.collisionResponse) == 0 {
					b.collisionResponse = append(b.collisionResponse, data...)
				}
			}
		}
	} else if len(data) == CollisionResponseSize-ResponsePacketMaxSize {
		// if this is the remaining part of the collision response,
		// then make sure the header and first part of data is already received
		if len(b.collisionResponse) == ResponsePacketMaxSize {
			b.collisionResponse = append(b.collisionResponse, data...)
		}
	} else {
		return // not collision event
	}

	// check expected sizes
	if len(b.collisionResponse) != CollisionResponseSize || b.collisionResponse[4] != CollisionDataSize {
		return
	}

	// confirm checksum
	size := len(b.collisionResponse)
	chk := b.collisionResponse[size-1] // last byte is checksum
	if chk != calculateChecksum(b.collisionResponse[2:size-1]) {
		return
	}

	var collision sphero.CollisionPacket
	buffer := bytes.NewBuffer(b.collisionResponse[5:]) // skip header
	binary.Read(buffer, binary.BigEndian, &collision)
	b.collisionResponse = nil // clear the current response

	b.Publish(Collision, collision)
}

func (b *Driver) calculateChecksum(packet *packet) uint8 {
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
