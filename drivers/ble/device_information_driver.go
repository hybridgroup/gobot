package ble

import (
	"bytes"

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
func NewDeviceInformationDriver(a gobot.BLEConnector, opts ...OptionApplier) *DeviceInformationDriver {
	n := &DeviceInformationDriver{
		Driver:  NewDriver(a, "DeviceInformation", nil, nil, opts...),
		Eventer: gobot.NewEventer(),
	}

	return n
}

// GetModelNumber returns the model number for the BLE Peripheral
func (d *DeviceInformationDriver) GetModelNumber() (string, error) {
	c, err := d.Adaptor().ReadCharacteristic(deviceInformationModelNumberCharaShort)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(c)
	model := buf.String()
	return model, nil
}

// GetFirmwareRevision returns the firmware revision for the BLE Peripheral
func (d *DeviceInformationDriver) GetFirmwareRevision() (string, error) {
	c, err := d.Adaptor().ReadCharacteristic(deviceInformationFirmwareRevisionCharaShort)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val, nil
}

// GetHardwareRevision returns the hardware revision for the BLE Peripheral
func (d *DeviceInformationDriver) GetHardwareRevision() (string, error) {
	c, err := d.Adaptor().ReadCharacteristic(deviceInformationHardwareRevisionCharaShort)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val, nil
}

// GetManufacturerName returns the manufacturer name for the BLE Peripheral
func (d *DeviceInformationDriver) GetManufacturerName() (string, error) {
	c, err := d.Adaptor().ReadCharacteristic(deviceInformationManufacturerNameCharaShort)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val, nil
}

// GetPnPId returns the PnP ID for the BLE Peripheral
func (d *DeviceInformationDriver) GetPnPId() (string, error) {
	c, err := d.Adaptor().ReadCharacteristic(deviceInformationPnPIdCharaShort)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(c)
	val := buf.String()
	return val, nil
}
