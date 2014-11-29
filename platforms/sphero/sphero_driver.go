package sphero

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*SpheroDriver)(nil)

type packet struct {
	header   []uint8
	body     []uint8
	checksum uint8
}

// Represents a Sphero
type SpheroDriver struct {
	name            string
	connection      gobot.Connection
	seq             uint8
	asyncResponse   [][]uint8
	syncResponse    [][]uint8
	packetChannel   chan *packet
	responseChannel chan []uint8
	gobot.Eventer
	gobot.Commander
}

type Collision struct {
	// Normalized impact components (direction of the collision event):
	X, Y, Z int16
	// Thresholds exceeded by X (1h) and/or Y (2h) axis (bitmask):
	Axis byte
	// Power that cross threshold Xt + Xs:
	XMagnitude, YMagnitude int16
	// Sphero's speed when impact detected:
	Speed uint8
	// Millisecond timer
	Timestamp uint32
}

// NewSpheroDriver returns a new SpheroDriver given a SpheroAdaptor and name.
//
// Adds the following API Commands:
// 	"Roll" - See SpheroDriver.Roll
// 	"Stop" - See SpheroDriver.Stop
// 	"GetRGB" - See SpheroDriver.GetRGB
// 	"SetBackLED" - See SpheroDriver.SetBackLED
// 	"SetHeading" - See SpheroDriver.SetHeading
// 	"SetStabilization" - See SpheroDriver.SetStabilization
func NewSpheroDriver(a *SpheroAdaptor, name string) *SpheroDriver {
	s := &SpheroDriver{
		name:            name,
		connection:      a,
		Eventer:         gobot.NewEventer(),
		Commander:       gobot.NewCommander(),
		packetChannel:   make(chan *packet, 1024),
		responseChannel: make(chan []uint8, 1024),
	}

	s.AddEvent("error")
	s.AddEvent("collision")
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

	s.AddCommand("SetBackLED", func(params map[string]interface{}) interface{} {
		level := uint8(params["level"].(float64))
		s.SetBackLED(level)
		return nil
	})

	s.AddCommand("SetHeading", func(params map[string]interface{}) interface{} {
		heading := uint16(params["heading"].(float64))
		s.SetHeading(heading)
		return nil
	})
	s.AddCommand("SetStabilization", func(params map[string]interface{}) interface{} {
		on := params["heading"].(bool)
		s.SetStabilization(on)
		return nil
	})

	return s
}

func (s *SpheroDriver) Name() string                 { return s.name }
func (s *SpheroDriver) Connection() gobot.Connection { return s.connection }

func (s *SpheroDriver) adaptor() *SpheroAdaptor {
	return s.Connection().(*SpheroAdaptor)
}

// Start starts the SpheroDriver and enables Collision Detection.
// Returns true on successful start.
//
// Emits the Events:
// 	"collision" SpheroDriver.Collision - On Collision Detected
func (s *SpheroDriver) Start() (errs []error) {
	go func() {
		for {
			packet := <-s.packetChannel
			err := s.write(packet)
			if err != nil {
				gobot.Publish(s.Event("error"), err)
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
			if header != nil && len(header) != 0 {
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
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	s.configureCollisionDetection()
	s.enableStopOnDisconnect()

	return
}

// Halt halts the SpheroDriver and sends a SpheroDriver.Stop command to the Sphero.
// Returns true on successful halt.
func (s *SpheroDriver) Halt() (errs []error) {
	gobot.Every(10*time.Millisecond, func() {
		s.Stop()
	})
	time.Sleep(1 * time.Second)
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

// SetBackLED sets the Sphero Back LED to the specified brightness
func (s *SpheroDriver) SetBackLED(level uint8) {
	s.packetChannel <- s.craftPacket([]uint8{level}, 0x02, 0x21)
}

// SetHeading sets the heading of the Sphero
func (s *SpheroDriver) SetHeading(heading uint16) {
	s.packetChannel <- s.craftPacket([]uint8{uint8(heading >> 8), uint8(heading & 0xFF)}, 0x02, 0x01)
}

// SetStabilization enables or disables the built-in auto stabilizing features of the Sphero
func (s *SpheroDriver) SetStabilization(on bool) {
	b := uint8(0x01)
	if on == false {
		b = 0x00
	}
	s.packetChannel <- s.craftPacket([]uint8{b}, 0x02, 0x02)
}

// Roll sends a roll command to the Sphero gives a speed and heading
func (s *SpheroDriver) Roll(speed uint8, heading uint16) {
	s.packetChannel <- s.craftPacket([]uint8{speed, uint8(heading >> 8), uint8(heading & 0xFF), 0x01}, 0x02, 0x30)
}

// Stop sets the Sphero to a roll speed of 0
func (s *SpheroDriver) Stop() {
	s.Roll(0, 0)
}

func (s *SpheroDriver) configureCollisionDetection() {
	s.packetChannel <- s.craftPacket([]uint8{0x01, 0x40, 0x40, 0x50, 0x50, 0x60}, 0x02, 0x12)
}

func (s *SpheroDriver) enableStopOnDisconnect() {
	s.packetChannel <- s.craftPacket([]uint8{0x00, 0x00, 0x00, 0x01}, 0x02, 0x37)
}

func (s *SpheroDriver) handleCollisionDetected(data []uint8) {
	// ensure data is the right length:
	if len(data) != 22 || data[4] != 17 {
		return
	}
	var collision Collision
	buffer := bytes.NewBuffer(data[5:]) // skip header
	binary.Read(buffer, binary.BigEndian, &collision)
	gobot.Publish(s.Event("collision"), collision)
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
	packet := new(packet)
	packet.body = body
	dlen := len(packet.body) + 1
	packet.header = []uint8{0xFF, 0xFF, did, cid, s.seq, uint8(dlen)}
	packet.checksum = s.calculateChecksum(packet)
	return packet
}

func (s *SpheroDriver) write(packet *packet) (err error) {
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
