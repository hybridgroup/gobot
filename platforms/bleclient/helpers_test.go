package bleclient

import (
	"fmt"
	"time"

	"tinygo.org/x/bluetooth"
)

type btTestAdapter struct {
	deviceAddress       string
	rssi                int16
	scanDelay           time.Duration
	payload             *btTestPayload
	simulateEnableErr   bool
	simulateScanErr     bool
	simulateStopScanErr bool
	simulateConnectErr  bool
}

func (bta *btTestAdapter) Enable() error {
	if bta.simulateEnableErr {
		return fmt.Errorf("adapter enable error")
	}

	return nil
}

func (bta *btTestAdapter) Scan(callback func(*bluetooth.Adapter, bluetooth.ScanResult)) error {
	if bta.simulateScanErr {
		return fmt.Errorf("adapter scan error")
	}

	devAddr, err := bluetooth.ParseMAC(bta.deviceAddress)
	if err != nil {
		// normally this error should not happen in test
		return err
	}
	time.Sleep(bta.scanDelay)

	a := bluetooth.Address{MACAddress: bluetooth.MACAddress{MAC: devAddr}}
	r := bluetooth.ScanResult{Address: a, RSSI: bta.rssi, AdvertisementPayload: bta.payload}
	callback(nil, r)

	return nil
}

func (bta *btTestAdapter) StopScan() error {
	if bta.simulateStopScanErr {
		return fmt.Errorf("adapter stop scan error")
	}

	return nil
}

func (bta *btTestAdapter) Connect(addr bluetooth.Address, _ bluetooth.ConnectionParams) (bluetooth.Device, error) {
	if bta.simulateConnectErr {
		return bluetooth.Device{}, fmt.Errorf("adapter connect error")
	}

	return bluetooth.Device{Address: addr}, nil
}

type btTestPayload struct {
	name string
}

func (ptp *btTestPayload) LocalName() string { return ptp.name }

func (*btTestPayload) HasServiceUUID(bluetooth.UUID) bool { return true }

func (*btTestPayload) Bytes() []byte { return nil }

func (*btTestPayload) ManufacturerData() []bluetooth.ManufacturerDataElement { return nil }

func (*btTestPayload) ServiceData() []bluetooth.ServiceDataElement { return nil }

type btTestDevice struct {
	simulateDiscoverServicesErr bool
	simulateDisconnectErr       bool
}

func (btd btTestDevice) DiscoverServices(_ []bluetooth.UUID) ([]bluetooth.DeviceService, error) {
	if btd.simulateDiscoverServicesErr {
		return nil, fmt.Errorf("device discover services error")
	}

	// for this test we can not return any []bluetooth.DeviceService
	return nil, nil
}

func (btd btTestDevice) Disconnect() error {
	if btd.simulateDisconnectErr {
		return fmt.Errorf("device disconnect error")
	}

	return nil
}

type btTestChara struct {
	readData         []byte
	writtenData      []byte
	notificationFunc func(buf []byte)
}

func (btc *btTestChara) Read(data []byte) (int, error) {
	copy(data, btc.readData)
	return len(btc.readData), nil
}

func (btc *btTestChara) WriteWithoutResponse(data []byte) (int, error) {
	btc.writtenData = append(btc.writtenData, data...)
	return len(data), nil
}

func (btc *btTestChara) EnableNotifications(callback func(buf []byte)) error {
	btc.notificationFunc = callback
	return nil
}
