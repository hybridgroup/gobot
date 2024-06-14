package bleclient

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var (
	_ gobot.Adaptor      = (*Adaptor)(nil)
	_ gobot.BLEConnector = (*Adaptor)(nil)
)

func TestNewAdaptor(t *testing.T) {
	a := NewAdaptor("D7:99:5A:26:EC:38")
	assert.Equal(t, "D7:99:5A:26:EC:38", a.Address())
	assert.True(t, strings.HasPrefix(a.Name(), "BLEClient"))
}

func TestName(t *testing.T) {
	a := NewAdaptor("D7:99:5A:26:EC:38")
	a.SetName("awesome")
	assert.Equal(t, "awesome", a.Name())
}

func TestConnect(t *testing.T) {
	const (
		scanTimeout   = 5 * time.Millisecond
		deviceName    = "hello"
		deviceAddress = "11:22:44:AA:BB:CC"
		rssi          = 56
	)
	tests := map[string]struct {
		identifier  string
		extAdapter  *btTestAdapter
		extDevice   btTestDevice
		wantAddress string
		wantName    string
		wantErr     string
	}{
		"connect_by_address": {
			identifier: deviceAddress,
			extAdapter: &btTestAdapter{
				deviceAddress: deviceAddress,
				rssi:          rssi,
				payload:       &btTestPayload{name: deviceName},
			},
			extDevice:   btTestDevice{},
			wantAddress: deviceAddress,
			wantName:    deviceName,
		},
		"connect_by_name": {
			identifier: deviceName,
			extAdapter: &btTestAdapter{
				deviceAddress: deviceAddress,
				rssi:          rssi,
				payload:       &btTestPayload{name: deviceName},
			},
			extDevice:   btTestDevice{},
			wantAddress: deviceAddress,
			wantName:    deviceName,
		},
		"error_enable": {
			extAdapter: &btTestAdapter{
				simulateEnableErr: true,
			},
			wantName: "BLEClient",
			wantErr:  "can't get adapter default: adapter enable error",
		},
		"error_scan": {
			extAdapter: &btTestAdapter{
				simulateScanErr: true,
			},
			wantName: "BLEClient",
			wantErr:  "scan error",
		},
		"error_stop_scan": {
			extAdapter: &btTestAdapter{
				deviceAddress:       deviceAddress,
				payload:             &btTestPayload{},
				simulateStopScanErr: true,
			},
			wantName: "BLEClient",
			wantErr:  "stop scan error",
		},
		"error_timeout_long_delay": {
			extAdapter: &btTestAdapter{
				deviceAddress: deviceAddress,
				payload:       &btTestPayload{},
				scanDelay:     2 * scanTimeout,
			},
			wantName: "BLEClient",
			wantErr:  "scan timeout (5ms) elapsed",
		},
		"error_timeout_bad_identifier": {
			identifier: "bad_identifier",
			extAdapter: &btTestAdapter{
				deviceAddress: deviceAddress,
				payload:       &btTestPayload{},
			},
			wantAddress: "bad_identifier",
			wantName:    "BLEClient",
			wantErr:     "scan timeout (5ms) elapsed",
		},
		"error_connect": {
			extAdapter: &btTestAdapter{
				deviceAddress:      deviceAddress,
				payload:            &btTestPayload{},
				simulateConnectErr: true,
			},
			wantName: "BLEClient",
			wantErr:  "adapter connect error",
		},
		"error_discovery_services": {
			identifier: "disco_err",
			extAdapter: &btTestAdapter{
				deviceAddress: deviceAddress,
				payload:       &btTestPayload{name: "disco_err"},
			},
			extDevice: btTestDevice{
				simulateDiscoverServicesErr: true,
			},
			wantAddress: deviceAddress,
			wantName:    "disco_err",
			wantErr:     "device discover services error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor(tc.identifier)
			btdc := func(_ bluetoothExtDevicer, address, name string) *btDevice {
				return &btDevice{extDevice: tc.extDevice, devAddress: address, devName: name}
			}
			btac := func(bluetoothExtAdapterer, bool) *btAdapter {
				return &btAdapter{extAdapter: tc.extAdapter, btDeviceCreator: btdc}
			}
			a.btAdptCreator = btac
			a.cfg.scanTimeout = scanTimeout // to speed up test
			// act
			err := a.Connect()
			// assert
			if tc.wantErr == "" {
				require.NoError(t, err)
				assert.Equal(t, tc.wantName, a.Name())
				assert.Equal(t, tc.wantAddress, a.Address())
				assert.Equal(t, rssi, a.RSSI())
				assert.True(t, a.connected)
			} else {
				require.ErrorContains(t, err, tc.wantErr)
				assert.Contains(t, a.Name(), tc.wantName)
				assert.Equal(t, tc.wantAddress, a.Address())
				assert.False(t, a.connected)
			}
		})
	}
}

func TestReconnect(t *testing.T) {
	const (
		scanTimeout   = 5 * time.Millisecond
		deviceName    = "hello"
		deviceAddress = "11:22:44:AA:BB:CC"
		rssi          = 56
	)
	tests := map[string]struct {
		extAdapter   *btTestAdapter
		extDevice    *btTestDevice
		wasConnected bool
		wantErr      string
	}{
		"reconnect_not_connected": {
			extAdapter: &btTestAdapter{
				deviceAddress: deviceAddress,
				rssi:          rssi,
				payload:       &btTestPayload{name: deviceName},
			},
			extDevice: &btTestDevice{},
		},
		"reconnect_was_connected": {
			extAdapter: &btTestAdapter{
				deviceAddress: deviceAddress,
				rssi:          rssi,
				payload:       &btTestPayload{name: deviceName},
			},
			extDevice:    &btTestDevice{},
			wasConnected: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor(deviceAddress)
			btdc := func(_ bluetoothExtDevicer, address, name string) *btDevice {
				return &btDevice{extDevice: tc.extDevice, devAddress: address, devName: name}
			}
			a.btAdpt = &btAdapter{extAdapter: tc.extAdapter, btDeviceCreator: btdc}
			a.cfg.scanTimeout = scanTimeout // to speed up test in case of errors
			a.cfg.sleepAfterDisconnect = 0  // to speed up test
			if tc.wasConnected {
				a.btDevice = btdc(nil, "", "")
				a.connected = tc.wasConnected
			}
			// act
			err := a.Reconnect()
			// assert
			if tc.wantErr == "" {
				require.NoError(t, err)
				assert.Equal(t, rssi, a.RSSI())
			} else {
				require.ErrorContains(t, err, tc.wantErr)
			}
			assert.True(t, a.connected)
		})
	}
}

func TestFinalize(t *testing.T) {
	// this also tests Disconnect()
	tests := map[string]struct {
		extDevice *btTestDevice
		wantErr   string
	}{
		"disconnect": {
			extDevice: &btTestDevice{},
		},
		"error_disconnect": {
			extDevice: &btTestDevice{
				simulateDisconnectErr: true,
			},
			wantErr: "device disconnect error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor("")
			a.cfg.sleepAfterDisconnect = 0 // to speed up test
			a.btDevice = &btDevice{extDevice: tc.extDevice}
			// act
			err := a.Finalize()
			// assert
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.wantErr)
			}
			assert.False(t, a.connected)
		})
	}
}

func TestReadCharacteristic(t *testing.T) {
	const uuid = "00001234-0000-1000-8000-00805f9b34fb"
	tests := map[string]struct {
		inUUID       string
		chara        *btTestChara
		notConnected bool
		want         []byte
		wantErr      string
	}{
		"read_ok": {
			inUUID: uuid,
			chara:  &btTestChara{readData: []byte{1, 2, 3}},
			want:   []byte{1, 2, 3},
		},
		"error_not_connected": {
			notConnected: true,
			wantErr:      "cannot read from BLE device until connected",
		},
		"error_bad_chara": {
			inUUID:  "gag1",
			wantErr: "'gag1' is not a valid 16-bit Bluetooth UUID",
		},
		"error_unknown_chara": {
			inUUID:  uuid,
			wantErr: "unknown characteristic: 00001234-0000-1000-8000-00805f9b34fb",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor("")
			if tc.chara != nil {
				a.characteristics[uuid] = tc.chara
			}
			a.connected = !tc.notConnected
			// act
			got, err := a.ReadCharacteristic(tc.inUUID)
			// assert
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.wantErr)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestWriteCharacteristic(t *testing.T) {
	const uuid = "00004321-0000-1000-8000-00805f9b34fb"
	tests := map[string]struct {
		inUUID       string
		inData       []byte
		notConnected bool
		chara        *btTestChara
		want         []byte
		wantErr      string
	}{
		"write_ok": {
			inUUID: uuid,
			inData: []byte{3, 2, 1},
			chara:  &btTestChara{},
			want:   []byte{3, 2, 1},
		},
		"error_not_connected": {
			notConnected: true,
			wantErr:      "cannot write to BLE device until connected",
		},
		"error_bad_chara": {
			inUUID:  "gag2",
			wantErr: "'gag2' is not a valid 16-bit Bluetooth UUID",
		},
		"error_unknown_chara": {
			inUUID:  uuid,
			wantErr: "unknown characteristic: 00004321-0000-1000-8000-00805f9b34fb",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor("")
			if tc.chara != nil {
				a.characteristics[uuid] = tc.chara
			}
			a.connected = !tc.notConnected
			// act
			err := a.WriteCharacteristic(tc.inUUID, tc.inData)
			// assert
			if tc.wantErr == "" {
				require.NoError(t, err)
				assert.Equal(t, tc.want, tc.chara.writtenData)
			} else {
				require.ErrorContains(t, err, tc.wantErr)
			}
		})
	}
}

func TestSubscribe(t *testing.T) {
	const uuid = "00004321-0000-1000-8000-00805f9b34fb"
	tests := map[string]struct {
		inUUID       string
		notConnected bool
		chara        *btTestChara
		want         []byte
		wantErr      string
	}{
		"subscribe_ok": {
			inUUID: uuid,
			chara:  &btTestChara{},
			want:   []byte{3, 4, 5},
		},
		"error_not_connected": {
			notConnected: true,
			wantErr:      "cannot subscribe to BLE device until connected",
		},
		"error_bad_chara": {
			inUUID:  "gag2",
			wantErr: "'gag2' is not a valid 16-bit Bluetooth UUID",
		},
		"error_unknown_chara": {
			inUUID:  uuid,
			wantErr: "unknown characteristic: 00004321-0000-1000-8000-00805f9b34fb",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor("")
			if tc.chara != nil {
				a.characteristics[uuid] = tc.chara
			}
			a.connected = !tc.notConnected
			var got []byte
			notificationFunc := func(data []byte) {
				got = append(got, data...)
			}
			// act
			err := a.Subscribe(tc.inUUID, notificationFunc)
			// assert
			if tc.wantErr == "" {
				require.NoError(t, err)
				tc.chara.notificationFunc([]byte{3, 4, 5})
			} else {
				require.ErrorContains(t, err, tc.wantErr)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
