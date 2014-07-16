package firmata

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/hybridgroup/gobot"
)

const (
	open                     byte = 1
	close                    byte = 0
	input                    byte = 0x00
	output                   byte = 0x01
	analog                   byte = 0x02
	pwm                      byte = 0x03
	servo                    byte = 0x04
	low                      byte = 0
	high                     byte = 1
	reportVersion            byte = 0xF9
	systemReset              byte = 0xFF
	digitalMessage           byte = 0x90
	digitalMessageRangeStart byte = 0x90
	digitalMessageRangeEnd   byte = 0x9F
	analogMessage            byte = 0xE0
	analogMessageRangeStart  byte = 0xE0
	analogMessageRangeEnd    byte = 0xEF
	reportAnalog             byte = 0xC0
	reportDigital            byte = 0xD0
	pinMode                  byte = 0xF4
	startSysex               byte = 0xF0
	endSysex                 byte = 0xF7
	capabilityQuery          byte = 0x6B
	capabilityResponse       byte = 0x6C
	pinStateQuery            byte = 0x6D
	pinStateResponse         byte = 0x6E
	analogMappingQuery       byte = 0x69
	analogMappingResponse    byte = 0x6A
	stringData               byte = 0x71
	i2CRequest               byte = 0x76
	i2CReply                 byte = 0x77
	i2CConfig                byte = 0x78
	firmwareQuery            byte = 0x79
	i2CModeWrite             byte = 0x00
	i2CModeRead              byte = 0x01
	i2CmodeContinuousRead    byte = 0x02
	i2CModeStopReading       byte = 0x03
)

type board struct {
	serial           io.ReadWriteCloser
	pins             []pin
	analogPins       []byte
	firmwareName     string
	majorVersion     byte
	minorVersion     byte
	connected        bool
	events           map[string]*gobot.Event
	initTimeInterval time.Duration
}

type pin struct {
	supportedModes []byte
	mode           byte
	value          int
	analogChannel  byte
}

func newBoard(sp io.ReadWriteCloser) *board {
	board := &board{
		majorVersion:     0,
		minorVersion:     0,
		serial:           sp,
		firmwareName:     "",
		pins:             []pin{},
		analogPins:       []byte{},
		connected:        false,
		events:           make(map[string]*gobot.Event),
		initTimeInterval: 1 * time.Second,
	}

	for _, s := range []string{
		"firmware_query",
		"capability_query",
		"analog_mapping_query",
		"report_version",
		"i2c_reply",
		"analog_mapping_query",
		"string_data",
		"firmware_query",
	} {
		board.events[s] = gobot.NewEvent()
	}

	return board
}

func (b *board) connect() {
	if b.connected == false {
		b.reset()
		b.initBoard()

		for {
			b.queryReportVersion()
			<-time.After(b.initTimeInterval)
			b.readAndProcess()
			if b.connected == true {
				break
			}
		}
	}
}

func (b *board) initBoard() {
	gobot.Once(b.events["firmware_query"], func(data interface{}) {
		b.queryCapabilities()
	})

	gobot.Once(b.events["capability_query"], func(data interface{}) {
		b.queryAnalogMapping()
	})

	gobot.Once(b.events["analog_mapping_query"], func(data interface{}) {
		b.togglePinReporting(0, high, reportDigital)
		b.togglePinReporting(1, high, reportDigital)
		b.connected = true
	})
}

func (b *board) readAndProcess() {
	b.process(b.read())
}

func (b *board) reset() {
	b.write([]byte{systemReset})
}

func (b *board) setPinMode(pin byte, mode byte) {
	b.pins[pin].mode = mode
	b.write([]byte{pinMode, pin, mode})
}

func (b *board) digitalWrite(pin byte, value byte) {
	port := byte(math.Floor(float64(pin) / 8))
	portValue := byte(0)

	b.pins[pin].value = int(value)

	for i := byte(0); i < 8; i++ {
		if b.pins[8*port+i].value != 0 {
			portValue = portValue | (1 << i)
		}
	}
	b.write([]byte{digitalMessage | port, portValue & 0x7F, (portValue >> 7) & 0x7F})
}

func (b *board) analogWrite(pin byte, value byte) {
	b.pins[pin].value = int(value)
	b.write([]byte{analogMessage | pin, value & 0x7F, (value >> 7) & 0x7F})
}

func (b *board) version() string {
	return fmt.Sprintf("%v.%v", b.majorVersion, b.minorVersion)
}

func (b *board) reportVersion() {
	b.write([]byte{reportVersion})
}

func (b *board) queryFirmware() {
	b.write([]byte{startSysex, firmwareQuery, endSysex})
}

func (b *board) queryPinState(pin byte) {
	b.write([]byte{startSysex, pinStateQuery, pin, endSysex})
}

func (b *board) queryReportVersion() {
	b.write([]byte{reportVersion})
}

func (b *board) queryCapabilities() {
	b.write([]byte{startSysex, capabilityQuery, endSysex})
}

func (b *board) queryAnalogMapping() {
	b.write([]byte{startSysex, analogMappingQuery, endSysex})
}

func (b *board) togglePinReporting(pin byte, state byte, mode byte) {
	b.write([]byte{mode | pin, state})
}

func (b *board) i2cReadRequest(slaveAddress byte, numBytes uint) {
	b.write([]byte{startSysex, i2CRequest, slaveAddress, (i2CModeRead << 3),
		byte(numBytes & 0x7F), byte(((numBytes >> 7) & 0x7F)), endSysex})
}

func (b *board) i2cWriteRequest(slaveAddress byte, data []byte) {
	ret := []byte{startSysex, i2CRequest, slaveAddress, (i2CModeWrite << 3)}
	for _, val := range data {
		ret = append(ret, byte(val&0x7F))
		ret = append(ret, byte((val>>7)&0x7F))
	}
	ret = append(ret, endSysex)
	b.write(ret)
}

func (b *board) i2cConfig(data []byte) {
	ret := []byte{startSysex, i2CConfig}
	for _, val := range data {
		ret = append(ret, byte(val&0xFF))
		ret = append(ret, byte((val>>8)&0xFF))
	}
	ret = append(ret, endSysex)
	b.write(ret)
}

func (b *board) write(commands []byte) {
	b.serial.Write(commands[:])
}

func (b *board) read() []byte {
	buf := make([]byte, 1024)
	b.serial.Read(buf)
	return buf
}

func (b *board) process(data []byte) {
	buf := bytes.NewBuffer(data)
	for {
		messageType, err := buf.ReadByte()
		if err != nil {
			break
		}
		switch {
		case reportVersion == messageType:
			b.majorVersion, _ = buf.ReadByte()
			b.minorVersion, _ = buf.ReadByte()
			gobot.Publish(b.events["report_version"], b.version())
		case analogMessageRangeStart <= messageType &&
			analogMessageRangeEnd >= messageType:

			leastSignificantByte, _ := buf.ReadByte()
			mostSignificantByte, _ := buf.ReadByte()

			value := uint(leastSignificantByte) | uint(mostSignificantByte)<<7
			pin := (messageType & 0x0F)

			b.pins[b.analogPins[pin]].value = int(value)
			gobot.Publish(b.events[fmt.Sprintf("analog_read_%v", pin)],
				[]byte{
					byte(value >> 24),
					byte(value >> 16),
					byte(value >> 8),
					byte(value & 0xff),
				},
			)
		case digitalMessageRangeStart <= messageType &&
			digitalMessageRangeEnd >= messageType:

			port := messageType & 0x0F
			firstBitmask, _ := buf.ReadByte()
			secondBitmask, _ := buf.ReadByte()
			portValue := firstBitmask | (secondBitmask << 7)

			for i := 0; i < 8; i++ {
				pinNumber := (8*byte(port) + byte(i))
				pin := b.pins[pinNumber]
				if byte(pin.mode) == input {
					pin.value = int((portValue >> (byte(i) & 0x07)) & 0x01)
					gobot.Publish(b.events[fmt.Sprintf("digital_read_%v", pinNumber)],
						[]byte{byte(pin.value & 0xff)})
				}
			}
		case startSysex == messageType:
			currentBuffer := []byte{messageType}
			for {
				b, err := buf.ReadByte()
				if err != nil {
					break
				}
				currentBuffer = append(currentBuffer, b)
				if currentBuffer[len(currentBuffer)-1] == endSysex {
					break
				}
			}
			command := currentBuffer[1]
			switch command {
			case capabilityResponse:
				supportedModes := 0
				n := 0

				for _, val := range currentBuffer[2:(len(currentBuffer) - 5)] {
					if val == 127 {
						modes := []byte{}
						for _, mode := range []byte{input, output, analog, pwm, servo} {
							if (supportedModes & (1 << mode)) != 0 {
								modes = append(modes, mode)
							}
						}
						b.pins = append(b.pins, pin{modes, output, 0, 0})
						b.events[fmt.Sprintf("digital_read_%v", len(b.pins)-1)] = gobot.NewEvent()
						b.events[fmt.Sprintf("pin_%v_state", len(b.pins)-1)] = gobot.NewEvent()
						supportedModes = 0
						n = 0
						continue
					}

					if n == 0 {
						supportedModes = supportedModes | (1 << val)
					}
					n ^= 1
				}
				gobot.Publish(b.events["capability_query"], nil)
			case analogMappingResponse:
				pinIndex := byte(0)

				for _, val := range currentBuffer[2 : len(b.pins)-1] {

					b.pins[pinIndex].analogChannel = val

					if val != 127 {
						b.analogPins = append(b.analogPins, pinIndex)
					}
					b.events[fmt.Sprintf("analog_read_%v", pinIndex)] = gobot.NewEvent()
					pinIndex++
				}

				gobot.Publish(b.events["analog_mapping_query"], nil)
			case pinStateResponse:
				pin := b.pins[currentBuffer[2]]
				pin.mode = currentBuffer[3]
				pin.value = int(currentBuffer[4])

				if len(currentBuffer) > 6 {
					pin.value = int(uint(pin.value) | uint(currentBuffer[5])<<7)
				}
				if len(currentBuffer) > 7 {
					pin.value = int(uint(pin.value) | uint(currentBuffer[6])<<14)
				}

				gobot.Publish(b.events[fmt.Sprintf("pin_%v_state", currentBuffer[2])],
					map[string]int{
						"pin":   int(currentBuffer[2]),
						"mode":  int(pin.mode),
						"value": int(pin.value),
					},
				)
			case i2CReply:
				i2cReply := map[string][]byte{
					"slave_address": []byte{byte(currentBuffer[2]) | byte(currentBuffer[3])<<7},
					"register":      []byte{byte(currentBuffer[4]) | byte(currentBuffer[5])<<7},
					"data":          []byte{byte(currentBuffer[6]) | byte(currentBuffer[7])<<7},
				}
				for i := 8; i < len(currentBuffer); i = i + 2 {
					if currentBuffer[i] == byte(0xF7) {
						break
					}
					if i+2 > len(currentBuffer) {
						break
					}
					i2cReply["data"] = append(i2cReply["data"],
						byte(currentBuffer[i])|byte(currentBuffer[i+1])<<7,
					)
				}
				gobot.Publish(b.events["i2c_reply"], i2cReply)
			case firmwareQuery:
				name := []byte{}
				for _, val := range currentBuffer[4:(len(currentBuffer) - 1)] {
					if val != 0 {
						name = append(name, val)
					}
				}
				b.firmwareName = string(name[:])
				gobot.Publish(b.events["firmware_query"], b.firmwareName)
			case stringData:
				str := currentBuffer[2:len(currentBuffer)]
				gobot.Publish(b.events["string_data"], string(str[:len(str)]))
			default:
				fmt.Println("bad byte", fmt.Sprintf("0x%x", command))
			}
		}
	}
}
