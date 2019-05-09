package spi

import (
	"image/color"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*APA102Driver)(nil)

func initTestDriver() *APA102Driver {
	d := NewAPA102Driver(&TestConnector{}, 10, 31)
	return d
}

func TestDriverStart(t *testing.T) {
	d := initTestDriver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestDriverHalt(t *testing.T) {
	d := initTestDriver()
	d.Start()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestDriverLEDs(t *testing.T) {
	d := initTestDriver()
	d.Start()

	d.SetRGBA(0, Color{color.RGBA{255, 255, 255, 0}, 15})
	d.SetRGBA(1, Color{color.RGBA{255, 255, 255, 0}, 15})
	d.SetRGBA(2, Color{color.RGBA{255, 255, 255, 0}, 15})
	d.SetRGBA(3, Color{color.RGBA{255, 255, 255, 0}, 15})

	gobottest.Assert(t, d.Draw(), nil)
}
