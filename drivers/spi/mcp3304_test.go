package spi

import (
	"fmt"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*MCP3304Driver)(nil)

// must implement the AnalogReader interface
var _ aio.AnalogReader = (*MCP3304Driver)(nil)

func initTestMCP3304DriverWithStubbedAdaptor(simRead []byte) (*MCP3304Driver, *spiTestAdaptor) {
	a := newSpiTestAdaptor(simRead)
	d := NewMCP3304Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewMCP3304Driver(t *testing.T) {
	var di interface{} = NewMCP3304Driver(newSpiTestAdaptor([]byte{}))
	d, ok := di.(*MCP3304Driver)
	if !ok {
		t.Errorf("NewMCP3304Driver() should have returned a *MCP3304Driver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "MCP3304"), true)
}

func TestMCP3304Start(t *testing.T) {
	d := NewMCP3304Driver(newSpiTestAdaptor([]byte{}))
	gobottest.Assert(t, d.Start(), nil)
}

func TestMCP3304Halt(t *testing.T) {
	d, _ := initTestMCP3304DriverWithStubbedAdaptor([]byte{})
	gobottest.Assert(t, d.Halt(), nil)
}

func TestMCP3304Read(t *testing.T) {
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
			wantWritten: []byte{0x0C, 0x00, 0x00},
			want:        0x0FFF,
		},
		"number_1_ok": {
			chanNum:     1,
			simRead:     []byte{0xFF, 0xFE, 0xFF},
			wantWritten: []byte{0x0C, 0x80, 0x00},
			want:        0x0EFF,
		},
		"number_7_ok": {
			chanNum:     7,
			simRead:     []byte{0xFF, 0xF7, 0x65},
			wantWritten: []byte{0x0F, 0x80, 0x00},
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
			d, a := initTestMCP3304DriverWithStubbedAdaptor(tc.simRead)
			copy(a.device.simRead, tc.simRead)
			// act
			got, err := d.Read(tc.chanNum)
			// assert
			gobottest.Assert(t, err, tc.wantErr)
			gobottest.Assert(t, got, tc.want)
			gobottest.Assert(t, a.device.written, tc.wantWritten)
		})
	}
}

func TestMCP3304ReadWithError(t *testing.T) {
	// arrange
	d, a := initTestMCP3304DriverWithStubbedAdaptor([]byte{})
	a.device.spiReadErr = true
	// act
	got, err := d.Read(0)
	// assert
	gobottest.Assert(t, err, fmt.Errorf("Error on SPI read in helper"))
	gobottest.Assert(t, got, 0)
}
