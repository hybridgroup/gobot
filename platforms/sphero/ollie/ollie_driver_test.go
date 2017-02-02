package ollie

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"

	"gobot.io/x/gobot/platforms/ble"
)

var _ gobot.Driver = (*Driver)(nil)

func initTestOllieDriver() *Driver {
	d := NewDriver(ble.NewClientAdaptor("D7:99:5A:26:EC:38"))
	return d
}

func TestOllieDriver(t *testing.T) {
	d := initTestOllieDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Ollie"), true)
}
