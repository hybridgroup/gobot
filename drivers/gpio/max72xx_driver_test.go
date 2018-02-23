package gpio

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*MAX72xxDriver)(nil)

// --------- HELPERS
func initTestMAX72xxDriver() (driver *MAX72xxDriver) {
	driver, _ = initTestMAX72xxDriverWithStubbedAdaptor()
	return
}

func initTestMAX72xxDriverWithStubbedAdaptor() (*MAX72xxDriver, *gpioTestAdaptor) {
	adaptor := newGpioTestAdaptor()
	return NewMAX72xxDriver(adaptor, "1", "2", "3", 1), adaptor
}

// --------- TESTS
func TestMAX72xxDriver(t *testing.T) {
	var a interface{} = initTestMAX72xxDriver()
	_, ok := a.(*MAX72xxDriver)
	if !ok {
		t.Errorf("NewMAX72xxDriver() should have returned a *MAX72xxDriver")
	}
}

func TestMAX72xxDriverStart(t *testing.T) {
	d := initTestMAX72xxDriver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestMAX72xxDriverHalt(t *testing.T) {
	d := initTestMAX72xxDriver()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestMAX72xxDriverDefaultName(t *testing.T) {
	d := initTestMAX72xxDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "MAX72xxDriver"), true)
}

func TestMAX72xxDriverSetName(t *testing.T) {
	d := initTestMAX72xxDriver()
	d.SetName("mybot")
	gobottest.Assert(t, d.Name(), "mybot")
}

