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
var _ gobot.Driver = (*LIDARLiteDriver)(nil)

func initTestLIDARLiteDriver() *LIDARLiteDriver {
	d, _ := initTestLIDARLiteDriverWithStubbedAdaptor()
	return d
}

func initTestLIDARLiteDriverWithStubbedAdaptor() (*LIDARLiteDriver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewLIDARLiteDriver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewLIDARLiteDriver(t *testing.T) {
	var di interface{} = NewLIDARLiteDriver(newI2cTestAdaptor())
	d, ok := di.(*LIDARLiteDriver)
	if !ok {
		require.Fail(t, "NewLIDARLiteDriver() should have returned a *LIDARLiteDriver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "LIDARLite"))
	assert.Equal(t, 0x62, d.defaultAddress)
}

func TestLIDARLiteDriverOptions(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewLIDARLiteDriver(newI2cTestAdaptor(), WithBus(2))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
}

func TestLIDARLiteDriverStart(t *testing.T) {
	d := NewLIDARLiteDriver(newI2cTestAdaptor())
	require.NoError(t, d.Start())
}

func TestLIDARLiteDriverHalt(t *testing.T) {
	d := initTestLIDARLiteDriver()
	require.NoError(t, d.Halt())
}

func TestLIDARLiteDriverDistance(t *testing.T) {
	// when everything is happy
	d, a := initTestLIDARLiteDriverWithStubbedAdaptor()
	first := true
	a.i2cReadImpl = func(b []byte) (int, error) {
		if first {
			first = false
			copy(b, []byte{99})
			return 1, nil
		}
		copy(b, []byte{1})
		return 1, nil
	}

	distance, err := d.Distance()

	require.NoError(t, err)
	assert.Equal(t, int(25345), distance)

	// when insufficient bytes have been read
	d, a = initTestLIDARLiteDriverWithStubbedAdaptor()
	a.i2cReadImpl = func([]byte) (int, error) {
		return 0, nil
	}

	distance, err = d.Distance()
	assert.Equal(t, int(0), distance)
	assert.Equal(t, ErrNotEnoughBytes, err)

	// when read error
	d, a = initTestLIDARLiteDriverWithStubbedAdaptor()
	a.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}

	distance, err = d.Distance()
	assert.Equal(t, int(0), distance)
	require.ErrorContains(t, err, "read error")
}

func TestLIDARLiteDriverDistanceError1(t *testing.T) {
	d, a := initTestLIDARLiteDriverWithStubbedAdaptor()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	distance, err := d.Distance()
	assert.Equal(t, int(0), distance)
	require.ErrorContains(t, err, "write error")
}

func TestLIDARLiteDriverDistanceError2(t *testing.T) {
	d, a := initTestLIDARLiteDriverWithStubbedAdaptor()
	a.i2cWriteImpl = func(b []byte) (int, error) {
		if b[0] == 0x0f {
			return 0, errors.New("write error")
		}
		return len(b), nil
	}

	distance, err := d.Distance()
	assert.Equal(t, int(0), distance)
	require.ErrorContains(t, err, "write error")
}

func TestLIDARLiteDriverDistanceError3(t *testing.T) {
	d, a := initTestLIDARLiteDriverWithStubbedAdaptor()
	a.i2cWriteImpl = func(b []byte) (int, error) {
		if b[0] == 0x10 {
			return 0, errors.New("write error")
		}
		return len(b), nil
	}
	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		buf.Write([]byte{0x03})
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}

	distance, err := d.Distance()
	assert.Equal(t, int(0), distance)
	require.ErrorContains(t, err, "write error")
}
