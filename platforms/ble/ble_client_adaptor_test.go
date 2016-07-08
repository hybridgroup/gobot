package ble

import (
	"testing"

	"github.com/hybridgroup/gobot/gobottest"
)

func initTestBLEAdaptor() *BLEAdaptor {
	a := NewBLEAdaptor("bot", "D7:99:5A:26:EC:38")
	// a.connect = func(n *BLEAdaptor) (io.ReadWriteCloser, error) {
	// 	return &NullReadWriteCloser{}, nil
	// }
	return a
}

func TestBLEAdaptor(t *testing.T) {
	a := NewBLEAdaptor("bot", "D7:99:5A:26:EC:38")
	gobottest.Assert(t, a.Name(), "bot")
	gobottest.Assert(t, a.UUID(), "D7:99:5A:26:EC:38")
}
