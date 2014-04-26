package gobotFirmata

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"time"
)

const (
	OPEN                        byte = 1
	CLOSE                       byte = 0
	INPUT                       byte = 0x00
	OUTPUT                      byte = 0x01
	ANALOG                      byte = 0x02
	PWM                         byte = 0x03
	SERVO                       byte = 0x04
	LOW                         byte = 0
	HIGH                        byte = 1
	REPORT_VERSION              byte = 0xF9
	SYSTEM_RESET                byte = 0xFF
	DIGITAL_MESSAGE             byte = 0x90
	DIGITAL_MESSAGE_RANGE_START byte = 0x90
	DIGITAL_MESSAGE_RANGE_END   byte = 0x9F
	ANALOG_MESSAGE              byte = 0xE0
	ANALOG_MESSAGE_RANGE_START  byte = 0xE0
	ANALOG_MESSAGE_RANGE_END    byte = 0xEF
	REPORT_ANALOG               byte = 0xC0
	REPORT_DIGITAL              byte = 0xD0
	PIN_MODE                    byte = 0xF4
	START_SYSEX                 byte = 0xF0
	END_SYSEX                   byte = 0xF7
	CAPABILITY_QUERY            byte = 0x6B
	CAPABILITY_RESPONSE         byte = 0x6C
	PIN_STATE_QUERY             byte = 0x6D
	PIN_STATE_RESPONSE          byte = 0x6E
	ANALOG_MAPPING_QUERY        byte = 0x69
	ANALOG_MAPPING_RESPONSE     byte = 0x6A
	STRING_DATA                 byte = 0x71
	I2C_REQUEST                 byte = 0x76
	I2C_REPLY                   byte = 0x77
	I2C_CONFIG                  byte = 0x78
	FIRMWARE_QUERY              byte = 0x79
	I2C_MODE_WRITE              byte = 0x00
	I2C_MODE_READ               byte = 0x01
	I2C_MODE_CONTINUOUS_READ    byte = 0x02
	I2C_MODE_STOP_READING       byte = 0x03
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
	I2cReply map[string][]uint16
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
	b.togglePinReporting(0, HIGH, REPORT_DIGITAL)
	time.Sleep(50 * time.Millisecond)
	b.togglePinReporting(1, HIGH, REPORT_DIGITAL)
	time.Sleep(50 * time.Millisecond)
}

func (b *board) findEvents(name string) []event {
	ret := make([]event, 0)
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
	b.write([]byte{SYSTEM_RESET})
}

func (b *board) setPinMode(pin byte, mode byte) {
	b.Pins[pin].Mode = mode
	b.write([]byte{PIN_MODE, pin, mode})
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
	b.write([]byte{DIGITAL_MESSAGE | port, portValue & 0x7F, (portValue >> 7) & 0x7F})
}

func (b *board) analogWrite(pin byte, value byte) {
	b.Pins[pin].Value = int(value)
	b.write([]byte{ANALOG_MESSAGE | pin, value & 0x7F, (value >> 7) & 0x7F})
}

func (b *board) version() string {
	return fmt.Sprintf("%v.%v", b.MajorVersion, b.MinorVersion)
}

func (b *board) reportVersion() {
	b.write([]byte{REPORT_VERSION})
}

func (b *board) queryFirmware() {
	b.write([]byte{START_SYSEX, FIRMWARE_QUERY, END_SYSEX})
}

func (b *board) queryPinState(pin byte) {
	b.write([]byte{START_SYSEX, PIN_STATE_QUERY, pin, END_SYSEX})
}

func (b *board) queryReportVersion() {
	b.write([]byte{REPORT_VERSION})
}

func (b *board) queryCapabilities() {
	b.write([]byte{START_SYSEX, CAPABILITY_QUERY, END_SYSEX})
}

func (b *board) queryAnalogMapping() {
	b.write([]byte{START_SYSEX, ANALOG_MAPPING_QUERY, END_SYSEX})
}

func (b *board) togglePinReporting(pin byte, state byte, mode byte) {
	b.write([]byte{mode | pin, state})
}

func (b *board) i2cReadRequest(slave_address byte, num_bytes uint16) {
	b.write([]byte{START_SYSEX, I2C_REQUEST, slave_address, (I2C_MODE_READ << 3), byte(num_bytes & 0x7F), byte(((num_bytes >> 7) & 0x7F)), END_SYSEX})
}

func (b *board) i2cWriteRequest(slave_address byte, data []uint16) {
	ret := []byte{START_SYSEX, I2C_REQUEST, slave_address, (I2C_MODE_WRITE << 3)}
	for _, val := range data {
		ret = append(ret, byte(val&0x7F))
		ret = append(ret, byte((val>>7)&0x7F))
	}
	ret = append(ret, END_SYSEX)
	b.write(ret)
}

func (b *board) i2cConfig(data []uint16) {
	ret := []byte{START_SYSEX, I2C_CONFIG}
	for _, val := range data {
		ret = append(ret, byte(val&0xFF))
		ret = append(ret, byte((val>>8)&0xFF))
	}
	ret = append(ret, END_SYSEX)
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

func (me *board) process(data []byte) {
	buf := bytes.NewBuffer(data)
	for {
		b, err := buf.ReadByte()
		if err != nil {
			break
		}
		switch b {
		case REPORT_VERSION:
			me.MajorVersion, _ = buf.ReadByte()
			me.MinorVersion, _ = buf.ReadByte()
			me.Events = append(me.Events, event{Name: "report_version"})
		case ANALOG_MESSAGE:
			least_significant_byte, _ := buf.ReadByte()
			most_significant_byte, _ := buf.ReadByte()

			value := uint(least_significant_byte) | uint(most_significant_byte)<<7
			pin := (b & 0x0F)

			me.Pins[me.AnalogPins[pin]].Value = int(value)
			me.Events = append(me.Events, event{Name: fmt.Sprintf("analog_read_%v", pin), Data: []byte{byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value & 0xff)}})

		case DIGITAL_MESSAGE:
			port := b & 0x0F
			first_bitmask, _ := buf.ReadByte()
			second_bitmask, _ := buf.ReadByte()
			port_value := first_bitmask | (second_bitmask << 7)

			for i := 0; i < 8; i++ {
				pin_number := (8*byte(port) + byte(i))
				pin := me.Pins[pin_number]
				if byte(pin.Mode) == INPUT {
					pin.Value = int((port_value >> (byte(i) & 0x07)) & 0x01)
					me.Events = append(me.Events, event{Name: fmt.Sprintf("digital_read_%v", pin_number), Data: []byte{byte(pin.Value & 0xff)}})
				}
			}

		case START_SYSEX:
			current_buffer := []byte{b}
			for {
				b, err := buf.ReadByte()
				if err != nil {
					break
				}
				current_buffer = append(current_buffer, b)
				if current_buffer[len(current_buffer)-1] == END_SYSEX {
					break
				}
			}
			command := current_buffer[1]
			switch command {
			case CAPABILITY_RESPONSE:
				supported_modes := 0
				n := 0

				for _, val := range current_buffer[2:(len(current_buffer) - 5)] {
					if val == 127 {
						modes := make([]byte, 0)
						for _, mode := range []byte{INPUT, OUTPUT, ANALOG, PWM, SERVO} {
							if (supported_modes & (1 << mode)) != 0 {
								modes = append(modes, mode)
							}
						}
						me.Pins = append(me.Pins, pin{modes, OUTPUT, 0, 0})
						supported_modes = 0
						n = 0
						continue
					}

					if n == 0 {
						supported_modes = supported_modes | (1 << val)
					}
					n ^= 1
				}
				me.Events = append(me.Events, event{Name: "capability_query"})

			case ANALOG_MAPPING_RESPONSE:
				pin_index := byte(0)

				for _, val := range current_buffer[2 : len(current_buffer)-1] {

					me.Pins[pin_index].AnalogChannel = val

					if val != 127 {
						me.AnalogPins = append(me.AnalogPins, pin_index)
					}

					pin_index += 1
				}

				me.Events = append(me.Events, event{Name: "analog_mapping_query"})

			case PIN_STATE_RESPONSE:
				pin := me.Pins[current_buffer[2]]
				pin.Mode = current_buffer[3]
				pin.Value = int(current_buffer[4])

				if len(current_buffer) > 6 {
					pin.Value = int(uint(pin.Value) | uint(current_buffer[5])<<7)
				}
				if len(current_buffer) > 7 {
					pin.Value = int(uint(pin.Value) | uint(current_buffer[6])<<14)
				}

				me.Events = append(me.Events, event{Name: fmt.Sprintf("pin_%v_state", current_buffer[2]), Data: []byte{byte(pin.Value & 0xff)}})
			case I2C_REPLY:
				i2c_reply := map[string][]uint16{
					"slave_address": []uint16{uint16(current_buffer[2]) | uint16(current_buffer[3])<<7},
					"register":      []uint16{uint16(current_buffer[4]) | uint16(current_buffer[5])<<7},
					"data":          []uint16{uint16(current_buffer[6]) | uint16(current_buffer[7])<<7},
				}
				for i := 8; i < len(current_buffer); i = i + 2 {
					if current_buffer[i] == byte(0xF7) {
						break
					}
					if i+2 > len(current_buffer) {
						break
					}
					i2c_reply["data"] = append(i2c_reply["data"], uint16(current_buffer[i])|uint16(current_buffer[i+1])<<7)
				}
				me.Events = append(me.Events, event{Name: "i2c_reply", I2cReply: i2c_reply})

			case FIRMWARE_QUERY:
				name := make([]byte, 0)
				for _, val := range current_buffer[4:(len(current_buffer) - 1)] {
					if val != 0 {
						name = append(name, val)
					}
				}
				me.FirmwareName = string(name[:])
				me.Events = append(me.Events, event{Name: "firmware_query"})
			case STRING_DATA:
				str := current_buffer[2 : len(current_buffer)-1]
				fmt.Println(string(str[:len(str)]))
				me.Events = append(me.Events, event{Name: "string_data", Data: str})
			default:
				fmt.Println("bad byte", fmt.Sprintf("0x%x", command))
			}
		}
	}
}
