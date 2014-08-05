package sphero

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
)

type packet struct {
	header   []uint8
	body     []uint8
	checksum uint8
}

type SpheroDriver struct {
	gobot.Driver
	seq             uint8
	asyncResponse   [][]uint8
	syncResponse    [][]uint8
	packetChannel   chan *packet
	responseChannel chan []uint8
}

func NewSpheroDriver(a *SpheroAdaptor, name string) *SpheroDriver {
	s := &SpheroDriver{
		Driver: *gobot.NewDriver(
			name,
			"SpheroDriver",
			a,
		),
		packetChannel:   make(chan *packet, 1024),
		responseChannel: make(chan []uint8, 1024),
	}

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

func (s *SpheroDriver) adaptor() *SpheroAdaptor {
	return s.Adaptor().(*SpheroAdaptor)
}

func (s *SpheroDriver) Init() bool {
	return true
}

func (s *SpheroDriver) Start() bool {
	go func() {
		for {
			packet := <-s.packetChannel
			s.write(packet)
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
				if header[1] == 0xFE {
					async := append(header, body...)
					s.asyncResponse = append(s.asyncResponse, async)
				} else {
					s.responseChannel <- append(header, body...)
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

	return true
}

func (s *SpheroDriver) Halt() bool {
	gobot.Every(10*time.Millisecond, func() {
		s.Stop()
	})
	time.Sleep(1 * time.Second)
	return true
}

func (s *SpheroDriver) SetRGB(r uint8, g uint8, b uint8) {
	s.packetChannel <- s.craftPacket([]uint8{r, g, b, 0x01}, 0x02, 0x20)
}

func (s *SpheroDriver) GetRGB() []uint8 {
	return s.getSyncResponse(s.craftPacket([]uint8{}, 0x02, 0x22))
}

func (s *SpheroDriver) SetBackLED(level uint8) {
	s.packetChannel <- s.craftPacket([]uint8{level}, 0x02, 0x21)
}

func (s *SpheroDriver) SetHeading(heading uint16) {
	s.packetChannel <- s.craftPacket([]uint8{uint8(heading >> 8), uint8(heading & 0xFF)}, 0x02, 0x01)
}

func (s *SpheroDriver) SetStabilization(on bool) {
	b := uint8(0x01)
	if on == false {
		b = 0x00
	}
	s.packetChannel <- s.craftPacket([]uint8{b}, 0x02, 0x02)
}

func (s *SpheroDriver) Roll(speed uint8, heading uint16) {
	s.packetChannel <- s.craftPacket([]uint8{speed, uint8(heading >> 8), uint8(heading & 0xFF), 0x01}, 0x02, 0x30)
}

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
	gobot.Publish(s.Event("collision"), data)
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
		time.Sleep(10 * time.Microsecond)
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

func (s *SpheroDriver) write(packet *packet) {
	buf := append(packet.header, packet.body...)
	buf = append(buf, packet.checksum)
	length, err := s.adaptor().sp.Write(buf)
	if err != nil {
		fmt.Println(s.Name, err)
		s.adaptor().Disconnect()
		fmt.Println("Reconnecting to SpheroDriver...")
		s.adaptor().Connect()
		return
	} else if length != len(buf) {
		fmt.Println("Not enough bytes written", s.Name)
	}
	s.seq++
}

func (s *SpheroDriver) calculateChecksum(packet *packet) uint8 {
	buf := append(packet.header, packet.body...)
	buf = buf[2:]
	var calculatedChecksum uint16
	for i := range buf {
		calculatedChecksum += uint16(buf[i])
	}
	return uint8(^(calculatedChecksum % 256))
}

func (s *SpheroDriver) readHeader() []uint8 {
	data := s.readNextChunk(5)
	if data == nil {
		return nil
	}
	return data
}

func (s *SpheroDriver) readBody(length uint8) []uint8 {
	data := s.readNextChunk(length)
	if data == nil {
		return nil
	}
	return data
}

func (s *SpheroDriver) readNextChunk(length uint8) []uint8 {
	time.Sleep(1000 * time.Microsecond)
	var read = make([]uint8, int(length))
	l, err := s.adaptor().sp.Read(read[:])
	if err != nil || length != uint8(l) {
		return nil
	}
	return read
}
