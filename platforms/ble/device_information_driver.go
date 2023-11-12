package ble

import (
	"bytes"
	"log"

	"gobot.io/x/gobot/v2"
)

// DeviceInformationDriver represents the Device Information Service for a BLE Peripheral
type DeviceInformationDriver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}

// NewDeviceInformationDriver creates a DeviceInformationDriver
func NewDeviceInformationDriver(a BLEConnector) *DeviceInformationDriver {
	n := &DeviceInformationDriver{
		name:       gobot.DefaultName("DeviceInformation"),
		connection: a,
		Eventer:    gobot.NewEventer(),
	}

	return n
}

// Connection returns the Driver's Connection to the associated Adaptor
func (b *DeviceInformationDriver) Connection() gobot.Connection { return b.connection }

// Name returns the Driver name
func (b *DeviceInformationDriver) Name() string { return b.name }

// SetName sets the Driver name
func (b *DeviceInformationDriver) SetName(n string) { b.name = n }

// adaptor returns BLE adaptor for this device
func (b *DeviceInformationDriver) adaptor() BLEConnector {
	//nolint:forcetypeassert // ok here
	return b.Connection().(BLEConnector)
}

// Start tells driver to get ready to do work
func (b *DeviceInformationDriver) Start() error { return nil }

// Halt stops driver (void)
func (b *DeviceInformationDriver) Halt() error { return nil }

// GetModelNumber returns the model number for the BLE Peripheral
func (b *DeviceInformationDriver) GetModelNumber() string {
	c, err := b.adaptor().ReadCharacteristic("2a24")
	if err != nil {
		log.Println(err)
		return ""
	}
	buf := bytes.NewBuffer(c)
	model := buf.String()
	return model
}

// GetFirmwareRevision returns the firmware revision for the BLE Peripheral
func (b *DeviceInformationDriver) GetFirmwareRevision() string {
	c, err := b.adaptor().ReadCharacteristic("2a26")
	if err != nil {
		log.Println(err)
		return ""
	}
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val
}

// GetHardwareRevision returns the hardware revision for the BLE Peripheral
func (b *DeviceInformationDriver) GetHardwareRevision() string {
	c, err := b.adaptor().ReadCharacteristic("2a27")
	if err != nil {
		log.Println(err)
		return ""
	}
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val
}

// GetManufacturerName returns the manufacturer name for the BLE Peripheral
func (b *DeviceInformationDriver) GetManufacturerName() string {
	c, err := b.adaptor().ReadCharacteristic("2a29")
	if err != nil {
		log.Println(err)
		return ""
	}
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val
}

// GetPnPId returns the PnP ID for the BLE Peripheral
func (b *DeviceInformationDriver) GetPnPId() string {
	c, err := b.adaptor().ReadCharacteristic("2a50")
	if err != nil {
		log.Println(err)
		return ""
	}
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val
}
