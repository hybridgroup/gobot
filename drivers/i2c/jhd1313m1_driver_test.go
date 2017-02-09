package i2c

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*JHD1313M1Driver)(nil)

// --------- HELPERS
func initTestJHD1313M1Driver() (driver *JHD1313M1Driver) {
	driver, _ = initTestJHD1313M1DriverWithStubbedAdaptor()
	return
}

func initTestJHD1313M1DriverWithStubbedAdaptor() (*JHD1313M1Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewJHD1313M1Driver(adaptor), adaptor
}

// --------- TESTS

func TestNewJHD1313M1Driver(t *testing.T) {
	// Does it return a pointer to an instance of JHD1313M1Driver?
	var mpl interface{} = NewJHD1313M1Driver(newI2cTestAdaptor())
	_, ok := mpl.(*JHD1313M1Driver)
	if !ok {
		t.Errorf("NewJHD1313M1Driver() should have returned a *JHD1313M1Driver")
	}
}

// Methods
func TestJHD1313M1Driver(t *testing.T) {
	jhd := initTestJHD1313M1Driver()

	gobottest.Refute(t, jhd.Connection(), nil)
	gobottest.Assert(t, strings.HasPrefix(jhd.Name(), "JHD1313M1"), true)
}
