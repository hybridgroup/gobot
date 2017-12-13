package spi

import (
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*MCP3304Driver)(nil)

// must implement the AnalogReader interface
var _ aio.AnalogReader = (*MCP3304Driver)(nil)

func initTestMCP3304Driver() *MCP3304Driver {
	d := NewMCP3304Driver(&TestConnector{})
	return d
}

func TestMCP3304DriverStart(t *testing.T) {
	d := initTestMCP3304Driver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestMCP3304DriverHalt(t *testing.T) {
	d := initTestMCP3304Driver()
	d.Start()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestMCP3304DriverRead(t *testing.T) {
	d := initTestMCP3304Driver()
	d.Start()

	// TODO: actual read test
}
