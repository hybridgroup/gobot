package i2c

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*HMC6352Driver)(nil)

func initTestHMC6352DriverWithStubbedAdaptor() (*HMC6352Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewHMC6352Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewHMC6352Driver(t *testing.T) {
	var di interface{} = NewHMC6352Driver(newI2cTestAdaptor())
	d, ok := di.(*HMC6352Driver)
	if !ok {
		t.Errorf("NewHMC6352Driver() should have returned a *HMC6352Driver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "HMC6352"), true)
	gobottest.Assert(t, d.defaultAddress, 0x21)
}

func TestHMC6352Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewHMC6352Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}

func TestHMC6352Start(t *testing.T) {
	d := NewHMC6352Driver(newI2cTestAdaptor())
	gobottest.Assert(t, d.Start(), nil)
}

func TestHMC6352Halt(t *testing.T) {
	d, _ := initTestHMC6352DriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestHMC6352Heading(t *testing.T) {
	// when len(data) is 2
	d, a := initTestHMC6352DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99, 1})
		return 2, nil
	}

	heading, _ := d.Heading()
	gobottest.Assert(t, heading, uint16(2534))

	// when len(data) is not 2
	d, a = initTestHMC6352DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99})
		return 1, nil
	}

	heading, err := d.Heading()
	gobottest.Assert(t, heading, uint16(0))
	gobottest.Assert(t, err, ErrNotEnoughBytes)

	// when read error
	d, a = initTestHMC6352DriverWithStubbedAdaptor()
	a.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}

	heading, err = d.Heading()
	gobottest.Assert(t, heading, uint16(0))
	gobottest.Assert(t, err, errors.New("read error"))

	// when write error
	d, a = initTestHMC6352DriverWithStubbedAdaptor()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	heading, err = d.Heading()
	gobottest.Assert(t, heading, uint16(0))
	gobottest.Assert(t, err, errors.New("write error"))
}
