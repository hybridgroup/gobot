package spi

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*Driver)(nil)

func initTestDriverWithStubbedAdaptor() (*Driver, *spiTestAdaptor) {
	a := newSpiTestAdaptor()
	d := NewDriver(a, "SPI_BASIC")
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewDriver(t *testing.T) {
	var di interface{} = NewDriver(newSpiTestAdaptor(), "SPI_BASIC")
	d, ok := di.(*Driver)
	if !ok {
		t.Errorf("NewDriver() should have returned a *Driver")
	}
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "SPI_BASIC"), true)
}

func TestStart(t *testing.T) {
	d := NewDriver(newSpiTestAdaptor(), "SPI_BASIC")
	gobottest.Assert(t, d.Start(), nil)
}

func TestHalt(t *testing.T) {
	d, _ := initTestDriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestSetName(t *testing.T) {
	// arrange
	d, _ := initTestDriverWithStubbedAdaptor()
	// act
	d.SetName("TESTME")
	// assert
	gobottest.Assert(t, d.Name(), "TESTME")
}

func TestConnection(t *testing.T) {
	d, _ := initTestDriverWithStubbedAdaptor()
	gobottest.Refute(t, d.Connection(), nil)
}
