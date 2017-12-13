package spi

import (
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*MCP3204Driver)(nil)

// must implement the AnalogReader interface
var _ aio.AnalogReader = (*MCP3204Driver)(nil)

func initTestMCP3204Driver() *MCP3204Driver {
	d := NewMCP3204Driver(&TestConnector{})
	return d
}

func TestMCP3204DriverStart(t *testing.T) {
	d := initTestMCP3204Driver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestMCP3204DriverHalt(t *testing.T) {
	d := initTestMCP3204Driver()
	d.Start()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestMCP3204DriverRead(t *testing.T) {
	d := initTestMCP3204Driver()
	d.Start()

	// TODO: actual read test
}
