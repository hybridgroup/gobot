package microbit

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"

	"gobot.io/x/gobot/platforms/ble"
)

var _ gobot.Driver = (*ButtonDriver)(nil)

func initTestButtonDriver() *ButtonDriver {
	d := NewButtonDriver(ble.NewClientAdaptor("D7:99:5A:26:EC:38"))
	return d
}

func TestButtonDriver(t *testing.T) {
	d := initTestButtonDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Microbit Button"), true)
}
