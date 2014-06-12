package firmata

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"testing"
)

var adaptor *FirmataAdaptor

func init() {
	adaptor = NewFirmataAdaptor("board", "/dev/null")
	adaptor.connect = func(f *FirmataAdaptor) {
		f.Board = newBoard(sp{})
		f.Board.Events = append(f.Board.Events, event{Name: "firmware_query"})
		f.Board.Events = append(f.Board.Events, event{Name: "capability_query"})
		f.Board.Events = append(f.Board.Events, event{Name: "analog_mapping_query"})
	}
	adaptor.Connect()
}

func TestFinalize(t *testing.T) {
	gobot.Expect(t, adaptor.Finalize(), true)
}
func TestConnect(t *testing.T) {
	gobot.Expect(t, adaptor.Connect(), true)
}
func TestInitServo(t *testing.T) {
	adaptor.InitServo()
}
func TestServoWrite(t *testing.T) {
	adaptor.ServoWrite("1", 50)
}
func TestPwmWrite(t *testing.T) {
	adaptor.PwmWrite("1", 50)
}
func TestDigitalWrite(t *testing.T) {
	adaptor.DigitalWrite("1", 1)
}
func TestDigitalRead(t *testing.T) {
	// -1 on no data
	gobot.Expect(t, adaptor.DigitalRead("1"), -1)

	pinNumber := "1"
	adaptor.Board.Events = append(adaptor.Board.Events, event{Name: fmt.Sprintf("digital_read_%v", pinNumber), Data: []byte{0x01}})
	gobot.Expect(t, adaptor.DigitalRead(pinNumber), 0x01)
}
func TestAnalogRead(t *testing.T) {
	// -1 on no data
	gobot.Expect(t, adaptor.AnalogRead("1"), -1)

	pinNumber := "1"
	value := 133
	adaptor.Board.Events = append(adaptor.Board.Events, event{Name: fmt.Sprintf("analog_read_%v", pinNumber), Data: []byte{byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value & 0xff)}})
	gobot.Expect(t, adaptor.AnalogRead(pinNumber), 133)
}
func TestI2cStart(t *testing.T) {
	adaptor.I2cStart(0x00)
}
func TestI2cRead(t *testing.T) {
	// [] on no data
	gobot.Expect(t, adaptor.I2cRead(1), []byte{})

	i := []byte{100}
	i2cReply := map[string][]byte{}
	i2cReply["data"] = i
	adaptor.Board.Events = append(adaptor.Board.Events, event{Name: "i2c_reply", I2cReply: i2cReply})
	gobot.Expect(t, adaptor.I2cRead(1), i)
}
func TestI2cWrite(t *testing.T) {
	adaptor.I2cWrite([]byte{0x00, 0x01})
}
