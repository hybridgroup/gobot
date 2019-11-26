package i2c

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/gobottest"
	"strings"
	"testing"
)

var _ gobot.Driver = (*GrovePiDriver)(nil)

// must implement the DigitalReader interface
var _ gpio.DigitalReader = (*GrovePiDriver)(nil)

// must implement the DigitalWriter interface
var _ gpio.DigitalWriter = (*GrovePiDriver)(nil)

// must implement the AnalogReader interface
var _ aio.AnalogReader = (*GrovePiDriver)(nil)

// must implement the Adaptor interface
var _ gobot.Adaptor = (*GrovePiDriver)(nil)

func initTestGrovePiDriver() (driver *GrovePiDriver, adaptor *i2cTestAdaptor) {
	return initGrovePiDriverWithStubbedAdaptor()
}

func initGrovePiDriverWithStubbedAdaptor() (*GrovePiDriver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewGrovePiDriver(adaptor), adaptor
}

func TestGrovePiDriverName(t *testing.T) {
	g, _ := initTestGrovePiDriver()
	gobottest.Refute(t, g.Connection(), nil)
	gobottest.Assert(t, strings.HasPrefix(g.Name(), "GrovePi"), true)
}

func TestGrovePiDriverOptions(t *testing.T) {
	g := NewGrovePiDriver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, g.GetBusOrDefault(1), 2)
}

func TestGrovePiDriver_UltrasonicRead(t *testing.T) {
	g, a := initTestGrovePiDriver()
	g.Start()

	fakePin := byte(1)
    fakeI2cResponse := []byte{CommandReadUltrasonic, 1, 2}

	expectedCommand := []byte{CommandReadUltrasonic, fakePin, 0, 0}
    expectedResult := 257

	resultCommand := make([]byte, 3)

    // capture i2c command
	a.i2cWriteImpl = func(bytes []byte) (i int, e error) {
		resultCommand = bytes
		return len(bytes), nil
	}

	// fake i2c response
	a.i2cReadImpl = func(bytes []byte) (i int, e error) {
		bytes[0] = fakeI2cResponse[0]
		bytes[1] = fakeI2cResponse[1]
		bytes[2] = fakeI2cResponse[2]
		return len(bytes), nil
	}

	result, _ := g.readUltrasonic(fakePin, 10)

	gobottest.Assert(t, resultCommand, expectedCommand)
	gobottest.Assert(t, result, expectedResult)
}

// Methods
func TestGrovePiDriverStart(t *testing.T) {
	g, _ := initTestGrovePiDriver()

	gobottest.Assert(t, g.Start(), nil)
}

func TestGrovePiDrivergetPin(t *testing.T) {
	gobottest.Assert(t, getPin("a1"), "1")
	gobottest.Assert(t, getPin("A16"), "16")
	gobottest.Assert(t, getPin("D3"), "3")
	gobottest.Assert(t, getPin("d22"), "22")
	gobottest.Assert(t, getPin("22"), "22")
}
