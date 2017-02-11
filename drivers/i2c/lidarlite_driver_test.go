package i2c

import (
	"errors"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*LIDARLiteDriver)(nil)

// --------- HELPERS
func initTestLIDARLiteDriver() (driver *LIDARLiteDriver) {
	driver, _ = initTestLIDARLiteDriverWithStubbedAdaptor()
	return
}

func initTestLIDARLiteDriverWithStubbedAdaptor() (*LIDARLiteDriver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewLIDARLiteDriver(adaptor), adaptor
}

// --------- TESTS

func TestNewLIDARLiteDriver(t *testing.T) {
	// Does it return a pointer to an instance of LIDARLiteDriver?
	var bm interface{} = NewLIDARLiteDriver(newI2cTestAdaptor())
	_, ok := bm.(*LIDARLiteDriver)
	if !ok {
		t.Errorf("NewLIDARLiteDriver() should have returned a *LIDARLiteDriver")
	}

	b := NewLIDARLiteDriver(newI2cTestAdaptor())
	gobottest.Refute(t, b.Connection(), nil)
}

// Methods
func TestLIDARLiteDriverStart(t *testing.T) {
	hmc, _ := initTestLIDARLiteDriverWithStubbedAdaptor()

	gobottest.Assert(t, hmc.Start(), nil)
}

func TestLIDARLiteDriverHalt(t *testing.T) {
	hmc := initTestLIDARLiteDriver()

	gobottest.Assert(t, hmc.Halt(), nil)
}

func TestLIDARLiteDriverDistance(t *testing.T) {
	// when everything is happy
	hmc, adaptor := initTestLIDARLiteDriverWithStubbedAdaptor()

	gobottest.Assert(t, hmc.Start(), nil)

	first := true
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		if first {
			first = false
			copy(b, []byte{99})
			return 1, nil
		}
		copy(b, []byte{1})
		return 1, nil
	}

	distance, err := hmc.Distance()

	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, distance, int(25345))

	// when insufficient bytes have been read
	hmc, adaptor = initTestLIDARLiteDriverWithStubbedAdaptor()

	gobottest.Assert(t, hmc.Start(), nil)

	adaptor.i2cReadImpl = func([]byte) (int, error) {
		return 0, nil
	}

	distance, err = hmc.Distance()
	gobottest.Assert(t, distance, int(0))
	gobottest.Assert(t, err, ErrNotEnoughBytes)

	// when read error
	hmc, adaptor = initTestLIDARLiteDriverWithStubbedAdaptor()

	gobottest.Assert(t, hmc.Start(), nil)

	adaptor.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}

	distance, err = hmc.Distance()
	gobottest.Assert(t, distance, int(0))
	gobottest.Assert(t, err, errors.New("read error"))

	// when write error
	hmc, adaptor = initTestLIDARLiteDriverWithStubbedAdaptor()

	gobottest.Assert(t, hmc.Start(), nil)

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	distance, err = hmc.Distance()
	gobottest.Assert(t, distance, int(0))
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestLIDARLiteDriverSetName(t *testing.T) {
	l := initTestLIDARLiteDriver()
	l.SetName("TESTME")
	gobottest.Assert(t, l.Name(), "TESTME")
}

func TestLIDARLiteDriverOptions(t *testing.T) {
	l := NewLIDARLiteDriver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, l.GetBusOrDefault(1), 2)
}
