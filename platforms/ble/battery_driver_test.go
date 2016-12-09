package ble

import (
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*BatteryDriver)(nil)

func initTestBatteryDriver() *BatteryDriver {
	d := NewBatteryDriver(NewClientAdaptor("D7:99:5A:26:EC:38"))
	return d
}

func TestBatteryDriver(t *testing.T) {
	d := initTestBatteryDriver()
	gobottest.Assert(t, d.Name(), "Battery")
}
