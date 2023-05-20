package spi

import (
	"fmt"
	"strings"
	"testing"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
	"gobot.io/x/gobot/v2/gobottest"
)

// this ensures that the implementation is based on spi.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*MCP3204Driver)(nil)

// must implement the AnalogReader interface
var _ aio.AnalogReader = (*MCP3204Driver)(nil)

func initTestMCP3204DriverWithStubbedAdaptor() (*MCP3204Driver, *spiTestAdaptor) {
	a := newSpiTestAdaptor()
	d := NewMCP3204Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewMCP3204Driver(t *testing.T) {
	var di interface{} = NewMCP3204Driver(newSpiTestAdaptor())
	d, ok := di.(*MCP3204Driver)
	if !ok {
		t.Errorf("NewMCP3204Driver() should have returned a *MCP3204Driver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "MCP3204"), true)
}

func TestMCP3204Read(t *testing.T) {
	var tests = map[string]struct {
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
			wantWritten: []byte{0x06, 0x00, 0x00},
			want:        0x0FFF,
		},
		"number_1_ok": {
			chanNum:     1,
			simRead:     []byte{0xFF, 0xFE, 0xFF},
			wantWritten: []byte{0x06, 0x40, 0x00},
			want:        0x0EFF,
		},
		"number_3_ok": {
			chanNum:     3,
			simRead:     []byte{0xFF, 0xF3, 0x21},
			wantWritten: []byte{0x06, 0xC0, 0x00},
			want:        0x0321,
		},
		"number_4_error": {
			chanNum: 4,
			wantErr: fmt.Errorf("Invalid channel '4' for read"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestMCP3204DriverWithStubbedAdaptor()
			a.spi.SetSimRead(tc.simRead)
			// act
			got, err := d.Read(tc.chanNum)
			// assert
			gobottest.Assert(t, err, tc.wantErr)
			gobottest.Assert(t, got, tc.want)
			gobottest.Assert(t, a.spi.Written(), tc.wantWritten)
		})
	}
}

func TestMCP3204ReadWithError(t *testing.T) {
	// arrange
	d, a := initTestMCP3204DriverWithStubbedAdaptor()
	a.spi.SetReadError(true)
	// act
	got, err := d.Read(0)
	// assert
	gobottest.Assert(t, err, fmt.Errorf("error while SPI read in mock"))
	gobottest.Assert(t, got, 0)
}
