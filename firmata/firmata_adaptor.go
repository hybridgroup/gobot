package gobotFirmata

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/tarm/goserial"
	"strconv"
)

type FirmataAdaptor struct {
	gobot.Adaptor
	Board      *board
	i2cAddress byte
}

var connect = func(fa *FirmataAdaptor) {
	sp, err := serial.OpenPort(&serial.Config{Name: fa.Port, Baud: 57600})
	if err != nil {
		panic(err)
	}
	fa.Board = newBoard(sp)
}

func (fa *FirmataAdaptor) Connect() bool {
	connect(fa)
	fa.Board.connect()
	fa.Connected = true
	return true
}

func (da *FirmataAdaptor) Reconnect() bool { return true }
func (da *FirmataAdaptor) Disconnect() bool {
	err := da.Board.Serial.Close()
	if err != nil {
		fmt.Println(err)
	}
	return true
}
func (da *FirmataAdaptor) Finalize() bool { return da.Disconnect() }

func (da *FirmataAdaptor) InitServo() {}
func (da *FirmataAdaptor) ServoWrite(pin string, angle byte) {
	p, _ := strconv.Atoi(pin)

	da.Board.setPinMode(byte(p), SERVO)
	da.Board.analogWrite(byte(p), angle)
}

func (da *FirmataAdaptor) PwmWrite(pin string, level byte) {
	p, _ := strconv.Atoi(pin)

	da.Board.setPinMode(byte(p), PWM)
	da.Board.analogWrite(byte(p), level)
}

func (da *FirmataAdaptor) DigitalWrite(pin string, level byte) {
	p, _ := strconv.Atoi(pin)

	da.Board.setPinMode(byte(p), OUTPUT)
	da.Board.digitalWrite(byte(p), level)
}

func (da *FirmataAdaptor) DigitalRead(pin string) int {
	p, _ := strconv.Atoi(pin)
	da.Board.setPinMode(byte(p), INPUT)
	da.Board.togglePinReporting(byte(p), HIGH, REPORT_DIGITAL)
	events := da.Board.findEvents(fmt.Sprintf("digital_read_%v", pin))
	if len(events) > 0 {
		return int(events[len(events)-1].Data[0])
	}
	return -1
}

// NOTE pins are numbered A0-A5, which translate to digital pins 14-19
func (da *FirmataAdaptor) AnalogRead(pin string) int {
	p, _ := strconv.Atoi(pin)
	p = da.digitalPin(p)
	da.Board.setPinMode(byte(p), ANALOG)
	da.Board.togglePinReporting(byte(p), HIGH, REPORT_ANALOG)
	events := da.Board.findEvents(fmt.Sprintf("analog_read_%v", pin))
	if len(events) > 0 {
		event := events[len(events)-1]
		return int(uint(event.Data[0])<<24 | uint(event.Data[1])<<16 | uint(event.Data[2])<<8 | uint(event.Data[3]))
	}
	return -1
}

func (da *FirmataAdaptor) digitalPin(pin int) int {
	return pin + 14
}

func (fa *FirmataAdaptor) I2cStart(address byte) {
	fa.i2cAddress = address
	fa.Board.i2cConfig([]uint16{0})
}

func (fa *FirmataAdaptor) I2cRead(size uint16) []uint16 {
	fa.Board.i2cReadRequest(fa.i2cAddress, size)

	events := fa.Board.findEvents("i2c_reply")
	if len(events) > 0 {
		return events[len(events)-1].I2cReply["data"]
	}
	return make([]uint16, 0)
}

func (fa *FirmataAdaptor) I2cWrite(data []uint16) {
	fa.Board.i2cWriteRequest(fa.i2cAddress, data)
}
