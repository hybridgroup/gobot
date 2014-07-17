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
		// arduino uno r3 firmware response "StandardFirmata.ino"
		f.board.process([]byte{240, 121, 2, 3, 83, 0, 116, 0, 97, 0, 110, 0, 100,
			0, 97, 0, 114, 0, 100, 0, 70, 0, 105, 0, 114, 0, 109, 0, 97, 0, 116, 0,
			97, 0, 46, 0, 105, 0, 110, 0, 111, 0, 247})
		// arduino uno r3 capabilities response
		f.board.process([]byte{240, 108, 127, 127, 0, 1, 1, 1, 4, 14, 127, 0, 1,
			1, 1, 3, 8, 4, 14, 127, 0, 1, 1, 1, 4, 14, 127, 0, 1, 1, 1, 3, 8, 4, 14,
			127, 0, 1, 1, 1, 3, 8, 4, 14, 127, 0, 1, 1, 1, 4, 14, 127, 0, 1, 1, 1,
			4, 14, 127, 0, 1, 1, 1, 3, 8, 4, 14, 127, 0, 1, 1, 1, 3, 8, 4, 14, 127,
			0, 1, 1, 1, 3, 8, 4, 14, 127, 0, 1, 1, 1, 4, 14, 127, 0, 1, 1, 1, 4, 14,
			127, 0, 1, 1, 1, 2, 10, 127, 0, 1, 1, 1, 2, 10, 127, 0, 1, 1, 1, 2, 10,
			127, 0, 1, 1, 1, 2, 10, 127, 0, 1, 1, 1, 2, 10, 6, 1, 127, 0, 1, 1, 1,
			2, 10, 6, 1, 127, 247})
		// arduino uno r3 analog mapping response
		f.board.process([]byte{240, 106, 127, 127, 127, 127, 127, 127, 127, 127,
			127, 127, 127, 127, 127, 127, 0, 1, 2, 3, 4, 5, 247})
	}
	a.Connect()
	return a
}

func TestFirmataAdaptorFinalize(t *testing.T) {
	a := initTestFirmataAdaptor()
	gobot.Assert(t, a.Finalize(), true)
}
func TestFirmataAdaptorConnect(t *testing.T) {
	a := initTestFirmataAdaptor()
	gobot.Assert(t, a.Connect(), true)
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
	pinNumber := "1"
	// -1 on no data
	gobot.Assert(t, a.DigitalRead(pinNumber), -1)

	go func() {
		<-time.After(5 * time.Millisecond)
		gobot.Publish(a.board.events[fmt.Sprintf("digital_read_%v", pinNumber)],
			[]byte{0x01})
	}()
	gobot.Assert(t, a.DigitalRead(pinNumber), 0x01)
}

func TestFirmataAdaptorAnalogRead(t *testing.T) {
	a := initTestFirmataAdaptor()
	pinNumber := "1"
	// -1 on no data
	gobot.Assert(t, a.AnalogRead(pinNumber), -1)

	value := 133
	go func() {
		<-time.After(5 * time.Millisecond)
		gobot.Publish(a.board.events[fmt.Sprintf("analog_read_%v", pinNumber)],
			[]byte{
				byte(value >> 24),
				byte(value >> 16),
				byte(value >> 8),
				byte(value & 0xff),
			},
		)
	}()
	gobot.Assert(t, a.AnalogRead(pinNumber), 133)
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
	gobot.Assert(t, a.I2cRead(1), []byte{})

	i := []byte{100}
	i2cReply := map[string][]byte{}
	i2cReply["data"] = i
	go func() {
		<-time.After(5 * time.Millisecond)
		gobot.Publish(a.board.events["i2c_reply"], i2cReply)
	}()
	gobot.Assert(t, a.I2cRead(1), i)
}
func TestFirmataAdaptorI2cWrite(t *testing.T) {
	a := initTestFirmataAdaptor()
	a.I2cWrite([]byte{0x00, 0x01})
}
