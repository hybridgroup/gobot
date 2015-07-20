package i2c

import (
	"errors"
	"testing"

	"github.com/hybridgroup/gobot"
)

// --------- HELPERS
func initTestLIDARLiteDriver() (driver *LIDARLiteDriver) {
	driver, _ = initTestLIDARLiteDriverWithStubbedAdaptor()
	return
}

func initTestLIDARLiteDriverWithStubbedAdaptor() (*LIDARLiteDriver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor("adaptor")
	return NewLIDARLiteDriver(adaptor, "bot"), adaptor
}

// --------- TESTS

func TestNewLIDARLiteDriver(t *testing.T) {
	// Does it return a pointer to an instance of LIDARLiteDriver?
	var bm interface{} = NewLIDARLiteDriver(newI2cTestAdaptor("adaptor"), "bot")
	_, ok := bm.(*LIDARLiteDriver)
	if !ok {
		t.Errorf("NewLIDARLiteDriver() should have returned a *LIDARLiteDriver")
	}

	b := NewLIDARLiteDriver(newI2cTestAdaptor("adaptor"), "bot")
	gobot.Assert(t, b.Name(), "bot")
	gobot.Assert(t, b.Connection().Name(), "adaptor")
}

// Methods
func TestLIDARLiteDriverStart(t *testing.T) {
	hmc, adaptor := initTestLIDARLiteDriverWithStubbedAdaptor()

	gobot.Assert(t, len(hmc.Start()), 0)

	adaptor.i2cStartImpl = func() error {
		return errors.New("start error")
	}
	err := hmc.Start()
	gobot.Assert(t, err[0], errors.New("start error"))

}

func TestLIDARLiteDriverHalt(t *testing.T) {
	hmc := initTestLIDARLiteDriver()

	gobot.Assert(t, len(hmc.Halt()), 0)
}

func TestLIDARLiteDriverDistance(t *testing.T) {
	// when everything is happy
	hmc, adaptor := initTestLIDARLiteDriverWithStubbedAdaptor()

	first := true
	adaptor.i2cReadImpl = func() ([]byte, error) {
		if first {
			first = false
			return []byte{99}, nil
		}
		return []byte{1}, nil
	}

	distance, err := hmc.Distance()

	gobot.Assert(t, err, nil)
	gobot.Assert(t, distance, int(25345))

	// when insufficient bytes have been read
	hmc, adaptor = initTestLIDARLiteDriverWithStubbedAdaptor()

	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{}, nil
	}

	distance, err = hmc.Distance()
	gobot.Assert(t, distance, int(0))
	gobot.Assert(t, err, ErrNotEnoughBytes)

	// when read error
	hmc, adaptor = initTestLIDARLiteDriverWithStubbedAdaptor()

	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{}, errors.New("read error")
	}

	distance, err = hmc.Distance()
	gobot.Assert(t, distance, int(0))
	gobot.Assert(t, err, errors.New("read error"))

	// when write error
	hmc, adaptor = initTestLIDARLiteDriverWithStubbedAdaptor()

	adaptor.i2cWriteImpl = func() error {
		return errors.New("write error")
	}

	distance, err = hmc.Distance()
	gobot.Assert(t, distance, int(0))
	gobot.Assert(t, err, errors.New("write error"))
}
