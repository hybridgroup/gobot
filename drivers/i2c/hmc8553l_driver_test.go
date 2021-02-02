package i2c

import (
	"errors"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*HMC8553LDriver)(nil)

// --------- HELPERS
func initTestHMC8553LDriver() (driver *HMC8553LDriver) {
	driver, _ = initTestHMC8553LDriverWithStubbedAdaptor()
	return
}

func initTestHMC8553LDriverWithStubbedAdaptor() (*HMC8553LDriver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewHMC8553LDriver(adaptor), adaptor
}

// --------- TESTS

func TestNewHMC8553LDriver(t *testing.T) {
	// Does it return a pointer to an instance of HMC8553LDriver?
	var bm interface{} = NewHMC8553LDriver(newI2cTestAdaptor())
	_, ok := bm.(*HMC8553LDriver)
	if !ok {
		t.Errorf("NewHMC8553LDriver() should have returned a *HMC8553LDriver")
	}

	b := NewHMC8553LDriver(newI2cTestAdaptor())
	gobottest.Refute(t, b.Connection(), nil)
}

// Methods
func TestHMC8553LDriverStart(t *testing.T) {
	hmc, adaptor := initTestHMC8553LDriverWithStubbedAdaptor()

	gobottest.Assert(t, hmc.Start(), nil)

	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	err := hmc.Start()
	gobottest.Assert(t, err, errors.New("write error"))
}

func Test8553LStartConnectError(t *testing.T) {
	d, adaptor := initTestHMC8553LDriverWithStubbedAdaptor()
	adaptor.Testi2cConnectErr(true)
	gobottest.Assert(t, d.Start(), errors.New("Invalid i2c connection"))
}

func TestHMC8553LDriverHalt(t *testing.T) {
	hmc := initTestHMC8553LDriver()

	gobottest.Assert(t, hmc.Halt(), nil)
}

func TestHMC8553LDriverSetName(t *testing.T) {
	d := initTestHMC8553LDriver()
	d.SetName("TESTME")
	gobottest.Assert(t, d.Name(), "TESTME")
}

func TestHMC8553LDriverOptions(t *testing.T) {
	d := NewHMC8553LDriver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}
