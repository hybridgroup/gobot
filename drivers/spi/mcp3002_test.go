package spi

import (
	"fmt"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*MCP3002Driver)(nil)

// must implement the AnalogReader interface
var _ aio.AnalogReader = (*MCP3002Driver)(nil)

func initTestMCP3002DriverWithStubbedAdaptor(simRead []byte) (*MCP3002Driver, *spiTestAdaptor) {
	a := newSpiTestAdaptor(simRead)
	d := NewMCP3002Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewMCP3002Driver(t *testing.T) {
	var di interface{} = NewMCP3002Driver(newSpiTestAdaptor([]byte{}))
	d, ok := di.(*MCP3002Driver)
	if !ok {
		t.Errorf("NewMCP3002Driver() should have returned a *MCP3002Driver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "MCP3002"), true)
}

func TestMCP3002Start(t *testing.T) {
	d := NewMCP3002Driver(newSpiTestAdaptor([]byte{}))
	gobottest.Assert(t, d.Start(), nil)
}

func TestMCP3002Halt(t *testing.T) {
	d, _ := initTestMCP3002DriverWithStubbedAdaptor([]byte{})
	gobottest.Assert(t, d.Halt(), nil)
}

func TestMCP3002Read(t *testing.T) {
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
			simRead:     []byte{0xFF, 0xFF},
			wantWritten: []byte{0x68, 0x00},
			want:        0x3FF,
		},
		"number_1_ok": {
			chanNum:     1,
			simRead:     []byte{0xF2, 0x22},
			wantWritten: []byte{0x78, 0x00},
			want:        0x222,
		},
		"number_2_error": {
			chanNum: 2,
			wantErr: fmt.Errorf("Invalid channel '2' for read"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestMCP3002DriverWithStubbedAdaptor(tc.simRead)
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

func TestMCP3002ReadWithError(t *testing.T) {
	// arrange
	d, a := initTestMCP3002DriverWithStubbedAdaptor([]byte{})
	a.device.spiReadErr = true
	// act
	got, err := d.Read(0)
	// assert
	gobottest.Assert(t, err, fmt.Errorf("Error on SPI read in helper"))
	gobottest.Assert(t, got, 0)
}
