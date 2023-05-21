package i2c

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/gobottest"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*MMA7660Driver)(nil)

func initTestMMA7660DriverWithStubbedAdaptor() (*MMA7660Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewMMA7660Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewMMA7660Driver(t *testing.T) {
	var di interface{} = NewMMA7660Driver(newI2cTestAdaptor())
	d, ok := di.(*MMA7660Driver)
	if !ok {
		t.Errorf("NewMMA7660Driver() should have returned a *MMA7660Driver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "MMA7660"), true)
	gobottest.Assert(t, d.defaultAddress, 0x4c)
}

func TestMMA7660Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewMMA7660Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}

func TestMMA7660Start(t *testing.T) {
	d := NewMMA7660Driver(newI2cTestAdaptor())
	gobottest.Assert(t, d.Start(), nil)
}

func TestMMA7660Halt(t *testing.T) {
	d, _ := initTestMMA7660DriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestMMA7660Acceleration(t *testing.T) {
	d, _ := initTestMMA7660DriverWithStubbedAdaptor()
	x, y, z := d.Acceleration(21.0, 21.0, 21.0)
	gobottest.Assert(t, x, 1.0)
	gobottest.Assert(t, y, 1.0)
	gobottest.Assert(t, z, 1.0)
}

func TestMMA7660NullXYZ(t *testing.T) {
	d, _ := initTestMMA7660DriverWithStubbedAdaptor()

	x, y, z, _ := d.XYZ()
	gobottest.Assert(t, x, 0.0)
	gobottest.Assert(t, y, 0.0)
	gobottest.Assert(t, z, 0.0)
}

func TestMMA7660XYZ(t *testing.T) {
	d, a := initTestMMA7660DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{0x11, 0x12, 0x13})
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}

	x, y, z, _ := d.XYZ()
	gobottest.Assert(t, x, 17.0)
	gobottest.Assert(t, y, 18.0)
	gobottest.Assert(t, z, 19.0)
}

func TestMMA7660XYZError(t *testing.T) {
	d, a := initTestMMA7660DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}

	_, _, _, err := d.XYZ()
	gobottest.Assert(t, err, errors.New("read error"))
}

func TestMMA7660XYZNotReady(t *testing.T) {
	d, a := initTestMMA7660DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{0x40, 0x40, 0x40})
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}

	_, _, _, err := d.XYZ()
	gobottest.Assert(t, err, ErrNotReady)
}
