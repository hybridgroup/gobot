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
	name               string
	connection         gobot.Connection
	seq                uint8
	mtx                sync.Mutex
	collisionResponse  []uint8
	packetChannel      chan *Packet
	asyncBuffer        []byte
	asyncMessage       []byte
	locatorCallback    func(p Point2D)
	powerstateCallback func(p PowerStatePacket)
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

// MotorModes is used to configure the motor
type MotorModes uint8

// MotorModes required for SetRawMotorValues command
const (
	Off MotorModes = iota
	Forward
	Reverse
	Brake
	Ignore
)

// Packet describes head, body and checksum for a data package to be sent to the sphero.
type Packet struct {
	Header   []uint8
	Body     []uint8
	Checksum uint8
}

// Point2D represents a koordinate in 2-Dimensional space
type Point2D struct {
	X int16
	Y int16
}

// NewDriver creates a Driver for a Sphero Ollie
func NewDriver(a ble.BLEConnector) *Driver {
	n := &Driver{
		name:          gobot.DefaultName("Ollie"),
		connection:    a,
		Eventer:       gobot.NewEventer(),
		packetChannel: make(chan *Packet, 1024),
	}

	n.AddEvent(Collision)

	return n
}

// PacketChannel returns the channel for packets to be sent to the sp
func (b *Driver) PacketChannel() chan *Packet { return b.packetChannel }

// Sequence returns the Sequence number of the current packet
func (b *Driver) Sequence() uint8 { return b.seq }

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

	//since packets can only be 20 bytes long, we have to puzzle them together
	newMessage := false

	//append message parts to existing
	if len(data) > 0 && data[0] != 0xFF {
		b.asyncBuffer = append(b.asyncBuffer, data...)
	}

	//clear message when new one begins (first byte is always 0xFF)
	if len(data) > 0 && data[0] == 0xFF {
		b.asyncMessage = b.asyncBuffer
		b.asyncBuffer = data
		newMessage = true
	}

	parts := b.asyncMessage
	//3 is the id of data streaming, located at index 2 byte
	if newMessage && len(parts) > 2 && parts[2] == 3 {
		b.handleDataStreaming(parts)
	}

	//index 1 is the type of the message, 0xFF being a direct response, 0xFE an asynchronous message
	if len(data) > 4 && data[1] == 0xFF && data[0] == 0xFF {
		//locator request
		if data[4] == 0x0B && len(data) == 16 {
			b.handleLocatorDetected(data)
		}

		if data[4] == 0x09 {
			b.handlePowerStateDetected(data)
		}
	}

	b.handleCollisionDetected(data)
}

// GetLocatorData calls the passed function with the data from the locator
func (b *Driver) GetLocatorData(f func(p Point2D)) {
	//CID 0x15 is the code for the locator request
	b.PacketChannel() <- b.craftPacket([]uint8{}, 0x02, 0x15)
	b.locatorCallback = f
}

// GetPowerState calls the passed function with the Power State information from the sphero
func (b *Driver) GetPowerState(f func(p PowerStatePacket)) {
	//CID 0x20 is the code for the power state
	b.PacketChannel() <- b.craftPacket([]uint8{}, 0x00, 0x20)
	b.powerstateCallback = f
}

func (b *Driver) handleDataStreaming(data []byte) {
	// ensure data is the right length:
	if len(data) != 88 {
		return
	}

	//data packet is the same as for the normal sphero, since the same communication api is used
	//only difference in communication is that the "newer" spheros use BLE for communinations
	var dataPacket DataStreamingPacket
	buffer := bytes.NewBuffer(data[5:]) // skip header
	binary.Read(buffer, binary.BigEndian, &dataPacket)

	b.Publish(SensorData, dataPacket)
}

// SetRGB sets the Ollie to the given r, g, and b values
func (b *Driver) SetRGB(r uint8, g uint8, bl uint8) {
	b.packetChannel <- b.craftPacket([]uint8{r, g, bl, 0x01}, 0x02, 0x20)
}

// Roll tells the Ollie to roll
func (b *Driver) Roll(speed uint8, heading uint16) {
	b.packetChannel <- b.craftPacket([]uint8{speed, uint8(heading >> 8), uint8(heading & 0xFF), 0x01}, 0x02, 0x30)
}

// Boost executes the boost macro from within the SSB which takes a
// 1 byte parameter which is either 01h to begin boosting or 00h to stop.
func (b *Driver) Boost(state bool) {
	s := uint8(0x01)
	if !state {
		s = 0x00
	}
	b.packetChannel <- b.craftPacket([]uint8{s}, 0x02, 0x31)
}

// SetStabilization enables or disables the built-in auto stabilizing features of the Ollie
func (b *Driver) SetStabilization(state bool) {
	s := uint8(0x01)
	if !state {
		s = 0x00
	}
	b.packetChannel <- b.craftPacket([]uint8{s}, 0x02, 0x02)
}

// SetRotationRate allows you to control the rotation rate that Sphero will use to meet new
// heading commands. A value of 255 jumps to the maximum (currently 400 degrees/sec).
// A value of zero doesn't make much sense so it's interpreted as 1, the minimum.
func (b *Driver) SetRotationRate(speed uint8) {
	b.packetChannel <- b.craftPacket([]uint8{speed}, 0x02, 0x03)
}

// SetRawMotorValues allows you to take over one or both of the motor output values,
// instead of having the stabilization system control them. Each motor (left and right)
// requires a mode and a power value from 0-255
func (b *Driver) SetRawMotorValues(lmode MotorModes, lpower uint8, rmode MotorModes, rpower uint8) {
	b.packetChannel <- b.craftPacket([]uint8{uint8(lmode), lpower, uint8(rmode), rpower}, 0x02, 0x33)
}

// SetBackLEDOutput allows you to control the brightness of the back(tail) LED.
func (b *Driver) SetBackLEDOutput(value uint8) {
	b.packetChannel <- b.craftPacket([]uint8{value}, 0x02, 0x21)
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

//SetDataStreamingConfig passes the config to the sphero to stream sensor data
func (b *Driver) SetDataStreamingConfig(d sphero.DataStreamingConfig) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, d)
	b.PacketChannel() <- b.craftPacket(buf.Bytes(), 0x02, 0x11)
}

func (b *Driver) write(packet *Packet) (err error) {
	buf := append(packet.Header, packet.Body...)
	buf = append(buf, packet.Checksum)
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

func (b *Driver) craftPacket(body []uint8, did byte, cid byte) *Packet {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	packet := new(Packet)
	packet.Body = body
	dlen := len(packet.Body) + 1
	packet.Header = []uint8{0xFF, 0xFF, did, cid, b.seq, uint8(dlen)}
	packet.Checksum = b.calculateChecksum(packet)
	return packet
}

func (b *Driver) handlePowerStateDetected(data []uint8) {

	var dataPacket PowerStatePacket
	buffer := bytes.NewBuffer(data[5:]) // skip header
	binary.Read(buffer, binary.BigEndian, &dataPacket)

	b.powerstateCallback(dataPacket)
}

func (b *Driver) handleLocatorDetected(data []uint8) {
	//read the unsigned raw values
	ux := binary.BigEndian.Uint16(data[5:7])
	uy := binary.BigEndian.Uint16(data[7:9])

	//convert to signed values
	var x, y int16

	if ux > 32255 {
		x = int16(ux - 65535)
	} else {
		x = int16(ux)
	}

	if uy > 32255 {
		y = int16(uy - 65535)
	} else {
		y = int16(uy)
	}

	//create point obj
	p := new(Point2D)
	p.X = x
	p.Y = y

	if b.locatorCallback != nil {
		b.locatorCallback(*p)
	}
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

func (b *Driver) calculateChecksum(packet *Packet) uint8 {
	buf := append(packet.Header, packet.Body...)
	return calculateChecksum(buf[2:])
}

func calculateChecksum(buf []byte) byte {
	var calculatedChecksum uint16
	for i := range buf {
		calculatedChecksum += uint16(buf[i])
	}
	return uint8(^(calculatedChecksum % 256))
}
