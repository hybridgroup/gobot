package bleclient

import (
	"fmt"
	"time"

	"tinygo.org/x/bluetooth"
)

// bluetoothExtDevicer is the interface usually implemented by bluetooth.Device
type bluetoothExtDevicer interface {
	DiscoverServices(uuids []bluetooth.UUID) ([]bluetooth.DeviceService, error)
	Disconnect() error
}

// bluetoothExtAdapterer is the interface usually implemented by bluetooth.Adapter
type bluetoothExtAdapterer interface {
	Enable() error
	Scan(callback func(*bluetooth.Adapter, bluetooth.ScanResult)) error
	StopScan() error
	Connect(address bluetooth.Address, params bluetooth.ConnectionParams) (bluetooth.Device, error)
}

type bluetoothExtCharacteristicer interface {
	Read(data []byte) (int, error)
	WriteWithoutResponse(p []byte) (n int, err error)
	EnableNotifications(callback func(buf []byte)) error
}

// btAdptCreatorFunc is just a convenience type, used in the BLE client to ensure testability
type btAdptCreatorFunc func(bluetoothExtAdapterer, bool) *btAdapter

// btAdapter is the wrapper for an external adapter implementation
type btAdapter struct {
	extAdapter      bluetoothExtAdapterer
	btDeviceCreator func(bluetoothExtDevicer, string, string) *btDevice
	debug           bool
}

// newBtAdapter creates a new wrapper around the given external implementation
func newBtAdapter(a bluetoothExtAdapterer, debug bool) *btAdapter {
	bta := btAdapter{
		extAdapter:      a,
		btDeviceCreator: newBtDevice,
		debug:           debug,
	}

	return &bta
}

// Enable configures the BLE stack. It must be called before any Bluetooth-related calls (unless otherwise indicated).
// It pass through the function of the external implementation.
func (bta *btAdapter) enable() error {
	return bta.extAdapter.Enable()
}

// StopScan stops any in-progress scan. It can be called from within a Scan callback to stop the current scan.
// If no scan is in progress, an error will be returned.
func (bta *btAdapter) stopScan() error {
	return bta.extAdapter.StopScan()
}

// Connect starts a connection attempt to the given peripheral device address.
//
// On Linux and Windows, the IsRandom part of the address is ignored.
func (bta *btAdapter) connect(address bluetooth.Address, devName string) (*btDevice, error) {
	extDev, err := bta.extAdapter.Connect(address, bluetooth.ConnectionParams{})
	if err != nil {
		return nil, err
	}

	return bta.btDeviceCreator(extDev, address.String(), devName), nil
}

// Scan starts a BLE scan for the given identifier (address or name).
func (bta *btAdapter) scan(identifier string, scanTimeout time.Duration) (*bluetooth.ScanResult, error) {
	resultChan := make(chan bluetooth.ScanResult, 1)
	errChan := make(chan error)

	go func() {
		callback := func(_ *bluetooth.Adapter, result bluetooth.ScanResult) {
			if bta.debug {
				fmt.Printf("[scan result]: address: '%s', rssi: %d, name: '%s', manufacturer: %v\n",
					result.Address, result.RSSI, result.LocalName(), result.ManufacturerData())
			}
			if result.Address.String() == identifier || result.LocalName() == identifier {
				resultChan <- result
			}
		}
		err := bta.extAdapter.Scan(callback)
		if err != nil {
			errChan <- err
		}
	}()

	select {
	case result := <-resultChan:
		if err := bta.stopScan(); err != nil {
			return nil, err
		}

		return &result, nil
	case err := <-errChan:
		return nil, err
	case <-time.After(scanTimeout):
		_ = bta.stopScan()
		return nil, fmt.Errorf("scan timeout (%s) elapsed", scanTimeout)
	}
}

// btDevice is the wrapper for an external device implementation
type btDevice struct {
	extDevice  bluetoothExtDevicer
	devAddress string
	devName    string
}

// newBtDevice creates a new wrapper around the given external implementation
func newBtDevice(d bluetoothExtDevicer, address, name string) *btDevice {
	return &btDevice{extDevice: d, devAddress: address, devName: name}
}

func (btd *btDevice) name() string { return btd.devName }

func (btd *btDevice) address() string { return btd.devAddress }

func (btd *btDevice) discoverServices(uuids []bluetooth.UUID) ([]bluetooth.DeviceService, error) {
	return btd.extDevice.DiscoverServices(uuids)
}

// Disconnect from the BLE device. This method is non-blocking and does not wait until the connection is fully gone.
func (btd *btDevice) disconnect() error {
	return btd.extDevice.Disconnect()
}

func readFromCharacteristic(chara bluetoothExtCharacteristicer) ([]byte, error) {
	buf := make([]byte, 255)
	n, err := chara.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func writeToCharacteristicWithoutResponse(chara bluetoothExtCharacteristicer, data []byte) error {
	if _, err := chara.WriteWithoutResponse(data); err != nil {
		return err
	}
	return nil
}

func enableNotificationsForCharacteristic(chara bluetoothExtCharacteristicer, f func(data []byte)) error {
	return chara.EnableNotifications(f)
}
