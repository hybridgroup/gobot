package bb8

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"

	"gobot.io/x/gobot/platforms/ble"
)

var _ gobot.Driver = (*BB8Driver)(nil)

func initTestBB8Driver() *BB8Driver {
	d := NewDriver(ble.NewClientAdaptor("D7:99:5A:26:EC:38"))
	return d
}

func TestBB8Driver(t *testing.T) {
	d := initTestBB8Driver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "BB8"), true)
}
