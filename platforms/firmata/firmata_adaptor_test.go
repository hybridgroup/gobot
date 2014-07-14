package firmata

import (
	"fmt"
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
)

func initTestFirmataAdaptor() *FirmataAdaptor {
	a := NewFirmataAdaptor("board", "/dev/null")
	a.connect = func(f *FirmataAdaptor) {
		f.board = newBoard(gobot.NullReadWriteCloser{})
		f.board.initTimeInterval = 0 * time.Second
		f.board.pins = make([]pin, 100)
		f.board.events["digital_read_1"] = gobot.NewEvent()
		f.board.events["analog_read_1"] = gobot.NewEvent()
		gobot.Publish(f.board.events["firmware_query"], nil)
		gobot.Publish(f.board.events["capability_query"], nil)
		gobot.Publish(f.board.events["analog_mapping_query"], nil)
	}
	a.Connect()
	return a
}

func TestFirmataAdaptorFinalize(t *testing.T) {
	a := initTestFirmataAdaptor()
	gobot.Expect(t, a.Finalize(), true)
}
func TestFirmataAdaptorConnect(t *testing.T) {
	a := initTestFirmataAdaptor()
	gobot.Expect(t, a.Connect(), true)
}

func TestFirmataAdaptorInitServo(t *testing.T) {
	a := initTestFirmataAdaptor()
	a.InitServo()
}

func TestFirmataAdaptorServoWrite(t *testing.T) {
	a := initTestFirmataAdaptor()
	a.ServoWrite("1", 50)
}

func TestFirmataAdaptorPwmWrite(t *testing.T) {
	a := initTestFirmataAdaptor()
	a.PwmWrite("1", 50)
}

func TestFirmataAdaptorDigitalWrite(t *testing.T) {
	a := initTestFirmataAdaptor()
	a.DigitalWrite("1", 1)
}

func TestFirmataAdaptorDigitalRead(t *testing.T) {
	a := initTestFirmataAdaptor()
	// -1 on no data
	gobot.Expect(t, a.DigitalRead("1"), -1)

	pinNumber := "1"
	gobot.Publish(a.board.events[fmt.Sprintf("digital_read_%v", pinNumber)],
		[]byte{0x01})
	gobot.Expect(t, a.DigitalRead(pinNumber), 0x01)
}

func TestFirmataAdaptorAnalogRead(t *testing.T) {
	a := initTestFirmataAdaptor()
	// -1 on no data
	gobot.Expect(t, a.AnalogRead("1"), -1)

	pinNumber := "1"
	value := 133
	gobot.Publish(a.board.events[fmt.Sprintf("analog_read_%v", pinNumber)],
		[]byte{
			byte(value >> 24),
			byte(value >> 16),
			byte(value >> 8),
			byte(value & 0xff),
		},
	)
	gobot.Expect(t, a.AnalogRead(pinNumber), 133)
}
func TestFirmataAdaptorAnalogWrite(t *testing.T) {
	a := initTestFirmataAdaptor()
	a.AnalogWrite("1", 50)
}
func TestFirmataAdaptorI2cStart(t *testing.T) {
	a := initTestFirmataAdaptor()
	a.I2cStart(0x00)
}
func TestFirmataAdaptorI2cRead(t *testing.T) {
	a := initTestFirmataAdaptor()
	// [] on no data
	gobot.Expect(t, a.I2cRead(1), []byte{})

	i := []byte{100}
	i2cReply := map[string][]byte{}
	i2cReply["data"] = i
	gobot.Publish(a.board.events["i2c_reply"], i2cReply)
	gobot.Expect(t, a.I2cRead(1), i)
}
func TestFirmataAdaptorI2cWrite(t *testing.T) {
	a := initTestFirmataAdaptor()
	a.I2cWrite([]byte{0x00, 0x01})
}
