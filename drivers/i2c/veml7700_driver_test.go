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
var _ gobot.Driver = (*VEML7700Driver)(nil)

func initTestVEML7700DriverWithStubbedAdaptor() (*VEML7700Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewVEML7700Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewVEML7700Driver(t *testing.T) {
	var di interface{} = NewVEML7700Driver(newI2cTestAdaptor())
	d, ok := di.(*VEML7700Driver)
	if !ok {
		t.Errorf("NewVEML7700Driver() should have returned a *VEML7700Driver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "VEML7700"))
	assert.Equal(t, 0x10, d.defaultAddress)
}

func TestVEML7700Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewVEML7700Driver(newI2cTestAdaptor(), WithBus(2))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
}

func TestVEML7700Start(t *testing.T) {
	d := NewVEML7700Driver(newI2cTestAdaptor())
	require.NoError(t, d.Start())
}

func TestVEML7700Halt(t *testing.T) {
	d, _ := initTestVEML7700DriverWithStubbedAdaptor()
	require.NoError(t, d.Halt())
}

func TestVEML7700NullLux(t *testing.T) {
	d, _ := initTestVEML7700DriverWithStubbedAdaptor()
	lux, _ := d.Lux()
	assert.Equal(t, 0.0, lux)
}

func TestVEML7700Lux(t *testing.T) {
	d, a := initTestVEML7700DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{0x05, 0xb0})
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}

	lux, _ := d.Lux()
	assert.Equal(t, 20764.108799999998, lux)
}

func TestVEML7700LuxError(t *testing.T) {
	d, a := initTestVEML7700DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("wrong number of bytes read")
	}

	_, err := d.Lux()
	require.ErrorContains(t, err, "wrong number of bytes read")
}
