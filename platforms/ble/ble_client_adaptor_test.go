package ble

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Adaptor = (*ClientAdaptor)(nil)

func initTestBLEClientAdaptor() *ClientAdaptor {
	a := NewClientAdaptor("D7:99:5A:26:EC:38")
	return a
}

func TestBLEClientAdaptor(t *testing.T) {
	a := NewClientAdaptor("D7:99:5A:26:EC:38")
	gobottest.Assert(t, a.Address(), "D7:99:5A:26:EC:38")
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "BLEClient"), true)
}

func TestBLEClientAdaptorName(t *testing.T) {
	a := NewClientAdaptor("D7:99:5A:26:EC:38")
	a.SetName("awesome")
	gobottest.Assert(t, a.Name(), "awesome")
}
