package i2c

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*MPL115A2Driver)(nil)

// --------- HELPERS
func initTestMPL115A2Driver() (driver *MPL115A2Driver) {
	driver, _ = initTestMPL115A2DriverWithStubbedAdaptor()
	return
}

func initTestMPL115A2DriverWithStubbedAdaptor() (*MPL115A2Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewMPL115A2Driver(adaptor), adaptor
}

// --------- TESTS

func TestNewMPL115A2Driver(t *testing.T) {
	// Does it return a pointer to an instance of MPL115A2Driver?
	var mpl interface{} = NewMPL115A2Driver(newI2cTestAdaptor())
	_, ok := mpl.(*MPL115A2Driver)
	if !ok {
		t.Errorf("NewMPL115A2Driver() should have returned a *MPL115A2Driver")
	}
}

// Methods
func TestMPL115A2Driver(t *testing.T) {
	mpl := initTestMPL115A2Driver()

	gobottest.Refute(t, mpl.Connection(), nil)
	gobottest.Assert(t, strings.HasPrefix(mpl.Name(), "MPL115A2"), true)
}

func TestMPL115A2DriverOptions(t *testing.T) {
	mpl := NewMPL115A2Driver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, mpl.GetBusOrDefault(1), 2)
}

func TestMPL115A2DriverSetName(t *testing.T) {
	mpl := initTestMPL115A2Driver()
	mpl.SetName("TESTME")
	gobottest.Assert(t, mpl.Name(), "TESTME")
}

func TestMPL115A2DriverStart(t *testing.T) {
	mpl, _ := initTestMPL115A2DriverWithStubbedAdaptor()

	gobottest.Assert(t, mpl.Start(), nil)
}

func TestMPL115A2StartConnectError(t *testing.T) {
	d, adaptor := initTestMPL115A2DriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestMPL115A2DriverStartWriteError(t *testing.T) {
	mpl, adaptor := initTestMPL115A2DriverWithStubbedAdaptor()
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, mpl.Start(), errors.New("write error"))
}

func TestMPL115A2DriverReadData(t *testing.T) {
	mpl, adaptor := initTestMPL115A2DriverWithStubbedAdaptor()

	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{0x00, 0x01, 0x02, 0x04})
		return 4, nil
	}
	mpl.Start()

	press, _ := mpl.Pressure()
	temp, _ := mpl.Temperature()
	gobottest.Assert(t, press, float32(50.007942))
	gobottest.Assert(t, temp, float32(116.58878))
}

func TestMPL115A2DriverReadDataError(t *testing.T) {
	mpl, adaptor := initTestMPL115A2DriverWithStubbedAdaptor()
	mpl.Start()

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	_, err := mpl.Pressure()

	gobottest.Assert(t, err, errors.New("write error"))
}

func TestMPL115A2DriverHalt(t *testing.T) {
	mpl := initTestMPL115A2Driver()

	gobottest.Assert(t, mpl.Halt(), nil)
}
