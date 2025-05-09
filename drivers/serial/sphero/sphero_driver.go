package sphero

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/common/spherocommon"
	"gobot.io/x/gobot/v2/drivers/serial"
)

type spheroSerialAdaptor interface {
	gobot.Adaptor
	serial.SerialReader
	serial.SerialWriter

	IsConnected() bool
}

type packet struct {
	header   []uint8
	body     []uint8
	checksum uint8
}

// SpheroDriver Represents a Sphero 2.0
type SpheroDriver struct {
	*serial.Driver
	gobot.Eventer
	seq              uint8
	asyncResponse    [][]uint8
	syncResponse     [][]uint8
	packetChannel    chan *packet
	responseChannel  chan []uint8
	originalColor    []uint8 // Only used for calibration.
	shutdownWaitTime time.Duration
}

// NewSpheroDriver returns a new SpheroDriver given a Sphero Adaptor.
//
// Adds the following API Commands:
//
//	"ConfigureLocator" - See SpheroDriver.ConfigureLocator
//	"Roll" - See SpheroDriver.Roll
//	"Stop" - See SpheroDriver.Stop
//	"GetRGB" - See SpheroDriver.GetRGB
//	"ReadLocator" - See SpheroDriver.ReadLocator
//	"SetBackLED" - See SpheroDriver.SetBackLED
//	"SetHeading" - See SpheroDriver.SetHeading
//	"SetStabilization" - See SpheroDriver.SetStabilization
//	"SetDataStreaming" - See SpheroDriver.SetDataStreaming
//	"SetRotationRate" - See SpheroDriver.SetRotationRate
func NewSpheroDriver(a spheroSerialAdaptor, opts ...serial.OptionApplier) *SpheroDriver {
	d := &SpheroDriver{
		Eventer:          gobot.NewEventer(),
		packetChannel:    make(chan *packet, 1024),
		responseChannel:  make(chan []uint8, 1024),
		shutdownWaitTime: 1 * time.Second,
	}
	d.Driver = serial.NewDriver(a, "Sphero", d.initialize, d.shutdown, opts...)

	d.AddEvent(spherocommon.ErrorEvent)
	d.AddEvent(spherocommon.CollisionEvent)
	d.AddEvent(spherocommon.SensorDataEvent)

	//nolint:forcetypeassert // ok here
	d.AddCommand("SetRGB", func(params map[string]interface{}) interface{} {
		r := uint8(params["r"].(float64))
		g := uint8(params["g"].(float64))
		b := uint8(params["b"].(float64))
		d.SetRGB(r, g, b)
		return nil
	})

	//nolint:forcetypeassert // ok here
	d.AddCommand("Roll", func(params map[string]interface{}) interface{} {
		speed := uint8(params["speed"].(float64))
		heading := uint16(params["heading"].(float64))
		d.Roll(speed, heading)
		return nil
	})

	d.AddCommand("Stop", func(_ map[string]interface{}) interface{} {
		d.Stop()
		return nil
	})

	d.AddCommand("GetRGB", func(_ map[string]interface{}) interface{} {
		return d.GetRGB()
	})

	d.AddCommand("ReadLocator", func(_ map[string]interface{}) interface{} {
		return d.ReadLocator()
	})

	//nolint:forcetypeassert // ok here
	d.AddCommand("SetBackLED", func(params map[string]interface{}) interface{} {
		level := uint8(params["level"].(float64))
		d.SetBackLED(level)
		return nil
	})
	//nolint:forcetypeassert // ok here
	d.AddCommand("SetRotationRate", func(params map[string]interface{}) interface{} {
		level := uint8(params["level"].(float64))
		d.SetRotationRate(level)
		return nil
	})
	//nolint:forcetypeassert // ok here
	d.AddCommand("SetHeading", func(params map[string]interface{}) interface{} {
		heading := uint16(params["heading"].(float64))
		d.SetHeading(heading)
		return nil
	})
	//nolint:forcetypeassert // ok here
	d.AddCommand("SetStabilization", func(params map[string]interface{}) interface{} {
		on := params["enable"].(bool)
		d.SetStabilization(on)
		return nil
	})
	//nolint:forcetypeassert // ok here
	d.AddCommand("SetDataStreaming", func(params map[string]interface{}) interface{} {
		N := uint16(params["N"].(float64))
		M := uint16(params["M"].(float64))
		Mask := uint32(params["Mask"].(float64))
		Pcnt := uint8(params["Pcnt"].(float64))
		Mask2 := uint32(params["Mask2"].(float64))

		d.SetDataStreaming(spherocommon.DataStreamingConfig{N: N, M: M, Mask2: Mask2, Pcnt: Pcnt, Mask: Mask})
		return nil
	})
	//nolint:forcetypeassert // ok here
	d.AddCommand("ConfigureLocator", func(params map[string]interface{}) interface{} {
		Flags := uint8(params["Flags"].(float64))
		X := int16(params["X"].(float64))
		Y := int16(params["Y"].(float64))
		YawTare := int16(params["YawTare"].(float64))

		d.ConfigureLocator(spherocommon.LocatorConfig{Flags: Flags, X: X, Y: Y, YawTare: YawTare})
		return nil
	})

	return d
}

// SetRGB sets the Sphero to the given r, g, and b values
func (d *SpheroDriver) SetRGB(r uint8, g uint8, b uint8) {
	d.sendCraftPacket([]uint8{r, g, b, 0x01}, 0x20)
}

// GetRGB returns the current r, g, b value of the Sphero
func (d *SpheroDriver) GetRGB() []uint8 {
	buf := d.getSyncResponse(d.craftPacket([]uint8{}, 0x22))
	if len(buf) == 9 {
		return []uint8{buf[5], buf[6], buf[7]}
	}
	return []uint8{}
}

// ReadLocator reads Sphero's current position (X,Y), component velocities and SOG (speed over ground).
func (d *SpheroDriver) ReadLocator() []int16 {
	buf := d.getSyncResponse(d.craftPacket([]uint8{}, 0x15))
	if len(buf) == 16 {
		vals := make([]int16, 5)
		_ = binary.Read(bytes.NewReader(buf[5:15]), binary.BigEndian, &vals)
		return vals
	}
	return []int16{}
}

// SetBackLED sets the Sphero Back LED to the specified brightness
func (d *SpheroDriver) SetBackLED(level uint8) {
	d.sendCraftPacket([]uint8{level}, 0x21)
}

// SetRotationRate sets the Sphero rotation rate
// A value of 255 jumps to the maximum (currently 400 degrees/sec).
func (d *SpheroDriver) SetRotationRate(level uint8) {
	d.sendCraftPacket([]uint8{level}, 0x03)
}

// SetHeading sets the heading of the Sphero
func (d *SpheroDriver) SetHeading(heading uint16) {
	//nolint:gosec // TODO: fix later
	d.sendCraftPacket([]uint8{uint8(heading >> 8), uint8(heading & 0xFF)}, 0x01)
}

// SetStabilization enables or disables the built-in auto stabilizing features of the Sphero
func (d *SpheroDriver) SetStabilization(on bool) {
	b := uint8(0x01)
	if !on {
		b = 0x00
	}
	d.sendCraftPacket([]uint8{b}, 0x02)
}

// Roll sends a roll command to the Sphero gives a speed and heading
func (d *SpheroDriver) Roll(speed uint8, heading uint16) {
	//nolint:gosec // TODO: fix later
	d.sendCraftPacket([]uint8{speed, uint8(heading >> 8), uint8(heading & 0xFF), 0x01}, 0x30)
}

// ConfigureLocator configures and enables the Locator
func (d *SpheroDriver) ConfigureLocator(lc spherocommon.LocatorConfig) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, lc); err != nil {
		panic(err)
	}

	d.sendCraftPacket(buf.Bytes(), 0x13)
}

// SetDataStreaming enables sensor data streaming
func (d *SpheroDriver) SetDataStreaming(dsc spherocommon.DataStreamingConfig) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, dsc); err != nil {
		panic(err)
	}

	d.sendCraftPacket(buf.Bytes(), 0x11)
}

// Stop sets the Sphero to a roll speed of 0
func (d *SpheroDriver) Stop() {
	d.Roll(0, 0)
}

// ConfigureCollisionDetection configures the sensitivity of the detection.
func (d *SpheroDriver) ConfigureCollisionDetection(cc spherocommon.CollisionConfig) {
	d.sendCraftPacket([]uint8{cc.Method, cc.Xt, cc.Yt, cc.Xs, cc.Ys, cc.Dead}, 0x12)
}

// SetCalibration sets up Sphero for manual heading calibration.
// It does this by turning on the tail light (so you can tell where it's
// facing) and disabling stabilization (so you can adjust the heading).
//
// When done, call FinishCalibration to set the new heading, and re-enable
// stabilization.
func (d *SpheroDriver) StartCalibration() {
	d.originalColor = d.GetRGB()
	d.SetRGB(0, 0, 0)
	d.SetBackLED(127)
	d.SetStabilization(false)
}

// FinishCalibration ends Sphero's calibration mode, by setting
// the new heading as current, and re-enabling normal defaults. This is a NOP
// in case StartCalibration was not called.
func (d *SpheroDriver) FinishCalibration() {
	if d.originalColor == nil {
		// Piggybacking on the original color being set to know if we are
		// calibrating or not.
		return
	}

	d.SetHeading(0)
	d.SetRGB(d.originalColor[0], d.originalColor[1], d.originalColor[2])
	d.SetBackLED(0)
	d.SetStabilization(true)
	d.originalColor = nil
}

// initialize starts the SpheroDriver and enables Collision Detection.
// Returns true on successful start.
//
// Emits the Events:
//
//	Collision  spherocommon.CollisionPacket - On Collision Detected
//	SensorData spherocommon.DataStreamingPacket - On Data Streaming event
//	Error      error- On error while processing asynchronous response
//
// TODO: stop the go routines gracefully on shutdown()
func (d *SpheroDriver) initialize() error {
	go func() {
		for {
			packet := <-d.packetChannel
			err := d.write(packet)
			if err != nil {
				d.Publish(spherocommon.ErrorEvent, err)
			}
		}
	}()

	go func() {
		for {
			response := <-d.responseChannel
			d.syncResponse = append(d.syncResponse, response)
		}
	}()

	go func() {
		for {
			header := d.readHeader()
			if len(header) > 0 {
				body := d.readBody(header[4])
				data := append(header, body...)
				checksum := data[len(data)-1]
				if checksum != spherocommon.CalculateChecksum(data[2:len(data)-1]) {
					continue
				}
				switch header[1] {
				case 0xFE:
					d.asyncResponse = append(d.asyncResponse, data)
				case 0xFF:
					d.responseChannel <- data
				}
			}
		}
	}()

	go func() {
		for {
			var evt []uint8
			for len(d.asyncResponse) != 0 {
				evt, d.asyncResponse = d.asyncResponse[len(d.asyncResponse)-1], d.asyncResponse[:len(d.asyncResponse)-1]
				if evt[2] == 0x07 {
					d.handleCollisionDetected(evt)
				} else if evt[2] == 0x03 {
					d.handleDataStreaming(evt)
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	d.ConfigureCollisionDetection(spheroDefaultCollisionConfig())
	d.enableStopOnDisconnect()

	return nil
}

// shutdown halts the SpheroDriver and sends a SpheroDriver.Stop command to the Sphero.
func (d *SpheroDriver) shutdown() error {
	if d.adaptor().IsConnected() {
		gobot.Every(10*time.Millisecond, func() {
			d.Stop()
		})
		time.Sleep(d.shutdownWaitTime)
	}
	return nil
}

func (d *SpheroDriver) enableStopOnDisconnect() {
	d.sendCraftPacket([]uint8{0x00, 0x00, 0x00, 0x01}, 0x37)
}

func (d *SpheroDriver) handleCollisionDetected(data []uint8) {
	// ensure data is the right length:
	if len(data) != 22 || data[4] != 17 {
		return
	}
	var collision spherocommon.CollisionPacket
	buffer := bytes.NewBuffer(data[5:]) // skip header
	if err := binary.Read(buffer, binary.BigEndian, &collision); err != nil {
		panic(err)
	}
	d.Publish(spherocommon.CollisionEvent, collision)
}

func (d *SpheroDriver) handleDataStreaming(data []uint8) {
	// ensure data is the right length:
	if len(data) != 90 {
		return
	}
	var dataPacket spherocommon.DataStreamingPacket
	buffer := bytes.NewBuffer(data[5:]) // skip header
	if err := binary.Read(buffer, binary.BigEndian, &dataPacket); err != nil {
		panic(err)
	}
	d.Publish(spherocommon.SensorDataEvent, dataPacket)
}

func (d *SpheroDriver) getSyncResponse(packet *packet) []byte {
	d.packetChannel <- packet
	for i := 0; i < 500; i++ {
		for key := range d.syncResponse {
			if d.syncResponse[key][3] == packet.header[4] && len(d.syncResponse[key]) > 6 {
				var response []byte
				response, d.syncResponse = d.syncResponse[len(d.syncResponse)-1], d.syncResponse[:len(d.syncResponse)-1]
				return response
			}
		}
		time.Sleep(100 * time.Microsecond)
	}

	return []byte{}
}

func (d *SpheroDriver) sendCraftPacket(body []uint8, cid byte) {
	d.packetChannel <- d.craftPacket(body, cid)
}

func (d *SpheroDriver) craftPacket(body []uint8, cid byte) *packet {
	dlen := len(body) + 1
	did := uint8(0x02)
	hdr := []uint8{0xFF, 0xFF, did, cid, d.seq, uint8(dlen)} //nolint:gosec // TODO: fix later
	buf := append(hdr, body...)

	packet := &packet{
		body:     body,
		header:   hdr,
		checksum: spherocommon.CalculateChecksum(buf[2:]),
	}

	return packet
}

func (d *SpheroDriver) write(packet *packet) error {
	d.Mutex().Lock()
	defer d.Mutex().Unlock()

	buf := append(packet.header, packet.body...)
	buf = append(buf, packet.checksum)
	length, err := d.adaptor().SerialWrite(buf)
	if err != nil {
		return err
	}

	if length != len(buf) {
		return errors.New("Not enough bytes written")
	}
	d.seq++
	return nil
}

func (d *SpheroDriver) readHeader() []uint8 {
	return d.readNextChunk(5)
}

func (d *SpheroDriver) readBody(length uint8) []uint8 {
	return d.readNextChunk(int(length))
}

func (d *SpheroDriver) readNextChunk(length int) []uint8 {
	read := make([]uint8, length)
	bytesRead := 0

	for bytesRead < length {
		time.Sleep(1 * time.Millisecond)
		n, err := d.adaptor().SerialRead(read[bytesRead:])
		if err != nil {
			return nil
		}
		bytesRead += n
	}
	return read
}

func (d *SpheroDriver) adaptor() spheroSerialAdaptor {
	if a, ok := d.Connection().(spheroSerialAdaptor); ok {
		return a
	}

	log.Printf("%s has no Sphero serial connector\n", d.Name())
	return nil
}

// spheroDefaultCollisionConfig returns a CollisionConfig with sensible collision defaults
func spheroDefaultCollisionConfig() spherocommon.CollisionConfig {
	return spherocommon.CollisionConfig{
		Method: 0x01,
		Xt:     0x80,
		Yt:     0x80,
		Xs:     0x80,
		Ys:     0x80,
		Dead:   0x60,
	}
}

// spheroDefaultLocatorConfig returns a LocatorConfig with defaults
func spheroDefaultLocatorConfig() spherocommon.LocatorConfig {
	return spherocommon.LocatorConfig{
		Flags:   0x01,
		X:       0x00,
		Y:       0x00,
		YawTare: 0x00,
	}
}
