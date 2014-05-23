package sphero

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"time"
)

type packet struct {
	header   []uint8
	body     []uint8
	checksum uint8
}

type SpheroDriver struct {
	gobot.Driver
	Adaptor          *SpheroAdaptor
	seq              uint8
	async_response   [][]uint8
	sync_response    [][]uint8
	packet_channel   chan *packet
	response_channel chan []uint8
}

func NewSpheroDriver(a *SpheroAdaptor, name string) *SpheroDriver {
	return &SpheroDriver{
		Driver: gobot.Driver{
			Name:   name,
			Events: make(map[string]chan interface{}),
			Commands: []string{
				"SetRGBC",
				"RollC",
				"StopC",
				"GetRGBC",
				"SetBackLEDC",
				"SetHeadingC",
				"SetStabilizationC",
			},
		},
		Adaptor:          a,
		packet_channel:   make(chan *packet, 1024),
		response_channel: make(chan []uint8, 1024),
	}
}
func (s *SpheroDriver) Init() bool {
	return true
}

func (s *SpheroDriver) Start() bool {
	go func() {
		for {
			packet := <-s.packet_channel
			s.write(packet)
		}
	}()

	go func() {
		for {
			response := <-s.response_channel
			s.sync_response = append(s.sync_response, response)
		}
	}()

	go func() {
		for {
			header := s.readHeader()
			if header != nil && len(header) != 0 {
				body := s.readBody(header[4])
				if header[1] == 0xFE {
					async := append(header, body...)
					s.async_response = append(s.async_response, async)
				} else {
					s.response_channel <- append(header, body...)
				}
			}
		}
	}()

	go func() {
		for {
			var evt []uint8
			for len(s.async_response) != 0 {
				evt, s.async_response = s.async_response[len(s.async_response)-1], s.async_response[:len(s.async_response)-1]
				if evt[2] == 0x07 {
					s.handleCollisionDetected(evt)
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	s.configureCollisionDetection()

	return true
}

func (s *SpheroDriver) Halt() bool {
	go func() {
		for {
			s.Stop()
		}
	}()
	time.Sleep(1 * time.Second)
	return true
}

func (s *SpheroDriver) SetRGB(r uint8, g uint8, b uint8) {
	s.packet_channel <- s.craftPacket([]uint8{r, g, b, 0x01}, 0x20)
}

func (s *SpheroDriver) GetRGB() []uint8 {
	return s.syncResponse(s.craftPacket([]uint8{}, 0x22))
}

func (s *SpheroDriver) SetBackLED(level uint8) {
	s.packet_channel <- s.craftPacket([]uint8{level}, 0x21)
}

func (s *SpheroDriver) SetHeading(heading uint16) {
	s.packet_channel <- s.craftPacket([]uint8{uint8(heading >> 8), uint8(heading & 0xFF)}, 0x01)
}

func (s *SpheroDriver) SetStabilization(on bool) {
	b := uint8(0x01)
	if on == false {
		b = 0x00
	}
	s.packet_channel <- s.craftPacket([]uint8{b}, 0x02)
}

func (s *SpheroDriver) Roll(speed uint8, heading uint16) {
	s.packet_channel <- s.craftPacket([]uint8{speed, uint8(heading >> 8), uint8(heading & 0xFF), 0x01}, 0x30)
}

func (s *SpheroDriver) Stop() {
	s.Roll(0, 0)
}

func (s *SpheroDriver) configureCollisionDetection() {
	s.Events["Collision"] = make(chan interface{})
	s.packet_channel <- s.craftPacket([]uint8{0x01, 0x40, 0x40, 0x50, 0x50, 0x60}, 0x12)
}

func (s *SpheroDriver) handleCollisionDetected(data []uint8) {
	gobot.Publish(s.Events["Collision"], data)
}

func (s *SpheroDriver) syncResponse(packet *packet) []byte {
	s.packet_channel <- packet
	for i := 0; i < 500; i++ {
		for key := range s.sync_response {
			if s.sync_response[key][3] == packet.header[4] && len(s.sync_response[key]) > 6 {
				var response []byte
				response, s.sync_response = s.sync_response[len(s.sync_response)-1], s.sync_response[:len(s.sync_response)-1]
				return response
			}
		}
		time.Sleep(10 * time.Microsecond)
	}

	return make([]byte, 0)
}

func (s *SpheroDriver) craftPacket(body []uint8, cid byte) *packet {
	packet := new(packet)
	packet.body = body
	dlen := len(packet.body) + 1
	packet.header = []uint8{0xFF, 0xFF, 0x02, cid, s.seq, uint8(dlen)}
	packet.checksum = s.calculateChecksum(packet)
	return packet
}

func (s *SpheroDriver) write(packet *packet) {
	buf := append(packet.header, packet.body...)
	buf = append(buf, packet.checksum)
	length, err := s.Adaptor.sp.Write(buf)
	if err != nil {
		fmt.Println(s.Name, err)
		s.Adaptor.Disconnect()
		fmt.Println("Reconnecting to SpheroDriver...")
		s.Adaptor.Connect()
		return
	} else if length != len(buf) {
		fmt.Println("Not enough bytes written", s.Name)
	}
	s.seq += 1
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
	} else {
		return data
	}
}

func (s *SpheroDriver) readBody(length uint8) []uint8 {
	data := s.readNextChunk(length)
	if data == nil {
		return nil
	} else {
		return data
	}
}

func (s *SpheroDriver) readNextChunk(length uint8) []uint8 {
	time.Sleep(1000 * time.Microsecond)
	var read = make([]uint8, int(length))
	l, err := s.Adaptor.sp.Read(read[:])
	if err != nil || length != uint8(l) {
		return nil
	} else {
		return read
	}
}
