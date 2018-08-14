package i2c

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/gobottest"
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

func initTestGrovePiDriver() (driver *GrovePiDriver) {
	driver, _ = initGrovePiDriverWithStubbedAdaptor()
	return
}

func initGrovePiDriverWithStubbedAdaptor() (*GrovePiDriver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewGrovePiDriver(adaptor), adaptor
}

func TestGrovePiDriverName(t *testing.T) {
	g := initTestGrovePiDriver()
	gobottest.Refute(t, g.Connection(), nil)
	gobottest.Assert(t, strings.HasPrefix(g.Name(), "GrovePi"), true)
}

func TestGrovePiDriverOptions(t *testing.T) {
	g := NewGrovePiDriver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, g.GetBusOrDefault(1), 2)
}

// Methods
func TestGrovePiDriverStart(t *testing.T) {
	g := initTestGrovePiDriver()

	gobottest.Assert(t, g.Start(), nil)
}

func TestGrovePiDrivergetPin(t *testing.T) {
	gobottest.Assert(t, getPin("a1"), "1")
	gobottest.Assert(t, getPin("A16"), "16")
	gobottest.Assert(t, getPin("D3"), "3")
	gobottest.Assert(t, getPin("d22"), "22")
	gobottest.Assert(t, getPin("22"), "22")
}
