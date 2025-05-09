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

const (
	// spheroBLEService    = "22bb746f2bb075542d6f726568705327"
	// robotControlService = "22bb746f2ba075542d6f726568705327"

	wakeChara     = "22bb746f2bbf75542d6f726568705327"
	txPowerChara  = "22bb746f2bb275542d6f726568705327"
	antiDosChara  = "22bb746f2bbd75542d6f726568705327"
	commandsChara = "22bb746f2ba175542d6f726568705327"
	responseChara = "22bb746f2ba675542d6f726568705327"

	// packet header size
	packetHeaderSize = 5

	// Response packet max size
	responsePacketMaxSize = 20

	// Collision packet data size: The number of bytes following the DLEN field through the end of the packet
	collisionDataSize = 17

	// Full size of the collision response
	collisionResponseSize = packetHeaderSize + collisionDataSize
)

// packet describes head, body and checksum for a data package to be sent
type packet struct {
	header   []uint8
	body     []uint8
	checksum uint8
}

// Point2D represents a coordinate in 2-Dimensional space, exposed because used in a callback
type Point2D struct {
	X int16
	Y int16
}

// OllieDriver is the Gobot driver for the Sphero Ollie robot
type OllieDriver struct {
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

// NewOllieDriver creates a driver for a Sphero Ollie
func NewOllieDriver(a gobot.BLEConnector, opts ...ble.OptionApplier) *OllieDriver {
	return newOllieBaseDriver(a, "Ollie", ollieDefaultCollisionConfig(), opts...)
}

func newOllieBaseDriver(
	a gobot.BLEConnector, name string,
	dcc spherocommon.CollisionConfig, opts ...ble.OptionApplier,
) *OllieDriver {
	d := &OllieDriver{
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
func (d *OllieDriver) SetTXPower(level int) error {
	buf := []byte{byte(level)}

	if err := d.Adaptor().WriteCharacteristic(txPowerChara, buf); err != nil {
		return err
	}

	return nil
}

// Wake wakes Ollie up so we can play
func (d *OllieDriver) Wake() error {
	buf := []byte{0x01}

	if err := d.Adaptor().WriteCharacteristic(wakeChara, buf); err != nil {
		return err
	}

	return nil
}

// ConfigureCollisionDetection configures the sensitivity of the detection.
func (d *OllieDriver) ConfigureCollisionDetection(cc spherocommon.CollisionConfig) {
	d.sendCraftPacket([]uint8{cc.Method, cc.Xt, cc.Yt, cc.Xs, cc.Ys, cc.Dead}, 0x02, 0x12)
}

// GetLocatorData calls the passed function with the data from the locator
func (d *OllieDriver) GetLocatorData(f func(p Point2D)) {
	// CID 0x15 is the code for the locator request
	d.sendCraftPacket([]uint8{}, 0x02, 0x15)
	d.locatorCallback = f
}

// GetPowerState calls the passed function with the Power State information from the sphero
func (d *OllieDriver) GetPowerState(f func(p spherocommon.PowerStatePacket)) {
	// CID 0x20 is the code for the power state
	d.sendCraftPacket([]uint8{}, 0x00, 0x20)
	d.powerstateCallback = f
}

// SetRGB sets the Ollie to the given r, g, and b values
func (d *OllieDriver) SetRGB(r uint8, g uint8, b uint8) {
	d.sendCraftPacket([]uint8{r, g, b, 0x01}, 0x02, 0x20)
}

// Roll tells the Ollie to roll
func (d *OllieDriver) Roll(speed uint8, heading uint16) {
	//nolint:gosec // TODO: fix later
	d.sendCraftPacket([]uint8{speed, uint8(heading >> 8), uint8(heading & 0xFF), 0x01}, 0x02, 0x30)
}

// Boost executes the boost macro from within the SSB which takes a 1 byte parameter which is
// either 01h to begin boosting or 00h to stop.
func (d *OllieDriver) Boost(state bool) {
	s := uint8(0x01)
	if !state {
		s = 0x00
	}
	d.sendCraftPacket([]uint8{s}, 0x02, 0x31)
}

// SetStabilization enables or disables the built-in auto stabilizing features of the Ollie
func (d *OllieDriver) SetStabilization(state bool) {
	s := uint8(0x01)
	if !state {
		s = 0x00
	}
	d.sendCraftPacket([]uint8{s}, 0x02, 0x02)
}

// SetRotationRate allows you to control the rotation rate that Sphero will use to meet new heading commands. A value
// of 255 jumps to the maximum (currently 400 degrees/sec). A value of zero doesn't make much sense so it's interpreted
// as 1, the minimum.
func (d *OllieDriver) SetRotationRate(speed uint8) {
	d.sendCraftPacket([]uint8{speed}, 0x02, 0x03)
}

// SetRawMotorValues allows you to take over one or both of the motor output values, instead of having the stabilization
// system control them. Each motor (left and right) requires a mode and a power value from 0-255.
func (d *OllieDriver) SetRawMotorValues(lmode MotorModes, lpower uint8, rmode MotorModes, rpower uint8) {
	d.sendCraftPacket([]uint8{uint8(lmode), lpower, uint8(rmode), rpower}, 0x02, 0x33)
}

// SetBackLEDBrightness allows you to control the brightness of the back(tail) LED.
func (d *OllieDriver) SetBackLEDBrightness(value uint8) {
	d.sendCraftPacket([]uint8{value}, 0x02, 0x21)
}

// Stop tells the Ollie to stop
func (d *OllieDriver) Stop() {
	d.Roll(0, 0)
}

// Sleep says Go to sleep
func (d *OllieDriver) Sleep() {
	d.sendCraftPacket([]uint8{0x00, 0x00, 0x00, 0x00, 0x00}, 0x00, 0x22)
}

// SetDataStreamingConfig passes the config to the sphero to stream sensor data
func (d *OllieDriver) SetDataStreamingConfig(dsc spherocommon.DataStreamingConfig) error {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, dsc); err != nil {
		return err
	}
	d.sendCraftPacket(buf.Bytes(), 0x02, 0x11)
	return nil
}

// initialize tells driver to get ready to do work
func (d *OllieDriver) initialize() error {
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
	if err := d.Adaptor().Subscribe(responseChara, d.handleResponses); err != nil {
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

// antiDOSOff turns off Anti-DOS code so we can control Ollie
func (d *OllieDriver) antiDOSOff() error {
	str := "011i3"
	buf := &bytes.Buffer{}
	buf.WriteString(str)

	if err := d.Adaptor().WriteCharacteristic(antiDosChara, buf.Bytes()); err != nil {
		return err
	}

	return nil
}

func (d *OllieDriver) writeCommand(packet *packet) error {
	d.Mutex().Lock()
	defer d.Mutex().Unlock()

	buf := append(packet.header, packet.body...)
	buf = append(buf, packet.checksum)
	if err := d.Adaptor().WriteCharacteristic(commandsChara, buf); err != nil {
		fmt.Println("async send command error:", err)
		return err
	}

	d.seq++
	return nil
}

// enableStopOnDisconnect auto-sends a Stop command after losing the connection
func (d *OllieDriver) enableStopOnDisconnect() {
	d.sendCraftPacket([]uint8{0x00, 0x00, 0x00, 0x01}, 0x02, 0x37)
}

// shutdown stops Ollie driver (void)
func (d *OllieDriver) shutdown() error {
	d.Sleep()
	time.Sleep(750 * time.Microsecond)
	return nil
}

// handleResponses handles responses returned from Ollie
func (d *OllieDriver) handleResponses(data []byte) {
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

func (d *OllieDriver) handleDataStreaming(data []byte) {
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

func (d *OllieDriver) handleLocatorDetected(data []uint8) {
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

func (d *OllieDriver) handlePowerStateDetected(data []uint8) {
	var dataPacket spherocommon.PowerStatePacket
	buffer := bytes.NewBuffer(data[5:]) // skip header
	if err := binary.Read(buffer, binary.BigEndian, &dataPacket); err != nil {
		panic(err)
	}

	d.powerstateCallback(dataPacket)
}

func (d *OllieDriver) handleCollisionDetected(data []uint8) {
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

func (d *OllieDriver) sendCraftPacket(body []uint8, did byte, cid byte) {
	d.packetChannel <- d.craftPacket(body, did, cid)
}

func (d *OllieDriver) craftPacket(body []uint8, did byte, cid byte) *packet {
	dlen := len(body) + 1
	hdr := []uint8{0xFF, 0xFF, did, cid, d.seq, uint8(dlen)} //nolint:gosec // TODO: fix later
	buf := append(hdr, body...)

	packet := &packet{
		body:     body,
		header:   hdr,
		checksum: spherocommon.CalculateChecksum(buf[2:]),
	}

	return packet
}

// ollieDefaultCollisionConfig returns a CollisionConfig with sensible collision defaults
func ollieDefaultCollisionConfig() spherocommon.CollisionConfig {
	return spherocommon.CollisionConfig{
		Method: 0x01,
		Xt:     0x20,
		Yt:     0x20,
		Xs:     0x20,
		Ys:     0x20,
		Dead:   0x60,
	}
}
