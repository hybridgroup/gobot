package gobotSphero

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
	SpheroAdaptor    *SpheroAdaptor
	seq              uint8
	async_response   [][]uint8
	sync_response    [][]uint8
	packet_channel   chan *packet
	response_channel chan []uint8
}

func NewSphero(sa *SpheroAdaptor) *SpheroDriver {
	s := new(SpheroDriver)
	s.Events = make(map[string]chan interface{})
	s.SpheroAdaptor = sa
	s.packet_channel = make(chan *packet, 1024)
	s.response_channel = make(chan []uint8, 1024)
	s.Commands = []string{
		"SetRGBC",
		"RollC",
		"StopC",
		"GetRGBC",
		"SetBackLEDC",
		"SetHeadingC",
		"SetStabilizationC",
	}
	return s
}
func (sd *SpheroDriver) Init() bool {
	return true
}

func (sd *SpheroDriver) Start() bool {
	go func() {
		for {
			packet := <-sd.packet_channel
			sd.write(packet)
		}
	}()

	go func() {
		for {
			response := <-sd.response_channel
			sd.sync_response = append(sd.sync_response, response)
		}
	}()

	go func() {
		for {
			header := sd.readHeader()
			if header != nil && len(header) != 0 {
				body := sd.readBody(header[4])
				if header[1] == 0xFE {
					async := append(header, body...)
					sd.async_response = append(sd.async_response, async)
				} else {
					sd.response_channel <- append(header, body...)
				}
			}
		}
	}()

	go func() {
		for {
			var evt []uint8
			for len(sd.async_response) != 0 {
				evt, sd.async_response = sd.async_response[len(sd.async_response)-1], sd.async_response[:len(sd.async_response)-1]
				if evt[2] == 0x07 {
					sd.handleCollisionDetected(evt)
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	sd.configureCollisionDetection()

	return true
}

func (sd *SpheroDriver) Halt() bool {
	go func() {
		for {
			sd.Stop()
		}
	}()
	time.Sleep(1 * time.Second)
	return true
}

func (sd *SpheroDriver) SetRGB(r uint8, g uint8, b uint8) {
	sd.packet_channel <- sd.craftPacket([]uint8{r, g, b, 0x01}, 0x20)
}

func (sd *SpheroDriver) GetRGB() []uint8 {
	return sd.syncResponse(sd.craftPacket([]uint8{}, 0x22))
}

func (sd *SpheroDriver) SetBackLED(level uint8) {
	sd.packet_channel <- sd.craftPacket([]uint8{level}, 0x21)
}

func (sd *SpheroDriver) SetHeading(heading uint16) {
	sd.packet_channel <- sd.craftPacket([]uint8{uint8(heading >> 8), uint8(heading & 0xFF)}, 0x01)
}

func (sd *SpheroDriver) SetStabilization(on bool) {
	b := uint8(0x01)
	if on == false {
		b = 0x00
	}
	sd.packet_channel <- sd.craftPacket([]uint8{b}, 0x02)
}

func (sd *SpheroDriver) Roll(speed uint8, heading uint16) {
	sd.packet_channel <- sd.craftPacket([]uint8{speed, uint8(heading >> 8), uint8(heading & 0xFF), 0x01}, 0x30)
}

func (sd *SpheroDriver) Stop() {
	sd.Roll(0, 0)
}

func (sd *SpheroDriver) configureCollisionDetection() {
	sd.Events["Collision"] = make(chan interface{})
	sd.packet_channel <- sd.craftPacket([]uint8{0x01, 0x40, 0x40, 0x50, 0x50, 0x60}, 0x12)
}

func (sd *SpheroDriver) handleCollisionDetected(data []uint8) {
	gobot.Publish(sd.Events["Collision"], data)
}

func (sd *SpheroDriver) syncResponse(packet *packet) []byte {
	sd.packet_channel <- packet
	for i := 0; i < 500; i++ {
		for key := range sd.sync_response {
			if sd.sync_response[key][3] == packet.header[4] && len(sd.sync_response[key]) > 6 {
				var response []byte
				response, sd.sync_response = sd.sync_response[len(sd.sync_response)-1], sd.sync_response[:len(sd.sync_response)-1]
				return response
			}
		}
		time.Sleep(10 * time.Microsecond)
	}

	return make([]byte, 0)
}

func (sd *SpheroDriver) craftPacket(body []uint8, cid byte) *packet {
	packet := new(packet)
	packet.body = body
	dlen := len(packet.body) + 1
	packet.header = []uint8{0xFF, 0xFF, 0x02, cid, sd.seq, uint8(dlen)}
	packet.checksum = sd.calculateChecksum(packet)
	return packet
}

func (sd *SpheroDriver) write(packet *packet) {
	buf := append(packet.header, packet.body...)
	buf = append(buf, packet.checksum)
	length, err := sd.SpheroAdaptor.sp.Write(buf)
	if err != nil {
		fmt.Println(sd.Name, err)
		sd.SpheroAdaptor.Disconnect()
		fmt.Println("Reconnecting to sphero...")
		sd.SpheroAdaptor.Connect()
		return
	} else if length != len(buf) {
		fmt.Println("Not enough bytes written", sd.Name)
	}
	sd.seq += 1
}

func (sd *SpheroDriver) calculateChecksum(packet *packet) uint8 {
	buf := append(packet.header, packet.body...)
	buf = buf[2:]
	var calculatedChecksum uint16
	for i := range buf {
		calculatedChecksum += uint16(buf[i])
	}
	return uint8(^(calculatedChecksum % 256))
}

func (sd *SpheroDriver) readHeader() []uint8 {
	data := sd.readNextChunk(5)
	if data == nil {
		return nil
	} else {
		return data
	}
}

func (sd *SpheroDriver) readBody(length uint8) []uint8 {
	data := sd.readNextChunk(length)
	if data == nil {
		return nil
	} else {
		return data
	}
}

func (sd *SpheroDriver) readNextChunk(length uint8) []uint8 {
	time.Sleep(1000 * time.Microsecond)
	var read = make([]uint8, int(length))
	l, err := sd.SpheroAdaptor.sp.Read(read[:])
	if err != nil || length != uint8(l) {
		return nil
	} else {
		return read
	}
}
