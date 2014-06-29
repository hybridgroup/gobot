package firmata

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"time"
)

const (
	Open                     byte = 1
	Close                    byte = 0
	Input                    byte = 0x00
	Output                   byte = 0x01
	Analog                   byte = 0x02
	PWM                      byte = 0x03
	Servo                    byte = 0x04
	Low                      byte = 0
	High                     byte = 1
	ReportVersion            byte = 0xF9
	SystemReset              byte = 0xFF
	DigitalMessage           byte = 0x90
	DigitalMessageRangeStart byte = 0x90
	DigitalMessageRangeEnd   byte = 0x9F
	AnalogMessage            byte = 0xE0
	AnalogMessageRangeStart  byte = 0xE0
	AnalogMessageRangeEnd    byte = 0xEF
	ReportAnalog             byte = 0xC0
	ReportDigital            byte = 0xD0
	PinMode                  byte = 0xF4
	StartSysex               byte = 0xF0
	EndSysex                 byte = 0xF7
	CapabilityQuery          byte = 0x6B
	CapabilityResponse       byte = 0x6C
	PinStateQuery            byte = 0x6D
	PinStateResponse         byte = 0x6E
	AnalogMappingQuery       byte = 0x69
	AnalogMappingResponse    byte = 0x6A
	StringData               byte = 0x71
	I2CRequest               byte = 0x76
	I2CReply                 byte = 0x77
	I2CConfig                byte = 0x78
	FirmwareQuery            byte = 0x79
	I2CModeWrite             byte = 0x00
	I2CModeRead              byte = 0x01
	I2CmodeContinuousRead    byte = 0x02
	I2CModeStopReading       byte = 0x03
)

type board struct {
	Serial       io.ReadWriteCloser
	Pins         []pin
	AnalogPins   []byte
	FirmwareName string
	MajorVersion byte
	MinorVersion byte
	Events       []event
	Connected    bool
}

type pin struct {
	SupportedModes []byte
	Mode           byte
	Value          int
	AnalogChannel  byte
}

type event struct {
	Name     string
	Data     []byte
	I2cReply map[string][]byte
}

func newBoard(sp io.ReadWriteCloser) *board {
	board := new(board)
	board.MajorVersion = 0
	board.MinorVersion = 0
	board.Serial = sp
	board.FirmwareName = ""
	board.Pins = make([]pin, 100)
	board.AnalogPins = make([]byte, 0)
	board.Connected = false
	board.Events = make([]event, 0)
	return board
}

func (b *board) connect() {
	if b.Connected == false {
		b.initBoard()
		b.Connected = true

		go func() {
			for {
				b.queryReportVersion()
				time.Sleep(50 * time.Millisecond)
				b.readAndProcess()
			}
		}()
	}
}

func (b *board) initBoard() {
	for {
		b.queryFirmware()
		time.Sleep(50 * time.Millisecond)
		b.readAndProcess()
		if len(b.findEvents("firmware_query")) > 0 {
			break
		}
	}
	for {
		b.queryCapabilities()
		time.Sleep(50 * time.Millisecond)
		b.readAndProcess()
		if len(b.findEvents("capability_query")) > 0 {
			break
		}
	}
	for {
		b.queryAnalogMapping()
		time.Sleep(50 * time.Millisecond)
		b.readAndProcess()
		if len(b.findEvents("analog_mapping_query")) > 0 {
			break
		}
	}
	b.togglePinReporting(0, High, ReportDigital)
	time.Sleep(50 * time.Millisecond)
	b.togglePinReporting(1, High, ReportDigital)
	time.Sleep(50 * time.Millisecond)
}

func (b *board) findEvents(name string) []event {
	ret := []event{}
	for key, val := range b.Events {
		if val.Name == name {
			ret = append(ret, val)
			if len(b.Events) > key+1 {
				b.Events = append(b.Events[:key], b.Events[key+1:]...)
			}
		}
	}
	return ret
}

func (b *board) readAndProcess() {
	b.process(b.read())
}

func (b *board) reset() {
	b.write([]byte{SystemReset})
}

func (b *board) setPinMode(pin byte, mode byte) {
	b.Pins[pin].Mode = mode
	b.write([]byte{PinMode, pin, mode})
}

func (b *board) digitalWrite(pin byte, value byte) {
	port := byte(math.Floor(float64(pin) / 8))
	portValue := byte(0)

	b.Pins[pin].Value = int(value)

	for i := byte(0); i < 8; i++ {
		if b.Pins[8*port+i].Value != 0 {
			portValue = portValue | (1 << i)
		}
	}
	b.write([]byte{DigitalMessage | port, portValue & 0x7F, (portValue >> 7) & 0x7F})
}

func (b *board) analogWrite(pin byte, value byte) {
	b.Pins[pin].Value = int(value)
	b.write([]byte{AnalogMessage | pin, value & 0x7F, (value >> 7) & 0x7F})
}

func (b *board) version() string {
	return fmt.Sprintf("%v.%v", b.MajorVersion, b.MinorVersion)
}

func (b *board) reportVersion() {
	b.write([]byte{ReportVersion})
}

func (b *board) queryFirmware() {
	b.write([]byte{StartSysex, FirmwareQuery, EndSysex})
}

func (b *board) queryPinState(pin byte) {
	b.write([]byte{StartSysex, PinStateQuery, pin, EndSysex})
}

func (b *board) queryReportVersion() {
	b.write([]byte{ReportVersion})
}

func (b *board) queryCapabilities() {
	b.write([]byte{StartSysex, CapabilityQuery, EndSysex})
}

func (b *board) queryAnalogMapping() {
	b.write([]byte{StartSysex, AnalogMappingQuery, EndSysex})
}

func (b *board) togglePinReporting(pin byte, state byte, mode byte) {
	b.write([]byte{mode | pin, state})
}

func (b *board) i2cReadRequest(slaveAddress byte, numBytes uint) {
	b.write([]byte{StartSysex, I2CRequest, slaveAddress, (I2CModeRead << 3), byte(numBytes & 0x7F), byte(((numBytes >> 7) & 0x7F)), EndSysex})
}

func (b *board) i2cWriteRequest(slaveAddress byte, data []byte) {
	ret := []byte{StartSysex, I2CRequest, slaveAddress, (I2CModeWrite << 3)}
	for _, val := range data {
		ret = append(ret, byte(val&0x7F))
		ret = append(ret, byte((val>>7)&0x7F))
	}
	ret = append(ret, EndSysex)
	b.write(ret)
}

func (b *board) i2cConfig(data []byte) {
	ret := []byte{StartSysex, I2CConfig}
	for _, val := range data {
		ret = append(ret, byte(val&0xFF))
		ret = append(ret, byte((val>>8)&0xFF))
	}
	ret = append(ret, EndSysex)
	b.write(ret)
}

func (b *board) write(commands []byte) {
	b.Serial.Write(commands[:])
}

func (b *board) read() []byte {
	buf := make([]byte, 1024)
	b.Serial.Read(buf)
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
		case ReportVersion == messageType:
			b.MajorVersion, _ = buf.ReadByte()
			b.MinorVersion, _ = buf.ReadByte()
			b.Events = append(b.Events, event{Name: "report_version"})
		case AnalogMessageRangeStart <= messageType && AnalogMessageRangeEnd >= messageType:
			leastSignificantByte, _ := buf.ReadByte()
			mostSignificantByte, _ := buf.ReadByte()

			value := uint(leastSignificantByte) | uint(mostSignificantByte)<<7
			pin := (messageType & 0x0F)

			b.Pins[b.AnalogPins[pin]].Value = int(value)
			b.Events = append(b.Events, event{Name: fmt.Sprintf("analog_read_%v", pin), Data: []byte{byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value & 0xff)}})

		case DigitalMessageRangeStart <= messageType && DigitalMessageRangeEnd >= messageType:
			port := messageType & 0x0F
			firstBitmask, _ := buf.ReadByte()
			secondBitmask, _ := buf.ReadByte()
			portValue := firstBitmask | (secondBitmask << 7)

			for i := 0; i < 8; i++ {
				pinNumber := (8*byte(port) + byte(i))
				pin := b.Pins[pinNumber]
				if byte(pin.Mode) == Input {
					pin.Value = int((portValue >> (byte(i) & 0x07)) & 0x01)
					b.Events = append(b.Events, event{Name: fmt.Sprintf("digital_read_%v", pinNumber), Data: []byte{byte(pin.Value & 0xff)}})
				}
			}

		case StartSysex == messageType:
			currentBuffer := []byte{messageType}
			for {
				b, err := buf.ReadByte()
				if err != nil {
					break
				}
				currentBuffer = append(currentBuffer, b)
				if currentBuffer[len(currentBuffer)-1] == EndSysex {
					break
				}
			}
			command := currentBuffer[1]
			switch command {
			case CapabilityResponse:
				supportedModes := 0
				n := 0

				for _, val := range currentBuffer[2:(len(currentBuffer) - 5)] {
					if val == 127 {
						modes := []byte{}
						for _, mode := range []byte{Input, Output, Analog, PWM, Servo} {
							if (supportedModes & (1 << mode)) != 0 {
								modes = append(modes, mode)
							}
						}
						b.Pins = append(b.Pins, pin{modes, Output, 0, 0})
						supportedModes = 0
						n = 0
						continue
					}

					if n == 0 {
						supportedModes = supportedModes | (1 << val)
					}
					n ^= 1
				}
				b.Events = append(b.Events, event{Name: "capability_query"})

			case AnalogMappingResponse:
				pinIndex := byte(0)

				for _, val := range currentBuffer[2 : len(currentBuffer)-1] {

					b.Pins[pinIndex].AnalogChannel = val

					if val != 127 {
						b.AnalogPins = append(b.AnalogPins, pinIndex)
					}

					pinIndex++
				}

				b.Events = append(b.Events, event{Name: "analog_mapping_query"})

			case PinStateResponse:
				pin := b.Pins[currentBuffer[2]]
				pin.Mode = currentBuffer[3]
				pin.Value = int(currentBuffer[4])

				if len(currentBuffer) > 6 {
					pin.Value = int(uint(pin.Value) | uint(currentBuffer[5])<<7)
				}
				if len(currentBuffer) > 7 {
					pin.Value = int(uint(pin.Value) | uint(currentBuffer[6])<<14)
				}

				b.Events = append(b.Events, event{Name: fmt.Sprintf("pin_%v_state", currentBuffer[2]), Data: []byte{byte(pin.Value & 0xff)}})
			case I2CReply:
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
					i2cReply["data"] = append(i2cReply["data"], byte(currentBuffer[i])|byte(currentBuffer[i+1])<<7)
				}
				b.Events = append(b.Events, event{Name: "i2c_reply", I2cReply: i2cReply})

			case FirmwareQuery:
				name := []byte{}
				for _, val := range currentBuffer[4:(len(currentBuffer) - 1)] {
					if val != 0 {
						name = append(name, val)
					}
				}
				b.FirmwareName = string(name[:])
				b.Events = append(b.Events, event{Name: "firmware_query"})
			case StringData:
				str := currentBuffer[2 : len(currentBuffer)-1]
				fmt.Println(string(str[:len(str)]))
				b.Events = append(b.Events, event{Name: "string_data", Data: str})
			default:
				fmt.Println("bad byte", fmt.Sprintf("0x%x", command))
			}
		}
	}
}
