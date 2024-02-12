package serialport

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

func initTestAdaptor() (*Adaptor, *nullReadWriteCloser) {
	a := NewAdaptor("/dev/null")
	rwc := newNullReadWriteCloser()

	a.connectFunc = func(string, int) (io.ReadWriteCloser, error) {
		return rwc, nil
	}

	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, rwc
}

func TestNewAdaptor(t *testing.T) {
	// arrange
	a := NewAdaptor("/dev/null")
	assert.Equal(t, "/dev/null", a.Port())
	require.NotNil(t, a.cfg)
	assert.Equal(t, 115200, a.cfg.baudRate)
	assert.True(t, strings.HasPrefix(a.Name(), "Serial"))
}

func TestSerialRead(t *testing.T) {
	tests := map[string]struct {
		readDataBuffer []byte
		simReadErr     bool
		wantCount      int
		wantErr        string
	}{
		"read_ok": {
			readDataBuffer: []byte{0, 0},
			wantCount:      2,
		},
		"error_read": {
			readDataBuffer: []byte{},
			simReadErr:     true,
			wantErr:        "read error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a, rwc := initTestAdaptor()
			rwc.simulateReadErr = tc.simReadErr
			// act
			gotCount, err := a.SerialRead(tc.readDataBuffer)
			// assert
			if tc.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantErr)
			}
			assert.Equal(t, tc.wantCount, gotCount)
		})
	}
}

func TestSerialWrite(t *testing.T) {
	tests := map[string]struct {
		writeDataBuffer []byte
		simWriteErr     bool
		wantCount       int
		wantWritten     []byte
		wantErr         string
	}{
		"write_ok": {
			writeDataBuffer: []byte{1, 3, 6},
			wantWritten:     []byte{1, 3, 6},
			wantCount:       3,
		},
		"error_write": {
			writeDataBuffer: []byte{},
			simWriteErr:     true,
			wantErr:         "write error",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a, rwc := initTestAdaptor()
			rwc.simulateWriteErr = tc.simWriteErr
			// act
			gotCount, err := a.SerialWrite(tc.writeDataBuffer)
			// assert
			if tc.wantErr == "" {
				require.NoError(t, err)
				assert.Equal(t, tc.wantWritten, rwc.written)
			} else {
				require.EqualError(t, err, tc.wantErr)
			}
			assert.Equal(t, tc.wantCount, gotCount)
		})
	}
}

func TestConnect(t *testing.T) {
	// arrange
	a, _ := initTestAdaptor()
	require.True(t, a.IsConnected())
	// act & assert
	require.EqualError(t, a.Connect(), "serial port is already connected, try reconnect or run disconnect first")
	// re-arrange error
	a.sp = nil
	a.connectFunc = func(string, int) (io.ReadWriteCloser, error) {
		return nil, errors.New("connect error")
	}
	// act & assert
	require.ErrorContains(t, a.Connect(), "connect error")
	assert.False(t, a.IsConnected())
}

func TestReconnect(t *testing.T) {
	// arrange
	a, _ := initTestAdaptor()
	require.NotNil(t, a.sp)
	// act & assert
	require.NoError(t, a.Reconnect())
	require.NotNil(t, a.sp)
	// act & assert
	require.NoError(t, a.Disconnect())
	require.Nil(t, a.sp)
	// act & assert
	require.NoError(t, a.Reconnect())
	require.NotNil(t, a.sp)
}

func TestFinalize(t *testing.T) {
	// arrange
	a, rwc := initTestAdaptor()
	// act & assert
	require.NoError(t, a.Finalize())
	assert.False(t, a.IsConnected())
	// re-arrange error
	rwc.simulateCloseErr = true
	require.NoError(t, a.Connect())
	// act & assert
	require.ErrorContains(t, a.Finalize(), "close error")
}
