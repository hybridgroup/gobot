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
var _ gobot.Driver = (*INA3221Driver)(nil)

func initTestINA3221DriverWithStubbedAdaptor() (*INA3221Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewINA3221Driver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewINA3221Driver(t *testing.T) {
	var di interface{} = NewINA3221Driver(newI2cTestAdaptor())
	d, ok := di.(*INA3221Driver)
	if !ok {
		t.Error("NewINA3221Driver() should return a *INA3221Driver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "INA3221"))
	assert.Equal(t, 0x40, d.defaultAddress)
}

func TestINA3221Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewINA3221Driver(newI2cTestAdaptor(), WithBus(2))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
}

func TestINA3221Start(t *testing.T) {
	d := NewINA3221Driver(newI2cTestAdaptor())
	require.NoError(t, d.Start())
}

func TestINA3221Halt(t *testing.T) {
	d, _ := initTestINA3221DriverWithStubbedAdaptor()
	require.NoError(t, d.Halt())
}

func TestINA3221GetBusVoltage(t *testing.T) {
	d, a := initTestINA3221DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		// bus voltage sensor values from 12V battery
		copy(b, []byte{0x36, 0x68})
		return 2, nil
	}

	v, err := d.GetBusVoltage(INA3221Channel1)
	assert.InDelta(t, float64(13.928), v, 0.0)
	require.NoError(t, err)
}

func TestINA3221GetBusVoltageReadError(t *testing.T) {
	d, a := initTestINA3221DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}

	_, err := d.GetBusVoltage(INA3221Channel1)
	require.ErrorContains(t, err, "read error")
}

func TestINA3221GetShuntVoltage(t *testing.T) {
	d, a := initTestINA3221DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		// shunt voltage sensor values from 12V battery
		copy(b, []byte{0x05, 0xD8})
		return 2, nil
	}

	v, err := d.GetShuntVoltage(INA3221Channel1)
	assert.InDelta(t, float64(7.48), v, 0.0)
	require.NoError(t, err)
}

func TestINA3221GetShuntVoltageReadError(t *testing.T) {
	d, a := initTestINA3221DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}

	_, err := d.GetShuntVoltage(INA3221Channel1)
	require.ErrorContains(t, err, "read error")
}

func TestINA3221GetCurrent(t *testing.T) {
	d, a := initTestINA3221DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		// shunt voltage sensor values from 12V battery
		copy(b, []byte{0x05, 0x0D8})
		return 2, nil
	}

	v, err := d.GetCurrent(INA3221Channel1)
	assert.InDelta(t, float64(74.8), v, 0.0)
	require.NoError(t, err)
}

func TestINA3221CurrentReadError(t *testing.T) {
	d, a := initTestINA3221DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}

	_, err := d.GetCurrent(INA3221Channel1)
	require.ErrorContains(t, err, "read error")
}

func TestINA3221GetLoadVoltage(t *testing.T) {
	d, a := initTestINA3221DriverWithStubbedAdaptor()
	i := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		// TODO: return test data as read from actual sensor
		copy(b, []byte{0x36, 0x68, 0x05, 0xd8}[i:i+2])
		i += 2
		return 2, nil
	}

	v, err := d.GetLoadVoltage(INA3221Channel2)
	assert.InDelta(t, float64(13.935480), v, 0.0)
	require.NoError(t, err)
}

func TestINA3221GetLoadVoltageReadError(t *testing.T) {
	d, a := initTestINA3221DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}

	_, err := d.GetLoadVoltage(INA3221Channel2)
	require.ErrorContains(t, err, "read error")
}
