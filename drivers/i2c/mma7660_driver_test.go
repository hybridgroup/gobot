package i2c

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*MMA7660Driver)(nil)

func initTestMMA7660DriverWithStubbedAdaptor() (*MMA7660Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewMMA7660Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewMMA7660Driver(t *testing.T) {
	var di interface{} = NewMMA7660Driver(newI2cTestAdaptor())
	d, ok := di.(*MMA7660Driver)
	if !ok {
		require.Fail(t, "NewMMA7660Driver() should have returned a *MMA7660Driver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "MMA7660"))
	assert.Equal(t, 0x4c, d.defaultAddress)
}

func TestMMA7660Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewMMA7660Driver(newI2cTestAdaptor(), WithBus(2))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
}

func TestMMA7660Start(t *testing.T) {
	d := NewMMA7660Driver(newI2cTestAdaptor())
	require.NoError(t, d.Start())
}

func TestMMA7660Halt(t *testing.T) {
	d, _ := initTestMMA7660DriverWithStubbedAdaptor()
	require.NoError(t, d.Halt())
}

func TestMMA7660Acceleration(t *testing.T) {
	d, _ := initTestMMA7660DriverWithStubbedAdaptor()
	x, y, z := d.Acceleration(21.0, 21.0, 21.0)
	assert.InDelta(t, 1.0, x, 0.0)
	assert.InDelta(t, 1.0, y, 0.0)
	assert.InDelta(t, 1.0, z, 0.0)
}

func TestMMA7660NullXYZ(t *testing.T) {
	d, _ := initTestMMA7660DriverWithStubbedAdaptor()

	x, y, z, _ := d.XYZ()
	assert.InDelta(t, 0.0, x, 0.0)
	assert.InDelta(t, 0.0, y, 0.0)
	assert.InDelta(t, 0.0, z, 0.0)
}

func TestMMA7660XYZ(t *testing.T) {
	d, a := initTestMMA7660DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{0x11, 0x12, 0x13})
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}

	x, y, z, _ := d.XYZ()
	assert.InDelta(t, 17.0, x, 0.0)
	assert.InDelta(t, 18.0, y, 0.0)
	assert.InDelta(t, 19.0, z, 0.0)
}

func TestMMA7660XYZError(t *testing.T) {
	d, a := initTestMMA7660DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}

	x, y, z, err := d.XYZ()
	require.ErrorContains(t, err, "read error")
	assert.InDelta(t, 0.0, x, 0.0)
	assert.InDelta(t, 0.0, y, 0.0)
	assert.InDelta(t, 0.0, z, 0.0)
}

func TestMMA7660XYZNotReady(t *testing.T) {
	d, a := initTestMMA7660DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{0x40, 0x40, 0x40})
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}

	x, y, z, err := d.XYZ()
	assert.Equal(t, ErrNotReady, err)
	assert.InDelta(t, 0.0, x, 0.0)
	assert.InDelta(t, 0.0, y, 0.0)
	assert.InDelta(t, 0.0, z, 0.0)
}
