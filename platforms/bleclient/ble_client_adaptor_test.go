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
		identifier                  string
		extAdapter                  *btTestAdapter
		simulateConnected           bool
		simulateDiscoverServicesErr bool
		wantAddress                 string
		wantName                    string
		wantErr                     string
	}{
		"connect_by_address": {
			identifier: deviceAddress,
			extAdapter: &btTestAdapter{
				deviceAddress: deviceAddress,
				rssi:          rssi,
				payload:       &btTestPayload{name: deviceName},
			},
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
			wantAddress: deviceAddress,
			wantName:    deviceName,
		},
		"error_already_connected": {
			extAdapter:        &btTestAdapter{},
			simulateConnected: true,
			wantName:          "BLEClient",
			wantErr:           "is already connected",
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
			simulateDiscoverServicesErr: true,
			wantAddress:                 deviceAddress,
			wantName:                    "disco_err",
			wantErr:                     "device discover services error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor(tc.identifier)
			extDevice := btTestDevice{simulateDiscoverServicesErr: tc.simulateDiscoverServicesErr}
			btdc := func(_ bluetoothExtDevicer, address, name string) *btDevice {
				return &btDevice{extDevice: extDevice, devAddress: address, devName: name}
			}
			btac := func(bluetoothExtAdapterer, bool) *btAdapter {
				return &btAdapter{extAdapter: tc.extAdapter, btDeviceCreator: btdc}
			}
			a.btAdptCreator = btac
			a.cfg.scanTimeout = scanTimeout // to speed up test
			a.connected = tc.simulateConnected
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
				assert.Equal(t, tc.simulateConnected, a.connected)
			}
		})
	}
}

func TestDisconnect(t *testing.T) {
	const (
		scanTimeout   = 5 * time.Millisecond
		deviceName    = "hello"
		deviceAddress = "11:22:44:AA:BB:CC"
	)
	tests := map[string]struct {
		connectBefore         bool
		dropCharacteristics   bool
		simulateDisconnectErr bool
		wantErr               string
	}{
		"disconnect_not_connectected": {},
		"disconnect_connectected_before": {
			connectBefore: true,
		},
		"disconnect_drop_charas": {
			connectBefore:       true,
			dropCharacteristics: true,
		},
		"error_disconnect": {
			simulateDisconnectErr: true,
			connectBefore:         true,
			wantErr:               "device disconnect error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor(deviceName)
			a.cfg.sleepAfterDisconnect = 0 // to speed up test
			a.cfg.dropCharacteristicsOnDisconnect = tc.dropCharacteristics
			var wantCharaLen int
			if tc.connectBefore {
				extDevice := btTestDevice{simulateDisconnectErr: tc.simulateDisconnectErr}
				btdc := func(_ bluetoothExtDevicer, address, name string) *btDevice {
					return &btDevice{extDevice: extDevice, devAddress: address, devName: name}
				}
				extAdapter := &btTestAdapter{
					deviceAddress: deviceAddress,
					payload:       &btTestPayload{name: deviceName},
				}
				btac := func(bluetoothExtAdapterer, bool) *btAdapter {
					return &btAdapter{extAdapter: extAdapter, btDeviceCreator: btdc}
				}
				a.btAdptCreator = btac
				a.cfg.scanTimeout = scanTimeout // to speed up test
				require.NoError(t, a.Connect())
				a.characteristics["charauuid"] = &btTestChara{}
				if !tc.dropCharacteristics {
					wantCharaLen = 1
				}
			}
			// act
			err := a.Disconnect()
			// assert
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.wantErr)
			}
			assert.False(t, a.connected)
			require.NotNil(t, a.characteristics)
			assert.Len(t, a.characteristics, wantCharaLen)
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
		wasConnected          bool
		dropCharacteristics   bool
		simulateDisconnectErr bool
		simulateConnectErr    bool
		wantErr               string
	}{
		"reconnect_not_connected": {},
		"reconnect_was_connected": {
			wasConnected: true,
		},
		"reconnect_drop_charas": {
			wasConnected:        true,
			dropCharacteristics: true,
		},
		"error_disconnect": {
			simulateDisconnectErr: true,
			wasConnected:          true,
			wantErr:               "device disconnect error",
		},
		"error_connect": {
			simulateConnectErr: true,
			wasConnected:       true,
			wantErr:            "adapter connect error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor(deviceAddress)
			extDevice := btTestDevice{simulateDisconnectErr: tc.simulateDisconnectErr}
			btdc := func(_ bluetoothExtDevicer, address, name string) *btDevice {
				return &btDevice{extDevice: extDevice, devAddress: address, devName: name}
			}
			extAdapter := &btTestAdapter{
				simulateConnectErr: tc.simulateConnectErr,
				deviceAddress:      deviceAddress,
				rssi:               rssi,
				payload:            &btTestPayload{name: deviceName},
			}
			a.btAdpt = &btAdapter{extAdapter: extAdapter, btDeviceCreator: btdc}
			a.cfg.scanTimeout = scanTimeout // to speed up test in case of errors
			a.cfg.sleepAfterDisconnect = 0  // to speed up test
			a.cfg.dropCharacteristicsOnDisconnect = tc.dropCharacteristics
			var wantCharaLen int
			if tc.wasConnected {
				a.btDevice = btdc(nil, "", "")
				a.connected = tc.wasConnected
				a.characteristics["charauuid"] = &btTestChara{}
				if !tc.dropCharacteristics {
					wantCharaLen = 1
				}
			}
			// act
			err := a.Reconnect()
			// assert
			if tc.wantErr == "" {
				require.NoError(t, err)
				assert.Equal(t, rssi, a.RSSI())
				assert.Len(t, a.characteristics, wantCharaLen)
				assert.True(t, a.connected)
			} else {
				require.ErrorContains(t, err, tc.wantErr)
			}
		})
	}
}

func TestFinalize(t *testing.T) {
	// this also tests Disconnect()
	tests := map[string]struct {
		simulateDisconnectErr bool
		wantErr               string
	}{
		"disconnect": {},
		"error_disconnect": {
			simulateDisconnectErr: true,
			wantErr:               "device disconnect error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor("")
			a.cfg.sleepAfterDisconnect = 0 // to speed up test
			extDevice := btTestDevice{simulateDisconnectErr: tc.simulateDisconnectErr}
			a.btDevice = &btDevice{extDevice: extDevice}
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

func TestUnsubscribe(t *testing.T) {
	const uuid = "00004321-0000-1000-8000-00805f9b34fb"
	tests := map[string]struct {
		inUUID       string
		notConnected bool
		chara        *btTestChara
		wantErr      string
	}{
		"unsubscribe_ok": {
			inUUID: uuid,
			chara:  &btTestChara{notificationFunc: func([]byte) {}},
		},
		"error_not_connected": {
			notConnected: true,
			wantErr:      "cannot unsubscribe from BLE device until connected",
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
			err := a.Unsubscribe(tc.inUUID)
			// assert
			if tc.wantErr == "" {
				require.NoError(t, err)
				require.NotNil(t, a.characteristics)
				require.NotNil(t, a.characteristics[uuid])
				assert.Nil(t, tc.chara.notificationFunc)
			} else {
				require.ErrorContains(t, err, tc.wantErr)
			}
		})
	}
}
