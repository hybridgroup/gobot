package gobot

import "testing"

func TestAdaptor(t *testing.T) {
	a := NewAdaptor("", "testBot", "/dev/null")
	Refute(t, a.Name(), "")
	a.SetPort("/dev/null1")
	Assert(t, a.Port(), "/dev/null1")
	a.SetName("myAdaptor")
	Assert(t, a.Name(), "myAdaptor")
	a.SetConnected(true)
	Assert(t, a.Connected(), true)
}
