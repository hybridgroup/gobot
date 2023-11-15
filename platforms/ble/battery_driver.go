package ble

import (
	"bytes"
	"log"

	"gobot.io/x/gobot/v2"
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
	//nolint:forcetypeassert // ok here
	return b.Connection().(BLEConnector)
}

// Start tells driver to get ready to do work
func (b *BatteryDriver) Start() error { return nil }

// Halt stops battery driver (void)
func (b *BatteryDriver) Halt() error { return nil }

// GetBatteryLevel reads and returns the current battery level
func (b *BatteryDriver) GetBatteryLevel() uint8 {
	c, err := b.adaptor().ReadCharacteristic("2a19")
	if err != nil {
		log.Println(err)
		return 0
	}
	buf := bytes.NewBuffer(c)
	val, _ := buf.ReadByte()
	level := val
	return level
}
