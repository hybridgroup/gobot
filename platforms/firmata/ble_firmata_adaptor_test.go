package firmata

import (
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Adaptor = (*BLEAdaptor)(nil)

func initTestBLEAdaptor() *BLEAdaptor {
	a := NewBLEAdaptor("DEVICE", "123", "456")
	return a
}

func TestFirmataBLEAdaptor(t *testing.T) {
	a := initTestBLEAdaptor()
	gobottest.Assert(t, a.Name(), "BLEFirmata")
}
