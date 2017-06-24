package gpio

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*NeopixelDriver)(nil)

func initTestNeopixelDriver(conn DigitalWriter) *NeopixelDriver {
	return NewNeopixelDriver(conn, "1", 5)
}

func TestNeopixelDriverDefaultName(t *testing.T) {
	g := initTestNeopixelDriver(newGpioTestAdaptor())
	gobottest.Assert(t, strings.HasPrefix(g.Name(), "Neopixel"), true)
}

func TestNeopixelDriverSetName(t *testing.T) {
	g := initTestNeopixelDriver(newGpioTestAdaptor())
	g.SetName("mybot")
	gobottest.Assert(t, g.Name(), "mybot")
}
