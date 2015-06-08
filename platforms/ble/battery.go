package ble

import (
//	"bytes"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*BLEBatteryDriver)(nil)


type BLEBatteryDriver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}


// NewBLEBatteryDriver creates a BLEBatteryDriver by name
func NewBLEBatteryDriver(a *BLEAdaptor, name string) *BLEBatteryDriver {
	n := &BLEBatteryDriver{
		name:       name,
		connection: a,
		Eventer:    gobot.NewEventer(),
	}

	// n.AddEvent("extended")
	// n.AddEvent("signal")
	// n.AddEvent("attention")
	// n.AddEvent("meditation")
	// n.AddEvent("blink")
	// n.AddEvent("wave")
	// n.AddEvent("eeg")
	// n.AddEvent("error")

	return n
}
func (b *BLEBatteryDriver) Connection() gobot.Connection { return b.connection }
func (b *BLEBatteryDriver) Name() string                 { return b.name }

// adaptor returns BLE adaptor
func (b *BLEBatteryDriver) adaptor() *BLEAdaptor {
	return b.Connection().(*BLEAdaptor)
}

// Start creates a go routine to listen from serial port
// and parse buffer readings
func (b *BLEBatteryDriver) Start() (errs []error) {
	return
}

// Halt stops battery driver (void)
func (b *BLEBatteryDriver) Halt() (errs []error) { return }

