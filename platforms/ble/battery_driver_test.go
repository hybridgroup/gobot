package ble

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*BatteryDriver)(nil)

func initTestBatteryDriver() *BatteryDriver {
	d := NewBatteryDriver(NewBleTestAdaptor())
	return d
}

func TestBatteryDriver(t *testing.T) {
	d := initTestBatteryDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Battery"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestBatteryDriverStartAndHalt(t *testing.T) {
	d := initTestBatteryDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}

func TestBatteryDriverRead(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewBatteryDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte{20}, nil
	})

	gobottest.Assert(t, d.GetBatteryLevel(), uint8(20))
}
