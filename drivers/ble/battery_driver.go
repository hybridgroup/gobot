package ble

import (
	"bytes"

	"gobot.io/x/gobot/v2"
)

const batteryCharaShort = "2a19"

// BatteryDriver represents the battery service for a BLE peripheral
type BatteryDriver struct {
	*Driver
	gobot.Eventer
}

// NewBatteryDriver creates a new driver
func NewBatteryDriver(a gobot.BLEConnector, opts ...OptionApplier) *BatteryDriver {
	d := &BatteryDriver{
		Driver:  NewDriver(a, "Battery", nil, nil, opts...),
		Eventer: gobot.NewEventer(),
	}

	return d
}

// GetBatteryLevel reads and returns the current battery level
func (d *BatteryDriver) GetBatteryLevel() (uint8, error) {
	c, err := d.Adaptor().ReadCharacteristic(batteryCharaShort)
	if err != nil {
		return 0, err
	}
	buf := bytes.NewBuffer(c)
	val, _ := buf.ReadByte()
	level := val
	return level, nil
}
