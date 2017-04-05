package ble

import (
	"bytes"

	"gobot.io/x/gobot"
)

// BatteryDriver represents the Battery Service for a BLE Peripheral
type BatteryDriver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}

// NewBatteryDriver creates a BatteryDriver
func NewBatteryDriver(a BLEConnector) *BatteryDriver {
	n := &BatteryDriver{
		name:       gobot.DefaultName("Battery"),
		connection: a,
		Eventer:    gobot.NewEventer(),
	}

	return n
}

// Connection returns the Driver's Connection to the associated Adaptor
func (b *BatteryDriver) Connection() gobot.Connection { return b.connection }

// Name returns the Driver name
func (b *BatteryDriver) Name() string { return b.name }

// SetName sets the Driver name
func (b *BatteryDriver) SetName(n string) { b.name = n }

// adaptor returns BLE adaptor
func (b *BatteryDriver) adaptor() BLEConnector {
	return b.Connection().(BLEConnector)
}

// Start tells driver to get ready to do work
func (b *BatteryDriver) Start() (err error) {
	return
}

// Halt stops battery driver (void)
func (b *BatteryDriver) Halt() (err error) { return }

// GetBatteryLevel reads and returns the current battery level
func (b *BatteryDriver) GetBatteryLevel() (level uint8) {
	var l uint8
	c, _ := b.adaptor().ReadCharacteristic("2a19")
	buf := bytes.NewBuffer(c)
	val, _ := buf.ReadByte()
	l = uint8(val)
	return l
}
