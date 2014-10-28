package i2c

import (
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
)

// --------- HELPERS
func initTestHMC6352Driver() (driver *HMC6352Driver) {
	driver, _ = initTestHMC6352DriverWithStubbedAdaptor()
	return
}

func initTestHMC6352DriverWithStubbedAdaptor() (*HMC6352Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor("adaptor")
	return NewHMC6352Driver(adaptor, "bot"), adaptor
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
func TestHMC6352DriverStart(t *testing.T) {
	sem := make(chan bool)
	// when len(data) is 2
	hmc, adaptor := initTestHMC6352DriverWithStubbedAdaptor()

	adaptor.i2cReadImpl = func() []byte {
		return []byte{99, 1}
	}

	numberOfCyclesForEvery := 3

	// Make sure "Heading" is set to 0 so that we assert
	// its new value after executing "Start()"
	gobot.Assert(t, hmc.Heading, uint16(0))

	hmc.SetInterval(1 * time.Millisecond)
	gobot.Assert(t, hmc.Start(), true)
	go func() {
		for {
			<-time.After(time.Duration(numberOfCyclesForEvery) * time.Millisecond)
			if hmc.Heading == uint16(2534) {
				sem <- true
			}
		}
	}()

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Heading not read correctly")
	}

	// when len(data) is not 2
	hmc, adaptor = initTestHMC6352DriverWithStubbedAdaptor()

	adaptor.i2cReadImpl = func() []byte {
		return []byte{99}
	}

	hmc.SetInterval(1 * time.Millisecond)
	gobot.Assert(t, hmc.Start(), true)
	go func() {
		for {
			<-time.After(time.Duration(numberOfCyclesForEvery) * time.Millisecond)
			if hmc.Heading == uint16(0) {
				sem <- true
			}
		}
	}()

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Heading not read correctly")
	}
}

func TestHMC6352DriverInit(t *testing.T) {
	hmc := initTestHMC6352Driver()

	gobot.Assert(t, hmc.Init(), true)
}

func TestHMC6352DriverHalt(t *testing.T) {
	hmc := initTestHMC6352Driver()

	gobot.Assert(t, hmc.Halt(), true)
}
