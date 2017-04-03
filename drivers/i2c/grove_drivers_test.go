package i2c

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*GroveLcdDriver)(nil)
var _ gobot.Driver = (*GroveAccelerometerDriver)(nil)

func initTestGroveLcdDriver() (driver *GroveLcdDriver) {
	driver, _ = initGroveLcdDriverWithStubbedAdaptor()
	return
}

func initGroveLcdDriverWithStubbedAdaptor() (*GroveLcdDriver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewGroveLcdDriver(adaptor), adaptor
}

func initTestGroveAccelerometerDriver() (driver *GroveAccelerometerDriver) {
	driver, _ = initGroveAccelerometerDriverWithStubbedAdaptor()
	return
}

func initGroveAccelerometerDriverWithStubbedAdaptor() (*GroveAccelerometerDriver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewGroveAccelerometerDriver(adaptor), adaptor
}

func TestGroveLcdDriverName(t *testing.T) {
	g := initTestGroveLcdDriver()
	gobottest.Refute(t, g.Connection(), nil)
	gobottest.Assert(t, strings.HasPrefix(g.Name(), "JHD1313M1"), true)
}

func TestLcdDriverWithAddress(t *testing.T) {
	adaptor := newI2cTestAdaptor()
	g := NewGroveLcdDriver(adaptor, WithAddress(0x66))
	gobottest.Assert(t, g.GetAddressOrDefault(0x33), 0x66)
}

func TestGroveAccelerometerDriverName(t *testing.T) {
	g := initTestGroveAccelerometerDriver()
	gobottest.Refute(t, g.Connection(), nil)
	gobottest.Assert(t, strings.HasPrefix(g.Name(), "MMA7660"), true)
}

func TestGroveAccelerometerDriverWithAddress(t *testing.T) {
	adaptor := newI2cTestAdaptor()
	g := NewGroveAccelerometerDriver(adaptor, WithAddress(0x66))
	gobottest.Assert(t, g.GetAddressOrDefault(0x33), 0x66)
}
