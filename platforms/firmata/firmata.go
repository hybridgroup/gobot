package firmata

import (
	"bytes"
	"errors"
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

var defaultInitTimeInterval time.Duration = 1 * time.Second

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

// newBoard creates a new board connected in specified serial port.
// Adds following events: "firmware_query", "capability_query",
// "analog_mapping_query", "report_version", "i2c_reply",
// "string_data", "firmware_query"
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
		initTimeInterval: defaultInitTimeInterval,
	}

	for _, s := range []string{
		"firmware_query",
		"capability_query",
		"analog_mapping_query",
		"report_version",
		"i2c_reply",
		"string_data",
		"firmware_query",
	} {
		board.events[s] = gobot.NewEvent()
	}

	return board
}

// connect starts connection to board.
// Queries report version until connected
func (b *board) connect() (err error) {
	if b.connected == false {
		if err = b.reset(); err != nil {
			return err
		}
		b.initBoard()
		for {
			if err = b.queryReportVersion(); err != nil {
				return err
			}
			<-time.After(b.initTimeInterval)
			if err = b.readAndProcess(); err != nil {
				return err
			}
			if b.connected == true {
				break
			}
		}
	}
	return
}

// initBoard initializes board by listening for "firware_query", "capability_query"
// and "analog_mapping_query" events
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

// readAndProcess reads from serial port and parses data.
func (b *board) readAndProcess() error {
	buf, err := b.read()
	if err != nil {
		return err
	}
	return b.process(buf)
}

// reset writes system reset bytes.
func (b *board) reset() error {
	return b.write([]byte{systemReset})
}

// setPinMode writes pin mode bytes for specified pin.
func (b *board) setPinMode(pin byte, mode byte) error {
	b.pins[pin].mode = mode
	return b.write([]byte{pinMode, pin, mode})
}

// digitalWrite is used to send a digital value to a specified pin.
func (b *board) digitalWrite(pin byte, value byte) error {
	port := byte(math.Floor(float64(pin) / 8))
	portValue := byte(0)

	b.pins[pin].value = int(value)

	for i := byte(0); i < 8; i++ {
		if b.pins[8*port+i].value != 0 {
			portValue = portValue | (1 << i)
		}
	}
	return b.write([]byte{digitalMessage | port, portValue & 0x7F, (portValue >> 7) & 0x7F})
}

// analogWrite writes value to specified pin
func (b *board) analogWrite(pin byte, value byte) error {
	b.pins[pin].value = int(value)
	return b.write([]byte{analogMessage | pin, value & 0x7F, (value >> 7) & 0x7F})
}

// version returns board version following MAYOR.minor convention.
func (b *board) version() string {
	return fmt.Sprintf("%v.%v", b.majorVersion, b.minorVersion)
}

// queryFirmware writes bytes to query firmware from board.
func (b *board) queryFirmware() error {
	return b.write([]byte{startSysex, firmwareQuery, endSysex})
}

// queryPinState writes bytes to retrieve pin state
func (b *board) queryPinState(pin byte) error {
	return b.write([]byte{startSysex, pinStateQuery, pin, endSysex})
}

// queryReportVersion sends query for report version
func (b *board) queryReportVersion() error {
	return b.write([]byte{reportVersion})
}

// queryCapabilities is used to retrieve board capabilities.
func (b *board) queryCapabilities() error {
	return b.write([]byte{startSysex, capabilityQuery, endSysex})
}

// queryAnalogMapping returns analog mapping for board.
func (b *board) queryAnalogMapping() error {
	return b.write([]byte{startSysex, analogMappingQuery, endSysex})
}

// togglePinReporting is used to change pin reporting mode.
func (b *board) togglePinReporting(pin byte, state byte, mode byte) error {
	return b.write([]byte{mode | pin, state})
}

// i2cReadRequest reads from slaveAddress.
func (b *board) i2cReadRequest(slaveAddress byte, numBytes uint) error {
	return b.write([]byte{startSysex, i2CRequest, slaveAddress, (i2CModeRead << 3),
		byte(numBytes & 0x7F), byte(((numBytes >> 7) & 0x7F)), endSysex})
}

// i2cWriteRequest writes to slaveAddress.
func (b *board) i2cWriteRequest(slaveAddress byte, data []byte) error {
	ret := []byte{startSysex, i2CRequest, slaveAddress, (i2CModeWrite << 3)}
	for _, val := range data {
		ret = append(ret, byte(val&0x7F))
		ret = append(ret, byte((val>>7)&0x7F))
	}
	ret = append(ret, endSysex)
	return b.write(ret)
}

// i2xConfig returns i2c configuration.
func (b *board) i2cConfig(data []byte) error {
	ret := []byte{startSysex, i2CConfig}
	for _, val := range data {
		ret = append(ret, byte(val&0xFF))
		ret = append(ret, byte((val>>8)&0xFF))
	}
	ret = append(ret, endSysex)
	return b.write(ret)
}

// write is used to send commands to serial port
func (b *board) write(commands []byte) (err error) {
	_, err = b.serial.Write(commands[:])
	return
}

// read returns buffer reading from serial port (1024 bytes)
func (b *board) read() (buf []byte, err error) {
	buf = make([]byte, 1024)
	_, err = b.serial.Read(buf)
	return
}

// process uses incoming data and executes actions depending on what is received.
// The following messages are processed: reportVersion, AnalogMessageRangeStart,
// digitalMessageRangeStart.
// And the following responses: capability, analog mapping, pin state,
// i2c, firmwareQuery, string data.
// If neither of those messages is received, then data is treated as "bad_byte"
func (b *board) process(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	for {
		messageType, err := buf.ReadByte()
		if err != nil {
			// we ran out of bytes so we break out of the process loop
			break
		}
		switch {
		case reportVersion == messageType:
			if b.majorVersion, err = buf.ReadByte(); err != nil {
				return err
			}
			if b.minorVersion, err = buf.ReadByte(); err != nil {
				return err
			}
			gobot.Publish(b.events["report_version"], b.version())
		case analogMessageRangeStart <= messageType &&
			analogMessageRangeEnd >= messageType:

			leastSignificantByte, err := buf.ReadByte()
			if err != nil {
				return err
			}
			mostSignificantByte, err := buf.ReadByte()
			if err != nil {
				return err
			}

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
			firstBitmask, err := buf.ReadByte()
			if err != nil {
				return err
			}
			secondBitmask, err := buf.ReadByte()
			if err != nil {
				return err
			}
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
					// we ran out of bytes before we reached the endSysex so we break
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
				return errors.New(fmt.Sprintf("bad byte: 0x%x", command))
			}
		}
	}
	return
}
