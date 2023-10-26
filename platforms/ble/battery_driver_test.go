package ble

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*BatteryDriver)(nil)

func initTestBatteryDriver() *BatteryDriver {
	d := NewBatteryDriver(NewBleTestAdaptor())
	return d
}

func TestBatteryDriver(t *testing.T) {
	d := initTestBatteryDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "Battery"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestBatteryDriverStartAndHalt(t *testing.T) {
	d := initTestBatteryDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.Halt())
}

func TestBatteryDriverRead(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewBatteryDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte{20}, nil
	})

	assert.Equal(t, uint8(20), d.GetBatteryLevel())
}
