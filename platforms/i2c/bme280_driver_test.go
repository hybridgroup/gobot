package i2c

import (
	"testing"
	"time"

	"github.com/hybridgroup/gobot/gobottest"
)

// --------- HELPERS
func initTestBME280Driver() (driver *BME280Driver) {
	driver, _ = initTestBME280DriverWithStubbedAdaptor()
	return
}

func initTestBME280DriverWithStubbedAdaptor() (*BME280Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor("adaptor")
	return NewBME280Driver(adaptor, "bot"), adaptor
}

// --------- TESTS

func TestNewBME280Driver(t *testing.T) {
	// Does it return a pointer to an instance of BME280Driver?
	var mpl interface{} = NewBME280Driver(newI2cTestAdaptor("adaptor"), "bot")
	_, ok := mpl.(*BME280Driver)
	if !ok {
		t.Errorf("NewBME280Driver() should have returned a *BME280Driver")
	}
}

// Methods
func TestBME280Driver(t *testing.T) {
	mpl := initTestBME280Driver()

	gobottest.Assert(t, mpl.Name(), "bot")
	gobottest.Assert(t, mpl.Connection().Name(), "adaptor")
	gobottest.Assert(t, mpl.interval, 10*time.Millisecond)

	mpl = NewBME280Driver(newI2cTestAdaptor("adaptor"), "bot", 100*time.Millisecond)
	gobottest.Assert(t, mpl.interval, 100*time.Millisecond)
}

func TestBME280DriverStart(t *testing.T) {
	mpl, adaptor := initTestBME280DriverWithStubbedAdaptor()

	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{0x00, 0x01, 0x02, 0x04}, nil
	}
	gobottest.Assert(t, len(mpl.Start()), 0)
	<-time.After(100 * time.Millisecond)
	gobottest.Assert(t, mpl.Pressure, float32(50.007942))
	gobottest.Assert(t, mpl.Temperature, float32(116.58878))
}

func TestBME280DriverHalt(t *testing.T) {
	mpl := initTestBME280Driver()

	gobottest.Assert(t, len(mpl.Halt()), 0)
}
