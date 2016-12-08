package ble

import (
	"testing"

	"gobot.io/x/gobot/gobottest"
)

func initTestBLEClientAdaptor() *ClientAdaptor {
	a := NewClientAdaptor("D7:99:5A:26:EC:38")
	return a
}

func TestBLEClientAdaptor(t *testing.T) {
	a := NewClientAdaptor("D7:99:5A:26:EC:38")
	gobottest.Assert(t, a.UUID(), "D7:99:5A:26:EC:38")
}
