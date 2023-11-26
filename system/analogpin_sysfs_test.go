package system

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newAnalogPinSysfs(t *testing.T) {
	// arrange
	m := &MockFilesystem{}
	sfa := sysfsFileAccess{fs: m, readBufLen: 2}
	const path = "/sys/whatever"
	// act
	pin := newAnalogPinSysfs(&sfa, path, true, false)
	// assert
	assert.Equal(t, path, pin.sysfsPath)
	assert.Equal(t, &sfa, pin.sfa)
	assert.True(t, pin.r)
	assert.False(t, pin.w)
}

func TestRead(t *testing.T) {
	const (
		sysfsPath = "/sys/testread"
		value     = "32"
	)
	tests := map[string]struct {
		readable bool
		simErr   bool
		wantVal  int
		wantErr  string
	}{
		"read_ok": {
			readable: true,
			wantVal:  32,
		},
		"error_not_readable": {
			readable: false,
			wantErr:  "the pin '/sys/testread' is not allowed to read",
		},
		"error_on_read": {
			readable: true,
			simErr:   true,
			wantErr:  "read error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			fs := newMockFilesystem([]string{sysfsPath})
			sfa := sysfsFileAccess{fs: fs, readBufLen: 2}
			pin := newAnalogPinSysfs(&sfa, sysfsPath, tc.readable, false)
			fs.Files[sysfsPath].Contents = value
			if tc.simErr {
				fs.Files[sysfsPath].simulateReadError = fmt.Errorf("read error")
			}
			// act
			got, err := pin.Read()
			// assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.wantVal, got)
		})
	}
}

func TestWrite(t *testing.T) {
	const (
		sysfsPath = "/sys/testwrite"
		oldVal    = "old_value"
	)
	tests := map[string]struct {
		writeable bool
		simErr    bool
		wantVal   string
		wantErr   string
	}{
		"write_ok": {
			writeable: true,
			wantVal:   "23",
		},
		"error_not_writeable": {
			writeable: false,
			wantErr:   "the pin '/sys/testwrite' is not allowed to write (val: 23)",
			wantVal:   oldVal,
		},
		"error_on_write": {
			writeable: true,
			simErr:    true,
			wantErr:   "write error",
			wantVal:   "23", // the mock is implemented in this way
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			fs := newMockFilesystem([]string{sysfsPath})
			sfa := sysfsFileAccess{fs: fs, readBufLen: 10}
			pin := newAnalogPinSysfs(&sfa, sysfsPath, false, tc.writeable)
			fs.Files[sysfsPath].Contents = oldVal
			if tc.simErr {
				fs.Files[sysfsPath].simulateWriteError = fmt.Errorf("write error")
			}
			// act
			err := pin.Write(23)
			// assert
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.wantVal, fs.Files[sysfsPath].Contents)
		})
	}
}
