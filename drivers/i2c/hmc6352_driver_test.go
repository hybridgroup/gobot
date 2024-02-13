package i2c

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*HMC6352Driver)(nil)

func initTestHMC6352DriverWithStubbedAdaptor() (*HMC6352Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewHMC6352Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewHMC6352Driver(t *testing.T) {
	var di interface{} = NewHMC6352Driver(newI2cTestAdaptor())
	d, ok := di.(*HMC6352Driver)
	if !ok {
		require.Fail(t, "NewHMC6352Driver() should have returned a *HMC6352Driver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "HMC6352"))
	assert.Equal(t, 0x21, d.defaultAddress)
}

func TestHMC6352Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewHMC6352Driver(newI2cTestAdaptor(), WithBus(2))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
}

func TestHMC6352Start(t *testing.T) {
	d := NewHMC6352Driver(newI2cTestAdaptor())
	require.NoError(t, d.Start())
}

func TestHMC6352Halt(t *testing.T) {
	d, _ := initTestHMC6352DriverWithStubbedAdaptor()
	require.NoError(t, d.Halt())
}

func TestHMC6352Heading(t *testing.T) {
	// when len(data) is 2
	d, a := initTestHMC6352DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99, 1})
		return 2, nil
	}

	heading, _ := d.Heading()
	assert.Equal(t, uint16(2534), heading)

	// when len(data) is not 2
	d, a = initTestHMC6352DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99})
		return 1, nil
	}

	heading, err := d.Heading()
	assert.Equal(t, uint16(0), heading)
	assert.Equal(t, ErrNotEnoughBytes, err)

	// when read error
	d, a = initTestHMC6352DriverWithStubbedAdaptor()
	a.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}

	heading, err = d.Heading()
	assert.Equal(t, uint16(0), heading)
	require.ErrorContains(t, err, "read error")

	// when write error
	d, a = initTestHMC6352DriverWithStubbedAdaptor()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	heading, err = d.Heading()
	assert.Equal(t, uint16(0), heading)
	require.ErrorContains(t, err, "write error")
}
