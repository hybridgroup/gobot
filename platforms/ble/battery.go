package ble

import (
	"bytes"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*BatteryDriver)(nil)

type BatteryDriver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}

// NewBatteryDriver creates a BatteryDriver
func NewBatteryDriver(a *ClientAdaptor, name string) *BatteryDriver {
	n := &BatteryDriver{
		connection: a,
		Eventer:    gobot.NewEventer(),
	}

	return n
}
func (b *BatteryDriver) Connection() gobot.Connection { return b.connection }
func (b *BatteryDriver) Name() string                 { return b.name }
func (b *BatteryDriver) SetName(n string)             { b.name = n }

// adaptor returns BLE adaptor
func (b *BatteryDriver) adaptor() *ClientAdaptor {
	return b.Connection().(*ClientAdaptor)
}

// Start tells driver to get ready to do work
func (b *BatteryDriver) Start() (errs []error) {
	return
}

// Halt stops battery driver (void)
func (b *BatteryDriver) Halt() (errs []error) { return }

func (b *BatteryDriver) GetBatteryLevel() (level uint8) {
	var l uint8
	c, _ := b.adaptor().ReadCharacteristic("180f", "2a19")
	buf := bytes.NewBuffer(c)
	val, _ := buf.ReadByte()
	l = uint8(val)
	return l
}
