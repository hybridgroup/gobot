package sphero

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble"
	"gobot.io/x/gobot/v2/drivers/common/spherocommon"
)

const (
	// spheroBLEService    = "22bb746f2bb075542d6f726568705327"
	// robotControlService = "22bb746f2ba075542d6f726568705327"

	r2WakeChara     = "22bb746f2bbf75542d6f726568705327"
	r2TxPowerChara  = "22bb746f2bb275542d6f726568705327"
	r2AntiDosChara  = "22bb746f2bbd75542d6f726568705327"
	r2CommandsChara = "00010002574f4f2053706865726f2121"
	r2ResponseChara = r2CommandsChara

	// start of packet
	sop = 0x8D

	// end of packet
	eop = 0xD8

	// flags
	isResponse                = 0b1
	requestsResponse          = 0b10
	requestsOnlyErrorResponse = 0b100
	isActivity                = 0b1000
	hasTargetId               = 0b10000
	hasSourceId               = 0b100000
	unused                    = 0b1000000
	extendedFlags             = 0b10000000

	// packet header size
	r2PacketHeaderSize = 5

	// Response packet max size
	r2ResponsePacketMaxSize = 20

	// Collision packet data size: The number of bytes following the DLEN field through the end of the packet
	r2CollisionDataSize = 17

	// Full size of the collision response
	r2CollisionResponseSize = packetHeaderSize + collisionDataSize
)

// R2D2Driver is the Gobot driver for the Sphero R2D2 robot
type R2D2Driver struct {
	*ble.Driver
	gobot.Eventer
	defaultCollisionConfig spherocommon.CollisionConfig
	seq                    uint8
	collisionResponse      []uint8
	packetChannel          chan *packet
	asyncBuffer            []byte
	asyncMessage           []byte
	locatorCallback        func(p Point2D)
	powerstateCallback     func(p spherocommon.PowerStatePacket)
}

// NewR2D2Driver creates a driver for a Sphero R2D2
func NewR2D2Driver(a gobot.BLEConnector, opts ...ble.OptionApplier) *R2D2Driver {
	return newR2D2BaseDriver(a, "R2D2", r2d2DefaultCollisionConfig(), opts...)
}

func newR2D2BaseDriver(
	a gobot.BLEConnector, name string,
	dcc spherocommon.CollisionConfig, opts ...ble.OptionApplier,
) *R2D2Driver {
	d := &R2D2Driver{
		defaultCollisionConfig: dcc,
		Eventer:                gobot.NewEventer(),
		packetChannel:          make(chan *packet, 1024),
	}
	d.Driver = ble.NewDriver(a, name, d.initialize, d.shutdown, opts...)

	d.AddEvent(spherocommon.ErrorEvent)
	d.AddEvent(spherocommon.CollisionEvent)

	return d
}

// SetTXPower sets transmit level
func (d *R2D2Driver) SetTXPower(level int) error {
	buf := []byte{byte(level)}

	if err := d.Adaptor().WriteCharacteristic(r2TxPowerChara, buf); err != nil {
		return err
	}

	return nil
}

// Wake wakes R2D2 up so we can play
func (d *R2D2Driver) Wake() error {
	buf := []byte{0x01}

	if err := d.Adaptor().WriteCharacteristic(r2WakeChara, buf); err != nil {
		return err
	}

	return nil
}

// ConfigureCollisionDetection configures the sensitivity of the detection.
func (d *R2D2Driver) ConfigureCollisionDetection(cc spherocommon.CollisionConfig) {
	d.sendCraftPacket([]uint8{cc.Method, cc.Xt, cc.Yt, cc.Xs, cc.Ys, cc.Dead}, 0x02, 0x12)
}

// GetLocatorData calls the passed function with the data from the locator
func (d *R2D2Driver) GetLocatorData(f func(p Point2D)) {
	// CID 0x15 is the code for the locator request
	d.sendCraftPacket([]uint8{}, 0x02, 0x15)
	d.locatorCallback = f
}

// GetPowerState calls the passed function with the Power State information from the sphero
func (d *R2D2Driver) GetPowerState(f func(p spherocommon.PowerStatePacket)) {
	// CID 0x20 is the code for the power state
	d.sendCraftPacket([]uint8{}, 0x00, 0x20)
	d.powerstateCallback = f
}

// SetRGB sets the R2D2 to the given r, g, and b values
func (d *R2D2Driver) SetRGB(r uint8, g uint8, b uint8) {
	d.sendCraftPacket([]uint8{r, g, b, 0x01}, 0x02, 0x20)
}

// Roll tells the R2D2 to roll
func (d *R2D2Driver) Roll(speed uint8, heading uint16) {
	//nolint:gosec // TODO: fix later
	d.sendCraftPacket([]uint8{speed, uint8(heading >> 8), uint8(heading & 0xFF), 0x01}, 0x02, 0x30)
}

// Boost executes the boost macro from within the SSB which takes a 1 byte parameter which is
// either 01h to begin boosting or 00h to stop.
func (d *R2D2Driver) Boost(state bool) {
	s := uint8(0x01)
	if !state {
		s = 0x00
	}
	d.sendCraftPacket([]uint8{s}, 0x02, 0x31)
}

// SetStabilization enables or disables the built-in auto stabilizing features of the R2D2
func (d *R2D2Driver) SetStabilization(state bool) {
	s := uint8(0x01)
	if !state {
		s = 0x00
	}
	d.sendCraftPacket([]uint8{s}, 0x02, 0x02)
}

// SetRotationRate allows you to control the rotation rate that Sphero will use to meet new heading commands. A value
// of 255 jumps to the maximum (currently 400 degrees/sec). A value of zero doesn't make much sense so it's interpreted
// as 1, the minimum.
func (d *R2D2Driver) SetRotationRate(speed uint8) {
	d.sendCraftPacket([]uint8{speed}, 0x02, 0x03)
}

// SetRawMotorValues allows you to take over one or both of the motor output values, instead of having the stabilization
// system control them. Each motor (left and right) requires a mode and a power value from 0-255.
func (d *R2D2Driver) SetRawMotorValues(lmode MotorModes, lpower uint8, rmode MotorModes, rpower uint8) {
	d.sendCraftPacket([]uint8{uint8(lmode), lpower, uint8(rmode), rpower}, 0x02, 0x33)
}

// SetBackLEDBrightness allows you to control the brightness of the back(tail) LED.
func (d *R2D2Driver) SetBackLEDBrightness(value uint8) {
	d.sendCraftPacket([]uint8{value}, 0x02, 0x21)
}

// Stop tells the R2D2 to stop
func (d *R2D2Driver) Stop() {
	d.Roll(0, 0)
}

// Sleep says Go to sleep
func (d *R2D2Driver) Sleep() {
	d.sendCraftPacket([]uint8{0x00, 0x00, 0x00, 0x00, 0x00}, 0x00, 0x22)
}

// SetDataStreamingConfig passes the config to the sphero to stream sensor data
func (d *R2D2Driver) SetDataStreamingConfig(dsc spherocommon.DataStreamingConfig) error {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, dsc); err != nil {
		return err
	}
	d.sendCraftPacket(buf.Bytes(), 0x02, 0x11)
	return nil
}

// initialize tells driver to get ready to do work
func (d *R2D2Driver) initialize() error {
	if err := d.antiDOSOff(); err != nil {
		return err
	}
	if err := d.SetTXPower(7); err != nil {
		return err
	}
	if err := d.Wake(); err != nil {
		return err
	}

	// subscribe to Sphero response notifications
	if err := d.Adaptor().Subscribe(r2ResponseChara, d.handleResponses); err != nil {
		return err
	}

	go func() {
		for {
			packet := <-d.packetChannel
			err := d.writeCommand(packet)
			if err != nil {
				d.Publish(d.Event(spherocommon.ErrorEvent), err)
			}
		}
	}()

	d.ConfigureCollisionDetection(d.defaultCollisionConfig)
	d.enableStopOnDisconnect()

	return nil
}

// antiDOSOff turns off Anti-DOS code so we can control R2D2
func (d *R2D2Driver) antiDOSOff() error {
	str := "011i3"
	buf := &bytes.Buffer{}
	buf.WriteString(str)

	if err := d.Adaptor().WriteCharacteristic(r2AntiDosChara, buf.Bytes()); err != nil {
		return err
	}

	return nil
}

func (d *R2D2Driver) writeCommand(packet *packet) error {
	d.Mutex().Lock()
	defer d.Mutex().Unlock()

	buf := append([]uint8{sop}, packet.header...)
	buf = append(buf, packet.body...)
	buf = append(buf, packet.checksum, eop)
	if err := d.Adaptor().WriteCharacteristic(r2CommandsChara, buf); err != nil {
		fmt.Println("async send command error:", err)
		return err
	}

	d.seq++
	return nil
}

// enableStopOnDisconnect auto-sends a Stop command after losing the connection
func (d *R2D2Driver) enableStopOnDisconnect() {
	d.sendCraftPacket([]uint8{0x00, 0x00, 0x00, 0x01}, 0x02, 0x37)
}

// shutdown stops R2D2 driver (void)
func (d *R2D2Driver) shutdown() error {
	d.Sleep()
	time.Sleep(750 * time.Microsecond)
	return nil
}

// handleResponses handles responses returned from R2D2
func (d *R2D2Driver) handleResponses(data []byte) {
	// since packets can only be 20 bytes long, we have to puzzle them together
	newMessage := false

	// append message parts to existing
	if len(data) > 0 && data[0] != 0xFF {
		d.asyncBuffer = append(d.asyncBuffer, data...)
	}

	// clear message when new one begins (first byte is always 0xFF)
	if len(data) > 0 && data[0] == 0xFF {
		d.asyncMessage = d.asyncBuffer
		d.asyncBuffer = data
		newMessage = true
	}

	parts := d.asyncMessage
	// 3 is the id of data streaming, located at index 2 byte
	if newMessage && len(parts) > 2 && parts[2] == 3 {
		d.handleDataStreaming(parts)
	}

	// index 1 is the type of the message, 0xFF being a direct response, 0xFE an asynchronous message
	if len(data) > 4 && data[1] == 0xFF && data[0] == 0xFF {
		// locator request
		if data[4] == 0x0B && len(data) == 16 {
			d.handleLocatorDetected(data)
		}

		if data[4] == 0x09 {
			d.handlePowerStateDetected(data)
		}
	}

	d.handleCollisionDetected(data)
}

func (d *R2D2Driver) handleDataStreaming(data []byte) {
	// ensure data is the right length:
	if len(data) != 88 {
		return
	}

	// data packet is the same as for the normal sphero, since the same communication api is used
	// only difference in communication is that the "newer" spheros use BLE for communications
	var dataPacket spherocommon.DataStreamingPacket
	buffer := bytes.NewBuffer(data[5:]) // skip header
	if err := binary.Read(buffer, binary.BigEndian, &dataPacket); err != nil {
		panic(err)
	}

	d.Publish(spherocommon.SensorDataEvent, dataPacket)
}

func (d *R2D2Driver) handleLocatorDetected(data []uint8) {
	if d.locatorCallback == nil {
		return
	}

	// read the unsigned raw values
	ux := binary.BigEndian.Uint16(data[5:7])
	uy := binary.BigEndian.Uint16(data[7:9])

	// convert to signed values
	var x, y int16

	if ux > 32255 {
		x = int16(ux - 65535) //nolint:gosec // ok here
	} else {
		x = int16(ux)
	}

	if uy > 32255 {
		y = int16(uy - 65535) //nolint:gosec // ok here
	} else {
		y = int16(uy)
	}

	d.locatorCallback(Point2D{X: x, Y: y})
}

func (d *R2D2Driver) handlePowerStateDetected(data []uint8) {
	var dataPacket spherocommon.PowerStatePacket
	buffer := bytes.NewBuffer(data[5:]) // skip header
	if err := binary.Read(buffer, binary.BigEndian, &dataPacket); err != nil {
		panic(err)
	}

	d.powerstateCallback(dataPacket)
}

func (d *R2D2Driver) handleCollisionDetected(data []uint8) {
	switch len(data) {
	case responsePacketMaxSize:
		// Check if this is the header of collision response. (i.e. first part of data)
		// Collision response is 22 bytes long. (individual packet size is maxed at 20)
		if data[1] == 0xFE && data[2] == 0x07 && len(d.collisionResponse) == 0 {
			// response code 7 is for a detected collision
			d.collisionResponse = append(d.collisionResponse, data...)
		}
	case collisionResponseSize - responsePacketMaxSize:
		// if this is the remaining part of the collision response,
		// then make sure the header and first part of data is already received
		if len(d.collisionResponse) == responsePacketMaxSize {
			d.collisionResponse = append(d.collisionResponse, data...)
		}
	default:
		return // not collision event
	}

	// check expected sizes
	if len(d.collisionResponse) != collisionResponseSize || d.collisionResponse[4] != collisionDataSize {
		return
	}

	// confirm checksum
	size := len(d.collisionResponse)
	chk := d.collisionResponse[size-1] // last byte is checksum
	if chk != spherocommon.CalculateChecksum(d.collisionResponse[2:size-1]) {
		return
	}

	var collision spherocommon.CollisionPacket
	buffer := bytes.NewBuffer(d.collisionResponse[5:]) // skip header
	if err := binary.Read(buffer, binary.BigEndian, &collision); err != nil {
		panic(err)
	}
	d.collisionResponse = nil // clear the current response

	d.Publish(spherocommon.CollisionEvent, collision)
}

func (d *R2D2Driver) sendCraftPacket(body []uint8, did byte, cid byte) {
	d.packetChannel <- d.craftPacket(body, nil, did, cid)
}

func (d *R2D2Driver) craftPacket(body []uint8, tid *byte, did byte, cid byte) *packet {
	//TODO handle flags, tid, sid, err better
	flags := byte(requestsResponse | isActivity)
	var sid byte
	if tid != nil {
		flags |= hasSourceId | hasTargetId
		sid = 0x1
	}
	hdr := []uint8{flags, *tid, sid, did, cid, d.seq} //nolint:gosec // TODO: fix later
	buf := append(hdr, body...)

	packet := &packet{
		body:     body,
		header:   hdr,
		checksum: spherocommon.CalculateChecksum(buf[2:]),
	}

	return packet
}

// r2d2DefaultCollisionConfig returns a CollisionConfig with sensible collision defaults
func r2d2DefaultCollisionConfig() spherocommon.CollisionConfig {
	return spherocommon.CollisionConfig{
		Method: 0x01,
		Xt:     0x20,
		Yt:     0x20,
		Xs:     0x20,
		Ys:     0x20,
		Dead:   0x60,
	}
}
