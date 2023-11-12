package i2c

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2"
)

var (
	_ gobot.Driver = (*GroveLcdDriver)(nil)
	_ gobot.Driver = (*GroveAccelerometerDriver)(nil)
)

func initTestGroveLcdDriver() *GroveLcdDriver {
	d, _ := initGroveLcdDriverWithStubbedAdaptor()
	return d
}

func initGroveLcdDriverWithStubbedAdaptor() (*GroveLcdDriver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewGroveLcdDriver(adaptor), adaptor
}

func initTestGroveAccelerometerDriver() *GroveAccelerometerDriver {
	d, _ := initGroveAccelerometerDriverWithStubbedAdaptor()
	return d
}

func initGroveAccelerometerDriverWithStubbedAdaptor() (*GroveAccelerometerDriver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewGroveAccelerometerDriver(adaptor), adaptor
}

func TestGroveLcdDriverName(t *testing.T) {
	g := initTestGroveLcdDriver()
	assert.NotNil(t, g.Connection())
	assert.True(t, strings.HasPrefix(g.Name(), "JHD1313M1"))
}

func TestLcdDriverWithAddress(t *testing.T) {
	adaptor := newI2cTestAdaptor()
	g := NewGroveLcdDriver(adaptor, WithAddress(0x66))
	assert.Equal(t, 0x66, g.GetAddressOrDefault(0x33))
}

func TestGroveAccelerometerDriverName(t *testing.T) {
	g := initTestGroveAccelerometerDriver()
	assert.NotNil(t, g.Connection())
	assert.True(t, strings.HasPrefix(g.Name(), "MMA7660"))
}

func TestGroveAccelerometerDriverWithAddress(t *testing.T) {
	adaptor := newI2cTestAdaptor()
	g := NewGroveAccelerometerDriver(adaptor, WithAddress(0x66))
	assert.Equal(t, 0x66, g.GetAddressOrDefault(0x33))
}
