package i2c

import (
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
)

// --------- HELPERS
func initTestMPL115A2Driver() (driver *MPL115A2Driver) {
	driver, _ = initTestMPL115A2DriverWithStubbedAdaptor()
	return
}

func initTestMPL115A2DriverWithStubbedAdaptor() (*MPL115A2Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor("adaptor")
	return NewMPL115A2Driver(adaptor, "bot"), adaptor
}

// --------- TESTS

func TestNewMPL115A2Driver(t *testing.T) {
	// Does it return a pointer to an instance of MPL115A2Driver?
	var mpl interface{} = NewMPL115A2Driver(newI2cTestAdaptor("adaptor"), "bot")
	_, ok := mpl.(*MPL115A2Driver)
	if !ok {
		t.Errorf("NewMPL115A2Driver() should have returned a *MPL115A2Driver")
	}
}

// Methods
func TestMPL115A2Driver(t *testing.T) {
	mpl := initTestMPL115A2Driver()

	gobot.Assert(t, mpl.Name(), "bot")
	gobot.Assert(t, mpl.Connection().Name(), "adaptor")
	gobot.Assert(t, mpl.interval, 10*time.Millisecond)

	mpl = NewMPL115A2Driver(newI2cTestAdaptor("adaptor"), "bot", 100*time.Millisecond)
	gobot.Assert(t, mpl.interval, 100*time.Millisecond)
}

func TestMPL115A2DriverStart(t *testing.T) {
	mpl, adaptor := initTestMPL115A2DriverWithStubbedAdaptor()

	adaptor.i2cReadImpl = func() ([]byte, error) {
		return []byte{0x00, 0x01, 0x02, 0x04}, nil
	}
	gobot.Assert(t, len(mpl.Start()), 0)
	<-time.After(100 * time.Millisecond)
	gobot.Assert(t, mpl.Pressure, float32(50.007942))
	gobot.Assert(t, mpl.Temperature, float32(116.58878))
}

func TestMPL115A2DriverHalt(t *testing.T) {
	mpl := initTestMPL115A2Driver()

	gobot.Assert(t, len(mpl.Halt()), 0)
}
