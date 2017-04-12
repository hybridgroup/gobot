package i2c

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*MMA7660Driver)(nil)

// --------- HELPERS
func initTestMMA7660Driver() (driver *MMA7660Driver) {
	driver, _ = initTestMMA7660DriverWithStubbedAdaptor()
	return
}

func initTestMMA7660DriverWithStubbedAdaptor() (*MMA7660Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewMMA7660Driver(adaptor), adaptor
}

// --------- TESTS

func TestNewMMA7660Driver(t *testing.T) {
	// Does it return a pointer to an instance of MMA7660Driver?
	var mma interface{} = NewMMA7660Driver(newI2cTestAdaptor())
	_, ok := mma.(*MMA7660Driver)
	if !ok {
		t.Errorf("NewMMA7660Driver() should have returned a *MMA7660Driver")
	}
}

// Methods
func TestMMA7660Driver(t *testing.T) {
	mma := initTestMMA7660Driver()

	gobottest.Refute(t, mma.Connection(), nil)
	gobottest.Assert(t, strings.HasPrefix(mma.Name(), "MMA7660"), true)
}

func TestMMA7660DriverSetName(t *testing.T) {
	d := initTestMMA7660Driver()
	d.SetName("TESTME")
	gobottest.Assert(t, d.Name(), "TESTME")
}

func TestMMA7660DriverOptions(t *testing.T) {
	d := NewMMA7660Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}

func TestMMA7660DriverStart(t *testing.T) {
	d := initTestMMA7660Driver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestMMA7660StartConnectError(t *testing.T) {
	d, adaptor := initTestMMA7660DriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestMMA7660DriverStartWriteError(t *testing.T) {
	mma, adaptor := initTestMMA7660DriverWithStubbedAdaptor()
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, mma.Start(), errors.New("write error"))
}

func TestMMA7660DriverHalt(t *testing.T) {
	d := initTestMMA7660Driver()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestMMA7660DriverAcceleration(t *testing.T) {
	d := initTestMMA7660Driver()
	x, y, z := d.Acceleration(21.0, 21.0, 21.0)
	gobottest.Assert(t, x, 1.0)
	gobottest.Assert(t, y, 1.0)
	gobottest.Assert(t, z, 1.0)
}

func TestMMA7660DriverNullXYZ(t *testing.T) {
	d, _ := initTestMMA7660DriverWithStubbedAdaptor()
	d.Start()
	x, y, z, _ := d.XYZ()
	gobottest.Assert(t, x, 0.0)
	gobottest.Assert(t, y, 0.0)
	gobottest.Assert(t, z, 0.0)
}

func TestMMA7660DriverXYZ(t *testing.T) {
	d, adaptor := initTestMMA7660DriverWithStubbedAdaptor()
	d.Start()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
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

func TestMMA7660DriverXYZError(t *testing.T) {
	d, adaptor := initTestMMA7660DriverWithStubbedAdaptor()
	d.Start()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}

	_, _, _, err := d.XYZ()
	gobottest.Assert(t, err, errors.New("read error"))
}

func TestMMA7660DriverXYZNotReady(t *testing.T) {
	d, adaptor := initTestMMA7660DriverWithStubbedAdaptor()
	d.Start()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{0x40, 0x40, 0x40})
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}

	_, _, _, err := d.XYZ()
	gobottest.Assert(t, err, ErrNotReady)
}
