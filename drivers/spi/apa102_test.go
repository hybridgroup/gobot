package spi

import (
	"image/color"
	"strings"
	"testing"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/gobottest"
)

// this ensures that the implementation is based on spi.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*APA102Driver)(nil)

func initTestAPA102DriverWithStubbedAdaptor() *APA102Driver {
	a := newSpiTestAdaptor()
	d := NewAPA102Driver(a, 10, 31)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d
}

func TestNewAPA102Driver(t *testing.T) {
	var di interface{} = NewAPA102Driver(newSpiTestAdaptor(), 10, 31)
	d, ok := di.(*APA102Driver)
	if !ok {
		t.Errorf("NewAPA102Driver() should have returned a *APA102Driver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "APA102"), true)
}

func TestDriverLEDs(t *testing.T) {
	d := initTestAPA102DriverWithStubbedAdaptor()

	d.SetRGBA(0, color.RGBA{255, 255, 255, 15})
	d.SetRGBA(1, color.RGBA{255, 255, 255, 15})
	d.SetRGBA(2, color.RGBA{255, 255, 255, 15})
	d.SetRGBA(3, color.RGBA{255, 255, 255, 15})

	gobottest.Assert(t, d.Draw(), nil)
}
