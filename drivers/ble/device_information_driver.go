package ble

import (
	"bytes"
	"log"

	"gobot.io/x/gobot/v2"
)

const (
	deviceInformationModelNumberCharaShort      = "2a24"
	deviceInformationFirmwareRevisionCharaShort = "2a26"
	deviceInformationHardwareRevisionCharaShort = "2a27"
	deviceInformationManufacturerNameCharaShort = "2a29"
	deviceInformationPnPIdCharaShort            = "2a50"
)

// DeviceInformationDriver represents the device information service for a BLE peripheral
type DeviceInformationDriver struct {
	*Driver
	gobot.Eventer
}

// NewDeviceInformationDriver creates a new driver
func NewDeviceInformationDriver(a gobot.BLEConnector) *DeviceInformationDriver {
	n := &DeviceInformationDriver{
		Driver:  NewDriver(a, "DeviceInformation", nil, nil),
		Eventer: gobot.NewEventer(),
	}

	return n
}

// GetModelNumber returns the model number for the BLE Peripheral
func (d *DeviceInformationDriver) GetModelNumber() string {
	c, err := d.Adaptor().ReadCharacteristic(deviceInformationModelNumberCharaShort)
	if err != nil {
		log.Println(err)
		return ""
	}
	buf := bytes.NewBuffer(c)
	model := buf.String()
	return model
}

// GetFirmwareRevision returns the firmware revision for the BLE Peripheral
func (d *DeviceInformationDriver) GetFirmwareRevision() string {
	c, err := d.Adaptor().ReadCharacteristic(deviceInformationFirmwareRevisionCharaShort)
	if err != nil {
		log.Println(err)
		return ""
	}
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val
}

// GetHardwareRevision returns the hardware revision for the BLE Peripheral
func (d *DeviceInformationDriver) GetHardwareRevision() string {
	c, err := d.Adaptor().ReadCharacteristic(deviceInformationHardwareRevisionCharaShort)
	if err != nil {
		log.Println(err)
		return ""
	}
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val
}

// GetManufacturerName returns the manufacturer name for the BLE Peripheral
func (d *DeviceInformationDriver) GetManufacturerName() string {
	c, err := d.Adaptor().ReadCharacteristic(deviceInformationManufacturerNameCharaShort)
	if err != nil {
		log.Println(err)
		return ""
	}
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val
}

// GetPnPId returns the PnP ID for the BLE Peripheral
func (d *DeviceInformationDriver) GetPnPId() string {
	c, err := d.Adaptor().ReadCharacteristic(deviceInformationPnPIdCharaShort)
	if err != nil {
		log.Println(err)
		return ""
	}
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val
}
