package ble

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*BatteryDriver)(nil)

func initTestBatteryDriver() *BatteryDriver {
	d := NewBatteryDriver(newBleTestAdaptor())
	return d
}

func TestBatteryDriver(t *testing.T) {
	d := initTestBatteryDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Battery"), true)
}

func TestBatteryDriverRead(t *testing.T) {
	a := newBleTestAdaptor()
	d := NewBatteryDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte{20}, nil
	})

	gobottest.Assert(t, d.GetBatteryLevel(), uint8(20))
}
