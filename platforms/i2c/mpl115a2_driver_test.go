package i2c

import (
	"testing"

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

func TestMPL115A2DriverDriver(t *testing.T) {
	// Does it implement gobot.DriverInterface?
	var _ gobot.Driver = (*MPL115A2Driver)(nil)

	// Does its adaptor implements the I2cInterface?
	driver := initTestMPL115A2Driver()
	var _ I2cInterface = driver.adaptor()
}

func TestNewMPL115A2Driver(t *testing.T) {
	// Does it return a pointer to an instance of MPL115A2Driver?
	var mpl interface{} = NewMPL115A2Driver(newI2cTestAdaptor("adaptor"), "bot")
	_, ok := mpl.(*MPL115A2Driver)
	if !ok {
		t.Errorf("NewMPL115A2Driver() should have returned a *MPL115A2Driver")
	}
}

// Methods
func TestMPL115A2DriverStart(t *testing.T) {
	mpl := initTestMPL115A2Driver()

	gobot.Assert(t, len(mpl.Start()), 0)
}

func TestMPL115A2DriverHalt(t *testing.T) {
	mpl := initTestMPL115A2Driver()

	gobot.Assert(t, len(mpl.Halt()), 0)
}
