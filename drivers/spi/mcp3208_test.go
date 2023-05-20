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
var _ gobot.Driver = (*MCP3208Driver)(nil)

// must implement the AnalogReader interface
var _ aio.AnalogReader = (*MCP3208Driver)(nil)

func initTestMCP3208DriverWithStubbedAdaptor() (*MCP3208Driver, *spiTestAdaptor) {
	a := newSpiTestAdaptor()
	d := NewMCP3208Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewMCP3208Driver(t *testing.T) {
	var di interface{} = NewMCP3208Driver(newSpiTestAdaptor())
	d, ok := di.(*MCP3208Driver)
	if !ok {
		t.Errorf("NewMCP3208Driver() should have returned a *MCP3208Driver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "MCP3208"), true)
}

func TestMCP3208Read(t *testing.T) {
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
		"number_7_ok": {
			chanNum:     7,
			simRead:     []byte{0xFF, 0xF7, 0x65},
			wantWritten: []byte{0x07, 0xC0, 0x00},
			want:        0x0765,
		},
		"number_8_error": {
			chanNum: 8,
			wantErr: fmt.Errorf("Invalid channel '8' for read"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestMCP3208DriverWithStubbedAdaptor()
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

func TestMCP3208ReadWithError(t *testing.T) {
	// arrange
	d, a := initTestMCP3208DriverWithStubbedAdaptor()
	a.spi.SetReadError(true)
	// act
	got, err := d.Read(0)
	// assert
	gobottest.Assert(t, err, fmt.Errorf("error while SPI read in mock"))
	gobottest.Assert(t, got, 0)
}
