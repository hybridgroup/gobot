package minidrone

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"

	"gobot.io/x/gobot/platforms/ble"
)

var _ gobot.Driver = (*Driver)(nil)

func initTestMinidroneDriver() *Driver {
	d := NewDriver(ble.NewClientAdaptor("D7:99:5A:26:EC:38"))
	return d
}

func TestMinidroneDriver(t *testing.T) {
	d := initTestMinidroneDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Minidrone"), true)
}
