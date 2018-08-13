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
