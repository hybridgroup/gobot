package i2c

import (
	"errors"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*HMC6352Driver)(nil)

// --------- HELPERS
func initTestHMC6352Driver() (driver *HMC6352Driver) {
	driver, _ = initTestHMC6352DriverWithStubbedAdaptor()
	return
}

func initTestHMC6352DriverWithStubbedAdaptor() (*HMC6352Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewHMC6352Driver(adaptor), adaptor
}

// --------- TESTS

func TestNewHMC6352Driver(t *testing.T) {
	// Does it return a pointer to an instance of HMC6352Driver?
	var bm interface{} = NewHMC6352Driver(newI2cTestAdaptor())
	_, ok := bm.(*HMC6352Driver)
	if !ok {
		t.Errorf("NewHMC6352Driver() should have returned a *HMC6352Driver")
	}

	b := NewHMC6352Driver(newI2cTestAdaptor())
	gobottest.Refute(t, b.Connection(), nil)
}

// Methods
func TestHMC6352DriverStart(t *testing.T) {
	hmc, adaptor := initTestHMC6352DriverWithStubbedAdaptor()

	gobottest.Assert(t, hmc.Start(), nil)

	adaptor.i2cWriteImpl = func() error {
		return errors.New("write error")
	}
	err := hmc.Start()
	gobottest.Assert(t, err, errors.New("write error"))

	adaptor.i2cStartImpl = func() error {
		return errors.New("start error")
	}
	err = hmc.Start()
	gobottest.Assert(t, err, errors.New("start error"))

}

func TestHMC6352DriverHalt(t *testing.T) {
	hmc := initTestHMC6352Driver()

	gobottest.Assert(t, hmc.Halt(), nil)
}

func TestHMC6352DriverHeading(t *testing.T) {
	// when len(data) is 2
	hmc, adaptor := initTestHMC6352DriverWithStubbedAdaptor()

	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{99, 1}, nil
	}

	heading, _ := hmc.Heading()
	gobottest.Assert(t, heading, uint16(2534))

	// when len(data) is not 2
	hmc, adaptor = initTestHMC6352DriverWithStubbedAdaptor()

	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{99}, nil
	}

	heading, err := hmc.Heading()
	gobottest.Assert(t, heading, uint16(0))
	gobottest.Assert(t, err, ErrNotEnoughBytes)

	// when read error
	hmc, adaptor = initTestHMC6352DriverWithStubbedAdaptor()

	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{}, errors.New("read error")
	}

	heading, err = hmc.Heading()
	gobottest.Assert(t, heading, uint16(0))
	gobottest.Assert(t, err, errors.New("read error"))

	// when write error
	hmc, adaptor = initTestHMC6352DriverWithStubbedAdaptor()

	adaptor.i2cWriteImpl = func() error {
		return errors.New("write error")
	}

	heading, err = hmc.Heading()
	gobottest.Assert(t, heading, uint16(0))
	gobottest.Assert(t, err, errors.New("write error"))
}
