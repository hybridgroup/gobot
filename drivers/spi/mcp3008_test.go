package spi

import (
	"fmt"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*MCP3008Driver)(nil)

// must implement the AnalogReader interface
var _ aio.AnalogReader = (*MCP3008Driver)(nil)

func initTestMCP3008DriverWithStubbedAdaptor(simRead []byte) (*MCP3008Driver, *spiTestAdaptor) {
	a := newSpiTestAdaptor(simRead)
	d := NewMCP3008Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewMCP3008Driver(t *testing.T) {
	var di interface{} = NewMCP3008Driver(newSpiTestAdaptor([]byte{}))
	d, ok := di.(*MCP3008Driver)
	if !ok {
		t.Errorf("NewMCP3008Driver() should have returned a *MCP3008Driver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "MCP3008"), true)
}

func TestMCP3008Start(t *testing.T) {
	d := NewMCP3008Driver(newSpiTestAdaptor([]byte{}))
	gobottest.Assert(t, d.Start(), nil)
}

func TestMCP3008Halt(t *testing.T) {
	d, _ := initTestMCP3008DriverWithStubbedAdaptor([]byte{})
	gobottest.Assert(t, d.Halt(), nil)
}

func TestMCP3008Read(t *testing.T) {
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
			wantWritten: []byte{0x01, 0x80, 0x00},
			want:        0x03FF,
		},
		"number_1_ok": {
			chanNum:     1,
			simRead:     []byte{0xFF, 0xF1, 0xFF},
			wantWritten: []byte{0x01, 0x90, 0x00},
			want:        0x01FF,
		},
		"number_7_ok": {
			chanNum:     7,
			simRead:     []byte{0xFF, 0xF2, 0x11},
			wantWritten: []byte{0x01, 0xF0, 0x00},
			want:        0x0211,
		},
		"number_8_error": {
			chanNum: 8,
			wantErr: fmt.Errorf("Invalid channel '8' for read"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestMCP3008DriverWithStubbedAdaptor(tc.simRead)
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

func TestMCP3008ReadWithError(t *testing.T) {
	// arrange
	d, a := initTestMCP3008DriverWithStubbedAdaptor([]byte{})
	a.device.spiReadErr = true
	// act
	got, err := d.Read(0)
	// assert
	gobottest.Assert(t, err, fmt.Errorf("Error on SPI read in helper"))
	gobottest.Assert(t, got, 0)
}
