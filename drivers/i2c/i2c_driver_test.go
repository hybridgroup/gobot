package i2c

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*Driver)(nil)

func initDriverWithStubbedAdaptor() (*Driver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewDriver(a, "I2C_BASIC", 0x15)
	return d, a
}

func initTestDriver() *Driver {
	d, _ := initDriverWithStubbedAdaptor()
	return d
}

func TestNewDriver(t *testing.T) {
	// arrange
	a := newI2cTestAdaptor()
	// act
	var di interface{} = NewDriver(a, "I2C_BASIC", 0x15)
	// assert
	d, ok := di.(*Driver)
	if !ok {
		require.Fail(t, "NewDriver() should have returned a *Driver")
	}
	assert.Contains(t, d.name, "I2C_BASIC")
	assert.Equal(t, 0x15, d.defaultAddress)
	assert.Equal(t, a, d.connector)
	assert.Nil(t, d.connection)
	require.NoError(t, d.afterStart())
	require.NoError(t, d.beforeHalt())
	assert.NotNil(t, d.Config)
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
}

func TestSetName(t *testing.T) {
	// arrange
	d := initTestDriver()
	// act
	d.SetName("TESTME")
	// assert
	assert.Equal(t, "TESTME", d.Name())
}

func TestConnection(t *testing.T) {
	// arrange
	d := initTestDriver()
	// act, assert
	assert.NotNil(t, d.Connection())
}

func TestStart(t *testing.T) {
	// arrange
	d, a := initDriverWithStubbedAdaptor()
	// act, assert
	require.NoError(t, d.Start())
	assert.Equal(t, 0x15, a.address)
}

func TestStartConnectError(t *testing.T) {
	// arrange
	d, a := initDriverWithStubbedAdaptor()
	a.Testi2cConnectErr(true)
	// act, assert
	require.ErrorContains(t, d.Start(), "Invalid i2c connection")
}

func TestHalt(t *testing.T) {
	// arrange
	d := initTestDriver()
	// act, assert
	require.NoError(t, d.Halt())
}

func TestWrite(t *testing.T) {
	// arrange
	const (
		address     = "82"
		wantAddress = uint8(0x52)
		value       = 0x25
	)
	d, a := initDriverWithStubbedAdaptor()
	_ = d.Start()
	// prepare all writes
	numCallsWrite := 0
	a.i2cWriteImpl = func([]byte) (int, error) {
		numCallsWrite++
		return 0, nil
	}
	// act
	err := d.Write(address, value)
	// assert
	require.NoError(t, err)
	assert.Equal(t, 1, numCallsWrite)
	assert.Equal(t, wantAddress, a.written[0])
	assert.Equal(t, uint8(value), a.written[1])
}

func TestRead(t *testing.T) {
	// arrange
	const (
		address     = "83"
		wantAddress = uint8(0x53)
		want        = uint8(0x44)
	)
	d, a := initDriverWithStubbedAdaptor()
	_ = d.Start()
	// prepare all writes
	numCallsWrite := 0
	a.i2cWriteImpl = func(b []byte) (int, error) {
		numCallsWrite++
		return 0, nil
	}
	// prepare all reads
	numCallsRead := 0
	a.i2cReadImpl = func(b []byte) (int, error) {
		numCallsRead++
		b[0] = want
		return len(b), nil
	}
	// act
	val, err := d.Read(address)
	// assert
	require.NoError(t, err)
	assert.Equal(t, int(want), val)
	assert.Equal(t, 1, numCallsWrite)
	assert.Equal(t, wantAddress, a.written[0])
	assert.Equal(t, 1, numCallsRead)
}
