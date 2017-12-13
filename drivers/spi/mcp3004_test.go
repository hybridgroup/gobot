package spi

import (
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*MCP3004Driver)(nil)

// must implement the AnalogReader interface
var _ aio.AnalogReader = (*MCP3004Driver)(nil)

func initTestMCP3004Driver() *MCP3004Driver {
	d := NewMCP3004Driver(&TestConnector{})
	return d
}

func TestMCP3004DriverStart(t *testing.T) {
	d := initTestMCP3004Driver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestMCP3004DriverHalt(t *testing.T) {
	d := initTestMCP3004Driver()
	d.Start()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestMCP3004DriverRead(t *testing.T) {
	d := initTestMCP3004Driver()
	d.Start()

	// TODO: actual read test
}
