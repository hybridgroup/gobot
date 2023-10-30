package gpio

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*MAX7219Driver)(nil)

// --------- HELPERS
func initTestMAX7219Driver() (driver *MAX7219Driver) {
	driver, _ = initTestMAX7219DriverWithStubbedAdaptor()
	return
}

func initTestMAX7219DriverWithStubbedAdaptor() (*MAX7219Driver, *gpioTestAdaptor) {
	adaptor := newGpioTestAdaptor()
	return NewMAX7219Driver(adaptor, "1", "2", "3", 1), adaptor
}

// --------- TESTS
func TestMAX7219Driver(t *testing.T) {
	var a interface{} = initTestMAX7219Driver()
	_, ok := a.(*MAX7219Driver)
	if !ok {
		t.Errorf("NewMAX7219Driver() should have returned a *MAX7219Driver")
	}
}

func TestMAX7219DriverStart(t *testing.T) {
	d := initTestMAX7219Driver()
	assert.NoError(t, d.Start())
}

func TestMAX7219DriverHalt(t *testing.T) {
	d := initTestMAX7219Driver()
	assert.NoError(t, d.Halt())
}

func TestMAX7219DriverDefaultName(t *testing.T) {
	d := initTestMAX7219Driver()
	assert.True(t, strings.HasPrefix(d.Name(), "MAX7219Driver"))
}

func TestMAX7219DriverSetName(t *testing.T) {
	d := initTestMAX7219Driver()
	d.SetName("mybot")
	assert.Equal(t, "mybot", d.Name())
}
