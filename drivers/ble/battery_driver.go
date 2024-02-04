package ble

import (
	"bytes"
	"log"

	"gobot.io/x/gobot/v2"
)

const batteryCharaShort = "2a19"

// BatteryDriver represents the battery service for a BLE peripheral
type BatteryDriver struct {
	*Driver
	gobot.Eventer
}

// NewBatteryDriver creates a new driver
func NewBatteryDriver(a gobot.BLEConnector) *BatteryDriver {
	d := &BatteryDriver{
		Driver:  NewDriver(a, "Battery", nil, nil),
		Eventer: gobot.NewEventer(),
	}

	return d
}

// GetBatteryLevel reads and returns the current battery level
func (d *BatteryDriver) GetBatteryLevel() uint8 {
	c, err := d.Adaptor().ReadCharacteristic(batteryCharaShort)
	if err != nil {
		log.Println(err)
		return 0
	}
	buf := bytes.NewBuffer(c)
	val, _ := buf.ReadByte()
	level := val
	return level
}
