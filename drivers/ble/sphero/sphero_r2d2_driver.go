package sphero

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"sync"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble"
	"gobot.io/x/gobot/v2/drivers/common/spherocommon"
)

type Playback uint8

type LegAction uint8

const (
	r2AntiDosChara  = "00020005574f4f2053706865726f2121" // handshake 1st transmission
	r2CommandsChara = "00010002574f4f2053706865726f2121" // send uuid
	r2ResponseChara = r2CommandsChara

	// command safe interval
	commandInterval = time.Duration(12) * time.Millisecond

	// start of packet
	sop = 0x8D
	// end of packet
	eop        = 0xD8
	escapeHex  = 0xAB
	escapeMask = 0x88

	// flags
	isResponse                = 0b1
	requestsResponse          = 0b10
	requestsOnlyErrorResponse = 0b100
	isActivity                = 0b1000
	hasTargetId               = 0b10000
	hasSourceId               = 0b100000
	unused                    = 0b1000000
	extendedFlags             = 0b10000000

	Immediate         Playback = 0x0
	IfNotPlaying      Playback = 0x1
	AfterCurrentSound Playback = 0x2

	Stop      LegAction = 0
	ThreeLegs LegAction = 1
	TwoLegs   LegAction = 2
	Waddle    LegAction = 3
)

// R2D2Driver is the Gobot driver for the Sphero R2D2 robot
type R2D2Driver struct {
	*ble.Driver
	gobot.Eventer

	defaultCollisionConfig spherocommon.CollisionConfig
	seq                    uint8
	seqMutex               sync.Mutex
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
		seqMutex:               sync.Mutex{},
	}
	d.Driver = ble.NewDriver(a, name, d.initialize, d.shutdown, opts...)

	d.AddEvent(spherocommon.ErrorEvent)
	d.AddEvent(spherocommon.CollisionEvent)

	return d
}

// Wake wakes R2D2 up so we can play
func (d *R2D2Driver) Wake() error {
	// Power.wake did: 19 cid: 13
	d.sendCraftPacket([]uint8{}, 0x13, 0x0D)

	return nil
}

func (d *R2D2Driver) ResetLocatorData() {
	// did: 24, cid: 19
	d.sendCraftPacket([]uint8{}, 0x18, 0x13)
}

// ConfigureCollisionDetection configures the sensitivity of the detection.
func (d *R2D2Driver) ConfigureCollisionDetection(cc spherocommon.CollisionConfig) {
	// did: 24, cid: 17
	d.sendCraftPacket([]uint8{cc.Method, cc.Xt, cc.Yt, cc.Xs, cc.Ys, cc.Dead}, 0x18, 0x11)
}

// GetLocatorData calls the passed function with the data from the locator
func (d *R2D2Driver) GetLocatorData(f func(p Point2D)) {
	// CID 0x15 is the code for the locator request
	d.sendCraftPacket([]uint8{}, 0x02, 0x15)
	d.locatorCallback = f
}

// GetPowerState calls the passed function with the Power State information from the sphero
func (d *R2D2Driver) GetPowerState(f func(p spherocommon.PowerStatePacket)) {
	// did: 19, cid: 4
	d.sendCraftPacket([]uint8{}, 0x13, 0x04)
	//	CHARGED = 0
	//	CHARGING = 1
	//	NOT_CHARGING = 2
	//	OK = 3
	//	LOW = 4
	//	CRITICAL = 5
	//	UNKNOWN = 255
	d.powerstateCallback = f
}

// SetRGB sets the R2D2 to the given r, g, and b values
func (d *R2D2Driver) SetRGB(r uint8, g uint8, b uint8) {
	// did: 26, cid: 14
	d.sendCraftPacket([]uint8{0, 119, r, g, b, r, g, b}, 0x1A, 0x0E)
}

// Roll tells the R2D2 to roll
func (d *R2D2Driver) Roll(speed uint8, heading uint16) {
	//nolint:gosec // TODO: fix later
	// did: 22, cid: 7
	// last data packet is DriveFlags, may not be supported on the R2
	//	FORWARD = 0x0  # 0b0
	//	BACKWARD = 0x1  # 0b1
	//	TURBO = 0x2  # 0b10
	//	FAST_TURN = 0x4  # 0b100
	//	LEFT_DIRECTION = 0x8  # 0b1000
	//	RIGHT_DIRECTION = 0x10  # 0b10000
	//	ENABLE_DRIFT = 0x20  # 0b100000
	d.sendCraftPacket([]uint8{speed, uint8(heading >> 8), uint8(heading & 0xFF), 0x00}, 0x16, 0x07)
}

// SetStabilization enables or disables the built-in auto stabilizing features of the R2D2
func (d *R2D2Driver) SetStabilization(state bool) {
	s := uint8(0x01)
	if !state {
		s = 0x00
	}
	// did: 22, cid: 12
	//	NO_CONTROL_SYSTEM = 0
	//	FULL_CONTROL_SYSTEM = 1
	//	PITCH_CONTROL_SYSTEM = 2
	//	ROLL_CONTROL_SYSTEM = 3
	//	YAW_CONTROL_SYSTEM = 4
	//	SPEED_AND_YAW_CONTROL_SYSTEM = 5
	d.sendCraftPacket([]uint8{s}, 0x16, 0x0C)
}

// SetRawMotorValues allows you to take over one or both of the motor output values, instead of having the stabilization
// system control them. Each motor (left and right) requires a mode and a power value from 0-255.
// MotorModes Brake and Ignore are not supported on the R2.
func (d *R2D2Driver) SetRawMotorValues(lmode MotorModes, lpower uint8, rmode MotorModes, rpower uint8) {
	// did: 22, cid: 1
	d.sendCraftPacket([]uint8{uint8(lmode), lpower, uint8(rmode), rpower}, 0x16, 0x01)
}

// SetBackRGB sets the back R2D2 dome LED to the given r, g, and b values
func (d *R2D2Driver) SetBackRGB(r uint8, g uint8, b uint8) {
	// did: 26, cid: 14
	//	def set_leds(self, mapping: Dict[IntEnum, int]):
	//	mask = 0
	//	led_values = []
	//	for e in self.__toy.LEDs:
	//		if e in mapping:
	//			mask |= 1 << e
	//			led_values.append(mapping[e])
	//	self.__toy.set_all_leds_with_16_bit_mask(mask, led_values)
	// IO._encode(toy, 14, proc, [*to_bytes(mask, 2), *values])
	d.sendCraftPacket([]uint8{0, 119, r, g, b, r, g, b}, 0x1A, 0x0E)
}

func (d *R2D2Driver) PerformLegAction(action LegAction) {
	// did: 23, cid: 13
	d.sendCraftPacket([]uint8{uint8(action)}, 0x17, 0x0D)
}

// SetDomePosition pos can be -160 to 180
func (d *R2D2Driver) SetDomePosition(pos float32) {
	// did: 23, cid: 15
	d.sendCraftPacket(spherocommon.Float32ToBytes(pos), 0x17, 0x0F)
}

// PlaySound where playback is PlaybackImmediate, PlaybackIfNotPlaying or PlaybackAfterCurrentSound
func (d *R2D2Driver) PlaySound(sound Audio, playback Playback) {
	// did: 26, cid: 7
	d.sendCraftPacket(append(spherocommon.Uint16ToBytes(uint16(sound)), uint8(playback)), 0x1A, 0x07)
}

func (d *R2D2Driver) PlayAnimation(anima Animation) {
	// did: 23, cid: 5
	d.sendCraftPacket(spherocommon.Uint16ToBytes(uint16(anima)), 0x17, 0x05)
}

// Stop tells the R2D2 to stop
func (d *R2D2Driver) Stop() {
	d.Roll(0, 0)
}

// Sleep says Go to sleep
func (d *R2D2Driver) Sleep() {
	// did: 19, cid: 1
	d.sendCraftPacket([]uint8{}, 0x13, 0x1)
}

// SetDataStreamingConfig passes the config to the sphero to stream sensor data
// TODO get this working for R2
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
	// enableStopOnDisconnect is a temporary option, not supported by the R2
	// d.enableStopOnDisconnect()

	return nil
}

// antiDOSOff turns off Anti-DOS code so we can control R2D2
func (d *R2D2Driver) antiDOSOff() error {
	str := "usetheforce...band"
	buf := &bytes.Buffer{}
	buf.WriteString(str)

	if err := d.Adaptor().WriteCharacteristic(r2AntiDosChara, buf.Bytes()); err != nil {
		return err
	}

	return nil
}

func (d *R2D2Driver) writeCommand(packet *packet) error {
	log.Printf("request %X % X % X %X %X\n", sop, packet.header, packet.body, packet.checksum, eop)

	buf := append([]uint8{sop}, escapeBytes(packet.header)...)
	buf = append(buf, escapeBytes(packet.body)...)
	buf = append(buf, escapeByte(packet.checksum)...)
	buf = append(buf, eop)

	if err := d.Adaptor().WriteCharacteristic(r2CommandsChara, buf); err != nil {
		fmt.Println("async send command error:", err)
		return err
	}

	// avoid ddos the r2
	time.Sleep(commandInterval)

	return nil
}

// shutdown stops R2D2 driver (void)
func (d *R2D2Driver) shutdown() error {
	d.Sleep()
	time.Sleep(750 * time.Microsecond)
	return nil
}

// handleResponses handles responses returned from R2D2
func (d *R2D2Driver) handleResponses(data []byte) {
	// log.Printf("handleResponse of %v bytes: % X", len(data), data)

	// v2 packets can be arbitrary length, we have to puzzle them together
	// they also are sent in 1 byte chunks and can be out of order

	// ignore end of packet and wait till we see another packet
	if len(data) > 0 && data[0] == eop {
		return
	}

	// append message parts to existing
	if len(data) > 0 && data[0] != sop {
		d.asyncBuffer = append(d.asyncBuffer, data...)
		return
	}

	// clear message when new one begins (first byte is always 0x8D)
	if len(data) > 0 {
		// append end of packet to complete message
		d.asyncMessage = append(d.asyncBuffer, eop)
		d.asyncBuffer = data
	}

	log.Printf("processing asyncMessage: % X", d.asyncMessage)

	// TODO get sensor data from a packet starting with 0x8d 0x0 0x18 0x2 0xff
}

func (d *R2D2Driver) sendCraftPacket(body []uint8, did byte, cid byte) {
	d.packetChannel <- d.craftPacket(body, did, cid)
}

func (d *R2D2Driver) craftPacket(body []uint8, did byte, cid byte) *packet {
	// packet protocol v2
	// [SOP, FLAGS, TID (optional), SID (optional), DID, CID, SEQ, ERR (at response), DATA..., CHK, EOP]
	// FLAGS to DATA... included in CHK
	// FLAGS to CHK are encoded
	d.seqMutex.Lock()
	defer d.seqMutex.Unlock()

	// TODO handle flags, tid, sid err better
	flags := byte(requestsResponse | isActivity)
	// var sid byte
	// if tid != nil {
	// 	flags |= hasSourceId | hasTargetId
	// 	sid = 0x1
	// }
	hdr := []uint8{flags, did, cid, d.seq}
	buf := append(hdr, body...)

	packet := &packet{
		body:     body,
		header:   hdr,
		checksum: spherocommon.CalculateChecksum(buf),
	}

	d.seq++

	return packet
}

// r2d2DefaultCollisionConfig returns a CollisionConfig with sensible collision defaults
func r2d2DefaultCollisionConfig() spherocommon.CollisionConfig {
	return spherocommon.CollisionConfig{
		Method: 0x01,
		Xt:     0x5A,
		Yt:     0x82,
		Xs:     0x5A,
		Ys:     0x82,
		Dead:   0x01,
	}
}

func escapeBytes(data []uint8) []uint8 {
	var escaped []uint8
	for _, b := range data {
		escaped = append(escaped, escapeByte(b)...)
	}
	return escaped
}

func escapeByte(b uint8) []uint8 {
	var escaped []uint8
	if b == sop || b == eop || b == escapeHex {
		return append(escaped, escapeHex, b^escapeMask)
	}

	return append(escaped, b)
}
