package ble

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/testutil"
)

var _ gobot.Driver = (*SerialPortDriver)(nil)

var _ io.ReadWriteCloser = (*SerialPortDriver)(nil)

func TestNewSerialPortDriver(t *testing.T) {
	d := NewSerialPortDriver(testutil.NewBleTestAdaptor(), "123", "456")
	assert.Equal(t, "01:02:03:0A:0B:0C", d.Address())
}

func TestNewSerialPortDriverWithName(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option.	Further
	// tests for options can also be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const newName = "new name"
	a := testutil.NewBleTestAdaptor()
	// act
	d := NewSerialPortDriver(a, "123", "456", WithName(newName))
	// assert
	assert.Equal(t, newName, d.Name())
}

func TestSerialPortOpen(t *testing.T) {
	const receiveCharacteristicUUID = "123"
	tests := map[string]struct {
		simConnectErr   bool
		simSubscribeErr bool
		wantErr         string
	}{
		"open_ok": {},
		"error_connect": {
			simConnectErr: true,
			wantErr:       "connect error",
		},
		"error_subscribe": {
			simSubscribeErr: true,
			wantErr:         "subscribe error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := testutil.NewBleTestAdaptor()
			a.SetSimulateConnectError(tc.simConnectErr)
			a.SetSimulateSubscribeError(tc.simSubscribeErr)
			d := NewSerialPortDriver(a, receiveCharacteristicUUID, "456")
			// act
			err := d.Open()
			// assert
			if tc.wantErr == "" {
				require.NoError(t, err)
				a.SendTestDataToSubscriber([]byte{3, 5, 7})
				assert.Equal(t, []byte{3, 5, 7}, d.responseData)
				assert.Equal(t, receiveCharacteristicUUID, a.SubscribeCharaUUID())
			} else {
				require.EqualError(t, err, tc.wantErr)
			}
		})
	}
}

func TestSerialPortRead(t *testing.T) {
	tests := map[string]struct {
		availableData  []byte
		readDataBuffer []byte
		wantCount      int
		wantData       []byte
		wantRemaining  []byte
	}{
		"no_data": {
			availableData:  []byte{},
			readDataBuffer: []byte{0, 0, 0},
			wantCount:      0,
			wantData:       []byte{0, 0, 0},
			wantRemaining:  nil,
		},
		"read_all": {
			availableData:  []byte{1, 2, 3},
			readDataBuffer: []byte{0, 0, 0},
			wantCount:      3,
			wantData:       []byte{1, 2, 3},
			wantRemaining:  nil,
		},
		"read_smaller": {
			availableData:  []byte{4, 6, 7},
			readDataBuffer: []byte{0, 0},
			wantCount:      2,
			wantData:       []byte{4, 6},
			wantRemaining:  []byte{7},
		},
		"read_bigger": {
			availableData:  []byte{7, 8},
			readDataBuffer: []byte{0, 0, 0},
			wantCount:      2,
			wantData:       []byte{7, 8, 0},
			wantRemaining:  nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d := NewSerialPortDriver(testutil.NewBleTestAdaptor(), "123", "456")
			d.responseData = append(d.responseData, tc.availableData...)
			// act
			gotCount, err := d.Read(tc.readDataBuffer)
			// assert
			require.NoError(t, err)
			assert.Equal(t, tc.wantCount, gotCount)
			assert.Equal(t, tc.wantData, tc.readDataBuffer)
			assert.Equal(t, tc.wantRemaining, d.responseData)
		})
	}
}

func TestSerialPortWrite(t *testing.T) {
	const transmitCharacteristicUUID = "456"
	tests := map[string]struct {
		writeData []byte
		simError  bool
		wantCount int
		wantData  []byte
		wantErr   string
	}{
		"write_ok": {
			writeData: []byte{1, 2, 3},
			wantCount: 3,
			wantData:  []byte{1, 2, 3},
		},
		"error_write": {
			writeData: []byte{1, 2, 3},
			simError:  true,
			wantCount: 3,
			wantData:  []byte{1, 2, 3},
			wantErr:   "write error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := testutil.NewBleTestAdaptor()
			var gotUUID string
			var gotData []byte
			a.SetWriteCharacteristicTestFunc(func(cUUID string, data []byte) error {
				gotUUID = cUUID
				gotData = append(gotData, data...)
				if tc.simError {
					return fmt.Errorf("write error")
				}
				return nil
			})
			d := NewSerialPortDriver(a, "123", transmitCharacteristicUUID)

			// act
			gotCount, err := d.Write(tc.writeData)
			// assert
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantErr)
			}
			assert.Equal(t, tc.wantCount, gotCount)
			assert.Equal(t, transmitCharacteristicUUID, gotUUID)
			assert.Equal(t, tc.wantData, gotData)
		})
	}
}

func TestSerialPortClose(t *testing.T) {
	tests := map[string]struct {
		simDisconnectErr bool
		wantErr          string
	}{
		"close_ok": {},
		"error_close": {
			simDisconnectErr: true,
			wantErr:          "disconnect error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := testutil.NewBleTestAdaptor()
			a.SetSimulateDisconnectError(tc.simDisconnectErr)
			d := NewSerialPortDriver(a, "123", "456")
			// act
			err := d.Close()
			// assert
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantErr)
			}
		})
	}
}
