package ble

import (
	"testing"

	"github.com/hybridgroup/gobot/gobottest"
)

func initTestBLEClientAdaptor() *BLEClientAdaptor {
	a := NewBLEClientAdaptor("bot", "D7:99:5A:26:EC:38")
	return a
}

func TestBLEClientAdaptor(t *testing.T) {
	a := NewBLEClientAdaptor("bot", "D7:99:5A:26:EC:38")
	gobottest.Assert(t, a.Name(), "bot")
	gobottest.Assert(t, a.UUID(), "D7:99:5A:26:EC:38")
}
