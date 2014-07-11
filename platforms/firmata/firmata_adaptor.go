package firmata

import (
	"fmt"
	"strconv"

	"github.com/hybridgroup/gobot"
	"github.com/tarm/goserial"
)

type FirmataAdaptor struct {
	gobot.Adaptor
	Board      *board
	i2cAddress byte
	connect    func(*FirmataAdaptor)
}

func NewFirmataAdaptor(name, port string) *FirmataAdaptor {
	return &FirmataAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"FirmataAdaptor",
			port,
		),
		connect: func(f *FirmataAdaptor) {
			sp, err := serial.OpenPort(&serial.Config{Name: f.Port(), Baud: 57600})
			if err != nil {
				panic(err)
			}
			f.Board = newBoard(sp)
		},
	}
}

func (f *FirmataAdaptor) Connect() bool {
	f.connect(f)
	f.Board.connect()
	f.SetConnected(true)
	return true
}

func (f *FirmataAdaptor) Disconnect() bool {
	err := f.Board.Serial.Close()
	if err != nil {
		fmt.Println(err)
	}
	return true
}
func (f *FirmataAdaptor) Finalize() bool { return f.Disconnect() }

func (f *FirmataAdaptor) InitServo() {}
func (f *FirmataAdaptor) ServoWrite(pin string, angle byte) {
	p, _ := strconv.Atoi(pin)

	f.Board.setPinMode(byte(p), Servo)
	f.Board.analogWrite(byte(p), angle)
}

func (f *FirmataAdaptor) PwmWrite(pin string, level byte) {
	p, _ := strconv.Atoi(pin)

	f.Board.setPinMode(byte(p), PWM)
	f.Board.analogWrite(byte(p), level)
}

func (f *FirmataAdaptor) DigitalWrite(pin string, level byte) {
	p, _ := strconv.Atoi(pin)

	f.Board.setPinMode(byte(p), Output)
	f.Board.digitalWrite(byte(p), level)
}

func (f *FirmataAdaptor) DigitalRead(pin string) int {
	p, _ := strconv.Atoi(pin)
	f.Board.setPinMode(byte(p), Input)
	f.Board.togglePinReporting(byte(p), High, ReportDigital)
	events := f.Board.findEvents(fmt.Sprintf("digital_read_%v", pin))
	if len(events) > 0 {
		return int(events[len(events)-1].Data[0])
	}
	return -1
}

// NOTE pins are numbered A0-A5, which translate to digital pins 14-19
func (f *FirmataAdaptor) AnalogRead(pin string) int {
	p, _ := strconv.Atoi(pin)
	p = f.digitalPin(p)
	f.Board.setPinMode(byte(p), Analog)
	f.Board.togglePinReporting(byte(p), High, ReportAnalog)
	events := f.Board.findEvents(fmt.Sprintf("analog_read_%v", pin))
	if len(events) > 0 {
		event := events[len(events)-1]
		return int(uint(event.Data[0])<<24 |
			uint(event.Data[1])<<16 |
			uint(event.Data[2])<<8 |
			uint(event.Data[3]))
	}
	return -1
}

func (f *FirmataAdaptor) AnalogWrite(pin string, level byte) {
	f.PwmWrite(pin, level)
}

func (f *FirmataAdaptor) digitalPin(pin int) int {
	return pin + 14
}

func (f *FirmataAdaptor) I2cStart(address byte) {
	f.i2cAddress = address
	f.Board.i2cConfig([]byte{0})
}

func (f *FirmataAdaptor) I2cRead(size uint) []byte {
	f.Board.i2cReadRequest(f.i2cAddress, size)

	events := f.Board.findEvents("i2c_reply")
	if len(events) > 0 {
		return events[len(events)-1].I2cReply["data"]
	}
	return make([]byte, 0)
}

func (f *FirmataAdaptor) I2cWrite(data []byte) {
	f.Board.i2cWriteRequest(f.i2cAddress, data)
}
