package onewire

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConnection(t *testing.T) {
	// arrange & act
	c := NewConnection(newOneWireTestSystemDevice("id"))
	// assert
	assert.IsType(t, &onewireConnection{}, c)
	assert.Equal(t, "id", c.ID())
}

func TestClose(t *testing.T) {
	tests := map[string]struct {
		simulateErr bool
		wantErr     string
	}{
		"close_ok": {},
		"error_close": {
			wantErr: "close error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			sysCon := newOneWireTestSystemDevice("id")
			if tc.wantErr != "" {
				sysCon.retErr = errors.New(tc.wantErr)
			}
			c := NewConnection(sysCon)
			// act
			err := c.Close()
			// assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestReadData(t *testing.T) {
	data := []byte{10, 11, 21, 32}
	tests := map[string]struct {
		data     []byte
		wantData []byte
		wantErr  string
	}{
		"read_ok": {
			// only to test the parameter passing
			data:     []byte{0, 0, 0},
			wantData: []byte{10, 11, 21},
		},
		"error_read": {
			data:     []byte{0, 0},
			wantData: []byte{10, 11},
			wantErr:  "read error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			const command = "read data command"
			sysCon := newOneWireTestSystemDevice("id")
			sysCon.lastData = data
			if tc.wantErr != "" {
				sysCon.retErr = errors.New(tc.wantErr)
			}
			c := NewConnection(sysCon)
			// act
			err := c.ReadData(command, tc.data)
			// assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, command, sysCon.lastCommand)
			assert.Equal(t, tc.wantData, tc.data)
		})
	}
}

func TestWriteData(t *testing.T) {
	tests := map[string]struct {
		data    []byte
		wantErr string
	}{
		"write_ok": {
			data: []byte{10, 11, 21},
		},
		"error_write": {
			data:    []byte{11, 32},
			wantErr: "read error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			const command = "write data command"
			sysCon := newOneWireTestSystemDevice("id")
			if tc.wantErr != "" {
				sysCon.retErr = errors.New(tc.wantErr)
			}
			c := NewConnection(sysCon)
			// act
			err := c.WriteData(command, tc.data)
			// assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, command, sysCon.lastCommand)
			assert.Equal(t, tc.data, sysCon.lastData)
		})
	}
}

func TestReadInteger(t *testing.T) {
	tests := map[string]struct {
		wantValue int
		wantErr   string
	}{
		"read_ok": {
			wantValue: 12,
		},
		"error_read": {
			wantErr: "read error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			const command = "read data command"
			sysCon := newOneWireTestSystemDevice("id")
			sysCon.lastValue = tc.wantValue
			if tc.wantErr != "" {
				sysCon.retErr = errors.New(tc.wantErr)
			}
			c := NewConnection(sysCon)
			// act
			got, err := c.ReadInteger(command)
			// assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, command, sysCon.lastCommand)
			assert.Equal(t, tc.wantValue, got)
		})
	}
}

func TestWriteInteger(t *testing.T) {
	tests := map[string]struct {
		value   int
		wantErr string
	}{
		"write_ok": {
			value: 21,
		},
		"error_write": {
			wantErr: "read error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			const command = "write data command"
			sysCon := newOneWireTestSystemDevice("id")
			if tc.wantErr != "" {
				sysCon.retErr = errors.New(tc.wantErr)
			}
			c := NewConnection(sysCon)
			// act
			err := c.WriteInteger(command, tc.value)
			// assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, command, sysCon.lastCommand)
			assert.Equal(t, tc.value, sysCon.lastValue)
		})
	}
}
