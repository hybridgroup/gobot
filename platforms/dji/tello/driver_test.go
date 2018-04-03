package tello

import (
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*Driver)(nil)

func TestTelloDriver(t *testing.T) {
	d := NewDriver("127.0.0.1:8888")

	gobottest.Assert(t, d.respAddr, "127.0.0.1:8888")
}
