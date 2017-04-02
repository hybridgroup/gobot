package sphero

import (
	"bytes"
	"encoding/binary"
	"errors"
	"sync"
	"time"

	"gobot.io/x/gobot"
)

const (
	// Error event when error encountered
	Error = "error"

	// SensorData event when sensor data is received
	SensorData = "sensordata"

	// Collision event when collision is detected
	Collision = "collision"
)

type packet struct {
	header   []uint8
	body     []uint8
	checksum uint8
}

// SpheroDriver Represents a Sphero 2.0
type SpheroDriver struct {
	name            string
	connection      gobot.Connection
	mtx             sync.Mutex
	seq             uint8
	asyncResponse   [][]uint8
	syncResponse    [][]uint8
	packetChannel   chan *packet
	responseChannel chan []uint8
	gobot.Eventer
	gobot.Commander
}

// NewSpheroDriver returns a new SpheroDriver given a Sphero Adaptor.
//
// Adds the following API Commands:
// 	"ConfigureLocator" - See SpheroDriver.ConfigureLocator
// 	"Roll" - See SpheroDriver.Roll
// 	"Stop" - See SpheroDriver.Stop
// 	"GetRGB" - See SpheroDriver.GetRGB
//	"ReadLocator" - See SpheroDriver.ReadLocator
// 	"SetBackLED" - See SpheroDriver.SetBackLED
// 	"SetHeading" - See SpheroDriver.SetHeading
// 	"SetStabilization" - See SpheroDriver.SetStabilization
//  "SetDataStreaming" - See SpheroDriver.SetDataStreaming
//  "SetRotationRate" - See SpheroDriver.SetRotationRate
func NewSpheroDriver(a *Adaptor) *SpheroDriver {
	s := &SpheroDriver{
		name:            gobot.DefaultName("Sphero"),
		connection:      a,
		Eventer:         gobot.NewEventer(),
		Commander:       gobot.NewCommander(),
		packetChannel:   make(chan *packet, 1024),
		responseChannel: make(chan []uint8, 1024),
	}

	s.AddEvent(Error)
	s.AddEvent(Collision)
	s.AddEvent(SensorData)

	s.AddCommand("SetRGB", func(params map[string]interface{}) interface{} {
		r := uint8(params["r"].(float64))
		g := uint8(params["g"].(float64))
		b := uint8(params["b"].(float64))
		s.SetRGB(r, g, b)
		return nil
	})

	s.AddCommand("Roll", func(params map[string]interface{}) interface{} {
		speed := uint8(params["speed"].(float64))
		heading := uint16(params["heading"].(float64))
		s.Roll(speed, heading)
		return nil
	})

	s.AddCommand("Stop", func(params map[string]interface{}) interface{} {
		s.Stop()
		return nil
	})

	s.AddCommand("GetRGB", func(params map[string]interface{}) interface{} {
		return s.GetRGB()
	})

	s.AddCommand("ReadLocator", func(params map[string]interface{}) interface{} {
		return s.ReadLocator()
	})

	s.AddCommand("SetBackLED", func(params map[string]interface{}) interface{} {
		level := uint8(params["level"].(float64))
		s.SetBackLED(level)
		return nil
	})

	s.AddCommand("SetRotationRate", func(params map[string]interface{}) interface{} {
		level := uint8(params["level"].(float64))
		s.SetRotationRate(level)
		return nil
	})

	s.AddCommand("SetHeading", func(params map[string]interface{}) interface{} {
		heading := uint16(params["heading"].(float64))
		s.SetHeading(heading)
		return nil
	})

	s.AddCommand("SetStabilization", func(params map[string]interface{}) interface{} {
		on := params["enable"].(bool)
		s.SetStabilization(on)
		return nil
	})

	s.AddCommand("SetDataStreaming", func(params map[string]interface{}) interface{} {
		N := uint16(params["N"].(float64))
		M := uint16(params["M"].(float64))
		Mask := uint32(params["Mask"].(float64))
		Pcnt := uint8(params["Pcnt"].(float64))
		Mask2 := uint32(params["Mask2"].(float64))

		s.SetDataStreaming(DataStreamingConfig{N: N, M: M, Mask2: Mask2, Pcnt: Pcnt, Mask: Mask})
		return nil
	})

	s.AddCommand("ConfigureLocator", func(params map[string]interface{}) interface{} {
		Flags := uint8(params["Flags"].(float64))
		X := int16(params["X"].(float64))
		Y := int16(params["Y"].(float64))
		YawTare := int16(params["YawTare"].(float64))

		s.ConfigureLocator(LocatorConfig{Flags: Flags, X: X, Y: Y, YawTare: YawTare})
		return nil
	})

	return s
}

// Name returns the Driver Name
func (s *SpheroDriver) Name() string { return s.name }

// SetName sets the Driver Name
func (s *SpheroDriver) SetName(n string) { s.name = n }

// Connection returns the Driver's Connection
func (s *SpheroDriver) Connection() gobot.Connection { return s.connection }

func (s *SpheroDriver) adaptor() *Adaptor {
	return s.Connection().(*Adaptor)
}

// Start starts the SpheroDriver and enables Collision Detection.
// Returns true on successful start.
//
// Emits the Events:
// 	Collision  sphero.CollisionPacket - On Collision Detected
// 	SensorData sphero.DataStreamingPacket - On Data Streaming event
// 	Error      error- On error while processing asynchronous response
func (s *SpheroDriver) Start() (err error) {
	go func() {
		for {
			packet := <-s.packetChannel
			err := s.write(packet)
			if err != nil {
				s.Publish(Error, err)
			}
		}
	}()

	go func() {
		for {
			response := <-s.responseChannel
			s.syncResponse = append(s.syncResponse, response)
		}
	}()

	go func() {
		for {
			header := s.readHeader()
			if len(header) > 0 {
				body := s.readBody(header[4])
				data := append(header, body...)
				checksum := data[len(data)-1]
				if checksum != calculateChecksum(data[2:len(data)-1]) {
					continue
				}
				switch header[1] {
				case 0xFE:
					s.asyncResponse = append(s.asyncResponse, data)
				case 0xFF:
					s.responseChannel <- data
				}
			}
		}
	}()

	go func() {
		for {
			var evt []uint8
			for len(s.asyncResponse) != 0 {
				evt, s.asyncResponse = s.asyncResponse[len(s.asyncResponse)-1], s.asyncResponse[:len(s.asyncResponse)-1]
				if evt[2] == 0x07 {
					s.handleCollisionDetected(evt)
				} else if evt[2] == 0x03 {
					s.handleDataStreaming(evt)
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	s.ConfigureCollisionDetection(DefaultCollisionConfig())
	s.enableStopOnDisconnect()

	return
}

// Halt halts the SpheroDriver and sends a SpheroDriver.Stop command to the Sphero.
// Returns true on successful halt.
func (s *SpheroDriver) Halt() (err error) {
	if s.adaptor().connected {
		gobot.Every(10*time.Millisecond, func() {
			s.Stop()
		})
		time.Sleep(1 * time.Second)
	}
	return
}

// SetRGB sets the Sphero to the given r, g, and b values
func (s *SpheroDriver) SetRGB(r uint8, g uint8, b uint8) {
	s.packetChannel <- s.craftPacket([]uint8{r, g, b, 0x01}, 0x02, 0x20)
}

// GetRGB returns the current r, g, b value of the Sphero
func (s *SpheroDriver) GetRGB() []uint8 {
	buf := s.getSyncResponse(s.craftPacket([]uint8{}, 0x02, 0x22))
	if len(buf) == 9 {
		return []uint8{buf[5], buf[6], buf[7]}
	}
	return []uint8{}
}

// ReadLocator reads Sphero's current position (X,Y), component velocities and SOG (speed over ground).
func (s *SpheroDriver) ReadLocator() []int16 {
	buf := s.getSyncResponse(s.craftPacket([]uint8{}, 0x02, 0x15))
	if len(buf) == 16 {
		vals := make([]int16, 5)
		_ = binary.Read(bytes.NewReader(buf[5:15]), binary.BigEndian, &vals)
		return vals
	}
	return []int16{}
}

// SetBackLED sets the Sphero Back LED to the specified brightness
func (s *SpheroDriver) SetBackLED(level uint8) {
	s.packetChannel <- s.craftPacket([]uint8{level}, 0x02, 0x21)
}

// SetRotationRate sets the Sphero rotation rate
// A value of 255 jumps to the maximum (currently 400 degrees/sec).
func (s *SpheroDriver) SetRotationRate(level uint8) {
	s.packetChannel <- s.craftPacket([]uint8{level}, 0x02, 0x03)
}

// SetHeading sets the heading of the Sphero
func (s *SpheroDriver) SetHeading(heading uint16) {
	s.packetChannel <- s.craftPacket([]uint8{uint8(heading >> 8), uint8(heading & 0xFF)}, 0x02, 0x01)
}

// SetStabilization enables or disables the built-in auto stabilizing features of the Sphero
func (s *SpheroDriver) SetStabilization(on bool) {
	b := uint8(0x01)
	if !on {
		b = 0x00
	}
	s.packetChannel <- s.craftPacket([]uint8{b}, 0x02, 0x02)
}

// Roll sends a roll command to the Sphero gives a speed and heading
func (s *SpheroDriver) Roll(speed uint8, heading uint16) {
	s.packetChannel <- s.craftPacket([]uint8{speed, uint8(heading >> 8), uint8(heading & 0xFF), 0x01}, 0x02, 0x30)
}

// ConfigureLocator configures and enables the Locator
func (s *SpheroDriver) ConfigureLocator(d LocatorConfig) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, d)

	s.packetChannel <- s.craftPacket(buf.Bytes(), 0x02, 0x13)
}

// SetDataStreaming enables sensor data streaming
func (s *SpheroDriver) SetDataStreaming(d DataStreamingConfig) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, d)

	s.packetChannel <- s.craftPacket(buf.Bytes(), 0x02, 0x11)
}

// Stop sets the Sphero to a roll speed of 0
func (s *SpheroDriver) Stop() {
	s.Roll(0, 0)
}

// ConfigureCollisionDetection configures the sensitivity of the detection.
func (s *SpheroDriver) ConfigureCollisionDetection(cc CollisionConfig) {
	s.packetChannel <- s.craftPacket([]uint8{cc.Method, cc.Xt, cc.Yt, cc.Xs, cc.Ys, cc.Dead}, 0x02, 0x12)
}

func (s *SpheroDriver) enableStopOnDisconnect() {
	s.packetChannel <- s.craftPacket([]uint8{0x00, 0x00, 0x00, 0x01}, 0x02, 0x37)
}

func (s *SpheroDriver) handleCollisionDetected(data []uint8) {
	// ensure data is the right length:
	if len(data) != 22 || data[4] != 17 {
		return
	}
	var collision CollisionPacket
	buffer := bytes.NewBuffer(data[5:]) // skip header
	binary.Read(buffer, binary.BigEndian, &collision)
	s.Publish(Collision, collision)
}

func (s *SpheroDriver) handleDataStreaming(data []uint8) {
	// ensure data is the right length:
	if len(data) != 90 {
		return
	}
	var dataPacket DataStreamingPacket
	buffer := bytes.NewBuffer(data[5:]) // skip header
	binary.Read(buffer, binary.BigEndian, &dataPacket)
	s.Publish(SensorData, dataPacket)
}

func (s *SpheroDriver) getSyncResponse(packet *packet) []byte {
	s.packetChannel <- packet
	for i := 0; i < 500; i++ {
		for key := range s.syncResponse {
			if s.syncResponse[key][3] == packet.header[4] && len(s.syncResponse[key]) > 6 {
				var response []byte
				response, s.syncResponse = s.syncResponse[len(s.syncResponse)-1], s.syncResponse[:len(s.syncResponse)-1]
				return response
			}
		}
		time.Sleep(100 * time.Microsecond)
	}

	return []byte{}
}

func (s *SpheroDriver) craftPacket(body []uint8, did byte, cid byte) *packet {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	packet := new(packet)
	packet.body = body
	dlen := len(packet.body) + 1
	packet.header = []uint8{0xFF, 0xFF, did, cid, s.seq, uint8(dlen)}
	packet.checksum = s.calculateChecksum(packet)
	return packet
}

func (s *SpheroDriver) write(packet *packet) (err error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	buf := append(packet.header, packet.body...)
	buf = append(buf, packet.checksum)
	length, err := s.adaptor().sp.Write(buf)
	if err != nil {
		return err
	} else if length != len(buf) {
		return errors.New("Not enough bytes written")
	}
	s.seq++
	return
}

func (s *SpheroDriver) calculateChecksum(packet *packet) uint8 {
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

func (s *SpheroDriver) readHeader() []uint8 {
	return s.readNextChunk(5)
}

func (s *SpheroDriver) readBody(length uint8) []uint8 {
	return s.readNextChunk(int(length))
}

func (s *SpheroDriver) readNextChunk(length int) []uint8 {
	read := make([]uint8, length)
	bytesRead := 0

	for bytesRead < length {
		time.Sleep(1 * time.Millisecond)
		n, err := s.adaptor().sp.Read(read[bytesRead:])
		if err != nil {
			return nil
		}
		bytesRead += n
	}
	return read
}
