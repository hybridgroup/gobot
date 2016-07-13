package ble

import (
	"bytes"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*BLEDeviceInformationDriver)(nil)

type BLEDeviceInformationDriver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}

// NewBLEDeviceInformationDriver creates a BLEDeviceInformationDriver
// by name
func NewBLEDeviceInformationDriver(a *BLEClientAdaptor, name string) *BLEDeviceInformationDriver {
	n := &BLEDeviceInformationDriver{
		name:       name,
		connection: a,
		Eventer:    gobot.NewEventer(),
	}

	return n
}
func (b *BLEDeviceInformationDriver) Connection() gobot.Connection { return b.connection }
func (b *BLEDeviceInformationDriver) Name() string                 { return b.name }

// adaptor returns BLE adaptor for this device
func (b *BLEDeviceInformationDriver) adaptor() *BLEClientAdaptor {
	return b.Connection().(*BLEClientAdaptor)
}

// Start tells driver to get ready to do work
func (b *BLEDeviceInformationDriver) Start() (errs []error) {
	return
}

// Halt stops driver (void)
func (b *BLEDeviceInformationDriver) Halt() (errs []error) { return }

func (b *BLEDeviceInformationDriver) GetModelNumber() (model string) {
	c, _ := b.adaptor().ReadCharacteristic("180a", "2a24")
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val
}

func (b *BLEDeviceInformationDriver) GetFirmwareRevision() (revision string) {
	c, _ := b.adaptor().ReadCharacteristic("180a", "2a26")
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val
}

func (b *BLEDeviceInformationDriver) GetHardwareRevision() (revision string) {
	c, _ := b.adaptor().ReadCharacteristic("180a", "2a27")
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val
}

func (b *BLEDeviceInformationDriver) GetManufacturerName() (manufacturer string) {
	c, _ := b.adaptor().ReadCharacteristic("180a", "2a29")
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val
}

func (b *BLEDeviceInformationDriver) GetPnPId() (model string) {
	c, _ := b.adaptor().ReadCharacteristic("180a", "2a50")
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val
}
