package i2c

import (
	"github.com/hybridgroup/gobot"
	"testing"
	"time"
)

// --------- HELPERS
func initTestHMC6352Driver() *HMC6352Driver {
	return NewHMC6352Driver(newI2cTestAdaptor("adaptor"), "bot")
}

func initTestMockedHMC6352Driver() (*HMC6352Driver, *I2cInterfaceClient) {
	inter := NewI2cInterfaceClient()
	return NewHMC6352Driver(inter, "bot"), inter
}

// --------- TESTS

func TestHMC6352Driver(t *testing.T) {
	// Does it implement gobot.DriverInterface?
	var _ gobot.DriverInterface = (*HMC6352Driver)(nil)

	// Does its adaptor implements the I2cInterface?
	driver := initTestHMC6352Driver()
	var _ I2cInterface = driver.adaptor()
}

func TestNewHMC6352Driver(t *testing.T) {
	// Does it return a pointer to an instance of HMC6352Driver?
	var bm interface{} = NewHMC6352Driver(newI2cTestAdaptor("adaptor"), "bot")
	_, ok := bm.(*HMC6352Driver)
	if !ok {
		t.Errorf("NewHMC6352Driver() should have returned a *HMC6352Driver")
	}
}

// Methods
func TestHMC6352DriverAdaptor(t *testing.T) {
	hmc := initTestHMC6352Driver()

	gobot.Assert(t, hmc.adaptor(), hmc.Adaptor())
}

func TestHMC6352DriverStart(t *testing.T) {
	// when length of data returned by I2cRead is 2
	hmc, inter := initTestMockedHMC6352Driver()

	numberOfCyclesForEvery := 3

	inter.When("I2cStart", uint8(0x21))
	inter.When("I2cWrite", []byte("A")).Times(numberOfCyclesForEvery)
	inter.When("I2cRead", uint(2)).Return([]byte{99, 1}).Times(numberOfCyclesForEvery - 1)

	// Make sure "Heading" is set to 0 so that we assert
	// its new value after executing "Start()"
	gobot.Assert(t, hmc.Heading, uint16(0))

	hmc.SetInterval(1 * time.Millisecond)
	gobot.Assert(t, hmc.Start(), true)
	<-time.After(time.Duration(numberOfCyclesForEvery) * time.Millisecond)

	if ok, err := inter.Verify(); !ok {
		t.Errorf("Error:", err)
	}

	gobot.Assert(t, hmc.Heading, uint16(2534))

	// when length of data returned by I2cRead is not 2
	hmc, inter = initTestMockedHMC6352Driver()

	inter.When("I2cStart", uint8(0x21))
	inter.When("I2cWrite", []byte("A")).Times(numberOfCyclesForEvery)
	inter.When("I2cRead", uint(2)).Return([]byte{99}).Times(numberOfCyclesForEvery - 1)

	gobot.Assert(t, hmc.Heading, uint16(0))

	hmc.SetInterval(1 * time.Millisecond)
	gobot.Assert(t, hmc.Start(), true)
	<-time.After(time.Duration(numberOfCyclesForEvery) * time.Millisecond)

	if ok, err := inter.Verify(); !ok {
		t.Errorf("Error:", err)
	}

	gobot.Assert(t, hmc.Heading, uint16(0))
}

func TestHMC6352DriverInit(t *testing.T) {
	hmc := initTestHMC6352Driver()

	gobot.Assert(t, hmc.Init(), true)
}

func TestHMC6352DriverHalt(t *testing.T) {
	hmc := initTestHMC6352Driver()

	gobot.Assert(t, hmc.Halt(), true)
}
