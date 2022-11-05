package i2c

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*GenericDriver)(nil)

func TestNewGenericDriver(t *testing.T) {
	// arrange
	a := newI2cTestAdaptor()
	// act
	var di interface{} = NewGenericDriver(a, "GenericI2C", 0x17)
	// assert
	d, ok := di.(*GenericDriver)
	if !ok {
		t.Errorf("NewGenericDriver() should have returned a *GenericDriver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "GenericI2C"), true)
}
