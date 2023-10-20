package spi

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
)

// this ensures that the implementation is based on spi.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*MCP3004Driver)(nil)

// must implement the AnalogReader interface
var _ aio.AnalogReader = (*MCP3004Driver)(nil)

func initTestMCP3004DriverWithStubbedAdaptor() (*MCP3004Driver, *spiTestAdaptor) {
	a := newSpiTestAdaptor()
	d := NewMCP3004Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewMCP3004Driver(t *testing.T) {
	var di interface{} = NewMCP3004Driver(newSpiTestAdaptor())
	d, ok := di.(*MCP3004Driver)
	if !ok {
		t.Errorf("NewMCP3004Driver() should have returned a *MCP3004Driver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "MCP3004"))
}

func TestMCP3004Read(t *testing.T) {
	tests := map[string]struct {
		chanNum     int
		simRead     []byte
		want        int
		wantWritten []byte
		wantErr     error
	}{
		"number_negative_error": {
			chanNum: -1,
			wantErr: fmt.Errorf("Invalid channel '-1' for read"),
		},
		"number_0_ok": {
			chanNum:     0,
			simRead:     []byte{0xFF, 0xFF, 0xFF},
			wantWritten: []byte{0x01, 0x80, 0x00},
			want:        0x03FF,
		},
		"number_1_ok": {
			chanNum:     1,
			simRead:     []byte{0xFF, 0xF1, 0xFF},
			wantWritten: []byte{0x01, 0x90, 0x00},
			want:        0x01FF,
		},
		"number_3_ok": {
			chanNum:     3,
			simRead:     []byte{0xFF, 0xF2, 0x11},
			wantWritten: []byte{0x01, 0xB0, 0x00},
			want:        0x0211,
		},
		"number_4_error": {
			chanNum: 4,
			wantErr: fmt.Errorf("Invalid channel '4' for read"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestMCP3004DriverWithStubbedAdaptor()
			a.spi.SetSimRead(tc.simRead)
			// act
			got, err := d.Read(tc.chanNum)
			// assert
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantWritten, a.spi.Written())
		})
	}
}

func TestMCP3004ReadWithError(t *testing.T) {
	// arrange
	d, a := initTestMCP3004DriverWithStubbedAdaptor()
	a.spi.SetReadError(true)
	// act
	got, err := d.Read(0)
	// assert
	assert.Error(t, err, "error while SPI read in mock")
	assert.Equal(t, 0, got)
}
