package ble

import (
	"bytes"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*DeviceInformationDriver)(nil)

type DeviceInformationDriver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}

// NewDeviceInformationDriver creates a DeviceInformationDriver
func NewDeviceInformationDriver(a *ClientAdaptor) *DeviceInformationDriver {
	n := &DeviceInformationDriver{
		name:       "Device Information",
		connection: a,
		Eventer:    gobot.NewEventer(),
	}

	return n
}
func (b *DeviceInformationDriver) Connection() gobot.Connection { return b.connection }
func (b *DeviceInformationDriver) Name() string                 { return b.name }
func (b *DeviceInformationDriver) SetName(n string)             { b.name = n }

// adaptor returns BLE adaptor for this device
func (b *DeviceInformationDriver) adaptor() *ClientAdaptor {
	return b.Connection().(*ClientAdaptor)
}

// Start tells driver to get ready to do work
func (b *DeviceInformationDriver) Start() (errs []error) {
	return
}

// Halt stops driver (void)
func (b *DeviceInformationDriver) Halt() (errs []error) { return }

func (b *DeviceInformationDriver) GetModelNumber() (model string) {
	c, _ := b.adaptor().ReadCharacteristic("180a", "2a24")
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val
}

func (b *DeviceInformationDriver) GetFirmwareRevision() (revision string) {
	c, _ := b.adaptor().ReadCharacteristic("180a", "2a26")
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val
}

func (b *DeviceInformationDriver) GetHardwareRevision() (revision string) {
	c, _ := b.adaptor().ReadCharacteristic("180a", "2a27")
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val
}

func (b *DeviceInformationDriver) GetManufacturerName() (manufacturer string) {
	c, _ := b.adaptor().ReadCharacteristic("180a", "2a29")
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val
}

func (b *DeviceInformationDriver) GetPnPId() (model string) {
	c, _ := b.adaptor().ReadCharacteristic("180a", "2a50")
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val
}
