package ble

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*GenericAccessDriver)(nil)

func initTestGenericAccessDriver() *GenericAccessDriver {
	d := NewGenericAccessDriver(NewClientAdaptor("D7:99:5A:26:EC:38"))
	return d
}

func TestGenericAccessDriver(t *testing.T) {
	d := initTestGenericAccessDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "GenericAccess"), true)
}
