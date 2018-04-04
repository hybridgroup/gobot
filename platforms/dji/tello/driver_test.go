package tello

import (
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*Driver)(nil)

func TestTelloDriver(t *testing.T) {
	d := NewDriver("8888")

	gobottest.Assert(t, d.respPort, "8888")
}
